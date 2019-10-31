package dotgithub

import (
	"context"
	"io"
	"testing"

	"github.com/google/go-github/v28/github"
	"github.com/src-d/go-git/storage"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var _ storage.Storer = &Store{}
var _ = &ReferenceIterator{}

func TestStore_Reference(t *testing.T) {
	client := github.NewClient(nil)
	owner, repo := "crhntr", "dotgithub"
	store := Store{
		Client:          client,
		Context:         context.Background(),
		RepositoryOwner: owner,
		RepositoryName:  repo,
	}

	ref, err := store.Reference("refs/heads/master")
	if err != nil {
		t.Error("it should not return an error")
		t.Logf("%[1]s %#[1]v", err)
		return
	}
	t.Logf("%#v", ref)
}

func TestStore_IterReferences(t *testing.T) {
	client := github.NewClient(nil)
	owner, repo := "crhntr", "dotgithub"
	store := Store{
		Client:          client,
		Context:         context.Background(),
		RepositoryOwner: owner,
		RepositoryName:  repo,
	}

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
