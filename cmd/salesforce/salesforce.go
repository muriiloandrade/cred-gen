/*
Copyright © 2023 Murilo Andrade <murilo@muriloandrade.dev>
*/
package salesforce

import (
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	sfEnv       string
	promptItems = []string{"DevRC", "DevTechRC", "Prodlike", "Production"}
	envMapping  = map[string]string{
		"devrc":      "DEVRC",
		"devtechrc":  "DEVTECHRC",
		"prodlike":   "PRODLIKE",
		"prod":       "PROD",
		"production": "PROD",
	}
)

type SalesforceResponse struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

// salesforceCmd represents the salesforce command
var SalesforceCmd = &cobra.Command{
	Use:   "salesforce",
	Short: "Generates Salesforce access tokens for a given environment",
	Run: func(cmd *cobra.Command, args []string) {
		if sfEnv == "" {
			sfEnvSelected, err := getSalesforceEnvFromPrompt()
			if err != nil {
				os.Exit(1)
				return
			}

			sfEnv = sfEnvSelected
		}

		selectedEnv, ok := envMapping[strings.ToLower(sfEnv)]
		if !ok {
			fmt.Printf("env value invalid\n")
			return
		}

		getSalesforceToken(selectedEnv)
	},
}

func init() {
	// Here you will define your flags and configuration settings.
	var allowedSfEnvs string = strings.Join(promptItems, ", ")

	var usageText string = fmt.Sprintf("The environment to get an access token/sessionId (required).\nAllowed values: %s", allowedSfEnvs)
	SalesforceCmd.Flags().StringVarP(&sfEnv, "sfEnv", "e", "", usageText)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// salesforceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getSalesforceEnvFromPrompt() (string, error) {
	prompt := promptui.Select{
		Label: "Select environment",
		Items: promptItems,
	}

	_, selected, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return selected, nil
}

func getSalesforceToken(sfEnv string) {
	sfUrl := viper.GetString(sfEnv+"_SF_DOMAIN_URL") + "/services/oauth2/token?"

	query := url.Values{}
	query.Add("grant_type", "password")
	query.Add("client_id", viper.GetString(sfEnv+"_SF_OAUTH_CLIENT_ID"))
	query.Add("client_secret", viper.GetString(sfEnv+"_SF_OAUTH_CLIENT_SECRET"))
	query.Add("username", viper.GetString(sfEnv+"_SF_USERNAME"))
	query.Add("password", viper.GetString(sfEnv+"_SF_PASSWORD")+viper.GetString(sfEnv+"_SF_SEC_TOKEN"))

	resp, err := http.Post(sfUrl+query.Encode(), "application/json", nil)
	if err != nil {
		fmt.Printf("Failed to get an access token for env: %s\n", sfEnv)
		return
	}

	defer resp.Body.Close()
	respBody, bodyReadErr := io.ReadAll(resp.Body)

	if bodyReadErr != nil {
		fmt.Printf("Error reading response body\n")
		return
	}

	sfResponse := SalesforceResponse{}
	if err := json.Unmarshal(respBody, &sfResponse); err != nil {
		fmt.Printf("Error unmarshalling salesforce response into variable\n")
		return
	}

	fmt.Printf("Session ID: %s\n", sfResponse.AccessToken)
}
