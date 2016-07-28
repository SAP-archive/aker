package endpoint

import (
	"errors"
	"fmt"
)

type InvalidPathError string

func (e InvalidPathError) Error() string {
	return fmt.Sprintf("invalid endpoint path: %q", string(e))
}

var NoPluginsErr = errors.New("no plugins specified")
