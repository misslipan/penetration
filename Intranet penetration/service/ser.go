package main

import (
	_ "bufio"
	_ "flag"
	"fmt"
	"io"
	"net"
	_ "os"
)

type MidServer struct {
	//客户端监听
	clientLis *net.TCPListener
	//后端服务连接
	transferLis *net.TCPListener
	//所有通道
	channels map[int]*Channel
	//当前通道ID
	curChannelId int
}

type Channel struct {
	//通道ID
	id int
	//客户端连接
	client net.Conn
	//后端服务连接
	transfer net.Conn
	//客户端接收消息
	clientRecvMsg chan []byte
	//后端服务发送消息
	transferSendMsg chan []byte
}

//创建一个服务器
func New() *MidServer {
	return &MidServer{
		channels:     make(map[int]*Channel),
		curChannelId: 0,
	}
}

//启动服务
func (m *MidServer) Start(clientPort int, transferPort int) error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", clientPort))
	if err != nil {
		return err
	}
	m.clientLis, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", transferPort))
	if err != nil {
		return err
	}
	m.transferLis, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	go m.AcceptLoop()
	// defer m.Stop()
	return nil
}

//关闭服务
func (m *MidServer) Stop() {

	//循环关闭通道连接
	for _, v := range m.channels {
		v.client.Close()
		v.transfer.Close()
	}
	m.clientLis.Close()
	m.transferLis.Close()
}

//删除通道
func (m *MidServer) DelChannel(id int) {
	chs := m.channels
	delete(chs, id)
	m.channels = chs
}

//处理连接
func (m *MidServer) AcceptLoop() {
	transfer, err := m.transferLis.Accept()
	if err != nil {
		fmt.Println("服务端连接======", err)
		return
	}

	for {
		//获取连接
		client, err := m.clientLis.Accept()
		if err != nil {
			fmt.Println("client连接======", err)
			continue
		}

		// //创建一个通道
		// ch := &Channel{
		// 	id:              m.curChannelId,
		// 	client:          client,
		// 	transfer:        transfer,
		// 	clientRecvMsg:   make(chan []byte, 4096),
		// 	transferSendMsg: make(chan []byte, 4096),
		// }
		// m.curChannelId++

		// //把通道加入channels中
		// chs := m.channels
		// chs[ch.id] = ch
		// m.channels = chs
		// fmt.Println(m.curChannelId)
		go m.MsgLoop(client, transfer)
		// //启一个goroutine处理客户端消息
		// go m.ClientMsgLoop(ch)
		// //启一个goroutine处理后端服务消息
		// go m.TransferMsgLoop(ch)

	}
}

//处理客户端消息
func (m *MidServer) ClientMsgLoop(ch *Channel) {
	defer func() {
		fmt.Println("ClientMsgLoop exit")
	}()
	for {
		select {
		case data, isClose := <-ch.transferSendMsg:
			{
				fmt.Println("客服端", isClose)
				//判断channel是否关闭，如果是则返回
				if !isClose {
					return
				}
				_, err := ch.client.Write(data)
				if err != nil {
					return
				}
			}
		}
	}
}

//处理后端服务消息
func (m *MidServer) TransferMsgLoop(ch *Channel) {
	defer func() {
		fmt.Println("TransferMsgLoop exit")
	}()
	for {
		select {
		case data, isClose := <-ch.clientRecvMsg:
			{
				fmt.Println("服务端", isClose)
				//判断channel是否关闭，如果是则返回
				if !isClose {
					return
				}
				_, err := ch.transfer.Write(data)
				if err != nil {
					return
				}
			}
		}
	}
}

//客户端与后端服务消息处理
// func (m *MidServer) MsgLoop(ch *Channel) {
func (m *MidServer) MsgLoop(conn1 net.Conn, conn2 net.Conn) {
	// defer func() {
	// 	//关闭channel，好让ClientMsgLoop与TransferMsgLoop退出
	// 	// close(ch.clientRecvMsg)
	// 	// close(ch.transferSendMsg)
	// 	// m.DelChannel(ch.id)
	// 	// cli.Close()
	// 	// sers.Close()
	// 	fmt.Println("MsgLoop exit")
	// }()
	// buf := make([]byte, 4096)
	// // for {
	// // 	n, err := ch.client.Read(buf)
	// // 	if err != nil {
	// // 		return
	// // 	}

	// // 	ch.clientRecvMsg <- buf[:n]
	// // 	n, err = ch.transfer.Read(buf)
	// // 	if err != nil {
	// // 		return
	// // 	}

	// // 	ch.transferSendMsg <- buf[:n]
	// // }

	// for {

	// 	n, err := cli.Read(buf)
	// 	fmt.Println("打印读到客户端的消息=========1", string(buf))
	// 	if err != nil {
	// 		if err != io.EOF {
	// 			fmt.Println("打印读到客户端的消息1asdasdasdasd", err)
	// 			return
	// 		}
	// 	}
	// 	// if 0 == n {
	// 	// 	break
	// 	// }

	// 	n, err = sers.Write(buf[:n])
	// 	if err != nil {
	// 		fmt.Println("打印读到客户端的消息1")
	// 		panic(err)

	// 		return
	// 	}

	// 	// ch.clientRecvMsg <- buf[:n]
	// 	fmt.Println("向服务器写数据=====", err)

	// 	n, err = sers.Read(buf)
	// 	fmt.Println("打印读到服务端的消息=========2", string(buf), len(buf[:n]))
	// 	if err != nil && err != io.EOF {
	// 		fmt.Println("打印读到客户端的消息2")
	// 		panic(err)

	// 		return
	// 	}
	// 	// if 0 == n {
	// 	// 	break
	// 	// }
	// 	n, err = cli.Write(buf[:n])
	// 	if err != nil {
	// 		fmt.Println("打印读到客户端的消息4")
	// 		panic(err)
	// 		return
	// 	}

	// 	// ch.transferSendMsg <- buf[:n]

	// 	// fmt.Println("打印服务端消息写到管道", string(buf[:n]))
	// }

	f := func(local net.Conn, remote net.Conn) {
		//defer保证close
		defer local.Close()
		defer remote.Close()
		//使用io.Copy传输两个tcp连接，
		_, err := io.Copy(local, remote)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("join Conn2 end")
	}
	go f(conn2, conn1)
	go f(conn1, conn2)
}

func main() {
	//参数解析
	// localPort := flag.Int("localPort", 8080, "客户端访问端口")
	// remotePort := flag.Int("remotePort", 8888, "服务访问端口")
	// flag.Parse()
	// if flag.NFlag() != 2 {
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }

	ms := New()
	//启动服务
	ms.Start(8080, 8888)
	//循环
	select {}
}

// func readByLine(filename string) ([][]byte, error) {
// 	fp, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer fp.Close()

// 	reader := bufio.NewReader(fp)
// 	lines := make([][]byte, 0)
// 	for {
// 		line, err := reader.ReadBytes('\n')
// 		// fmt.Println(string(line), err)
// 		if err != nil || io.EOF == err {
// 			if string(line) == "" {
// 				break
// 			}
// 		}
// 		lines = append(lines, line)

// 	}
// 	return lines, nil
// }

// func main() {
// 	filename := "config.yaml"

// 	content, _ := readByLine(filename)

// 	for _, k := range content {
// 		fmt.Printf("%s\n", k)
// 	}

// }
