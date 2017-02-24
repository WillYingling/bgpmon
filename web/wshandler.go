package main

import (
	ws "github.com/gorilla/websocket"
	"net/http"
)

type WSHandler struct {
	m           Message
	startWSConn func(*ws.Conn, Message, chan bool)
	quit        chan bool
}

func NewWSHandler(msg Message, handler func(*ws.Conn, Message, chan bool), q chan bool) *WSHandler {
	return &WSHandler{msg, handler, q}
}

var upgrader = ws.Upgrader{}

func (wsh *WSHandler) Start(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		wsh.startWSConn(conn, wsh.m, wsh.quit)
	} else {
		//Print an error
	}
}

func GetAsByPrefix(wsconn *ws.Conn, m Message, bchan chan bool) {
	bchan <- true
	wsconn.Close()
}

func GetPrefixByAs(wsconn *ws.Conn, m Message, bchan chan bool) {
	bchan <- true
	wsconn.Close()
}
