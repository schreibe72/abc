package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// containerCmd represents the container command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save Config and Bash Complition",
}

var saveCredentialCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Save Key and Account",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if key == "" {
			return errors.New("Key should not be empty!")
		}
		if account == "" {
			return errors.New("Account shoud not be empty!")
		}
		return nil
	},
	RunE: saveCredential,
}

var saveBashComplitionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Save Bash Complition",
	RunE:  bashCompletion,
}

func bashCompletion(cmd *cobra.Command, args []string) error {
	if filename != "" {
		err := RootCmd.GenBashCompletionFile(filename)
		return err
	}
	err := RootCmd.GenBashCompletion(os.Stdout)
	return err

}

func saveCredential(cmd *cobra.Command, args []string) error {

	if filename == "" {
		filename = cfgFile
	}

	c := Config{
		Key:     key,
		Account: account,
	}
	if verbose {
		logger.Printf("Writing to Filename: %s\n", filename)
	}

	buf, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = createDirIfNeeded(filename)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, buf, 0600)
	return err

}

func init() {
	RootCmd.AddCommand(saveCmd)
	saveCmd.AddCommand(saveCredentialCmd)
	saveCmd.AddCommand(saveBashComplitionCmd)
	saveCredentialCmd.Flags().StringVarP(&filename, "filename", "f", cfgFile, "Filename to store Config")
	saveBashComplitionCmd.Flags().StringVarP(&filename, "filename", "f", "", "Filename to store Bash Complitition")
}
