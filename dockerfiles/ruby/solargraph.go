package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
)

// lazyObjectStream wraps a jsonrpc2.ObjectStream to only start the connection
// once we have received an initialize request.
type lazyObjectStream struct {
	Connect func(rootURI string) (jsonrpc2.ObjectStream, error)

	once   *sync.Once    // protects writes to stream, err and closing ready.
	ready  chan struct{} // protects reads to stream and err.
	stream jsonrpc2.ObjectStream
	err    error
}

func (t *lazyObjectStream) WriteObject(v interface{}) error {
	t.once.Do(func() {
		defer close(t.ready)

		// minimal json to get at rootUri
		var msg struct {
			Method string `json:"method"`
			Params struct {
				RootURI string `json:"rootUri"`
			} `json:"params"`
		}
		b, err := json.Marshal(v)
		if err != nil {
			t.err = err
			return
		}
		if err = json.Unmarshal(b, &msg); err != nil {
			t.err = err
			return
		}

		if msg.Method != "initialize" {
			t.err = fmt.Errorf("expected first message to be initialize, got %s", msg.Method)
			return
		}

		t.stream, t.err = t.Connect(msg.Params.RootURI)
	})
	if t.err != nil {
		return t.err
	}
	return t.stream.WriteObject(v)
}

func (t *lazyObjectStream) ReadObject(v interface{}) error {
	// we wait for the initialize request (via WriteObject) to happen first.
	<-t.ready
	if t.err != nil {
		return t.err
	}
	return t.stream.ReadObject(v)
}

func (t *lazyObjectStream) Close() error {
	t.once.Do(func() {
		t.err = jsonrpc2.ErrClosed
		close(t.ready)
	})
	if t.err != nil {
		return nil
	}
	return t.stream.Close()
}

func startSolargraph(ctx context.Context, dir string) (*exec.Cmd, int, error) {
	var (
		stderr io.Reader
		err    error
	)
	cmd := exec.CommandContext(ctx, "solargraph", "socket", "--port", "0")
	cmd.Dir = dir // set CWD env var?
	cmd.Stdout = os.Stderr
	stderr, err = cmd.StderrPipe()
	if err != nil {
		return nil, 0, err
	}
	stderr = io.TeeReader(stderr, os.Stderr)
	if err := cmd.Start(); err != nil {
		return nil, 0, err
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, "PORT=") {
			// read the rest of stderr to keep the tee reader going
			go io.Copy(ioutil.Discard, stderr)

			port, err := strconv.Atoi(word[len("PORT="):])
			if err != nil {
				break
			}
			return cmd, port, nil
		}
	}
	// if we get to this point we didn't find the port. cleanup
	cmd.Process.Kill()
	cmd.Wait()
	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}
	return nil, 0, errors.New("did not find port in stderr")
}

type cmd struct {
	*exec.Cmd

	// Reader and Writer do not need to be Closers since they are StdoutPipe
	// and StdinPipe respectively. Both of those will be closed by Cmd.Wait.
	io.Reader
	io.Writer
}

func (c *cmdRWCloser) Close() error {
	if err := c.Cmd.Process.Kill(); err != nil {
		return errors.Wrap(err, "unable to kill process during cmdRWCloser.Close()")
	}

	if err := c.Cmd.Wait(); err != nil {
		return errors.Wrap(err, "unable to wait on cmd to finish during cmdRWCloser.Close()")
	}

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := &lazyObjectStream{
		once:  new(sync.Once),
		ready: make(chan struct{}),

		Connect: func(rootURI string) (jsonrpc2.ObjectStream, error) {
			u, err := url.Parse(rootURI)
			if err != nil {
				return nil, err
			}
			if u.Scheme != "file" {
				return nil, fmt.Errorf("rootURI %s does not have a file scheme", rootURI)
			}

			cmd, port, err := startSolargraph(ctx, filepath.FromSlash(u.Path))
			if err != nil {
				return nil, err
			}

			conn, err := net.Dial("tcp", "127.0.0.1:"+stconv.Itoa(port))
			// TODO finish
		},
	}

	done := make(chan int, 2)
	stdin := bufio.NewReader(os.Stdin)
	codec := jsonrpc2.VSCodeObjectCodec{}
	go func() {
		for {
			var b json.RawMessage
			if err := codec.ReadObject(stdin, &b); err != nil {
				break
			}
			if err := stream.WriteObject(&b); err != nil {
				break
			}
		}
		done <- 1
	}()
	go func() {
		for {
			var b json.RawMessage
			if err := stream.ReadObject(&b); err != nil {
				break
			}
			if err := codec.WriteObject(os.Stdout, &b); err != nil {
				break
			}
		}
		done <- 1
	}()

	<-done
	stream.Close()
}
