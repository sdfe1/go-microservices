package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// 并发安全
type Connections struct {
	conns map[string]*websocket.Conn
	sync.Mutex
}

var connections = Connections{
	conns: make(map[string]*websocket.Conn),
}

func handler(w http.ResponseWriter, r *http.Request) {
	//1. 将http连接升级为websocket，获得了一个conn链接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("升级失败")
	}

	//读取第一条消息作为用户名
	_, nameBytes, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("读取用户名失败")
		return
	}
	name := string(nameBytes) // 这里获取了用户名

	//2.将新连接添加到全局连接列表
	connections.Lock()
	connections.conns[string(nameBytes)] = conn
	connections.Unlock()

	//3.确保连接最终会关闭
	defer conn.Close()

	//4.主循环，持续处理来自此连接的消息
	for {
		//5.读取客户端发送的消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		//7.在服务端控制台打印消息
		fmt.Printf(string(message))
		//6.向所有连接的客户端广播消息
		broadcast(string(message), name)

	}
}
func broadcast(m string, name string) {
	connections.Lock()
	//遍历每一个客户
	for i, _ := range connections.conns {
		//除了它自己
		if i == name {
			continue
		}
		err := connections.conns[i].WriteMessage(websocket.TextMessage, []byte(name+":"+m))
		if err != nil {
			fmt.Println("发送消息失败:", err)
			break
		}
	}
	connections.Unlock()
}
func main() {
	//注册处理函数到根路径
	http.HandleFunc("/", handler)
	//启动Http服务器监听8080端口
	http.ListenAndServe(":8080", nil)
}
