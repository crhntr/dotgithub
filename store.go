package dotgithub

import (
	"context"
	"io"

	"github.com/google/go-github/v28/github"
	"github.com/src-d/go-git/storage"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/index"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

type Store struct {
	Client          *github.Client
	Context         context.Context
	RepositoryName  string
	RepositoryOwner string
}

// storer.EncodedObjectStorer

// NewEncodedObject returns a new plumbing.EncodedObject, the real type
// of the object can be a custom implementation or the default one,
// plumbing.MemoryObject.
func (store *Store) NewEncodedObject() plumbing.EncodedObject {
	return &plumbing.MemoryObject{}
}

// SetEncodedObject saves an object into the storage, the object should
// be create with the NewEncodedObject, method, and file if the type is
// not supported.
func (store *Store) SetEncodedObject(plumbing.EncodedObject) (plumbing.Hash, error) {
	return plumbing.ZeroHash, nil
}

// EncodedObject gets an object by hash with the given
// plumbing.ObjectType. Implementors should return
// (nil, plumbing.ErrObjectNotFound) if an object doesn't exist with
// both the given hash and object type.
//
// Valid plumbing.ObjectType values are CommitObject, BlobObject, TagObject,
// TreeObject and AnyObject. If plumbing.AnyObject is given, the object must
// be looked up regardless of its type.
func (store *Store) EncodedObject(plumbing.ObjectType, plumbing.Hash) (plumbing.EncodedObject, error) {
	return nil, nil
}

// IterEncodedObjects returns a custom EncodedObjectStorer over all the object
// on the storage.
//
// Valid plumbing.ObjectType values are CommitObject, BlobObject, TagObject,
func (store *Store) IterEncodedObjects(plumbing.ObjectType) (storer.EncodedObjectIter, error) {
	return nil, nil
}

// HasEncodedObject returns ErrObjNotFound if the object doesn't
// exist.  If the object does exist, it returns nil.
func (store *Store) HasEncodedObject(plumbing.Hash) error { return nil }

// EncodedObjectSize returns the plaintext size of the encoded object.
func (store *Store) EncodedObjectSize(plumbing.Hash) (int64, error) { return 0, nil }

// storer.ReferenceStorer

func (store *Store) SetReference(*plumbing.Reference) error { return nil }

// CheckAndSetReference sets the reference `new`, but if `old` is
// not `nil`, it first checks that the current stored value for
// `old.Name()` matches the given reference value in `old`.  If
// not, it returns an error and doesn't update `new`.
func (store *Store) CheckAndSetReference(new, old *plumbing.Reference) error { return nil }
func (store *Store) Reference(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	ref, _, err := store.Client.Git.GetRef(store.Context, store.RepositoryOwner, store.RepositoryName, string(name))
	if err != nil {
		return nil, err
	}
	return plumbing.NewHashReference(
		plumbing.ReferenceName(ref.GetRef()),
		plumbing.NewHash(ref.Object.GetSHA()),
	), nil
}

func (store *Store) IterReferences() (storer.ReferenceIter, error) {
	done := make(chan struct{})
	refs := make(chan *plumbing.Reference)
	errs := make(chan error)

	go func() {
		defer close(errs)
		defer close(refs)
		// limit of 0 is omited by marshaler
		page, limit := 0, 100

		for {
			select {
			case <-done:
				return
			default:
				opts := &github.ReferenceListOptions{}
				opts.Page = page
				opts.PerPage = limit
				refSlice, _, err := store.Client.Git.ListRefs(store.Context, store.RepositoryOwner, store.RepositoryName, opts)
				page++
				if err != nil {
					errs <- err
					return
				}
				for _, ref := range refSlice {
					refs <- plumbing.NewHashReference(
						plumbing.ReferenceName(ref.GetRef()),
						plumbing.NewHash(ref.Object.GetSHA()),
					)
				}
				if len(refSlice) < limit {
					errs <- io.EOF
					return
				}
			}
		}
	}()

	return ReferenceIterator{done: done, refs: refs, errs: errs}, nil
}

type ReferenceIterator struct {
	done chan struct{}
	refs chan *plumbing.Reference
	errs chan error
}

func (iter ReferenceIterator) Next() (*plumbing.Reference, error) {
	select {
	case err := <-iter.errs:
		return nil, err
	case ref := <-iter.refs:
		return ref, nil
	}
}
func (iter ReferenceIterator) ForEach(fn func(*plumbing.Reference) error) error {
	for {
		ref, err := iter.Next()
		if err != nil {
			return err
		}
		if err := fn(ref); err != nil {
			return err
		}
	}
}

func (iter ReferenceIterator) Close() {
	close(iter.done)
	for range iter.refs {
	}
}

func (store *Store) RemoveReference(plumbing.ReferenceName) error { return nil }
func (store *Store) CountLooseRefs() (int, error)                 { return 0, nil }
func (store *Store) PackRefs() error                              { return nil }

// storer.ShallowStorer
func (store *Store) SetShallow([]plumbing.Hash) error  { return nil }
func (store *Store) Shallow() ([]plumbing.Hash, error) { return nil, nil }

// storer.IndexStorer

func (store *Store) SetIndex(*index.Index) error  { return nil }
func (store *Store) Index() (*index.Index, error) { return nil, nil }

// config.ConfigStorer

func (store *Store) Config() (*config.Config, error) { return nil, nil }
func (store *Store) SetConfig(*config.Config) error  { return nil }

// storage.ModuleStorer

func (store *Store) Module(name string) (storage.Storer, error) { return nil, nil }
