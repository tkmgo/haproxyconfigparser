package haproxyconfigparser

import (
	"github.com/mitchellh/hashstructure"
)

func (u Socket) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}

func (u Global) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}

func (u Acl) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}

func (u UseBackendClauses) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}

func (u UseBackend) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}

func (u Bind) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}

func (u Frontend) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}

func (u Server) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}

func (u Backend) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}

func (u Services) Hash() (uint64, error) {
	return hashstructure.Hash(u, nil)
}
