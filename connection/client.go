package connection

import (
	"errors"
	"net/url"
	"os"

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
	output := struct {
		Content []byte `json:"content"`
	}{
		Content: data,
	}
	p.configStream.Push(&output)

	input := struct {
		Error string `json:"error"`
	}{}
	if ok := p.configStream.Pop(&input); !ok {
		os.Exit(1)
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

	output := struct {
		RequestChannelId  int `json:"request_channel_id"`
		ResponseChannelId int `json:"response_channel_id"`
		DataChannelId     int `json:"data_channel_id"`
	}{
		RequestChannelId:  requestChannelId,
		ResponseChannelId: responseChannelId,
		DataChannelId:     dataChannelId,
	}
	p.processStream.Push(&output)
	input := struct{}{}
	if ok := p.processStream.Pop(&input); !ok {
		os.Exit(1)
	}

	p.peer.CloseChannel(requestChannel)
	p.peer.CloseChannel(responseChannel)
	p.peer.CloseChannel(dataChannel)
	return false
}

func NewRequestClient(channel Channel) api.Request {
	return &requestClient{
		channel:      channel,
		urlStream:    channel.GetStream("url"),
		methodStream: channel.GetStream("method"),
		headerStream: channel.GetStream("header"),
		readStream:   channel.GetStream("read"),
		closeStream:  channel.GetStream("close"),
	}
}

type requestClient struct {
	channel      Channel
	urlStream    Stream
	methodStream Stream
	headerStream Stream
	readStream   Stream
	closeStream  Stream
}

func (c *requestClient) URL() *url.URL {
	output := struct{}{}
	c.urlStream.Push(&output)

	input := struct {
		URL string `json:"url"`
	}{}
	if ok := c.urlStream.Pop(&input); !ok {
		os.Exit(1)
	}
	url, err := url.Parse(input.URL)
	if err != nil {
		panic(err)
	}
	return url
}

func (c *requestClient) Method() string {
	output := struct{}{}
	c.methodStream.Push(&output)

	input := struct {
		Method string `json:"method"`
	}{}
	if ok := c.methodStream.Pop(&input); !ok {
		os.Exit(1)
	}
	return input.Method
}

func (c *requestClient) Header(name string) string {
	output := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	c.headerStream.Push(&output)

	input := struct {
		Value string `json:"value"`
	}{}
	if ok := c.headerStream.Pop(&input); !ok {
		os.Exit(1)
	}
	return input.Value
}

func (c *requestClient) Read(data []byte) (int, error) {
	output := struct {
		Length int `json:"length"`
	}{
		Length: len(data),
	}
	c.readStream.Push(&output)

	input := struct {
		Content []byte `json:"content"`
		Error   string `json:"error"`
	}{}
	if ok := c.readStream.Pop(&input); !ok {
		os.Exit(1)
	}
	count := copy(data, input.Content)
	if count < len(input.Content) {
		panic("Message will drop content.")
	}
	if input.Error != "" {
		return len(input.Content), errors.New(input.Error)
	}
	return len(input.Content), nil
}

func (c *requestClient) Close() error {
	output := struct{}{}
	c.closeStream.Push(&output)

	input := struct {
		Error string `json:"error"`
	}{}
	if ok := c.closeStream.Pop(&input); !ok {
		os.Exit(1)
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
	output := struct {
		Name   string   `json:"name"`
		Values []string `json:"values"`
	}{
		Name:   name,
		Values: values,
	}
	c.setHeaderStream.Push(&output)
}

func (c *responseClient) WriteStatus(status int) {
	output := struct {
		Status int `json:"status"`
	}{
		Status: status,
	}
	c.writeStatusStream.Push(&output)
}

func (c *responseClient) Write(data []byte) (int, error) {
	output := struct {
		Content []byte `json:"content"`
	}{
		Content: data,
	}
	c.writeStream.Push(&output)

	input := struct {
		Count int    `json:"count"`
		Error string `json:"error"`
	}{}
	if ok := c.writeStream.Pop(&input); !ok {
		os.Exit(1)
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
	output := struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{
		Name:  name,
		Value: value,
	}
	w.setStringStream.Push(&output)
}

func (w *dataClient) String(name string) string {
	output := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	w.stringStream.Push(&output)

	input := struct {
		Value string `json:"value"`
	}{}
	if ok := w.stringStream.Pop(&input); !ok {
		os.Exit(1)
	}
	return input.Value
}
