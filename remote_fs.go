package main

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/neelance/parallel"
	"github.com/pkg/errors"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/go-langserver/pkg/lspext"
	"github.com/sourcegraph/jsonrpc2"
)

type remoteFS struct {
	conn *jsonrpc2.Conn
}

// BatchOpen opens all of the content for the specified paths.
func (fs *remoteFS) BatchOpen(ctx context.Context, fileURIs []lsp.DocumentURI) ([]batchFile, error) {
	par := parallel.NewRun(8)

	var mut sync.Mutex
	var batchFiles []batchFile

	for _, fileURI := range fileURIs {
		par.Acquire()

		go func(uri lsp.DocumentURI) {
			defer par.Release()

			text, err := fs.Open(ctx, uri)
			if err != nil {
				par.Error(err)
				return
			}

			mut.Lock()
			defer mut.Unlock()

			batchFiles = append(batchFiles, batchFile{uri: uri, content: text})

		}(fileURI)
	}

	if err := par.Wait(); err != nil {
		return nil, err
	}

	return batchFiles, nil
}

type batchFile struct {
	uri     lsp.DocumentURI
	content string
}

// Open returns the content of the text file for the given file uri path.
func (fs *remoteFS) Open(ctx context.Context, fileURI lsp.DocumentURI) (string, error) {
	params := lspext.ContentParams{TextDocument: lsp.TextDocumentIdentifier{URI: fileURI}}
	var res lsp.TextDocumentItem

	if err := fs.conn.Call(ctx, "textDocument/xcontent", params, &res); err != nil {
		return "", errors.Wrap(err, "calling textDocument/xcontent failed")
	}

	return res.Text, nil
}

// Walk returns a list of all file uris that are children of "base".
func (fs *remoteFS) Walk(ctx context.Context, base string) ([]lsp.DocumentURI, error) {
	params := lspext.FilesParams{Base: base}
	var res []lsp.TextDocumentIdentifier

	if err := fs.conn.Call(ctx, "workspace/xfiles", &params, &res); err != nil {
		return nil, errors.Wrap(err, "calling workspace/xfiles failed")
	}

	var fileURIs []lsp.DocumentURI
	for _, ident := range res {
		fileURIs = append(fileURIs, ident.URI)
	}

	return fileURIs, nil
}

func (fs *remoteFS) Clone(ctx context.Context, baseDir string) error {
	filePaths, err := fs.Walk(ctx, "/")
	if err != nil {
		return errors.Wrap(err, "failed to fetch all filePaths during clone")
	}

	files, err := fs.BatchOpen(ctx, filePaths)
	if err != nil {
		return errors.Wrap(err, "failed to batch open files during clone")
	}

	for _, file := range files {
		parsedFileURI, err := url.Parse(string(file.uri))
		if err != nil {
			errors.Wrapf(err, "failed to parse raw file uri %s for Clone", file.uri)
		}

		newFilePath := filepath.Join(baseDir, filepath.FromSlash(parsedFileURI.Path))

		// There is an assumption here that all paths returned from Walk()
		// point to files, not directories
		parentDir := filepath.Dir(newFilePath)

		if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to make parent dirs for %s")
		}

		if err := ioutil.WriteFile(newFilePath, []byte(file.content), os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to write file content for %s", newFilePath)
		}
	}
	return nil
}
