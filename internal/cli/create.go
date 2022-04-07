package cli

import (
	"fmt"

	"github.com/dwin/fastmask/pkg/fastmail"
	"github.com/spf13/cobra"
)

const (
	flagDescription = "description"
	flagDisabled    = "disabled"
)

func (f *fastmask) loadCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <domain>",
		Short: "Create masked email.",
		Long:  "Create a new masked email address using the given domain and optional description.",
		RunE:  f.runCreate,
	}

	cmd.SetUsageTemplate("fastmask create <domain> [flags]\n\n")

	cmd.Args = cobra.ExactArgs(1)

	cmd.Flags().StringP(flagDescription, "d", "", "Description of the masked email.")
	cmd.Flags().Bool(flagDisabled, false, "Create the masked email in disabled state, messages will go to trash.")

	return cmd
}

func (f *fastmask) runCreate(cmd *cobra.Command, args []string) error {
	domain := args[0]

	enabled, err := cmd.Flags().GetBool(flagDisabled)
	if err != nil {
		return fmt.Errorf("failed to get flag %s: %w", flagDisabled, err)
	}

	description, err := cmd.Flags().GetString(flagDescription)
	if err != nil {
		return fmt.Errorf("failed to get flag %s: %w", flagDescription, err)
	}

	m := fastmail.MaskedEmail{
		ForDomain:   domain,
		Description: description,
	}

	client := fastmail.NewClient(f.config.AppName)
	client.SetTokenAuthCredentials(f.config.accountID, f.config.accessToken)

	resp, err := client.CreateMaskedEmail(cmd.Context(), &m, !enabled) // must invert disabled to enabled
	if err != nil {
		if reauthNeeded(err) {
			fmt.Println("🛑 Authentication failed. Please run 'fastmask auth'.")
			return nil
		}

		return fmt.Errorf("failed to create masked email: %w", err)
	}

	return writeOutput(resp)
}
