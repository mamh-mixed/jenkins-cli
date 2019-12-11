package cmd

import (
	"fmt"
	"github.com/jenkins-zh/jenkins-cli/app/i18n"
	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/spf13/cobra"
	"net/http"
)

// CredentialDeleteOption option for credential delete command
type CredentialDeleteOption struct {
	BatchOption

	ID    string
	Store string

	RoundTripper http.RoundTripper
}

var credentialDeleteOption CredentialDeleteOption

func init() {
	credentialCmd.AddCommand(credentialDeleteCmd)
	credentialDeleteCmd.Flags().StringVarP(&credentialDeleteOption.Store, "store", "", "system",
		i18n.T("The store name of Jenkins credentials"))
	credentialDeleteCmd.Flags().StringVarP(&credentialDeleteOption.ID, "id", "", "",
		i18n.T("The ID of Jenkins credentials"))
	credentialDeleteOption.SetFlag(credentialDeleteCmd)
}

var credentialDeleteCmd = &cobra.Command{
	Use:     "delete [store] [id]",
	Aliases: []string{"remove", "del"},
	Short:   i18n.T("Delete a credential from Jenkins"),
	Long:    i18n.T("Delete a credential from Jenkins"),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) >= 1 {
			credentialDeleteOption.Store = args[0]
		}

		if len(args) >= 2 {
			credentialDeleteOption.ID = args[1]
		}

		if credentialDeleteOption.Store == "" || credentialDeleteOption.ID == "" {
			err = fmt.Errorf("the store or id of target credential is empty")
		}
		return
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if !credentialDeleteOption.Confirm(fmt.Sprintf("Are you sure to delete credential %s", credentialDeleteOption.ID)) {
			return
		}

		jClient := &client.CredentialsManager{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: credentialDeleteOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jClient.JenkinsCore))

		err = jClient.Delete(credentialDeleteOption.Store, credentialDeleteOption.ID)
		return
	},
}
