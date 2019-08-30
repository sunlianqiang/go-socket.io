package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/socket.io", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		server.JoinRoom("ceres", s)
		return nil
	})

	go func(server *socketio.Server) {
		for i := 0; i < 500; i++ {

			msg := fmt.Sprintf("broadcast progress %v/5", i)
			fmt.Println(msg)
			server.BroadcastToRoom("ceres", "notice", msg)
			time.Sleep(time.Second * 5)
		}
	}(server)

	server.OnEvent("/socket.io", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("notice", "have "+msg)
		server.BroadcastToRoom("ceres", "notice", fmt.Sprintf("send all msg:%v", msg))
	})

	server.OnError("/socket.io", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/socket.io", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
