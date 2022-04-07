package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

const (
	flagNoConfirm = "no-confirm"
	flagConfig    = "config"
)

type fastmask struct {
	cmd    *cobra.Command
	config *config
}

func (f *fastmask) LoadConfig() error {
	return f.loadConfig()
}

// nolint:golint,revive // internal pkg ok to return unexported type.
func LoadFastmask(date, version, commit string) *fastmask {
	// Set version details
	if version == "" {
		version = "dev"
	}

	if commit == "" {
		commit = "none"
	}

	f := &fastmask{}

	cmd := &cobra.Command{
		Use:              appName,
		Short:            "Un-Official CLI for interacting with Fastmail Masked Emails.",
		Long:             "Un-Official CLI for interacting with Fastmail Masked Emails.\n\nNot endorsed or supported by Fastmail.",
		TraverseChildren: true,
		Example: heredoc.Doc(`
			# Login with Fastmail. This also stores your credentials in a config file at ~/.fastmask/config.yaml.
			$ fastmask login -u me@you.com -p abc123 -m 012345 <- MFA is required only if enabled on account.

			# Mask your email address.
			$ fastmask create example.com -d "avoiding endless newsletters."

			{
				"createdBy": "",
				"lastMessageAt": null,
				"description": "still here.",
				"email": "fun.times1234@fastmail.com",
				"createdAt": "2000-01-01T00:00:01Z",
				"url": null,
				"id": "masked-12345678"
			}
			`),
		Version: fmt.Sprintf("%s (commit %s) at %s", version, commit, date),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return f.loadConfig()
		},
	}

	// Flags
	cmd.PersistentFlags().BoolP(flagNoConfirm, "y", false, "Disable confirmation prompt.")

	// Sub-Commands
	cmd.AddCommand(f.loadLoginCmd())
	cmd.AddCommand(f.loadCreateCmd())
	cmd.AddCommand(f.loadDeleteCmd())
	cmd.AddCommand(loadLicenseCmd())

	return &fastmask{
		cmd: cmd,
	}
}

func (f *fastmask) Execute() error {
	// nolint:wrapcheck // cobra.Command.Execute() ok unwrapped.
	return f.cmd.Execute()
}
