package key_cli

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/99designs/keyring"
)

var (
	//go:embed VERSION
	_v      string
	Version = strings.TrimSpace(_v)
)

func Some[T comparable, C ~[]T](coll C) C {
	var zero T
	nonZero := make(C, 0)
	for _, v := range coll {
		if v != zero {
			nonZero = append(nonZero, v)
		}
	}
	return nonZero
}

func Coalesce[T comparable](els ...T) T {
	var zero T
	for _, el := range els {
		if el != zero {
			return el
		}
	}
	return zero
}

func Fatal(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func MustDo[T any](f func(T) error, i T) {
	if err := f(i); err != nil {
		Fatal(err)
	}
}

func MustReturn[T any](f func() (T, error)) T {
	res, err := f()
	if err != nil {
		Fatal(err)
	}
	return res
}

func Must[Arg, Ret any](f func(Arg) (Ret, error), i Arg) Ret {
	r, err := f(i)
	if err != nil {
		Fatal(err)
	}
	return r
}

type Chain struct {
	ring keyring.Keyring
}

func (c *Chain) Get(key string) []byte {
	return Must(c.ring.Get, key).Data
}

func (c *Chain) Set(key string, value []byte) {
	MustDo(c.ring.Set, keyring.Item{
		Key:  key,
		Data: value,
	})
}

func (c *Chain) Remove(key string) {
	MustDo(c.ring.Remove, key)
}

func (c *Chain) Keys() []string {
	return MustReturn(c.ring.Keys)
}

func Open(name string) (*Chain, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName:                    name,
		KeychainName:                   name,
		KeychainTrustApplication:       true,
		KeychainSynchronizable:         false,
		KeychainAccessibleWhenUnlocked: false,
		LibSecretCollectionName:        name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring %s: %w", name, err)
	}
	return &Chain{ring}, nil
}
