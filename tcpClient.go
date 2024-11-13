package main

import (
	"fmt"
	"io"
	"net"
)

type TcpClient struct {
	chat     *Chat
	conn     *net.TCPConn
	isClosed bool
}

func NewTcpClient(chat *Chat) *TcpClient {
	return &TcpClient{
		chat:     chat,
		isClosed: true,
	}
}

func (tc *TcpClient) Start() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", tc.chat.address)

	if err != nil {
		return err
	}

  fmt.Println("Resolved tcp address: ", tcpAddr.String())

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return err
	}

  fmt.Println("Created tcp dial on:", tcpAddr.String())
	tc.isClosed = false
	tc.conn = conn
	go tc.readLoop()

	return nil
}

func (tc *TcpClient) Close() {
	tc.isClosed = true
	tc.conn.Close()
}

func (tc *TcpClient) readLoop() {
	defer tc.conn.Close()

	buffer := make([]byte, 2048)
	for {
		n, err := tc.conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed:", tc.conn.RemoteAddr().String())
				tc.chat.ReciveSystemMsg(fmt.Sprintf("Connection closed: %v", tc.conn.RemoteAddr().String()))

				tc.chat.sceneManager.SetScene(NewHostClosedScene())
				break
			}
			tc.chat.ReciveSystemMsg(fmt.Sprintf("Error reading from connection: %v", tc.conn.RemoteAddr().String()))
			fmt.Println("Error reading connection;", err)
			break
		}

		chatMsg := NewChatMsgFromBytes(buffer[:n])

		tc.chat.ReciveChatMsg(chatMsg)
	}
}

func (tc *TcpClient) SendMsg(chatMsg ChatMsg) {
	_, err := chatMsg.WriteTo(tc.conn)

	if err != nil {
		tc.chat.ReciveSystemMsg(fmt.Sprintf("Error writing to connection: %v", tc.conn.RemoteAddr().String()))
		fmt.Println("Error writing to connection;", err)
		return
	}

	tc.chat.messages = append([]ChatMsg{chatMsg}, tc.chat.messages...)
}
