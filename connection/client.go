package connection

import (
	"errors"
	"io"
	"net/url"

	"github.wdf.sap.corp/I061150/aker/api"
)

func NewPluginClient(peer Peer) api.Plugin {
	channel := peer.OpenChannel(0)
	return &pluginClient{
		peer:          peer,
		channel:       channel,
		configStream:  channel.GetStream("config"),
		processStream: channel.GetStream("process"),
	}
}

type pluginClient struct {
	peer          Peer
	channel       Channel
	configStream  Stream
	processStream Stream
}

func (p *pluginClient) Config(data []byte) error {
	output := pluginConfigInput{
		Content: data,
	}
	p.configStream.Push(&output)

	input := pluginConfigOutput{}
	if ok := p.configStream.Pop(&input); !ok {
		panic("Channel closed!")
	}

	if input.Error != "" {
		return errors.New(input.Error)
	}
	return nil
}

func (p *pluginClient) Process(context api.Context) bool {
	requestChannelId := p.peer.GetFreeChannelId()
	requestChannel := p.peer.OpenChannel(requestChannelId)
	ServeRequest(requestChannel, context.Request)

	responseChannelId := p.peer.GetFreeChannelId()
	responseChannel := p.peer.OpenChannel(responseChannelId)
	ServeResponse(responseChannel, context.Response)

	dataChannelId := p.peer.GetFreeChannelId()
	dataChannel := p.peer.OpenChannel(dataChannelId)
	ServeData(dataChannel, context.Data)

	output := pluginProcessInput{
		RequestChannelId:  requestChannelId,
		ResponseChannelId: responseChannelId,
		DataChannelId:     dataChannelId,
	}
	p.processStream.Push(&output)

	input := pluginProcessOutput{}
	if ok := p.processStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}

	p.peer.CloseChannel(requestChannel)
	p.peer.CloseChannel(responseChannel)
	p.peer.CloseChannel(dataChannel)
	return input.Done
}

func NewRequestClient(channel Channel) api.Request {
	return &requestClient{
		channel:             channel,
		urlStream:           channel.GetStream("url"),
		methodStream:        channel.GetStream("method"),
		hostStream:          channel.GetStream("host"),
		contentLengthStream: channel.GetStream("contentlength"),
		headersStream:       channel.GetStream("headers"),
		headerStream:        channel.GetStream("header"),
		readStream:          channel.GetStream("read"),
		closeStream:         channel.GetStream("close"),
	}
}

type requestClient struct {
	channel             Channel
	urlStream           Stream
	methodStream        Stream
	hostStream          Stream
	contentLengthStream Stream
	headersStream       Stream
	headerStream        Stream
	readStream          Stream
	closeStream         Stream
}

func (c *requestClient) URL() *url.URL {
	output := requestURLInput{}
	c.urlStream.Push(&output)

	input := requestURLOutput{}
	if ok := c.urlStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
	url, err := url.Parse(input.URL)
	if err != nil {
		panic(err)
	}
	return url
}

func (c *requestClient) Method() string {
	output := requestMethodInput{}
	c.methodStream.Push(&output)

	input := requestMethodOutput{}
	if ok := c.methodStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
	return input.Method
}

func (c *requestClient) Host() string {
	output := requestHostInput{}
	c.hostStream.Push(&output)

	input := requestHostOutput{}
	if ok := c.hostStream.Pop(&output); !ok {
		panic("Channel is closed!")
	}
	return input.Host
}

func (c *requestClient) ContentLength() int {
	output := requestContentLengthInput{}
	c.contentLengthStream.Push(&output)

	input := requestContentLengthOutput{}
	if ok := c.contentLengthStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
	return input.ContentLength
}

func (c *requestClient) Headers() map[string][]string {
	output := requestHeadersInput{}
	c.headersStream.Push(&output)

	input := requestHeadersOutput{}
	if ok := c.headersStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
	return input.Headers
}

func (c *requestClient) Header(name string) string {
	output := requestHeaderInput{
		Name: name,
	}
	c.headerStream.Push(&output)

	input := requestHeaderOutput{}
	if ok := c.headerStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
	return input.Value
}

func (c *requestClient) Read(data []byte) (int, error) {
	output := requestReadInput{
		Length: len(data),
	}
	c.readStream.Push(&output)

	input := requestReadOutput{}
	if ok := c.readStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
	count := copy(data, input.Content)
	if count < len(input.Content) {
		panic("Message will drop content.")
	}
	if input.EOF {
		return len(input.Content), io.EOF
	}
	if input.Error != "" {
		return len(input.Content), errors.New(input.Error)
	}
	return len(input.Content), nil
}

func (c *requestClient) Close() error {
	output := requestCloseInput{}
	c.closeStream.Push(&output)

	input := requestCloseOutput{}
	if ok := c.closeStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
	if input.Error != "" {
		return errors.New(input.Error)
	}
	return nil
}

func NewResponseClient(channel Channel) api.Response {
	return &responseClient{
		channel:           channel,
		setHeaderStream:   channel.GetStream("setheader"),
		writeStatusStream: channel.GetStream("writestatus"),
		writeStream:       channel.GetStream("write"),
	}
}

type responseClient struct {
	channel           Channel
	setHeaderStream   Stream
	writeStatusStream Stream
	writeStream       Stream
}

func (c *responseClient) SetHeader(name string, values []string) {
	output := responseSetHeaderInput{
		Name:   name,
		Values: values,
	}
	c.setHeaderStream.Push(&output)

	input := responseSetHeaderOutput{}
	if ok := c.setHeaderStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
}

func (c *responseClient) WriteStatus(status int) {
	output := responseWriteStatusInput{
		Status: status,
	}
	c.writeStatusStream.Push(&output)

	input := responseWriteStatusOutput{}
	if ok := c.writeStatusStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
}

func (c *responseClient) Write(data []byte) (int, error) {
	output := responseWriteInput{
		Content: data,
	}
	c.writeStream.Push(&output)

	input := responseWriteOutput{}
	if ok := c.writeStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
	if input.Error != "" {
		return input.Count, errors.New(input.Error)
	}
	return input.Count, nil
}

func NewDataClient(channel Channel) api.Data {
	return &dataClient{
		channel:         channel,
		setStringStream: channel.GetStream("setstring"),
		stringStream:    channel.GetStream("string"),
	}
}

type dataClient struct {
	channel         Channel
	setStringStream Stream
	stringStream    Stream
}

func (w *dataClient) SetString(name, value string) {
	output := dataSetStringInput{
		Name:  name,
		Value: value,
	}
	w.setStringStream.Push(&output)

	input := dataSetStringOutput{}
	if ok := w.setStringStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
}

func (w *dataClient) String(name string) string {
	output := dataStringInput{
		Name: name,
	}
	w.stringStream.Push(&output)

	input := dataStringOutput{}
	if ok := w.stringStream.Pop(&input); !ok {
		panic("Channel is closed!")
	}
	return input.Value
}
