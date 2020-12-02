package create

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghrepo"
)

func getOrgPublicKey(client *api.Client, host, orgName string) (*PubKey, error) {
	return getPubKey(client, host, fmt.Sprintf("orgs/%s/actions/secrets/public-key", orgName))
}

func getRepoPubKey(client *api.Client, repo ghrepo.Interface) (*PubKey, error) {
	return getPubKey(client, repo.RepoHost(), fmt.Sprintf("repos/%s/actions/secrets/public-key",
		ghrepo.FullName(repo)))
}

// TODO support key_id

type PubKey struct {
	Key     [32]byte
	ID      string
	encoded string
}

func (pk *PubKey) String() string {
	return pk.encoded
}

func NewPubKey(encodedKey, keyID string) (*PubKey, error) {
	pk, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	pka := [32]byte{}
	copy(pka[:], pk[0:32])
	return &PubKey{
		Key:     pka,
		ID:      keyID,
		encoded: encodedKey,
	}, nil
}

func getPubKey(client *api.Client, host, path string) (*PubKey, error) {
	result := struct {
		Key string
		ID  string `json:"key_id"`
	}{}

	err := client.REST(host, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	if result.Key == "" {
		return nil, fmt.Errorf("failed to find public key at %s/%s", host, path)
	}

	return NewPubKey(result.Key, result.ID)
}

type SecretPayload struct {
	EncryptedValue string   `json:"encrypted_value"`
	Visibility     string   `json:"visibility,omitempty"`
	Repositories   []string `json:"selected_repository_ids,omitempty"`
	KeyID          string   `json:"key_id"`
}

func putOrgSecret(client *api.Client, pk *PubKey, host, secretName, eValue string) error {
	// TODO handle repository names / visibility
	return nil
}

func putRepoSecret(client *api.Client, pk *PubKey, repo ghrepo.Interface, secretName, eValue string) error {
	payload := SecretPayload{
		EncryptedValue: eValue,
		KeyID:          pk.ID,
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
