package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// runHook runs the specified "program" after the contents of the repository are cloned
// to the workspace cache directory, but before the language server receives the "initialize"
// request from Sourcegraph.
//
// The workspace cache directory is both:
//    - passed as an argument to "program"
//    - used as the cwd for "program"
func (p *cloneProxy) runHook(ctx context.Context, program string) error {
	cmd := exec.CommandContext(ctx, program, p.workspaceCacheDir())
	cmd.Dir = p.workspaceCacheDir()
	cmd.Stdout = os.Stdout

	log.Printf("Running pre-init hook: '%s %s'\n", program, p.workspaceCacheDir())
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "When running pre-init hook: '%s %s'", program, p.workspaceCacheDir())
	}

	return nil
}
