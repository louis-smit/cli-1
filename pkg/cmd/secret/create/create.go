package create

import (
	"errors"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/internal/ghrepo"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   func() (ghrepo.Interface, error)

	OrgName         string
	Body            string
	Visibility      string
	RepositoryNames []string
}

func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "create <secret name>",
		Short: "Create secrets",
		Long:  "Locally encrypt a new secret and send it to GitHub for storage.",
		Example: heredoc.Doc(`
			$ gh secret create NEW_SECRET
			$ gh secret create NEW_SECRET -b"some literal value"
			$ gh secret create NEW_SECRET -b"@file.json"
			$ gh secret create ORG_SECRET --org=myOrg --visibility="repo1,repo2,repo3"
			$ gh secret create ORG_SECRET --org=myOrg --visibility="all"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			// support `-R, --repo` override
			opts.BaseRepo = f.BaseRepo

			// TODO process arguments

			if runF != nil {
				return runF(opts)
			}

			return createRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.OrgName, "org", "o", "", "List secrets for an organization")
	cmd.Flags().Lookup("org").NoOptDefVal = "@owner"
	cmd.Flags().StringVarP(&opts.Visibility, "visibility", "v", "private", "Set visibility for an organization secret.")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Provide either a literal string or a file path; prepend file paths with an @. Reads from STDIN if not provided.")

	return cmd
}

func createRun(opts *CreateOptions) error {
	// TODO

	return errors.New("not implemented")
}
