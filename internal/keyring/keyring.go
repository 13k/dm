package keyring

import (
	"github.com/zalando/go-keyring"

	"github.com/13k/dm/internal/meta"
)

func Get(key string) (string, error) {
	return keyring.Get(meta.AppID, key)
}

func Set(key, value string) error {
	return keyring.Set(meta.AppID, key, value)
}
