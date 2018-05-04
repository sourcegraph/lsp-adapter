// Command lsp-record provides an lsp proxy which records the session. It is intended
package main

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

// Request is the subset of a JSONRPC2 request payload we want to record
type Request struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params,omitempty"`
}

// Encoder wraps an io.Writer, but additionally provides Encode.
type Encoder struct {
	io.Writer
	enc *json.Encoder
}

// Encode writes the JSON encoding of v to the stream.
func (e *Encoder) Encode(v interface{}) error {
	if e.enc == nil {
		e.enc = json.NewEncoder(e.Writer)
		e.enc.SetIndent("", "  ")
	} else {
		if _, err := e.Writer.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	return e.enc.Encode(v)
}

func writeJSONRPC2Requests(r io.Reader, w io.Writer) error {
	stream := bufio.NewReader(r)
	codec := jsonrpc2.VSCodeObjectCodec{}
	enc := &Encoder{Writer: w}

	for {
		var req Request
		if err := codec.ReadObject(stream, &req); err != nil {
			return err
		}
		if req.Method != "" {
			if err := enc.Encode(req); err != nil {
				return err
			}
		}
	}
}

func retryDial(network, address string) (net.Conn, error) {
	conn, err := net.DialTimeout(network, address, time.Second)
	for i := 0; err != nil && i < 5; i++ {
		time.Sleep(time.Second)
		conn, err = net.DialTimeout(network, address, time.Second)
	}
	return conn, err
}

// massageGitHubArchive rewrites filenames to match what we expect. GitHub
// archives always have a top-level dir, so strip it out.
//
// before: dockerfile-language-server-nodejs-3083f51108b5e5ddfd440e6fe3da415d10b9c69c/src/server.ts
// after:  /src/server.ts
func massageGitHubArchive(r *zip.ReadCloser) {
	for i, file := range r.File {
		r.File[i].Name = file.Name[strings.Index(file.Name, "/"):]
	}
}

func fetchArchiveForRootURI(originalRootURI string) (*zip.ReadCloser, error) {
	dst := filepath.Join(os.TempDir(), "lsp-record", url.QueryEscape(originalRootURI)+".zip")
	if r, err := zip.OpenReader(dst); err == nil {
		massageGitHubArchive(r)
		return r, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	u, err := url.Parse(originalRootURI)
	if err != nil {
		return nil, err
	}
	if u.Host != "github.com" {
		return nil, errors.Errorf("Unsupported originalRootUri %s (only github supported)", originalRootURI)
	}

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return nil, err
	}

	repo := path.Join(u.Host, strings.TrimPrefix(u.Path, ".git"))
	rev := u.RawQuery
	url := fmt.Sprintf("https://codeload.%s/zip/%s", repo, rev)

	log.Println("fetching", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dstPart := dst + ".part"
	fd, err := os.Create(dstPart)
	if err != nil {
		return nil, err
	}
	defer os.Remove(dstPart)

	_, err = io.Copy(fd, resp.Body)
	if err2 := fd.Close(); err2 != nil {
		return nil, err2
	}
	if err != nil {
		return nil, err
	}

	if err := os.Rename(dstPart, dst); err != nil {
		return nil, err
	}
	r, err := zip.OpenReader(dst)
	if err != nil {
		return nil, err
	}
	massageGitHubArchive(r)
	return r, nil
}

type jsonrpc2HandlerFunc func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error)

func newVFSHandler(ar *zip.ReadCloser) jsonrpc2HandlerFunc {
	return func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		switch req.Method {
		case "workspace/xfiles":
			var results []lsp.TextDocumentIdentifier
			for _, f := range ar.File {
				results = append(results, lsp.TextDocumentIdentifier{URI: lsp.DocumentURI("file://" + f.Name)})
			}
			return results, nil
		case "textDocument/xcontent":
			var params struct {
				TextDocument lsp.TextDocumentIdentifier `json:"textDocument"`
			}
			if err := json.Unmarshal(*req.Params, &params); err != nil {
				return nil, err
			}
			u, err := url.Parse(string(params.TextDocument.URI))
			if err != nil {
				return nil, err
			}
			for _, f := range ar.File {
				if f.Name == u.Path {
					rc, err := f.Open()
					if err != nil {
						return nil, err
					}
					defer rc.Close()
					b, err := ioutil.ReadAll(rc)
					if err != nil {
						return nil, err
					}
					return lsp.TextDocumentItem{
						URI:  params.TextDocument.URI,
						Text: string(b),
					}, nil
				}
			}
			msg := fmt.Sprintf("URI %s does not exist", params.TextDocument.URI)
			log.Println(msg)
			return nil, &jsonrpc2.Error{
				Code:    jsonrpc2.CodeInvalidParams,
				Message: msg,
			}
		}

		log.Println("ignoring server->client request:", req.Method)
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("client handler: method not found: %q", req.Method)}
	}
}

func record() error {
	lis, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		return err
	}
	defer lis.Close()

	src, err := lis.Accept()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := retryDial("tcp", "127.0.0.1:8080")
	if err != nil {
		return err
	}
	defer dst.Close()

	done := make(chan error, 3)
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, os.Kill)
		<-stop
		done <- nil
	}()
	go func() {
		// src -> dst
		pr, pw := io.Pipe()
		go writeJSONRPC2Requests(pr, os.Stdout)
		_, err := io.Copy(dst, io.TeeReader(src, pw))
		pw.CloseWithError(err)
		done <- err
	}()
	go func() {
		// dst <- src
		_, err := io.Copy(src, dst)
		done <- err
	}()

	err = <-done
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func test() error {
	var mu sync.Mutex
	var vfsHandler jsonrpc2HandlerFunc
	handle := func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		mu.Lock()
		h := vfsHandler
		mu.Unlock()
		if h == nil {
			return nil, errors.Errorf("archive has not been fetched")
		}
		return h(ctx, conn, req)
	}

	conn, err := retryDial("tcp", "127.0.0.1:8080")
	if err != nil {
		return err
	}
	defer conn.Close()

	c := jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}),
		jsonrpc2.AsyncHandler(jsonrpc2.HandlerWithError(handle)))
	defer c.Close()

	dec := json.NewDecoder(os.Stdin)
	enc := &Encoder{Writer: os.Stdout}
	for {
		var req Request
		if err := dec.Decode(&req); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if req.Method == "initialize" {
			var params struct {
				OriginalRootURI string `json:"originalRootUri"`
			}
			if err := json.Unmarshal(*req.Params, &params); err != nil {
				return err
			}
			reader, err := fetchArchiveForRootURI(params.OriginalRootURI)
			if err != nil {
				return err
			}
			defer reader.Close()
			mu.Lock()
			vfsHandler = newVFSHandler(reader)
			mu.Unlock()
		}

		var res interface{}
		if err := c.Call(context.Background(), req.Method, req.Params, &res); err != nil {
			return err
		}

		if err := enc.Encode(res); err != nil {
			return err
		}
	}

	return nil
}

func mainErr() error {
	if len(os.Args) != 3 || (os.Args[1] != "record" && os.Args[1] != "test") {
		return errors.Errorf("USAGE: %s [record|test] language", os.Args[0])
	}
	action := os.Args[1]
	lang := filepath.Base(os.Args[2])
	image := fmt.Sprintf("sgtest/codeintel-%s", lang)
	dockerfile := filepath.Join("dockerfiles", lang, "Dockerfile")
	if _, err := os.Stat(dockerfile); os.IsNotExist(err) {
		return errors.Errorf("%s could not be found. Ensure you are running from github.com/sourcegraph/lsp-adapter directory and that %s integration exists", dockerfile, lang)
	}

	cmd := exec.Command("docker", "build", "-t", image, "-f", dockerfile, ".")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Println(strings.Join(cmd.Args, " "))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("docker", "run", "--rm=true", "-p", "8080:8080", image)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Println(strings.Join(cmd.Args, " "))
	if err := cmd.Start(); err != nil {
		return err
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	if action == "record" {
		return record()
	} else {
		return test()
	}
}

func main() {
	err := mainErr()
	if err != nil {
		log.Fatal(err)
	}
}
