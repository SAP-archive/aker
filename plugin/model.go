package plugin

type setup struct {
	SocketPath        string `json:"socket_path"`
	ForwardSocketPath string `json:"forward_socket_path"`
	Configuration     []byte `json:"configuration"`
}
