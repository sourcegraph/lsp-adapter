package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

var (
	proxyAddr = flag.String("proxyAddress", "127.0.0.1:8080", "proxy server listen address (tcp)")
	lspAddr   = flag.String("lspAddress", "", "language server listen address (tcp)")
	cacheDir  = flag.String("cacheDirectory", filepath.Join(os.TempDir(), "proxy-cache"), "cache directory location")
	trace     = flag.Bool("trace", false, "trace logs to stderr")
)

type cloneProxy struct {
	client *jsonrpc2.Conn // connection to the browser
	server *jsonrpc2.Conn // connection to the language server

	ready chan struct{} // barrier to block handling requests until proxy is fully initalized
	id    uuid.UUID     // unique ID for this session
	ctx   context.Context
}

func (p *cloneProxy) start() {
	close(p.ready)
}

type jsonrpc2HandlerFunc func(context.Context, *jsonrpc2.Conn, *jsonrpc2.Request)

func (h jsonrpc2HandlerFunc) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	h(ctx, conn, req)
}

func main() {
	flag.Parse()
	log.SetFlags(log.Flags() | log.Lshortfile)

	lspBin := flag.Args()
	if len(lspBin) > 0 && *lspAddr != "" {
		log.Fatalf("Both an LSP command (arguments %v) and an LSP address (-lspAddress %v) are specified. Please only specify one", lspBin, *lspAddr)
	}
	if len(lspBin) == 0 && *lspAddr == "" {
		log.Fatal("Specify either an LSP command (positional arguments) or an LSP address (-lspAddress flag)")
	}

	lis, err := net.Listen("tcp", *proxyAddr)
	if err != nil {
		err = errors.Wrap(err, "setting up proxy listener failed")
		log.Fatal(err)
	}
	defer lis.Close()
	log.Printf("CloneProxy: accepting connections at %s", lis.Addr())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go trapSignalsForShutdown(func() {
		cancel()
		lis.Close()
	})

	var wg sync.WaitGroup
	for {
		clientNetConn, err := lis.Accept()
		if err != nil {
			if ctx.Err() != nil { // shutdown
				break
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Println("error when accepting client connection: ", err.Error())
				continue
			}
			log.Fatal(err)
		}

		wg.Add(1)
		go func(clientNetConn net.Conn) {
			defer wg.Done()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			var lsConn io.ReadWriteCloser
			if *lspAddr != "" {
				lsConn, err = dialLanguageServer(ctx, *lspAddr)
				if err != nil {
					log.Println("dialing language server failed", err.Error())
					return
				}

			} else {
				lsConn, err = mkStdIoLSConn(ctx, lspBin[0], lspBin[1:]...)
				if err != nil {
					log.Println("connecting to language server over stdio failed", err.Error())
					return
				}
			}

			proxy := &cloneProxy{
				ready: make(chan struct{}),
				ctx:   ctx,
				id:    uuid.New(),
			}

			var serverConnOpts []jsonrpc2.ConnOpt
			if *trace {
				serverConnOpts = append(serverConnOpts, jsonrpc2.LogMessages(log.New(os.Stderr, fmt.Sprintf("TRACE %s ", proxy.id.String()), log.Ltime)))
			}
			proxy.client = jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(clientNetConn, jsonrpc2.VSCodeObjectCodec{}), jsonrpc2.AsyncHandler(jsonrpc2HandlerFunc(proxy.handleClientRequest)))
			proxy.server = jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(lsConn, jsonrpc2.VSCodeObjectCodec{}), jsonrpc2.AsyncHandler(jsonrpc2HandlerFunc(proxy.handleServerRequest)), serverConnOpts...)

			proxy.start()

			// When one side of the connection disconnects, close the other side.
			select {
			case <-proxy.client.DisconnectNotify():
				proxy.server.Close()
			case <-proxy.server.DisconnectNotify():
				proxy.client.Close()
			}
		}(clientNetConn)
	}

	wg.Wait()
}

// dialLanguageServer creates a connection to the language server specified at addr.
func dialLanguageServer(ctx context.Context, addr string) (net.Conn, error) {
	if addr == "" {
		return nil, errors.Errorf("language server not found at addr: %s", addr)
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
}

func (p *cloneProxy) handleServerRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	<-p.ready

	rTripper := roundTripper{
		req: req,

		src:  p.server,
		dest: p.client,

		updateURIFromSrc:  p.serverToClientURI,
		updateURIFromDest: p.clientToServerURI,
	}

	if err := rTripper.roundTrip(ctx); err != nil {
		log.Println("CloneProxy.handleServerRequest(): roundTrip failed", err)
	}
}

func (p *cloneProxy) handleClientRequest(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	<-p.ready

	if req.Method == "initialize" {
		if err := p.cloneWorkspaceToCache(); err != nil {
			log.Println("CloneProxy.handleClientRequest(): cloning workspace failed during initialize", err)
			return
		}
	}

	rTripper := roundTripper{
		req: req,

		src:  p.client,
		dest: p.server,

		updateURIFromSrc:  p.clientToServerURI,
		updateURIFromDest: p.serverToClientURI,
	}

	if err := rTripper.roundTrip(ctx); err != nil {
		log.Println("CloneProxy.handleClientRequest(): roundTrip failed", err)
	}
}

type roundTripper struct {
	req *jsonrpc2.Request

	src  *jsonrpc2.Conn
	dest *jsonrpc2.Conn

	updateURIFromSrc  func(lsp.DocumentURI) lsp.DocumentURI
	updateURIFromDest func(lsp.DocumentURI) lsp.DocumentURI
}

// roundTrip passes requests from one side of the connection to the other.
func (r *roundTripper) roundTrip(ctx context.Context) error {
	var params interface{}
	if r.req.Params != nil {
		if err := json.Unmarshal(*r.req.Params, &params); err != nil {
			return errors.Wrap(err, "unmarshling request parameters failed")
		}
	}

	WalkURIFields(params, nil, r.updateURIFromSrc)

	if r.req.Notif {
		err := r.dest.Notify(ctx, r.req.Method, params)
		if err != nil {
			err = errors.Wrap(err, "sending notification to dest failed")
		}
		// Don't send responses back to src for Notification requests
		return err
	}

	callOpts := []jsonrpc2.CallOption{
		// Proxy the ID used. Otherwise we assign our own ID, breaking
		// calls that depend on controlling the ID such as
		// $/cancelRequest and $/partialResult.
		jsonrpc2.PickID(r.req.ID),
	}

	var rawResult *json.RawMessage
	err := r.dest.Call(ctx, r.req.Method, params, &rawResult, callOpts...)

	if err != nil {
		var respErr *jsonrpc2.Error
		if e, ok := err.(*jsonrpc2.Error); ok {
			respErr = e
		} else {
			respErr = &jsonrpc2.Error{Message: err.Error()}
		}

		var multiErr error = respErr

		if err = r.src.ReplyWithError(ctx, r.req.ID, respErr); err != nil {
			multiErr = multierror.Append(multiErr, errors.Wrap(err, "when sending error reply back to src"))
		}

		return errors.Wrapf(multiErr, "calling method %s on dest failed", r.req.Method)
	}

	var result interface{}
	if rawResult != nil {
		if err := json.Unmarshal(*rawResult, &result); err != nil {
			return errors.Wrap(err, "unmarshling result failed")
		}
	}

	WalkURIFields(result, nil, r.updateURIFromDest)

	if err = r.src.Reply(ctx, r.req.ID, &result); err != nil {
		return errors.Wrap(err, "sending reply to back to src failed")
	}

	return nil
}

func trapSignalsForShutdown(shutdown func()) {
	// Listen for shutdown signals. When we receive one attempt to clean up,
	// but do an insta-shutdown if we receive more than one signal.
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP)
	<-c
	go func() {
		<-c
		os.Exit(0)
	}()

	shutdown()
}
