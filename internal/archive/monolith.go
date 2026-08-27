package archive

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// monolithTimeout bounds how long a single monolith run may take, so a
// slow or unresponsive site can't hang the archive command forever.
const monolithTimeout = 30 * time.Second

// monolithArgs are the fixed monolith flags used for every archive:
// strip JavaScript, keeping monolith's default behavior of inlining
// images, CSS, and fonts into one self-contained HTML file.
var monolithArgs = []string{"--no-js"}

// fetchHTML runs monolith against rawURL and returns the resulting
// HTML entirely in memory. monolith writes to stdout when no -o flag
// is given, so no temporary file is involved.
func fetchHTML(ctx context.Context, rawURL string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, monolithTimeout)
	defer cancel()

	args := append(append([]string{}, monolithArgs...), rawURL)
	cmd := exec.CommandContext(ctx, "monolith", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("archiving timed out after %s", monolithTimeout)
		}
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return nil, fmt.Errorf("monolith failed: %s", msg)
		}
		return nil, fmt.Errorf("monolith failed: %w", err)
	}

	if stdout.Len() == 0 {
		return nil, fmt.Errorf("monolith produced no output")
	}

	return stdout.Bytes(), nil
}
