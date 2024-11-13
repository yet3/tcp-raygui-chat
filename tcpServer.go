package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type TcpServer struct {
	listener net.Listener
	isClosed bool
	chat     *Chat

	connections      map[net.Conn]bool
	connectionsMutes sync.Mutex
}

func NewTcpServer(chat *Chat) *TcpServer {
	return &TcpServer{
		chat:        chat,
		connections: make(map[net.Conn]bool),
		isClosed:    true,
	}
}

func (ts *TcpServer) Start() error {
	listener, err := net.Listen("tcp", ts.chat.address)

	if err != nil {
		return err
	}

	fmt.Println("Created tcp listener on:", ts.chat.address)

	ts.isClosed = false
	ts.listener = listener

	go ts.acceptLoop()

	return nil
}

func (ts *TcpServer) Close() {
	fmt.Println("Closing listener on:", ts.chat.address)

	ts.connectionsMutes.Lock()
	for conn := range ts.connections {
		fmt.Println("Closing and removing connection", conn.RemoteAddr().String())
		delete(ts.connections, conn)
		conn.Close()
	}
	ts.connectionsMutes.Unlock()

	ts.isClosed = true
	ts.listener.Close()
}

func (ts *TcpServer) acceptLoop() {
	fmt.Println("Listening for connections on:", ts.chat.address)

	for {
		conn, err := ts.listener.Accept()
		if err != nil {
			if ts.isClosed {
				break
			}

			fmt.Printf("Error accepting connection: %v\n", err)
			ts.chat.ReciveSystemMsg(fmt.Sprintf("Error acceptin connection: %v", err.Error()))
			continue
		}

		ts.connectionsMutes.Lock()
		ts.connections[conn] = true
		ts.connectionsMutes.Unlock()

		fmt.Println("Connected to:", conn.RemoteAddr().String())

		ts.chat.ReciveSystemMsg(fmt.Sprintf("User conntected: %v", conn.RemoteAddr().String()))

		go ts.readLoop(conn)
	}
}

func (ts *TcpServer) readLoop(conn net.Conn) {
	defer func() {
		ts.connectionsMutes.Lock()
		if _, ok := ts.connections[conn]; ok {
			fmt.Println("Closing and removing connection", conn.RemoteAddr().String())
			delete(ts.connections, conn)
			conn.Close()
		}
		ts.connectionsMutes.Unlock()
	}()

	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed:", ts.chat.address)
				ts.chat.ReciveSystemMsg(fmt.Sprintf("Connection closed: %v", ts.chat.address))
				break
			}
			ts.chat.ReciveSystemMsg(fmt.Sprintf("Error reading from connection: %v", ts.chat.address))
			fmt.Println("Error reading from connection:", err)
			break
		}

		chatMsg := NewChatMsgFromBytes(buffer[:n])

		ts.chat.ReciveChatMsg(chatMsg)
		ts.BroadcastMessage(chatMsg, &conn)
	}
}

func (ts *TcpServer) BroadcastMessage(chatMsg ChatMsg, excludeConn *net.Conn) {
	ts.connectionsMutes.Lock()
	for conn := range ts.connections {
		if excludeConn != nil && conn == *excludeConn {
			continue
		}

		_, err := chatMsg.WriteTo(conn)

		if err != nil {
			ts.chat.ReciveSystemMsg(fmt.Sprintf("Error writing to connection: %v", conn.RemoteAddr().String()))
			fmt.Println("Error writing to connection:", err)
			continue
		}
	}
	ts.connectionsMutes.Unlock()
}
