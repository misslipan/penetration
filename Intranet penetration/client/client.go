package main

import (
	"bufio"
	_ "flag"
	"fmt"
	"io"
	"net"
	_ "os"
)

func handler(r net.Conn, localPort int) {
	// buf := make([]byte, 4096)
	local, _ := net.Dial("tcp", fmt.Sprintf(":%d", localPort))
	joinConn(local, r)
	// fmt.Printf("buf的类型%T======", buf)
	// fmt.Println("buf的长度======", len(buf))
	// for {

	// 	//先从远程读数据
	// 	n, err := r.Read(buf)

	// 	if err != nil && err != io.EOF {
	// 		panic(err)
	// 	}
	// 	if 0 == n {
	// 		fmt.Println("执行")
	// 		break
	// 	}

	// 	if err != nil {
	// 		fmt.Println("请求完成.....")
	// 		//修改了 这个地方 break

	// 		continue
	// 	}

	// 	data := buf[:n]
	// 	fmt.Println("读到服务器的数据buf========11", string(data))
	// 	//建立与本地80服务的连接
	// 	local, err := net.Dial("tcp", fmt.Sprintf(":%d", localPort))

	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	//向80服务写数据
	// 	fmt.Println(string(n))
	// 	n, err = local.Write(data)
	// 	fmt.Println("打印写到本地的buf====222222", string(n))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	//读取80服务返回的数据
	// 	// reader := bufio.NewReader(local)
	// 	// n, err = reader.Read(buf)
	// 	// chanss := rederser(local)
	// 	_, err = local.Read(buf)

	// 	if err != nil && err != io.EOF {
	// 		panic(err)
	// 	}
	// 	if 0 == n {
	// 		break
	// 	}
	// 	//关闭80服务，因为本地80服务是http服务，不是持久连接
	// 	//一个请求结束，就会自动断开。所以在for循环里我们要不断Dial，然后关闭。
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	defer local.Close()

	// 	a := buf[:n]
	// 	fmt.Println("buf==============3", string(a))

	// 	//向远程写数据
	// 	n, err = r.Write(a)
	// 	fmt.Println("打印向远程写数据的bu=======4f", string(a))
	// 	if err != nil {
	// 		fmt.Println("写完退出.....", err)
	// 		continue
	// 	}
	// 	fmt.Println("err=====", err)
	// 	// defer fmt.Println("连接请求中.....")
	// 	// writeser(r, local, chanss)

	// }

}
func joinConn(conn1 net.Conn, conn2 net.Conn) {
	f := func(local net.Conn, remotes net.Conn) {
		defer local.Close()
		defer remotes.Close()
		_, err := io.Copy(local, remotes)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("end")
	}
	go f(conn1, conn2)
	go f(conn2, conn1)
}

func rederser(local net.Conn) <-chan []byte {

	out := make(chan []byte, 1024)
	go func() {
		arr := make([]byte, 1024)
		byteReader := 0
		reader := bufio.NewReader(local)
		for {

			n, err := reader.Read(arr)

			byteReader += n
			if n > 0 {

				out <- arr
			}

			if err != nil {
				break
			}
		}

		close(out)

	}()
	return out

}

func writeser(r net.Conn, local net.Conn, ac <-chan []byte) {

	go func() {

		w := bufio.NewWriter(r)
		for v := range ac {
			fmt.Println(string(v))
			n, err := w.Write(v)
			fmt.Println(n)
			if err != nil {
				panic(err)
				return
			}
		}

		defer local.Close()

		defer w.Flush()

	}()

}

func main() {
	//参数解析
	// host := flag.String("host", "127.0.0.1", "服务器地址")
	// remotePort := flag.Int("remotePort", 8888, "服务器端口")
	// localPort := flag.Int("localPort", 80, "本地端口")
	// flag.Parse()
	// if flag.NFlag() != 3 {
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }
	//建立与服务器的连接
	// remote, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *remotePort))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// go handler(remote, *localPort)

	// select {}

	romte, err := net.Dial("tcp", "192.168.0.111:8888")
	if err != nil {
		fmt.Println(err)
	}
	go handler(romte, 80)
	defer romte.Close()
	select {}
}
