package keyring

import (
	"fmt"

	"github.com/99designs/keyring"

	"github.com/13k/dm/pkg/meta"
)

func open() (keyring.Keyring, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: meta.AppID,
	})

	if err != nil {
		err = fmt.Errorf("could not open keyring for service %q: %w", meta.AppID, err)
	}

	return ring, err
}

func Get(key string) (string, error) {
	ring, err := open()

	if err != nil {
		return "", err
	}

	item, err := ring.Get(key)

	if err != nil {
		return "", fmt.Errorf("could not get keyring item for key %q: %w", key, err)
	}

	return string(item.Data), nil
}

func Set(key, value string) error {
	ring, err := open()

	if err != nil {
		return err
	}

	item := keyring.Item{
		Key:  key,
		Data: []byte(value),
	}

	if err := ring.Set(item); err != nil {
		return fmt.Errorf("could not set keyring item for key %q: %w", key, err)
	}

	return nil
}
