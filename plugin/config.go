package plugin

import "gopkg.in/yaml.v2"

// UnmarshalConfig parses the plugin configuration data and stores the result
// in the value pointed by v.
func UnmarshalConfig(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// MarshalConfig returns the plugin configuration data encoding of v.
func MarshalConfig(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
