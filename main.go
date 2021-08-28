package main

import (
	"fmt"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type Roominfo struct {
	roomname string
	chat_id  string
}

type ServerMsg struct {
	chat_id string
	message string
}

type ChatMsg struct {
	roomname string
	chat_id  string
	message  string
}

type UserInfo struct {
	chat_id  string
	roomname string
}

func main() {

	Socket := socketio.NewServer(nil)

	// chat_socket namespace
	Socket.OnConnect("/chat", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected : ", s.ID())
		return nil
	})

	Socket.OnEvent("/chat", "intoroom", func(s socketio.Conn, roominfo *Roominfo) {
		s.Join(roominfo.roomname)
		Socket.BroadcastToRoom("/chat", roominfo.roomname, "server_msg", &ServerMsg{
			chat_id: "[서버]",
			message: fmt.Sprintf("%s님이 들어왔습니다. 환영해주세요", roominfo.chat_id),
		})
	})

	Socket.OnEvent("/chat", "createroom", func(s socketio.Conn, roominfo *Roominfo) {
		s.Join(roominfo.roomname)
	})

	Socket.OnEvent("/chat", "chatmsg", func(s socketio.Conn, chatmsg *ChatMsg) {
		Socket.BroadcastToRoom("/chat", chatmsg.roomname, "chatmsg", chatmsg)
	})

	//Socket.OnEvent("/chat_room_list", )

	Socket.OnEvent("/chat", "leaveroom", func(s socketio.Conn, userinfo *UserInfo) {
		Socket.BroadcastToRoom("/chat", userinfo.roomname, "leavemsg", &ServerMsg{
			chat_id: "[서버]",
			message: fmt.Sprintf("%s님이 방을 나갔습니다.", userinfo.chat_id),
		})
		s.Leave(userinfo.roomname)
	})

	Socket.OnDisconnect("/chat", func(c socketio.Conn, s string) {

	})

	//online_socket namespace
	Socket.OnConnect("/online", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected : ", s.ID())
		return nil
	})

	//Socket.OnEvent("/online", "online")

	//Socket.OnEvent("/online", "online_user_list")

	go Socket.Serve()
	defer Socket.Close()

	http.Handle("/", Socket)
	http.ListenAndServe(":8081", nil)
}
