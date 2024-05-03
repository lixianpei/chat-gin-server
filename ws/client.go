package ws

import (
	"GoChatServer/dal/structs"
	"GoChatServer/helper"
	"GoChatServer/service"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	//设置消息写入writer对象的时长，超出后报错：w.Close() write tcp 127.0.0.1:8081->127.0.0.1:60355: i/o timeout
	writeWait = 30 * time.Second

	//服务端会定期向客户端发送 ping 消息，以维持连接的活跃状态。
	//客户端在收到 ping 消息后，需要及时发送 pong 消息作为响应。
	//而 pongWait 参数指定了服务端等待客户端响应 pong 消息的时间，
	//即客户端在收到 ping 消息后需要在一定时间内发送 pong 消息，
	//否则服务端可能会认为客户端已断开连接。
	//pingPeriod 参数的设置应该小于 pongWait 参数，以确保客户端有足够的时间来响应 ping 消息。
	//这样可以确保服务端能够及时检测到客户端的连接状态，并维持连接的活跃性。

	//pongWait 参数的作用是设置服务端等待客户端响应 pong 消息的时间。
	//如果客户端在超过该时间后仍然没有发送 pong 消息，服务端可能会认为客户端已经断开连接，并采取相应的处理措施。
	// Time allowed to read the next pong message from the peer. pong相应消息
	//pongWait = 10 * 60 * time.Second
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	// pingPeriod 参数用于设置服务端定期发送 ping 消息给客户端的时间间隔。这个时间间隔必须小于客户端响应 ping 消息的等待时间 pongWait。
	//pingPeriod = (pongWait * 9) / 10
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer. 超过最大长度将会自动断开连接
	maxMessageSize = 5120
)

type Client struct {
	//im在线管理，后期可以
	clientManager *ClientManager

	//websocket connection
	conn *websocket.Conn

	//Buffered channel of outbound messages. 出站消息的缓冲通道
	send chan []byte

	userId   int64
	userInfo *structs.UserItem
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {

		return true
	}, //是否允许跨域
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.clientManager.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fmt.Println("SetReadDeadlineError:", err.Error())
	}
	c.conn.SetPongHandler(func(string) error {
		err = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			fmt.Println("SetPongHandlerSetReadDeadlineError:", err.Error())
		}
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		fmt.Println("服务器接收到客户端的消息：", string(message))
		if err != nil {
			fmt.Printf("读客户端消息发现错误：%s \r\n", err.Error())
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		//原始消息不处理
		if string(message) == "ping" {
			c.send <- []byte("pong")
			break
		}

		//消息发送和保存
		_, _ = HandleMessageSaveAndSend(string(message), c.userId)

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Printf("服务端监听客户端【%d】的writePump发现错误：%s \r\n", c.userId, ok)
				return
			}

			//在 WebSocket 中，NextWriter 返回一个 io.WriteCloser 接口，该接口用于向连接写入数据。
			//每次调用 NextWriter 都会返回一个新的写入器，用于向客户端发送消息。
			//当消息发送完成后，必须调用 Close 方法关闭该写入器，以确保将消息刷新到连接并释放相关资源
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			writeN, err := w.Write(message)
			if err != nil {
				//TODO
			}
			fmt.Printf("服务端发送给客户端[%d]消息，消息内容：%s \r\n", c.userId, string(message))

			//为了区分消息独立性，此处不建议全部刷数据给客户端，除非同客户端协商处理每个消息的分隔符
			// Add queued chat messages to the current websocket message.
			// 将chan中的其他未发送的消息全部发出去，若使用rang会阻塞进程，导致当前write无法真正发送出去
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline) //使用换行符分割每一条消息
			//	w.Write(<-c.send)
			//}

			//w.Close() 的作用是关闭当前的写入器，将消息刷新到连接并释放资源。
			if err := w.Close(); err != nil {
				fmt.Println("w.Close()", err.Error())
				fmt.Printf("服务端发送给客户端[%d]消息，消息长度[%d]，关闭当前的写入器状态失败：[%s]，消息内容：%s \r\n", c.userId, writeN, err.Error(), string(message))
				return
			}

		case <-ticker.C:
			//触发服务端ping客户端的定时器
			fmt.Printf("服务端主动发送ping消息给客户端：%d ,time:%s \r\n", c.userId, time.Now().Local().Format(time.DateTime))
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("服务端主动发送ping消息给客户端：%d ，发现错误：%s \r\n", c.userId, err.Error())
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(manager *ClientManager, w http.ResponseWriter, r *http.Request) {
	// 从请求头中获取认证信息
	token := strings.TrimSpace(r.Header.Get("token"))
	if len(token) == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		helper.Logger.Error("token为空")
		return
	}
	claims, err := helper.JwtParseChecking(token)
	if err != nil {
		helper.Logger.Error("token解析失败：" + token)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//获取用户信息
	userInfo, _ := service.User.GetMessageUserById(claims.UserId)
	if userInfo == nil || userInfo.ID == 0 {
		helper.Logger.Infof("ws.获取用户[%d]失败", claims.UserId)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		clientManager: manager,
		conn:          conn,
		send:          make(chan []byte, 256),
		userId:        claims.UserId,
		userInfo:      userInfo,
	}
	client.clientManager.register <- client
	helper.Logger.Infof("新用户连接ws：%+v", client)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
