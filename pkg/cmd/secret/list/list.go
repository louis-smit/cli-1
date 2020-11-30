package list

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghrepo"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/cli/cli/utils"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   func() (ghrepo.Interface, error)

	OrgName string
}

func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List secrets",
		Long:  "List secrets for a repository or organization",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// support `-R, --repo` override
			opts.BaseRepo = f.BaseRepo

			if runF != nil {
				return runF(opts)
			}

			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.OrgName, "org", "o", "", "List secrets for an organization")

	return cmd
}

func listRun(opts *ListOptions) error {
	c, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("could not create http client: %w", err)
	}
	client := api.NewClientFromHTTP(c)

	baseRepo, err := opts.BaseRepo()
	if err != nil {
		return fmt.Errorf("could not determine base repo: %w", err)
	}

	secrets, err := getSecrets(client, baseRepo)
	if err != nil {
		return fmt.Errorf("failed to get secrets: %w", err)
	}

	tp := utils.NewTablePrinter(opts.IO)
	for _, secret := range secrets {
		tp.AddField(secret.Name, nil, nil)
		updatedAt := fmt.Sprintf("updated %s", secret.UpdatedAt.Format("2006-Jan-02"))
		tp.AddField(updatedAt, nil, nil)
		tp.EndRow()
	}

	tp.Render()

	return nil
}

type Secret struct {
	Name      string
	UpdatedAt time.Time `json:"updated_at"`
}

func getSecrets(client *api.Client, repo ghrepo.Interface) ([]Secret, error) {
	url := fmt.Sprintf("repos/%s/%s/actions/secrets", repo.RepoOwner(), repo.RepoName())
	result := struct {
		Secrets []Secret
	}{}

	err := client.REST(repo.RepoHost(), "GET", url, nil, &result)
	if err != nil {
		return nil, err
	}

	return result.Secrets, nil
}
