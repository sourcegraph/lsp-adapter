// Command solargraph-wrapper starts up solargraph in the workspace root
// (after reading the initialize request). It also proxies the TCP connection
// of solargraph over stdin/stdout.
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

// extractRootURIFromRequest will take the JSONRPC2 initialize request and
// return the value of rootURI.
func extractRootURIFromRequest(v interface{}) (string, error) {
	// minimal json to get at rootUri
	var msg struct {
		Method string `json:"method"`
		Params struct {
			RootURI string `json:"rootUri"`
		} `json:"params"`
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(b, &msg); err != nil {
		return "", err
	}

	if msg.Method != "initialize" {
		return "", fmt.Errorf("expected first message to be initialize, got %s", msg.Method)
	}

	return msg.Params.RootURI, nil
}

// WriteObject proxies WriteObject to the underlying stream. It treats the
// first request specially and assumes it is an LSP initialize request. It
// uses that to initialize the underlying stream via Connect.
func (t *lazyObjectStream) WriteObject(v interface{}) error {
	t.once.Do(func() {
		defer close(t.ready)
		rootURI, err := extractRootURIFromRequest(v)
		if err != nil {
			t.err = err
			return
		}
		t.stream, t.err = t.Connect(rootURI)
	})
	if t.err != nil {
		return t.err
	}
	return t.stream.WriteObject(v)
}

// ReadObject proxies ReadObject to the underlying stream.
func (t *lazyObjectStream) ReadObject(v interface{}) error {
	// we wait for the initialize request (via WriteObject) to happen first.
	<-t.ready
	if t.err != nil {
		return t.err
	}
	return t.stream.ReadObject(v)
}

// Close closes the underlying stream. If the stream has not been created, it
// prevents future requests from opening the stream.
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

// startSolargraph starts up a TCP socket based solargraph server. It parses
// its stderr to discover the port it is running on. On success it returns a
// close function and port.
func startSolargraph(ctx context.Context, dir string) (func() error, int, error) {
	var (
		stderr io.Reader
		err    error
	)
	cmd := exec.CommandContext(ctx, "solargraph", "socket", "--port", "0")
	cmd.Dir = dir
	cmd.Stdout = os.Stderr
	stderr, err = cmd.StderrPipe()
	if err != nil {
		return nil, 0, err
	}
	// We want to parse stderr, but use a teereader so the user can also see
	// stderr.
	stderr = io.TeeReader(stderr, os.Stderr)
	if err := cmd.Start(); err != nil {
		return nil, 0, err
	}
	closer := func() error {
		// best-effort
		cmd.Process.Kill()
		cmd.Wait()
		return nil
	}

	// We are looking for \bPORT=(\d+)\b. We can achieve this nicely via using
	// a Scanner splitting on words.
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, "PORT=") {
			// We are done with scanning, but we need to keep reading stderr
			// to prevent it blocking the process.
			go io.Copy(ioutil.Discard, stderr)

			port, err := strconv.Atoi(word[len("PORT="):])
			if err != nil {
				break
			}
			return closer, port, nil
		}
	}
	// if we get to this point we didn't find the port. cleanup
	closer()
	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}
	return nil, 0, errors.New("did not find port in stderr")
}

type readWriteCloser struct {
	Closers []func() error
	io.ReadWriter
}

func (c *readWriteCloser) Close() error {
	for _, fn := range c.Closers {
		fn()
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

			closer, port, err := startSolargraph(ctx, filepath.FromSlash(u.Path))
			if err != nil {
				return nil, err
			}

			conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(port))
			if err != nil {
				closer()
				return nil, err
			}

			rwc := &readWriteCloser{
				Closers:    []func() error{conn.Close, closer},
				ReadWriter: conn,
			}

			return jsonrpc2.NewBufferedStream(rwc, jsonrpc2.VSCodeObjectCodec{}), nil
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
