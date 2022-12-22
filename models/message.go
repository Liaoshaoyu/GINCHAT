package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FromId   uint
	TargetId uint
	Type     int // 消息类型 私聊 群聊 广播
	Media    int // 音频 图片 文字等
	Content  string
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他消息统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射
var clientMap = make(map[int64]*Node, 0)

// 读写锁
var rwLocker = sync.RWMutex{}

func Chat(writer http.ResponseWriter, request *http.Request) {
	// 1. 获取并校验参数
	var query = request.URL.Query()
	//var token = query.Get("token")
	var userId, _ = strconv.ParseInt(query.Get("userId"), 10, 64)
	//var msgType = query.Get("token")
	//var targetId = query.Get("targetId")
	//var content = query.Get("content")

	var isValid = true // to-do:checkToken()
	var conn, err = (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 2. 获取连接
	var node = &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	clientMap[0] = node
	// 3. 用户关系
	// 4.userid与node绑定，并加锁
	// 5. 发送逻辑
	sendMsg(userId, []byte("欢迎进行聊天室！"))
	sendProc(node)
	// 6. 接收逻辑
	recvProc(node)
}
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue: // congDataQueue读取数据，给data
			var err = node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		default:
			print("")
		}
	}
}

func recvProc(node *Node) {
	for {
		var _, data, err = node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] <<<<< ", data)
	}
}

var udpsendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data // 写入数据到udpsendChan
}

func init() {
	go udpSendProc()
	go udpRecvProc()

}

func udpSendProc() {
	var con, err = net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			var _, err = con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		default:
			print("")
		}
	}
}
func udpRecvProc() {
	var con, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		var n, err = con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}
func dispatch(data []byte) {
	var msg = Message{}
	var err = json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1:
		sendMsg(int64(msg.TargetId), data)
	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	var node, ok = clientMap[userId]
	rwLocker.RUnlock()
	//ok = true
	fmt.Println(userId)
	if ok {
		node.DataQueue <- msg
	}
}
