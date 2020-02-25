package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ContractCommand is
type ContractCommand struct {
	Tag         string
	CreateInput *vmcommon.ContractCreateInput
	CallInput   *vmcommon.ContractCallInput
}

func (command *ContractCommand) String() string {
	return fmt.Sprintf("Command [%s]", command.Tag)
}

// HookCallRequest is
type HookCallRequest struct {
	Tag       string
	Hook      string
	Function  string
	Arguments []interface{}
}

// HookCallResponse is
type HookCallResponse struct {
	Tag    string
	Result []interface{}
	Error  error
}

// Messenger is
type Messenger struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

// NewMessenger is
func NewMessenger(reader *bufio.Reader, writer *bufio.Writer) *Messenger {
	return &Messenger{
		reader: reader,
		writer: writer,
	}
}

// WaitContractCommand waits
func (messenger *Messenger) WaitContractCommand() *ContractCommand {
	command := &ContractCommand{}

	err := messenger.receive(command)
	if err != nil {
		log.Fatalf("wait contract command error: %v", err)
	}

	return command
}

// CallFunction calls
func (messenger *Messenger) CallFunction(request *HookCallRequest) *HookCallResponse {
	response := &HookCallResponse{}

	err := messenger.send(request)
	if err != nil {
		log.Fatal("incorrect remote function call: cannot receive")
	}

	err = messenger.receive(response)
	if err != nil {
		log.Fatal("incorrect remote function call: cannot receive")
	}

	if response.Tag != request.Tag {
		log.Fatal("incorrect remote function call")
	}

	return response
}

func (messenger *Messenger) send(messageToNode interface{}) error {
	jsonData, err := messenger.marshal(messageToNode)
	if err != nil {
		return err
	}

	err = messenger.sendMessageLength(jsonData)
	if err != nil {
		return err
	}

	_, err = messenger.writer.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (messenger *Messenger) sendMessageLength(marshalizedMessage []byte) error {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(len(marshalizedMessage)))
	_, err := messenger.writer.Write(buffer)
	return err
}

func (messenger *Messenger) receive(messageFromNode interface{}) error {
	// Wait for the start of a message
	messenger.blockingPeek(4)

	length, err := messenger.receiveMessageLength()
	if err != nil {
		return err
	}

	// Now read the body of [length]
	messenger.blockingPeek(length)
	buffer := make([]byte, length)
	_, err = io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return err
	}
	return nil
}

func (messenger *Messenger) receiveMessageLength() (int, error) {
	buffer := make([]byte, 4)
	_, err := io.ReadFull(messenger.reader, buffer)
	if err != nil {
		return 0, err
	}

	length := binary.LittleEndian.Uint32(buffer)
	return int(length), nil
}

func (messenger *Messenger) blockingPeek(n int) {
	fmt.Printf("blockingPeek %d bytes\n", n)
	for {
		_, err := messenger.reader.Peek(n)
		if err == nil {
			break
		}
	}
}

func (messenger *Messenger) marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (messenger *Messenger) unmarshal(jsonData []byte, data interface{}) error {
	return json.Unmarshal(jsonData, data)
}
