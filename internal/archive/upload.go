package archive

import (
	"fmt"
	"io"
	"mime/multipart"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// FromUpload saves a page snapshot that was already captured client-side
// (e.g. by the SingleFile browser extension) as a "done" archive record,
// skipping the monolith fetch in monolith.go entirely.
func FromUpload(app core.App, rawURL string, fileHeader *multipart.FileHeader) (*core.Record, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	pageHTML, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("read uploaded file: %w", err)
	}

	collection, err := app.FindCollectionByNameOrId(collectionName)
	if err != nil {
		return nil, fmt.Errorf("find %q collection: %w", collectionName, err)
	}

	file, err := filesystem.NewFileFromBytes(pageHTML, fileHeader.Filename)
	if err != nil {
		return nil, fmt.Errorf("prepare archive file: %w", err)
	}

	record := core.NewRecord(collection)
	record.Set("url", rawURL)
	record.Set("title", ExtractTitle(pageHTML, fileHeader.Filename))
	record.Set("file", file)
	record.Set("status", StatusDone)
	if err := app.Save(record); err != nil {
		return nil, fmt.Errorf("save archive record: %w", err)
	}

	return record, nil
}
