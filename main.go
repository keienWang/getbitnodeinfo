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
		Version:     70001, // Protocol version
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
	fmt.Println(response)

}

func parseBitcoinResponse(response []byte) {
	if len(response) < 24 {
		fmt.Println("Response too short to be a valid message")
		return
	}

	// 解析消息头
	magic := binary.LittleEndian.Uint32(response[0:4])
	command := string(bytes.Trim(response[4:16], "\x00"))
	length := binary.LittleEndian.Uint32(response[16:20])
	checksum := response[20:24]

	fmt.Printf("Magic: %x\n", magic)
	fmt.Printf("Command: %s\n", command)
	fmt.Printf("Length: %d\n", length)
	fmt.Printf("Checksum: %x\n", checksum)

	if len(response) < int(24+length) {
		fmt.Println("Response body is incomplete")
		return
	}

	// 提取消息体
	body := response[24 : 24+length]

	// 解析 version 消息体
	if command == "version" {
		versionMsg, err := parseVersionMessage(body)
		if err != nil {
			fmt.Printf("Failed to parse version message: %v\n", err)
			return
		}

		// 打印解析后的 version 消息内容
		fmt.Printf("Version: %d\n", versionMsg.Version)
		fmt.Printf("Services: %d\n", versionMsg.Services)
		fmt.Printf("Timestamp: %s\n", time.Unix(versionMsg.Timestamp, 0))
		fmt.Printf("UserAgent: %s\n", versionMsg.UserAgent)
		fmt.Printf("StartHeight: %d\n", versionMsg.StartHeight)
		fmt.Printf("Relay: %v\n", versionMsg.Relay)
	}
}

// 解析 version 消息体的函数
func parseVersionMessage(body []byte) (*VersionMessage, error) {
	reader := bytes.NewReader(body)

	var versionMsg VersionMessage

	// 按顺序解析各个字段
	if err := binary.Read(reader, binary.LittleEndian, &versionMsg.Version); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &versionMsg.Services); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &versionMsg.Timestamp); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &versionMsg.AddrRecv); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &versionMsg.AddrFrom); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &versionMsg.Nonce); err != nil {
		return nil, err
	}

	// 读取 UserAgent 字符串的长度
	var userAgentLength uint8
	if err := binary.Read(reader, binary.LittleEndian, &userAgentLength); err != nil {
		return nil, err
	}

	// 读取 UserAgent 字符串
	userAgent := make([]byte, userAgentLength)
	if err := binary.Read(reader, binary.LittleEndian, &userAgent); err != nil {
		return nil, err
	}
	versionMsg.UserAgent = userAgent

	// 解析 StartHeight
	if err := binary.Read(reader, binary.LittleEndian, &versionMsg.StartHeight); err != nil {
		return nil, err
	}

	// 解析 Relay 字段
	var relay uint8
	if err := binary.Read(reader, binary.LittleEndian, &relay); err != nil {
		return nil, err
	}
	versionMsg.Relay = relay != 0

	return &versionMsg, nil
}

// 连接到比特币节点
func connectToNode() (net.Conn, error) {
	conn, err := net.Dial("tcp", "")
	if err != nil {
		return nil, err
	}
	return conn, nil
}
