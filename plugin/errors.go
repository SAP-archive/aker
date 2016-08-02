package plugin

import "fmt"

type ConfigDecodeError struct {
	original error
}

func (e *ConfigDecodeError) Error() string {
	return fmt.Sprintf("error decoding plugin config: %v", e.original.Error())
}
