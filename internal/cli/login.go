package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dwin/fastmask/pkg/fastmail"
)

var errAccountIDNotFound = errors.New("no account ID found in response")

func (f *fastmask) loadLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login with Fastmail",
		Long:  "Login with Fastmail and store auth token in config file.",
		RunE:  f.runLogin,
	}

	cmd.Flags().StringP("username", "u", "", "Fastmail email address.")
	cmd.Flags().StringP("password", "p", "", "Fastmail password.")
	cmd.Flags().StringP("mfa-code", "m", "", "Fastmail MFA code.")

	return cmd
}

func reauthNeeded(err error) bool {
	if err == nil {
		return false
	}

	var apiError fastmail.APIError

	if errors.As(err, &apiError) {
		if apiError.Status == "401 Unauthorized" {
			return true
		}
	}

	return false
}

func (f *fastmask) runLogin(cmd *cobra.Command, args []string) error {
	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	mfaCode, err := cmd.Flags().GetString("mfa-code")
	if err != nil {
		return fmt.Errorf("invalid mfa-code: %w", err)
	}

	client := fastmail.NewClient(f.config.AppName)

	resp, err := client.LoginUsernamePasswordMFA(cmd.Context(), username, password, mfaCode)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	accountID, ok := resp.GetMailAccountID()
	if !ok {
		return errAccountIDNotFound
	}

	f.config.setAccountID(accountID)
	f.config.setAccessToken(resp.GetAccessToken())

	fmt.Println("🟢 Login success. Stored access token in config file.")

	if err := f.config.Save(); err != nil {
		return err
	}

	return nil
}
