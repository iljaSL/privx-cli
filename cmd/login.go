package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login either user or client to PrivX",
	Long:  `login commands fetches access token for consequent calls of the client`,
	Example: `
export SESSION=$(privx-cli login [access flags])
privx-cli -s $SESSION ...
	`,
	SilenceUsage: true,
	RunE:         login,
}

func login(cmd *cobra.Command, args []string) error {
	token, err := auth().AccessToken()
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write([]byte(token))
	return err
}
