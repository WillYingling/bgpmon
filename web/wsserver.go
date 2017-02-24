package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Message struct {
	MType string
}

type WSServer struct {
	currReq    int
	waitlist   map[string]*WSHandler
	maxWaiting int
	basepath   string
	quitChan   chan bool
}

func NewWSServer(base string) *WSServer {
	srv := WSServer{}
	srv.waitlist = make(map[string]*WSHandler)
	srv.maxWaiting = 32
	srv.currReq = 0
	srv.basepath = base

	return &srv
}

func (s *WSServer) HandlePaths() {
	http.HandleFunc(s.basepath, s.Serve)
	http.HandleFunc(s.basepath+"ws", s.ServeWS)

	bchan := make(chan bool)
	s.quitChan = bchan
	go monitor(bchan)
}

func monitor(bchan chan bool) {
	ticker := time.Tick(15 * time.Minute)
	for {
		select {
		case <-bchan:
			// Decrement currReq
		case <-ticker:
			//Search waitlist for stale requests
		}
	}
}

func (s *WSServer) Add(m Message) (key string, err error) {
	if s.currReq < s.maxWaiting {
		key = s.generateKey()
		switch m.MType {
		case "AsByPrefix":
			s.waitlist[key] = NewWSHandler(m, GetAsByPrefix, s.quitChan)
		case "PrefixByAs":
			s.waitlist[key] = NewWSHandler(m, GetPrefixByAs, s.quitChan)
		default:
			//Not support message type
		}
		s.currReq++
	} else {
		err = fmt.Errorf("Too many waiting requests, ignoring")
	}
	return
}

func (s *WSServer) generateKey() string {
	//TODO: make this generate unique string keys
	return fmt.Sprintf("key%d%s", s.currReq, time.Now().Format("15-04-05.000"))
}

func (s *WSServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	keyS := v.Get("key")
	if handle, ok := s.waitlist[keyS]; ok {
		handle.Start(w, r)
	} else {
		w.Write([]byte("Key not regiestered"))
	}
}

func (s *WSServer) Serve(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Recieved connection!\n")
	switch r.Method {
	case "POST":
		dec := json.NewDecoder(r.Body)
		var m Message
		err := dec.Decode(&m)
		if err != nil {
			w.Write([]byte("Error parsing JSON"))
		} else {
			url, err := s.Add(m)
			fmt.Printf("Establishing key: %s:%v\n", url, err)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Hello %v", err)))
			} else {
				w.Write([]byte(url))
			}
		}
	default:
		w.Write([]byte("Not a post request"))
	}
}
