package plugin

import "fmt"

func sockPath(name string) string {
	return fmt.Sprintf("/tmp/%s-plugin.sock", name)
}
