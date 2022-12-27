//go:build unix

package util_test

import (
	"github.com/13k/dm/internal/util"
)

func init() {
	rootTestFS = util.NewPath("/")
}
