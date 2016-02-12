package connection

import "github.wdf.sap.corp/I061150/aker/api"

func ServePlugin(peer Peer, delegate api.Plugin) {
	channel := peer.OpenChannel(0)
	server := &pluginServer{
		delegate:      delegate,
		peer:          peer,
		channel:       channel,
		configStream:  channel.GetStream("config"),
		processStream: channel.GetStream("process"),
	}
	server.listenAndServe()
}

type pluginServer struct {
	delegate      api.Plugin
	peer          Peer
	channel       Channel
	configStream  Stream
	processStream Stream
}

func (s *pluginServer) listenAndServe() {
	go s.serveConfig()
	go s.serveProcess()
}

func (s *pluginServer) serveConfig() {
	var input struct {
		Content []byte `json:"content"`
	}
	var output struct {
		Error string `json:"error"`
	}
	for s.configStream.Pop(&input) {
		err := s.delegate.Config(input.Content)
		if err != nil {
			output.Error = err.Error()
		} else {
			output.Error = ""
		}
		s.configStream.Push(&output)
	}
}

func (s *pluginServer) serveProcess() {
	var input struct {
		RequestChannelId  int `json:"request_channel_id"`
		ResponseChannelId int `json:"response_channel_id"`
		DataChannelId     int `json:"data_channel_id"`
	}
	var output struct{}
	for s.processStream.Pop(&input) {
		requestChannel := s.peer.OpenChannel(input.RequestChannelId)
		request := NewRequestClient(requestChannel)
		responseChannel := s.peer.OpenChannel(input.ResponseChannelId)
		response := NewResponseClient(responseChannel)
		dataChannel := s.peer.OpenChannel(input.DataChannelId)
		data := NewDataClient(dataChannel)
		s.delegate.Process(api.Context{
			Request:  request,
			Response: response,
			Data:     data,
		})
		s.processStream.Push(&output)
		s.peer.CloseChannel(requestChannel)
		s.peer.CloseChannel(responseChannel)
		s.peer.CloseChannel(dataChannel)
	}
}

func ServeRequest(channel Channel, delegate api.Request) {
	server := &requestServer{
		delegate:     delegate,
		channel:      channel,
		urlStream:    channel.GetStream("url"),
		methodStream: channel.GetStream("method"),
		headerStream: channel.GetStream("header"),
		readStream:   channel.GetStream("read"),
		closeStream:  channel.GetStream("close"),
	}
	server.listenAndServe()
}

type requestServer struct {
	delegate     api.Request
	channel      Channel
	urlStream    Stream
	methodStream Stream
	headerStream Stream
	readStream   Stream
	closeStream  Stream
}

func (s *requestServer) listenAndServe() {
	go s.serveURL()
	go s.serveMethod()
	go s.serveHeader()
	go s.serveRead()
	go s.serveClose()
}

func (s *requestServer) serveURL() {
	var input struct{}
	var output struct {
		URL string `json:"url"`
	}
	for s.urlStream.Pop(&input) {
		output.URL = s.delegate.URL().String()
		s.urlStream.Push(&output)
	}
}

func (s *requestServer) serveMethod() {
	var input struct{}
	var output struct {
		Method string `json:"method"`
	}
	for s.methodStream.Pop(&input) {
		output.Method = s.delegate.Method()
		s.methodStream.Push(&output)
	}
}

func (s *requestServer) serveHeader() {
	var input struct {
		Name string `json:"name"`
	}
	var output struct {
		Value string `json:"value"`
	}
	for s.headerStream.Pop(&input) {
		output.Value = s.delegate.Header(input.Name)
		s.headerStream.Push(&output)
	}
}

func (s *requestServer) serveRead() {
	var input struct {
		Length int `json:"length"`
	}
	var output struct {
		Content []byte `json:"content"`
		Error   string `json:"error"`
	}
	for s.readStream.Pop(&input) {
		data := make([]byte, input.Length)
		count, err := s.delegate.Read(data)
		output.Content = data[:count]
		if err != nil {
			output.Error = err.Error()
		} else {
			output.Error = ""
		}
		s.readStream.Push(&output)
	}
}

func (s *requestServer) serveClose() {
	var input struct{}
	var output struct {
		Error string `json:"error"`
	}
	for s.closeStream.Pop(&input) {
		err := s.delegate.Close()
		if err != nil {
			output.Error = err.Error()
		} else {
			output.Error = ""
		}
		s.closeStream.Push(&output)
	}
}

func ServeResponse(channel Channel, delegate api.Response) {
	server := &responseServer{
		delegate:          delegate,
		channel:           channel,
		setHeaderStream:   channel.GetStream("setheader"),
		writeStatusStream: channel.GetStream("writeStatus"),
		writeStream:       channel.GetStream("write"),
	}
	server.listenAndServe()
}

type responseServer struct {
	delegate          api.Response
	channel           Channel
	setHeaderStream   Stream
	writeStatusStream Stream
	writeStream       Stream
}

func (s *responseServer) listenAndServe() {
	go s.serveSetHeader()
	go s.serveWriteStatus()
	go s.serveWrite()
}

func (s *responseServer) serveSetHeader() {
	var input struct {
		Name   string   `json:"name"`
		Values []string `json:"values"`
	}
	for s.setHeaderStream.Pop(&input) {
		s.delegate.SetHeader(input.Name, input.Values)
	}
}

func (s *responseServer) serveWriteStatus() {
	var input struct {
		Status int `json:"status"`
	}
	for s.writeStatusStream.Pop(&input) {
		s.delegate.WriteStatus(input.Status)
	}
}

func (s *responseServer) serveWrite() {
	var input struct {
		Content []byte `json:"content"`
	}
	var output struct {
		Count int    `json:"count"`
		Error string `json:"error"`
	}
	for s.writeStream.Pop(&input) {
		count, err := s.delegate.Write(input.Content)

		output.Count = count
		if err != nil {
			output.Error = err.Error()
		} else {
			output.Error = ""
		}
		s.writeStream.Push(&output)
	}
}

func ServeData(channel Channel, delegate api.Data) {
	server := &dataServer{
		delegate:        delegate,
		channel:         channel,
		setStringStream: channel.GetStream("setstring"),
		stringStream:    channel.GetStream("string"),
	}
	server.listenAndServe()
}

type dataServer struct {
	delegate        api.Data
	channel         Channel
	setStringStream Stream
	stringStream    Stream
}

func (s *dataServer) listenAndServe() {
	go s.serveSetString()
	go s.serveString()
}

func (s *dataServer) serveSetString() {
	var input struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	for s.setStringStream.Pop(&input) {
		s.delegate.SetString(input.Name, input.Value)
	}
}

func (s *dataServer) serveString() {
	var input struct {
		Name string `json:"name"`
	}
	var output struct {
		Value string `json:"value"`
	}
	for s.stringStream.Pop(&input) {
		output.Value = s.delegate.String(input.Name)
		s.stringStream.Push(&output)
	}
}
