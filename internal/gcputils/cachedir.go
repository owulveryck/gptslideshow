package gcputils

import (
	"os"
	"os/user"
	"path/filepath"
)

// GetCredentialCacheDir returns $HOME/.credentials or an error if it does not exists
func GetCredentialCacheDir() (string, error) {
	// Get the current user's home directory
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	// Define the .credentials directory path
	credentialsDir := filepath.Join(usr.HomeDir, ".credentials")

	// Check if the directory exists
	if _, err := os.Stat(credentialsDir); os.IsNotExist(err) {
		return "", err
	}
	return credentialsDir, nil
}
