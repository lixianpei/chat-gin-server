package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

// ClientManager [客户端管理中心]-维护连接的客户端信息，同时支持广播消息等
type ClientManager struct {
	//已注册的客户端
	clients     map[*Client]bool
	clientsLock sync.Mutex //并发锁

	//广播消息：需要发送给全站用户的消息接收通道
	broadcast chan []byte

	//接收新注册连接的客户端信息接收通道，注册成功后会将客户端信息保存在clients中
	register chan *Client

	//接收离线的客户端信息接收通道，离线处理成功后会将客户端信息从clients中删除
	unregister chan *Client

	//全局递增的客户端编号
	clientId int64
}

// 实例化一个客户端管理中心
func newClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
	}
}

// ClientRegister 用户新连接事件处理
func (manager *ClientManager) ClientRegister(client *Client) {
	manager.clientsLock.Lock()
	defer func() {
		manager.clientsLock.Unlock()
	}()

	//注册用户
	manager.clients[client] = true

	fmt.Println("EventClientRegister 用户建立连接：", client.id)

	//发送广播消息
	helloMessage := NewMessageTextHello(fmt.Sprintf("欢迎“%d”进入聊天", client.id))
	manager.SendBroadcastMessage([]byte(helloMessage))
}

// ClientUnregister 用户离线事件处理
func (manager *ClientManager) ClientUnregister(client *Client) {
	manager.clientsLock.Lock()
	defer func() {
		manager.clientsLock.Unlock()
	}()

	if _, ok := manager.clients[client]; ok {
		delete(manager.clients, client) //删除用户的在线列表
		close(client.send)              //关闭消息接收的通道
		fmt.Println("用户已下线：", client.id)
	} else {
		fmt.Println("ClientUnregister：未找到用户", client.id)
	}
}

// SendBroadcastMessage 发送广播消息
func (manager *ClientManager) SendBroadcastMessage(message []byte) {
	for client := range manager.clients {
		select {
		case client.send <- message:
			//将广播消息发送给客户端的chan，再由客户端通过conn发送给客户端
			fmt.Printf("给客户端【%d】发送消息: %s \r\n", client.id, message)
		default:
			//没有找到客户端则表示已离线
			manager.ClientUnregister(client)
		}
	}
}

func (manager *ClientManager) run() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ClientManager-recover: ", err)
		}
	}()
	for {
		select {
		case client, _ := <-manager.register:
			//接收连接的客户端
			manager.ClientRegister(client)

		case client, _ := <-manager.unregister:
			//离线客户端处理
			manager.ClientUnregister(client)

		case message, _ := <-manager.broadcast:
			//处理广播消息
			manager.SendBroadcastMessage(message)

		}
	}
}

func InitWebsocket(engine *gin.Engine) {
	manager := newClientManager()
	go manager.run()

	engine.Handle(http.MethodGet, "/chat", func(ctx *gin.Context) {
		serveWs(manager, ctx.Writer, ctx.Request)
	})
}
