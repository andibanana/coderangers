package notifications

import (
	".././cookies"
	".././templating"
	"fmt"
	"net/http"
	// "time"
)

var handler *ConnsHandler

type MsgChan chan []byte
type SSEConn MsgChan

type Conn struct {
	SSEConn SSEConn
	UserID  int
	URL     string
}

type Message struct {
	Message []byte
	To      int
	URL     string
}

type ConnsHandler struct {
	openingConns chan Conn
	closingConns chan Conn
	connSet      map[int]map[Conn]bool
	broadcasts   chan Message
}

const (
	Notifications = "notifications"
	Submissions   = "submissions-listener"
)

func NewSSEConnsHandler() (handler *ConnsHandler) {
	handler = &ConnsHandler{
		broadcasts:   make(chan Message, 1),
		openingConns: make(chan Conn),
		closingConns: make(chan Conn),
		connSet:      make(map[int]map[Conn]bool),
	}

	go handler.handleConns()

	return
}

func (handler *ConnsHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		templating.ErrorPage(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	messageChan := make(SSEConn)
	var conn Conn
	conn.UserID, ok = cookies.GetUserID(req)
	conn.URL = req.URL.Path[1:]
	if !ok {
		templating.ErrorPage(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	conn.SSEConn = messageChan

	handler.openingConns <- conn

	defer func() {
		handler.closingConns <- conn
	}()

	closeNotif := rw.(http.CloseNotifier).CloseNotify()

	go func() {
		<-closeNotif
		handler.closingConns <- conn
	}()

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	for {
		fmt.Fprintf(rw, "data: %s\n\n", <-conn.SSEConn)
		flusher.Flush()
	}
}

func (handler *ConnsHandler) handleConns() {
	for {
		select {
		case conn := <-handler.openingConns:
			if handler.connSet[conn.UserID] == nil {
				handler.connSet[conn.UserID] = make(map[Conn]bool)
			}
			handler.connSet[conn.UserID][conn] = true
		case conn := <-handler.closingConns:
			delete(handler.connSet[conn.UserID], conn)
		case msg := <-handler.broadcasts:
			for conn, _ := range handler.connSet[msg.To] {
				if conn.URL == msg.URL {
					conn.SSEConn <- msg.Message
				}
			}
		}
	}

}

func InitHandler() *ConnsHandler {
	if handler == nil {
		handler = NewSSEConnsHandler()
	}

	return handler
}

func SendMessageTo(userID int, stringMsg string, url string) {
	var message Message
	message.Message = []byte(stringMsg)
	message.To = userID
	message.URL = url
	handler.broadcasts <- message
}
