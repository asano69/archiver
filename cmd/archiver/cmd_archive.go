package main

import (
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"

	"github.com/asano69/archiver/internal/archive"
)

// archiveCmd defines the "archiver archive <url>" command. It runs
// synchronously: archive.Run blocks until the page has been fetched via
// monolith and the resulting record saved, so the command's own exit
// code and output directly reflect the outcome.
func archiveCmd(app *pocketbase.PocketBase) *cobra.Command {
	return &cobra.Command{
		Use:   "archive <url>",
		Short: "Archive a URL as a self-contained HTML snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rawURL := args[0]

			record, err := archive.Run(cmd.Context(), app, rawURL)
			if err != nil {
				return err
			}
			if record.GetString("status") == archive.StatusError {
				return fmt.Errorf("archive failed: %s", record.GetString("errorMessage"))
			}

			fmt.Printf("archived %s -> %s\n", rawURL, record.Id)
			return nil
		},
	}
}
