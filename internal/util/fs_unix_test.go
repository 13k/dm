//go:build linux || darwin

package util_test

import (
	"github.com/13k/dm/internal/util"
)

func init() {
	rootTestFS = util.NewPath("/")
}
