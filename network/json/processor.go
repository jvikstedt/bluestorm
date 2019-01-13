package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/jvikstedt/bluestorm/network"
)

// MsgInfo contains information about the msg
type MsgInfo struct {
	msgType reflect.Type
	handler MsgHandler
}

// Processor contains list of messages
type Processor struct {
	msgInfo map[string]*MsgInfo
}

// MsgHandler callback function that msg will be sent to
type MsgHandler func(*network.Agent, interface{})

// NewProcessor initialize new Processor
func NewProcessor() *Processor {
	return &Processor{
		msgInfo: map[string]*MsgInfo{},
	}
}

// Register stores new MsgInfo with msg type and handler
func (p *Processor) Register(msg interface{}, handler MsgHandler) error {
	t := reflect.TypeOf(msg)
	if t == nil || t.Kind() != reflect.Ptr {
		return fmt.Errorf("Invalid message %s", t)
	}

	id := t.Elem().Name()
	if _, ok := p.msgInfo[id]; ok {
		return fmt.Errorf("message %v is already registered", id)
	}

	log.Printf("Registering msg %s\n", id)

	p.msgInfo[id] = &MsgInfo{
		msgType: t,
		handler: handler,
	}

	return nil
}

// Marshal converts msg to json []byte
func (p *Processor) Marshal(msg interface{}) ([]byte, error) {
	t := reflect.TypeOf(msg)

	if t == nil || t.Kind() != reflect.Ptr {
		return []byte{}, fmt.Errorf("Invalid message %s", t)
	}

	id := t.Elem().Name()
	_, ok := p.msgInfo[id]
	if !ok {
		return []byte{}, fmt.Errorf("Could not find msg by id %s", id)
	}

	m := map[string]interface{}{id: msg}
	return json.Marshal(m)
}

// Unmarshal takes json []byte and converts it to interface{}
// Expects format: {"Greet": {"name": "Bob"}}
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
	var m map[string]json.RawMessage
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	if len(m) != 1 {
		return nil, errors.New("invalid json data")
	}

	for msgID, data := range m {
		i, ok := p.msgInfo[msgID]
		if !ok {
			return nil, fmt.Errorf("message %v not registered", msgID)
		}

		msg := reflect.New(i.msgType.Elem()).Interface()
		return msg, json.Unmarshal(data, msg)
	}

	return nil, nil
}

// Route forwards msg to correct handler
func (p *Processor) Route(a *network.Agent, msg interface{}) error {
	t := reflect.TypeOf(msg)
	if t == nil || t.Kind() != reflect.Ptr {
		return fmt.Errorf("Invalid message %s", t)
	}

	id := t.Elem().Name()
	i, ok := p.msgInfo[id]
	if !ok {
		return fmt.Errorf("Could not find msg handler by id %s", id)
	}

	i.handler(a, msg)

	return nil
}
