package create

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghrepo"
)

func getOrgPublicKey(client *api.Client, host, orgName string) (string, error) {
	return getPubKey(client, host, fmt.Sprintf("orgs/%s/actions/secrets/public-key", orgName))
}

func getRepoPubKey(client *api.Client, repo ghrepo.Interface) (string, error) {
	return getPubKey(client, repo.RepoHost(), fmt.Sprintf("repos/%s/actions/secrets/public-key",
		ghrepo.FullName(repo)))
}

// TODO support key_id

func getPubKey(client *api.Client, host, path string) (string, error) {
	result := struct {
		Key string
	}{}

	err := client.REST(host, "GET", path, nil, &result)
	if err != nil {
		return "", err
	}

	if result.Key == "" {
		return "", fmt.Errorf("failed to find public key at %s/%s", host, path)
	}

	return result.Key, nil
}

type SecretPayload struct {
	EncryptedValue string   `json:"encrypted_value"`
	Visibility     string   `json:"visibility,omitempty"`
	Repositories   []string `json:"selected_repository_ids,omitempty"`
}

func putOrgSecret(client *api.Client, host, secretName, eValue string) error {
	// TODO handle repository names / visibility
	return nil
}

func putRepoSecret(client *api.Client, repo ghrepo.Interface, secretName, eValue string) error {
	payload := SecretPayload{
		EncryptedValue: eValue,
	}
	path := fmt.Sprintf("repos/%s/actions/secrets/%s", ghrepo.FullName(repo), secretName)
	return putSecret(client, repo.RepoHost(), path, payload)
}

func putSecret(client *api.Client, host, path string, payload SecretPayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize: %w", err)
	}
	requestBody := bytes.NewReader(payloadBytes)

	return client.REST(host, "PUT", path, requestBody, nil)
}
