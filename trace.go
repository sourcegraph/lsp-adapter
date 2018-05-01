package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
	nettrace "golang.org/x/net/trace"
)

func traceRequests(family, sessionID string) jsonrpc2.ConnOpt {
	return func(c *jsonrpc2.Conn) {
		var (
			mu     sync.Mutex
			traces = map[jsonrpc2.ID]nettrace.Trace{}
		)

		go func() {
			<-c.DisconnectNotify()
			mu.Lock()
			defer mu.Unlock()
			for _, tr := range traces {
				tr.LazyPrintf("client disconnected")
				tr.SetError()
				tr.Finish()
			}
			traces = map[jsonrpc2.ID]nettrace.Trace{}
		}()

		jsonrpc2.OnRecv(func(req *jsonrpc2.Request, resp *jsonrpc2.Response) {
			if req == nil || req.Notif || resp != nil {
				return
			}

			mu.Lock()
			tr, ok := traces[req.ID]
			if ok {
				// misbehaving clients
				tr.LazyPrintf("error ID repeated")
				tr.SetError()
			} else {
				tr = nettrace.New(family, req.Method)
				traces[req.ID] = tr
			}
			mu.Unlock()

			tr.LazyPrintf("id: %s", lazyMarshal{req.ID})
			tr.LazyPrintf("params: %s", lazyMarshal{req.Params})
			tr.LazyPrintf("session: %s", sessionID)
		})(c)
		jsonrpc2.OnSend(func(req *jsonrpc2.Request, resp *jsonrpc2.Response) {
			if req != nil || resp == nil {
				return
			}

			mu.Lock()
			tr, ok := traces[resp.ID]
			delete(traces, resp.ID)
			mu.Unlock()
			if !ok {
				// bad client, ignore will be logged elsewhere
				return
			}

			if resp.Result != nil {
				tr.LazyPrintf("result: %s", lazyMarshal{resp.Result})
			} else if resp.Error != nil {
				tr.LazyPrintf("error: %s", lazyMarshal{resp.Error})
				tr.SetError()
			}
			tr.Finish()
		})(c)
	}
}

func traceEventLog(family, title string) jsonrpc2.ConnOpt {
	return func(c *jsonrpc2.Conn) {
		el := &finishOnceEventLog{EventLog: nettrace.NewEventLog(family, title)}

		go func() {
			<-c.DisconnectNotify()
			el.Finish()
		}()

		jsonrpc2.OnRecv(func(req *jsonrpc2.Request, resp *jsonrpc2.Response) {
			switch {
			case req != nil && resp == nil:
				params, _ := json.Marshal(req.Params)
				if req.Notif {
					el.Printf("--> notif: %s: %s", req.Method, params)
				} else {
					el.Printf("--> request #%s: %s: %s", req.ID, req.Method, params)
				}

			case resp != nil:
				switch {
				case resp.Result != nil:
					result, _ := json.Marshal(resp.Result)
					el.Printf("--> result #%s: %s", resp.ID, result)
				case resp.Error != nil:
					err, _ := json.Marshal(resp.Error)
					el.Errorf("--> error #%s: %s", resp.ID, err)
				}
			}
		})(c)
		jsonrpc2.OnSend(func(req *jsonrpc2.Request, resp *jsonrpc2.Response) {
			switch {
			case req != nil:
				params, _ := json.Marshal(req.Params)
				if req.Notif {
					el.Printf("<-- notif: %s: %s", req.Method, params)
				} else {
					el.Printf("<-- request #%s: %s: %s", req.ID, req.Method, params)
				}

			case resp != nil:
				if resp.Result != nil {
					result, _ := json.Marshal(resp.Result)
					el.Printf("<-- result #%s: %s", resp.ID, result)
				} else {
					err, _ := json.Marshal(resp.Error)
					el.Errorf("<-- error #%s: %s", resp.ID, err)
				}
			}
		})(c)
	}
}

func debugServer(addr string) {
	if addr == "" {
		return
	}

	nettrace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}

	pp := http.NewServeMux()
	index := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
				<a href="/debug/pprof/">PProf</a><br>
				<a href="/debug/requests">Requests</a><br>
				<a href="/debug/events">Events</a><br>
			`))
	})
	pp.Handle("/", index)
	pp.Handle("/debug", index)
	pp.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	pp.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	pp.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	pp.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	pp.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	pp.Handle("/debug/requests", http.HandlerFunc(nettrace.Traces))
	pp.Handle("/debug/events", http.HandlerFunc(nettrace.Events))
	log.Println("warning: could not start debug HTTP server:", http.ListenAndServe(addr, pp))
}

type lazyMarshal struct {
	w interface{}
}

func (m lazyMarshal) String() string {
	b, err := json.Marshal(m.w)
	if err != nil {
		return "error: " + err.Error()
	}
	return string(b)
}

type finishOnceEventLog struct {
	nettrace.EventLog
	mu     sync.Mutex
	closed bool
}

func (w *finishOnceEventLog) Finish() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.closed = true
	w.EventLog.Finish()
}

func (w *finishOnceEventLog) Printf(format string, a ...interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return
	}
	w.EventLog.Printf(format, a...)
}

func (w *finishOnceEventLog) Errorf(format string, a ...interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return
	}
	w.EventLog.Errorf(format, a...)
}
