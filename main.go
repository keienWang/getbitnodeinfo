package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

// Version message structure
type VersionMessage struct {
	Version     int32
	Services    uint64
	Timestamp   int64
	AddrRecv    [26]byte
	AddrFrom    [26]byte
	Nonce       uint64
	UserAgent   []byte
	StartHeight int32
	Relay       bool
}

// Serialize version message to bytes
func (msg *VersionMessage) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write version
	if err := binary.Write(buf, binary.LittleEndian, msg.Version); err != nil {
		return nil, err
	}
	// Write services
	if err := binary.Write(buf, binary.LittleEndian, msg.Services); err != nil {
		return nil, err
	}
	// Write timestamp
	if err := binary.Write(buf, binary.LittleEndian, msg.Timestamp); err != nil {
		return nil, err
	}
	// Write address of receiving node (dummy data for now)
	if err := binary.Write(buf, binary.LittleEndian, msg.AddrRecv); err != nil {
		return nil, err
	}
	// Write address of sending node (dummy data for now)
	if err := binary.Write(buf, binary.LittleEndian, msg.AddrFrom); err != nil {
		return nil, err
	}
	// Write nonce
	if err := binary.Write(buf, binary.LittleEndian, msg.Nonce); err != nil {
		return nil, err
	}
	// Write user agent
	if err := binary.Write(buf, binary.LittleEndian, byte(len(msg.UserAgent))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.UserAgent); err != nil {
		return nil, err
	}
	// Write start height
	if err := binary.Write(buf, binary.LittleEndian, msg.StartHeight); err != nil {
		return nil, err
	}
	// Write relay flag
	relayFlag := byte(0)
	if msg.Relay {
		relayFlag = 1
	}
	if err := binary.Write(buf, binary.LittleEndian, relayFlag); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {
	// Connect to a Bitcoin node
	conn, err := net.Dial("tcp", "203.11.72.110:8333")
	if err != nil {
		log.Fatalf("Failed to connect to node: %v", err)
	}
	defer conn.Close()

	// Prepare a version message
	versionMsg := &VersionMessage{
		Version:     70015, // Protocol version
		Services:    0,     // No services
		Timestamp:   time.Now().Unix(),
		AddrRecv:    [26]byte{},
		AddrFrom:    [26]byte{},
		Nonce:       12345, // Random nonce
		UserAgent:   []byte("/my-go-client:0.1/"),
		StartHeight: 0,
		Relay:       true,
	}

	// Serialize and send the version message
	versionData, err := versionMsg.Serialize()
	if err != nil {
		log.Fatalf("Failed to serialize version message: %v", err)
	}

	_, err = conn.Write(versionData)
	if err != nil {
		log.Fatalf("Failed to send version message: %v", err)
	}

	fmt.Println("Version message sent, waiting for response...")

	// Wait for the response (verack or version)
	// Here, for simplicity, we just read a fixed number of bytes.
	// In a full implementation, you'd parse and handle the incoming message properly.
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

}
