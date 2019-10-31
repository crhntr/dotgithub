package dotgithub

import "github.com/src-d/go-git/storage"

var _ storage.Storer = Store{}
