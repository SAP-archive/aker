package plugin

import "github.com/cloudfoundry-incubator/candiedyaml"

// UnmarshalConfig parses the plugin configuration data and stores the result
// in the value pointed by v.
func UnmarshalConfig(data []byte, v interface{}) error {
	return candiedyaml.Unmarshal(data, v)
}

// MarshalConfig returns the plugin configuration data encoding of v.
func MarshalConfig(v interface{}) ([]byte, error) {
	return candiedyaml.Marshal(v)
}
