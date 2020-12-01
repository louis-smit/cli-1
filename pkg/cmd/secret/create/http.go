package create

import (
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
