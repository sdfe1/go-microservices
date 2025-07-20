package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"strings"
)

func main() {
	//读取输入信息，模拟聊天框
	reader := bufio.NewReader(os.Stdin)
	//获取用户名
	fmt.Printf("请输入用户名：")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	//创建客户端
	dl2 := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	//连接
	conn, _, err := dl2.Dial("ws://127.0.0.1:8080", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	//连接成功后首先发送用户名
	err = conn.WriteMessage(websocket.TextMessage, []byte(name))
	if err != nil {
		fmt.Printf("用户名发送失败：%v", err)
		return
	}
	fmt.Printf("%s欢迎来到树洞聊天，祝你每天开心，不开心也没关系！\n", name)
	//持续监听来自服务端的消息
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			//打印收到的来自服务端的消息
			fmt.Println(string(message))
		}
	}()

	//读取用户输入并发送
	for {
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Printf("消息发送失败")
			return
		}
	}
}
