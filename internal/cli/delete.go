package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dwin/fastmask/pkg/fastmail"
	"github.com/spf13/cobra"
)

func (f *fastmask) loadDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>...",
		Short: "Delete masked emails.",
		Long:  "Delete masked email addresses.",
		RunE:  f.delete,
	}

	cmd.Args = cobra.MinimumNArgs(1)

	return cmd
}

var ErrOperationCancelled = errors.New("operation canceled")

func confirmDelete(cmd *cobra.Command, _ []string) error {
	skipConfirm, err := cmd.Flags().GetBool(flagNoConfirm)
	if err != nil {
		return fmt.Errorf("failed to get flag %s: %w", flagNoConfirm, err)
	}

	if skipConfirm {
		return nil
	}

	fmt.Printf("Confirm deletion of masked emails: (y/N): ")

	var confirm string

	fmt.Scanln(&confirm)

	if s := strings.ToLower(confirm); s == "y" || s == "yes" {
		return nil
	}

	fmt.Println("Deletion cancelled.")
	os.Exit(0)

	return nil
}

func (f *fastmask) delete(cmd *cobra.Command, args []string) error {
	if err := confirmDelete(cmd, args); err != nil {
		return err
	}

	client := fastmail.NewClient(f.config.AppName)
	client.SetTokenAuthCredentials(f.config.accountID, f.config.accessToken)

	if err := client.DeleteMaskedEmails(cmd.Context(), args...); err != nil {
		return fmt.Errorf("failed to delete masked emails: %w", err)
	}

	fmt.Println("Masked emails deleted.")

	return nil
}
