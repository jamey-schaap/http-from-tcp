package server

import (
	"fmt"
	"http-from-tcp/internal/request"
	"http-from-tcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
		handler:  handler,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	w := response.NewWriter(conn)

	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		headers := response.GetDefaultHeaders(len(err.Error()))
		w.WriteHeaders(headers)

		//hErr := &HandlerError{
		//	StatusCode: response.StatusCodeBadRequest,
		//	Message:    err.Error(),
		//}
		//hErr.Write(conn)
	} else {
		w.WriteStatusLine(response.StatusCodeSuccess)
	}

	response.WriteStatusLine(conn, response.StatusCodeSuccess)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)

	s.handler(w, req)

	return
}
