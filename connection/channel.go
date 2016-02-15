package connection

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type Peer interface {
	GetFreeChannelId() int
	OpenChannel(id int) Channel
	CloseChannel(Channel)
}

type Channel interface {
	Id() int
	GetStream(name string) Stream
}

type Stream interface {
	Push(payload interface{})
	Pop(payload interface{}) bool
}

func NewPeer(in io.Reader, out io.Writer, idOffset int) Peer {
	p := &peer{
		channelId: idOffset,
		decoder:   json.NewDecoder(in),
		encoder:   json.NewEncoder(out),
		channels:  make(map[int]*channel),
	}
	go p.streamPeerInboundMessages()
	return p
}

type peerMessage struct {
	ChannelId      int            `json:"channel_id"`
	ChannelMessage channelMessage `json:"channel_message"`
}

type peer struct {
	channelIdMutex sync.Mutex
	channelId      int

	decoder *json.Decoder

	encoderMutex sync.Mutex
	encoder      *json.Encoder

	channelsMutex sync.Mutex
	channels      map[int]*channel
}

func (p *peer) GetFreeChannelId() int {
	p.channelIdMutex.Lock()
	defer p.channelIdMutex.Unlock()
	p.channelId++
	return p.channelId
}

func (p *peer) OpenChannel(id int) Channel {
	ch := &channel{
		id:      id,
		input:   make(chan channelMessage, 1000),
		output:  make(chan channelMessage, 1000),
		streams: make(map[string]*stream),
	}
	p.channelsMutex.Lock()
	p.channels[id] = ch
	p.channelsMutex.Unlock()
	go ch.streamChannelInboundMessages()
	go p.streamPeerOutboundMessages(id, ch)
	return ch
}

func (p *peer) CloseChannel(c Channel) {
	p.channelsMutex.Lock()
	defer p.channelsMutex.Unlock()

	ch := p.channels[c.Id()]
	close(ch.input)
	close(ch.output)
	delete(p.channels, ch.id)
	ch.close()
}

func (p *peer) streamPeerInboundMessages() {
	for {
		msg := peerMessage{}
		err := p.decoder.Decode(&msg)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				panic(err)
			}
		}
		p.channelsMutex.Lock()
		ch, ok := p.channels[msg.ChannelId]
		p.channelsMutex.Unlock()
		if !ok {
			panic(fmt.Sprintf("Trying to write to missing channel '%d'", msg.ChannelId))
		}
		ch.input <- msg.ChannelMessage
	}
}

func (p *peer) streamPeerOutboundMessages(id int, ch *channel) {
	for msg := range ch.output {
		outMsg := peerMessage{
			ChannelId:      id,
			ChannelMessage: msg,
		}
		p.encoderMutex.Lock()
		err := p.encoder.Encode(&outMsg)
		p.encoderMutex.Unlock()
		if err != nil {
			panic(err)
		}
	}
}

type channelMessage struct {
	Key  string `json:"key"`
	Body []byte `json:"body"`
}

type channel struct {
	id int

	input  chan channelMessage
	output chan channelMessage

	streamsMutex sync.Mutex
	streams      map[string]*stream
}

func (c *channel) Id() int {
	return c.id
}

func (c *channel) GetStream(name string) Stream {
	return c.getStream(name)
}

func (c *channel) close() {
	c.streamsMutex.Lock()
	for _, str := range c.streams {
		close(str.input)
		close(str.output)
	}
	c.streamsMutex.Unlock()
}

func (c *channel) getStream(name string) *stream {
	c.streamsMutex.Lock()
	defer c.streamsMutex.Unlock()

	if str, ok := c.streams[name]; ok {
		return str
	}
	str := &stream{
		input:  make(chan []byte),
		output: make(chan []byte),
	}
	c.streams[name] = str
	go c.streamChannelOutboundMessages(name, str)
	return str
}

func (c *channel) streamChannelInboundMessages() {
	for msg := range c.input {
		str := c.getStream(msg.Key)
		str.input <- msg.Body
	}
}

func (c *channel) streamChannelOutboundMessages(name string, str *stream) {
	for msg := range str.output {
		c.output <- channelMessage{
			Key:  name,
			Body: msg,
		}
	}
}

type stream struct {
	input  chan []byte
	output chan []byte
}

func (s *stream) Push(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	s.output <- data
}

func (s *stream) Pop(payload interface{}) bool {
	data, ok := <-s.input
	if !ok {
		return false
	}
	json.Unmarshal(data, payload)
	return true
}
