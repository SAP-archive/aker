package connection

import (
	"io"

	"github.wdf.sap.corp/I061150/aker/api"
)

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
	var input pluginConfigInput
	var output pluginConfigOutput
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
	var input pluginProcessInput
	var output pluginProcessOutput
	for s.processStream.Pop(&input) {
		requestChannel := s.peer.OpenChannel(input.RequestChannelId)
		request := NewRequestClient(requestChannel)
		responseChannel := s.peer.OpenChannel(input.ResponseChannelId)
		response := NewResponseClient(responseChannel)
		dataChannel := s.peer.OpenChannel(input.DataChannelId)
		data := NewDataClient(dataChannel)

		output.Done = s.delegate.Process(api.Context{
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
		delegate:            delegate,
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
	server.listenAndServe()
}

type requestServer struct {
	delegate            api.Request
	channel             Channel
	urlStream           Stream
	hostStream          Stream
	contentLengthStream Stream
	methodStream        Stream
	headersStream       Stream
	headerStream        Stream
	readStream          Stream
	closeStream         Stream
}

func (s *requestServer) listenAndServe() {
	go s.serveURL()
	go s.serveMethod()
	go s.serveHost()
	go s.serveContentLength()
	go s.serveHeaders()
	go s.serveHeader()
	go s.serveRead()
	go s.serveClose()
}

func (s *requestServer) serveURL() {
	var input requestURLInput
	var output requestURLOutput
	for s.urlStream.Pop(&input) {
		output.URL = s.delegate.URL().String()
		s.urlStream.Push(&output)
	}
}

func (s *requestServer) serveMethod() {
	var input requestMethodInput
	var output requestMethodOutput
	for s.methodStream.Pop(&input) {
		output.Method = s.delegate.Method()
		s.methodStream.Push(&output)
	}
}

func (s *requestServer) serveHost() {
	var input requestHostInput
	var output requestHostOutput
	for s.hostStream.Pop(&input) {
		output.Host = s.delegate.Host()
		s.hostStream.Push(&output)
	}
}

func (s *requestServer) serveContentLength() {
	var input requestContentLengthInput
	var output requestContentLengthOutput
	for s.contentLengthStream.Pop(&input) {
		output.ContentLength = s.delegate.ContentLength()
		s.contentLengthStream.Push(&output)
	}
}

func (s *requestServer) serveHeaders() {
	var input requestHeadersInput
	var output requestHeadersOutput
	for s.headersStream.Pop(&input) {
		output.Headers = s.delegate.Headers()
		s.headersStream.Push(&output)
	}
}

func (s *requestServer) serveHeader() {
	var input requestHeaderInput
	var output requestHeaderOutput
	for s.headerStream.Pop(&input) {
		output.Value = s.delegate.Header(input.Name)
		s.headerStream.Push(&output)
	}
}

func (s *requestServer) serveRead() {
	var input requestReadInput
	var output requestReadOutput
	for s.readStream.Pop(&input) {
		data := make([]byte, input.Length)
		count, err := s.delegate.Read(data)
		output.Content = data[:count]
		if err != nil {
			if err == io.EOF {
				output.Error = ""
				output.EOF = true
			} else {
				output.Error = err.Error()
				output.EOF = false
			}
		} else {
			output.Error = ""
		}
		s.readStream.Push(&output)
	}
}

func (s *requestServer) serveClose() {
	var input requestCloseInput
	var output requestCloseOutput
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
		writeStatusStream: channel.GetStream("writestatus"),
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
	var input responseSetHeaderInput
	var output responseSetHeaderOutput
	for s.setHeaderStream.Pop(&input) {
		s.delegate.SetHeader(input.Name, input.Values)
		s.setHeaderStream.Push(&output)
	}
}

func (s *responseServer) serveWriteStatus() {
	var input responseWriteStatusInput
	var output responseWriteStatusOutput
	for s.writeStatusStream.Pop(&input) {
		s.delegate.WriteStatus(input.Status)
		s.writeStatusStream.Push(&output)
	}
}

func (s *responseServer) serveWrite() {
	var input responseWriteInput
	var output responseWriteOutput
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
	var input dataSetStringInput
	var output dataSetStringOutput
	for s.setStringStream.Pop(&input) {
		s.delegate.SetString(input.Name, input.Value)
		s.setStringStream.Push(&output)
	}
}

func (s *dataServer) serveString() {
	var input dataStringInput
	var output dataStringOutput
	for s.stringStream.Pop(&input) {
		output.Value = s.delegate.String(input.Name)
		s.stringStream.Push(&output)
	}
}
