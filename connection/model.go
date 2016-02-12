package connection

type pluginConfigInput struct {
	Content []byte `json:"content"`
}

type pluginConfigOutput struct {
	Error string `json:"error"`
}

type pluginProcessInput struct {
	RequestChannelId  int `json:"request_channel_id"`
	ResponseChannelId int `json:"response_channel_id"`
	DataChannelId     int `json:"data_channel_id"`
}

type pluginProcessOutput struct {
	Done bool `json:"done"`
}

type requestURLInput struct {
}

type requestURLOutput struct {
	URL string `json:"url"`
}

type requestMethodInput struct {
}

type requestMethodOutput struct {
	Method string `json:"method"`
}

type requestHostInput struct {
}

type requestHostOutput struct {
	Host string `json:"host"`
}

type requestContentLengthInput struct {
}

type requestContentLengthOutput struct {
	ContentLength int `json:"content_length"`
}

type requestHeadersInput struct {
}

type requestHeadersOutput struct {
	Headers map[string][]string `json:"headers"`
}

type requestHeaderInput struct {
	Name string `json:"name"`
}

type requestHeaderOutput struct {
	Value string `json:"value"`
}

type requestReadInput struct {
	Length int `json:"length"`
}

type requestReadOutput struct {
	Content []byte `json:"content"`
	EOF     bool   `json:"eof"`
	Error   string `json:"error"`
}

type requestCloseInput struct {
}

type requestCloseOutput struct {
	Error string `json:"error"`
}

type responseSetHeaderInput struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type responseSetHeaderOutput struct {
}

type responseWriteStatusInput struct {
	Status int `json:"status"`
}

type responseWriteStatusOutput struct {
}

type responseWriteInput struct {
	Content []byte `json:"content"`
}

type responseWriteOutput struct {
	Count int    `json:"count"`
	Error string `json:"error"`
}

type dataSetStringInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type dataSetStringOutput struct {
}

type dataStringInput struct {
	Name string `json:"name"`
}

type dataStringOutput struct {
	Value string `json:"value"`
}
