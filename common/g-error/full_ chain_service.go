package g_error

import "github.com/pkg/errors"

var (
	ErrNotMineMaster = errors.New("current node is not mine master")
)
