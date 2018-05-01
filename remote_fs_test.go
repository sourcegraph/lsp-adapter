package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/sourcegraph/go-langserver/pkg/lsp"

	"github.com/sourcegraph/go-langserver/pkg/lspext"
	"github.com/sourcegraph/jsonrpc2"
)

func TestClone(t *testing.T) {
	fileList := []batchFile{
		{
			uri:     "file:///a.py",
			content: "This is file A.",
		},
		{
			uri:     "/b.py",
			content: "This is file B.",
		},
		{
			uri:     "file:///dir/c.py",
			content: "This is file C.",
		},
		{
			uri:     "file:///dir/d.go",
			content: "This is file D.",
		},
	}

	files := make(map[string]string)

	for _, aFile := range fileList {
		parsedFileURI, err := url.Parse(string(aFile.uri))
		if err != nil {
			t.Fatalf("unable to parse uri for batchFile %v: %v", aFile, err)
		}
		files[parsedFileURI.Path] = aFile.content
	}

	cases := []struct {
		Name  string
		Globs []string
		Want  []string
	}{
		{
			Name:  "all",
			Globs: nil,
			Want:  []string{"/a.py", "/b.py", "/dir/c.py", "/dir/d.go"},
		},
		{
			Name:  "subset",
			Globs: []string{"*.py"},
			Want:  []string{"/a.py", "/b.py", "/dir/c.py"},
		},
		{
			Name:  "multi",
			Globs: []string{"a*", "b*"},
			Want:  []string{"/a.py", "/b.py"},
		},
		{
			Name:  "none",
			Globs: []string{"NOMATCH"},
			Want:  []string{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			baseDir, err := ioutil.TempDir("", uuid.New().String()+"testClone")
			if err != nil {
				t.Fatalf("when creating temp directory for clone test, err: %v", err)
			}

			defer os.Remove(baseDir)

			runTest(t, files, func(ctx context.Context, fs *remoteFS) {
				want := make(map[string]string)
				for _, k := range tt.Want {
					want[k] = files[k]
				}

				err := fs.Clone(ctx, baseDir, tt.Globs)
				if err != nil {
					t.Errorf("when calling clone(baseDir=%s): %v", baseDir, err)
				}

				found, err := findAll(baseDir)
				if err != nil {
					t.Errorf("when calling Walk for baseDir %s: %v", baseDir, err)
				}

				if !reflect.DeepEqual(found, want) {
					t.Errorf("for clone(baseDir=%s) expected %v, actual %v", baseDir, want, found)
				}
			})
		})
	}
}

func findAll(baseDir string) (map[string]string, error) {
	found := make(map[string]string)
	err := filepath.Walk(baseDir, func(currPath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		content, err := ioutil.ReadFile(currPath)
		if err != nil {
			return err
		}

		currPath = path.Join("/", filepath.ToSlash(filepathTrimPrefix(currPath, baseDir)))
		found[currPath] = string(content)
		return nil
	})
	return found, err
}

func TestBatchOpen(t *testing.T) {
	fileList := []batchFile{
		{
			uri:     "file:///a.py",
			content: "This is file A.",
		},
		{
			uri:     "/b.py",
			content: "This is file B.",
		},
		{
			uri:     "file:///dir/c.py",
			content: "This is file C.",
		},
	}

	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].uri < fileList[j].uri
	})

	files := make(map[string]string)

	for _, aFile := range fileList {
		parsedFileURI, err := url.Parse(string(aFile.uri))
		if err != nil {
			t.Fatalf("unable to parse uri for batchFile %v: %v", aFile, err)
		}
		files[parsedFileURI.Path] = aFile.content
	}

	// open single file
	for _, aFile := range fileList {
		runTest(t, files, func(ctx context.Context, fs *remoteFS) {
			results, err := fs.BatchOpen(ctx, []lsp.DocumentURI{aFile.uri})

			if err != nil {
				t.Errorf("when calling batchOpen on uri %s: %v", aFile.uri, err)
			}

			if !reflect.DeepEqual(results, []batchFile{aFile}) {
				t.Errorf("for batchOpen(paths=%v) expected %v, actual %v", []lsp.DocumentURI{aFile.uri}, []batchFile{aFile}, results)
			}
		})
	}

	// open multiple files
	runTest(t, files, func(ctx context.Context, fs *remoteFS) {
		var allURIs []lsp.DocumentURI

		for _, aFile := range fileList {
			allURIs = append(allURIs, aFile.uri)
		}

		results, err := fs.BatchOpen(ctx, allURIs)

		if err != nil {
			t.Errorf("when calling batchOpen on paths: %v, err: %v", allURIs, err)
		}

		sort.Slice(results, func(i, j int) bool {
			return results[i].uri < results[j].uri
		})

		if !reflect.DeepEqual(results, fileList) {
			t.Errorf("for batchOpen(paths=%v) expected %v, actual %v", allURIs, fileList, results)
		}
	})

	// open single invalid file
	runTest(t, files, func(ctx context.Context, fs *remoteFS) {
		_, err := fs.BatchOpen(ctx, []lsp.DocumentURI{"/non/existent/file.py"})

		if err == nil {
			t.Error("expected error when trying to batchOpen non-existent file '/non/existent/file.py'")
		}
	})

	// open multiple valid files and one invalid file
	runTest(t, files, func(ctx context.Context, fs *remoteFS) {
		allURIs := []lsp.DocumentURI{"non/existent/file.py"}

		for _, aFile := range fileList {
			allURIs = append(allURIs, aFile.uri)
		}

		_, err := fs.BatchOpen(ctx, allURIs)

		if err == nil {
			t.Errorf("expected error when trying to batchOpen(paths=%v) which includes non-existent file '/non/existent/file.py'", allURIs)
		}
	})

	// open zero files
	runTest(t, files, func(ctx context.Context, fs *remoteFS) {
		results, err := fs.BatchOpen(ctx, []lsp.DocumentURI{})

		if err != nil {
			t.Errorf("when calling batchOpen on zero paths: %v", err)
		}

		if len(results) > 0 {
			t.Error("expected zero results when trying to batchOpen zero paths")
		}
	})
}

func TestOpen(t *testing.T) {
	fileList := []batchFile{
		{
			uri:     "/a.py",
			content: "This is file A.",
		},
		{
			uri:     "/b.py",
			content: "This is file B.",
		},
		{
			uri:     "/dir/c.py",
			content: "This is file C.",
		},
	}

	files := make(map[string]string)

	for _, aFile := range fileList {
		files[string(aFile.uri)] = aFile.content
	}

	for _, aFile := range fileList {
		runTest(t, files, func(ctx context.Context, fs *remoteFS) {
			actualFileContent, err := fs.Open(ctx, aFile.uri)

			if err != nil {
				t.Errorf("when calling open on uri: %s, err: %v", aFile.uri, err)
			}

			if actualFileContent != aFile.content {
				t.Errorf("for open(path=%s) expected %v, actual %v", aFile.uri, aFile.content, actualFileContent)
			}
		})
	}

	runTest(t, files, func(ctx context.Context, fs *remoteFS) {
		_, err := fs.Open(ctx, "/c.py")
		if err == nil {
			t.Errorf("expected error when trying to open non-existent file '/c.py'")
		}
	})
}

func TestWalk(t *testing.T) {
	type testCase struct {
		fileNames        []string
		expectedFileURIs []string
	}

	tests := []testCase{
		{
			fileNames:        []string{"/a.py", "/b.py", "/dir/c.py"},
			expectedFileURIs: []string{"file:///a.py", "file:///b.py", "file:///dir/c.py"},
		},
	}

	for _, test := range tests {
		files := make(map[string]string)

		for _, fileName := range test.fileNames {
			files[fileName] = ""
		}

		runTest(t, files, func(ctx context.Context, fs *remoteFS) {
			actualFileURIs, err := fs.Walk(ctx)
			if err != nil {
				t.Errorf("when calling walk: %v", err)
			}

			var actualFileNames []string

			for _, uri := range actualFileURIs {
				actualFileNames = append(actualFileNames, string(uri))
			}

			sort.Strings(actualFileNames)
			sort.Strings(test.expectedFileURIs)

			if len(actualFileNames) == 0 && len(test.expectedFileURIs) == 0 {
				// special case empty slice versus nil comparison below?
				return
			}

			if !reflect.DeepEqual(actualFileNames, test.expectedFileURIs) {
				t.Errorf("for walk expected %v, actual %v", test.expectedFileURIs, actualFileNames)
			}
		})
	}
}

func runTest(t *testing.T, files map[string]string, checkFunc func(ctx context.Context, fs *remoteFS)) {
	ctx := context.Background()

	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()

	clientConn := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(a, jsonrpc2.VSCodeObjectCodec{}), &testFS{
		t:     t,
		files: files,
	})

	serverConn := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(b, jsonrpc2.VSCodeObjectCodec{}), &noopHandler{})
	defer clientConn.Close()
	defer serverConn.Close()

	fs := &remoteFS{
		conn: serverConn,
	}

	checkFunc(ctx, fs)
}

type testFS struct {
	t     *testing.T
	files map[string]string // map of file names to content
}

func (client *testFS) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if req.Notif {
		return
	}

	switch req.Method {
	case "textDocument/xcontent":
		var contentParams lspext.ContentParams
		if err := json.Unmarshal(*req.Params, &contentParams); err != nil {
			client.t.Fatalf("unable to unmarshal params %v for textdocument/xcontent, err: %v", req.Params, err)
		}

		filePathRawURI := string(contentParams.TextDocument.URI)
		filePathURI, err := url.Parse(filePathRawURI)
		if err != nil {
			client.t.Fatalf("unable to parse URI %vfor textdocument/xcontent, err: %v", filePathRawURI, err)
		}

		content, present := client.files[filePathURI.Path]

		if !present {
			err := &jsonrpc2.Error{
				Code:    jsonrpc2.CodeInvalidParams,
				Message: fmt.Sprintf("requested file path %s does not exist", filePathURI),
				Data:    nil,
			}
			if replyErr := conn.ReplyWithError(ctx, req.ID, err); replyErr != nil {
				client.t.Fatalf("error when sending back error reply for document %s, err: %v", filePathURI, replyErr)
			}
			return
		}

		document := lsp.TextDocumentItem{
			URI:  contentParams.TextDocument.URI,
			Text: content,
		}

		if replyErr := conn.Reply(ctx, req.ID, document); replyErr != nil {
			client.t.Fatalf("error when sending back content reply for document %v, err:%v", document, replyErr)
		}

	case "workspace/xfiles":
		var results []lsp.TextDocumentIdentifier
		for filePath := range client.files {
			fileURI, err := url.Parse(filePath)
			if err != nil {
				client.t.Fatalf("unable to parse filePath %s as URI for workspace/xfiles, err: %v", filePath, err)
			}
			fileURI.Scheme = "file"

			results = append(results, lsp.TextDocumentIdentifier{
				URI: lsp.DocumentURI(fileURI.String()),
			})
		}

		if replyErr := conn.Reply(ctx, req.ID, results); replyErr != nil {
			client.t.Fatalf("error when sending back files reply, err: %v", replyErr)
		}

	default:
		err := &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: fmt.Sprintf("method %s is invalid - only textdocument/xcontent and workspace/xfiles are supported", req.Method),
			Data:    nil,
		}

		if replyErr := conn.ReplyWithError(ctx, req.ID, err); replyErr != nil {
			client.t.Fatalf("error when sending back error reply for invalid method %s, err: %v", req.Method, replyErr)
		}
	}
}

type noopHandler struct{}

func (noopHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {}
