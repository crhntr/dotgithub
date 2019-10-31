package dotgithub

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/google/go-github/v28/github"
	"github.com/src-d/go-git/storage"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var _ storage.Storer = &Store{}
var _ = &ReferenceIterator{}

var _githubClient *github.Client

func init() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	_githubClient = github.NewClient(tc)
}

func setupStore() Store {
	owner, repo := "crhntr", "dotgithub"
	return Store{
		Client:          _githubClient,
		Context:         context.Background(),
		RepositoryOwner: owner,
		RepositoryName:  repo,
	}
}

func TestStore_Reference(t *testing.T) {
	store := setupStore()

	ref, err := store.Reference("refs/heads/master")
	if err != nil {
		t.Error("it should not return an error")
		t.Logf("%[1]s %#[1]v", err)
		return
	}
	t.Logf("%#v", ref)
}

func TestStore_IterReferences(t *testing.T) {
	store := setupStore()

	iter, err := store.IterReferences()
	if err != nil {
		t.Error("it should never error")
		t.Logf("got: %s", err)
	}

	gotMaster := false
	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() == "refs/heads/master" {
			gotMaster = true
		}
		return nil
	}); err != nil && err != io.EOF {
		t.Error("it should never error")
		t.Logf("got: %s", err)
	}
	iter.Close()
	if !gotMaster {
		t.Error("it should get the master ref")
	}
}
