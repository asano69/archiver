// Package archive turns a URL into a self-contained HTML snapshot
// using the monolith CLI and stores it as a record in the "archives"
// PocketBase collection.
//
// The package is split across three files:
//   - archive.go:   orchestration (this file)
//   - monolith.go:  running the monolith binary
//   - title.go:     extracting a title from the archived HTML
package archive

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// Status values stored in the archives collection's "status" field.
const (
	StatusPending = "pending"
	StatusDone    = "done"
	StatusError   = "error"
)

const collectionName = "archives"

// Run archives rawURL: it creates a "pending" record right away, then
// runs monolith and updates that same record to "done" (with the
// captured HTML attached) or "error" (with the failure message saved
// in its description) once finished.
//
// Run always blocks until archiving is done or fails, so a caller that
// wants this to happen in the background (e.g. a future HTTP handler)
// should invoke it inside its own goroutine.
func Run(ctx context.Context, app core.App, rawURL string) (*core.Record, error) {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, fmt.Errorf("invalid URL: %q", rawURL)
	}

	collection, err := app.FindCollectionByNameOrId(collectionName)
	if err != nil {
		return nil, fmt.Errorf("find %q collection: %w", collectionName, err)
	}

	record := core.NewRecord(collection)
	record.Set("url", rawURL)
	record.Set("title", rawURL)
	record.Set("status", StatusPending)
	if err := app.Save(record); err != nil {
		return nil, fmt.Errorf("create archive record: %w", err)
	}

	pageHTML, err := fetchHTML(ctx, rawURL)
	if err != nil {
		record.Set("status", StatusError)
		record.Set("errorMessage", err.Error())
		if saveErr := app.Save(record); saveErr != nil {
			return nil, fmt.Errorf("save failed archive record: %w", saveErr)
		}
		return record, nil
	}

	if err := finalize(app, record, pageHTML, rawURL); err != nil {
		return nil, err
	}

	return record, nil
}

// finalize saves the fetched page onto record as a "done" archive. If
// that save is rejected by the collection's schema (e.g. the file
// exceeds the configured maxSize, or "status" doesn't accept "done" --
// check the archives collection's fields in the PocketBase admin UI),
// it falls back to an "error" record instead of leaving record stuck
// in "pending" forever. The file is dropped for the fallback save so
// an oversized file doesn't also block it.
func finalize(app core.App, record *core.Record, pageHTML []byte, rawURL string) error {
	file, err := filesystem.NewFileFromBytes(pageHTML, "archive.html")
	if err != nil {
		return fmt.Errorf("prepare archive file: %w", err)
	}

	record.Set("title", ExtractTitle(pageHTML, rawURL))
	record.Set("file", file)
	record.Set("status", StatusDone)

	saveErr := app.Save(record)
	if saveErr == nil {
		return nil
	}

	wrappedErr := fmt.Errorf("save archive record: %w", saveErr)
	record.Set("file", nil)
	record.Set("status", StatusError)
	record.Set("errorMessage", wrappedErr.Error())
	if fallbackErr := app.Save(record); fallbackErr != nil {
		return fmt.Errorf("%w (fallback to error status also failed: %v)", wrappedErr, fallbackErr)
	}
	return nil
}
