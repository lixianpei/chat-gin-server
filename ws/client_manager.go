package ws

import (
	"GoChatServer/helper"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

// ClientManager [客户端管理中心]-维护连接的客户端信息，同时支持广播消息等
type ClientManager struct {
	//已注册的客户端
	clients          map[*Client]bool
	lock             sync.Mutex        //并发锁
	clientsUserIdMap map[int64]*Client //保存客户与连接的关系：暂时用手机号

	//广播消息：需要发送给全站用户的消息接收通道
	broadcast chan []byte

	//接收新注册连接的客户端信息接收通道，注册成功后会将客户端信息保存在clients中
	register chan *Client

	//接收离线的客户端信息接收通道，离线处理成功后会将客户端信息从clients中删除
	unregister chan *Client

	//客户端编号
	clientId string
}

// 实例化一个客户端管理中心
func newClientManager() *ClientManager {
	return &ClientManager{
		clients:          make(map[*Client]bool),
		clientsUserIdMap: make(map[int64]*Client),
		broadcast:        make(chan []byte, 256),
		register:         make(chan *Client, 256),
		unregister:       make(chan *Client, 256),
	}
}

// IM 全局唯一的im管理实例
var IM *ClientManager

func InitWebsocket(engine *gin.Engine) {
	IM = newClientManager()
	go IM.run()

	engine.Handle(http.MethodGet, "/chat", func(ctx *gin.Context) {
		serveWs(IM, ctx.Writer, ctx.Request)
	})
}

// ClientRegister 用户新连接事件处理
func (manager *ClientManager) ClientRegister(client *Client) {
	manager.lock.Lock()
	defer func() {
		manager.lock.Unlock()
	}()

	//注册用户
	manager.clients[client] = true
	manager.clientsUserIdMap[client.userId] = client

	//发送广播消息
	helloMessage := NewEntryGroupMessage(fmt.Sprintf("欢迎“%s”进入聊天", client.userInfo.Nickname))
	manager.SendBroadcastMessage([]byte(helloMessage))

	//发送消息给好友
}

// ClientUnregister 用户离线事件处理
func (manager *ClientManager) ClientUnregister(client *Client) {
	manager.lock.Lock()
	defer func() {
		manager.lock.Unlock()
		if err := recover(); err != nil {
			helper.Logger.Errorf("ClientUnregister: %T", err)
		}
	}()

	if _, ok := manager.clients[client]; ok {
		delete(manager.clients, client) //删除用户的在线列表
		close(client.send)              //关闭消息接收的通道
		fmt.Println("用户已下线：", client.userId)
	} else {
		fmt.Println("ClientUnregister：未找到用户", client.userId)
	}
}

// SendBroadcastMessage 发送广播消息
func (manager *ClientManager) SendBroadcastMessage(message []byte) {
	defer func() {
		if err := recover(); err != nil {
			helper.Logger.Errorf("SendBroadcastMessage: %T", err)
		}
	}()
	for client := range manager.clients {
		select {
		case client.send <- message:
			//将广播消息发送给客户端的chan，再由客户端通过conn发送给客户端
			fmt.Printf("给客户端【%d】发送消息: %s \r\n", client.userId, message)
		default:
			//没有找到客户端则表示已离线
			manager.ClientUnregister(client)
		}
	}
}

// SendMessageByUserId 发送消息
func (manager *ClientManager) SendMessageByUserId(message []byte, userId int64) {
	defer func() {
		if err := recover(); err != nil {
			helper.Logger.Errorf("SendMessageByUserId: %T", err)
		}
	}()
	client, ok := manager.clientsUserIdMap[userId]
	if !ok {
		errMessage := fmt.Sprintf("客户端【%d】发送的消息【%s】未能转发出去，客户端未在线", userId, string(message))
		helper.Logger.Errorf(errMessage)
		return
	}
	client.send <- message
}

func (manager *ClientManager) run() {
	defer func() {
		if err := recover(); err != nil {
			helper.Logger.Errorf("ClientManagerRun: %T", err)
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

// OnlineClients 获取在线的所有客户端
func (manager *ClientManager) OnlineClients() []int64 {
	clients := make([]int64, 0)
	for client, _ := range manager.clients {
		clients = append(clients, client.userId)
	}
	return clients
}

// GetClientByUserId 根据userid获取在线客户端
func (manager *ClientManager) GetClientByUserId(userId int64) *Client {
	manager.lock.Lock()
	defer func() {
		manager.lock.Unlock()
	}()
	client, ok := manager.clientsUserIdMap[userId]
	if ok {
		return client
	}
	return nil
}
