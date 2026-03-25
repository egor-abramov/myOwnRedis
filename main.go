package main

import (
	"fmt"
	"io"
	"myOwnRedis/handler"
	"myOwnRedis/resp"
	"myOwnRedis/storage"
	"net"
	"os"
	"os/signal"
)

func main() {
	s := storage.NewStorage("dump.json")
	err := s.Load()
	if err != nil {
		fmt.Println(err)
	}
	h := handler.NewHandler(s)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	go func() {
		<-sigChan
		if err := s.Save(); err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}()

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("redis server listening on", listener.Addr())
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn, h)
	}
}

func handleConn(conn net.Conn, h *handler.Handler) {
	defer conn.Close()
	respReader := resp.NewResp(conn)
	for {
		value, err := respReader.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Println("reading error: ", err)
			}
			break
		}
		response := h.Handle(value)
		_, err = conn.Write(response.Marshal())
		if err != nil {
			return
		}
	}
}
