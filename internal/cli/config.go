package cli

import (
	"errors"
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	appName                    = "fastmask"
	configFormat               = "yaml"
	configDirectoryPermissions = 0o700
)

type config struct {
	v           *viper.Viper
	AppName     string
	AppVersion  string
	accountID   string
	accessToken string
}

func (f *fastmask) loadConfig() error {
	v := viper.New()
	v.SetEnvPrefix(appName) // look for env vars prefixed as 'FASTMASK', will be uppercased automatically.

	configFilepath := v.GetString(flagConfig)

	if configFilepath != "" {
		v.SetConfigFile(configFilepath)
	} else {
		// If config filepath flag is not set, look for a config
		// file in home directory and current directory.
		home, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf("failed to determine home directory location: %w", err)
		}
		v.SetConfigName("config") // name of config file (without extension)
		v.SetConfigType(configFormat)
		v.AddConfigPath(path.Join(home, ".fastmask"))
	}

	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) {
			// Config file not found; ignore and create.
			if errCreate := createConfig(); errCreate != nil {
				return errCreate
			}
		} else {
			return fmt.Errorf("failed to read config file or input: %w", err)
		}
	}

	config := &config{
		v:       v,
		AppName: appName,
		// AppVersion:  appVersion,
		accountID:   v.GetString("account_id"),
		accessToken: v.GetString("access_token"),
	}

	f.config = config

	return nil
}

func (c *config) setAccountID(accountID string) {
	c.accountID = accountID
	c.v.Set("account_id", accountID)
}

func (c *config) setAccessToken(accessToken string) {
	c.accessToken = accessToken
	c.v.Set("access_token", accessToken)
}

func (c *config) Save() error {
	if err := c.v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func createConfig() error {
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to determine home directory location: %w", err)
	}

	directory := path.Join(home, ".fastmask")

	if err := os.MkdirAll(directory, configDirectoryPermissions); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := path.Join(directory, ".config.yaml")

	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	defer f.Close()

	return nil
}
