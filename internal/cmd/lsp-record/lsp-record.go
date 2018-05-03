// Command lsp-record provides an lsp proxy which records the session. It is intended
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sourcegraph/jsonrpc2"
)

func writeJSONRPC2Requests(r io.Reader, w io.Writer) error {
	stream := bufio.NewReader(r)
	codec := jsonrpc2.VSCodeObjectCodec{}

	for {
		var req struct {
			Method string           `json:"method"`
			Params *json.RawMessage `json:"params,omitempty"`
		}
		if err := codec.ReadObject(stream, &req); err != nil {
			return err
		}
		if req.Method != "" {
			b, err := json.Marshal(req)
			if err != nil {
				return err
			}
			if _, err := w.Write(b); err != nil {
				return err
			}
			if _, err := w.Write([]byte{'\n'}); err != nil {
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

func mainErr() error {
	if len(os.Args) != 2 {
		return errors.Errorf("USAGE: %s language", os.Args[0])
	}
	lang := filepath.Base(os.Args[1])
	image := fmt.Sprintf("sgtest/codeintel-%s", lang)
	dockerfile := filepath.Join("dockerfiles", lang, "Dockerfile")
	if _, err := os.Stat(dockerfile); os.IsNotExist(err) {
		return errors.Errorf("%s could not be found. Ensure you are running from github.com/sourcegraph/lsp-adapter directory and that %s integration exists", dockerfile, lang)
	}

	cmd := exec.Command("docker", "build", "-t", image, "-f", dockerfile, ".")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Println("running ", strings.Join(cmd.Args, " "))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("docker", "run", "--rm=true", "-p", "8080:8080", image)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	log.Println("running ", strings.Join(cmd.Args, " "))
	if err := cmd.Start(); err != nil {
		return err
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

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

func main() {
	err := mainErr()
	if err != nil {
		log.Fatal(err)
	}
}
