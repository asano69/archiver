// Package archive stores captured page snapshots as records in the
// "archives" PocketBase collection.
//
// The package is split across two files:
//   - archive.go: shared constants (this file)
//   - title.go:   extracting a title from the archived HTML
package archive

// Status values stored in the archives collection's "status" field.
const (
	StatusPending = "pending"
	StatusDone    = "done"
	StatusError   = "error"
)

const collectionName = "archives"
