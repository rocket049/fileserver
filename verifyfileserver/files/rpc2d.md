# rpc2d - go语言双向 RPC 调用的库
用数据流重定向的方法实现双向 RPC 调用，高效的实现从服务器 CALLBACK 客户端 API，基于 "net/rpc" 原生库。

安装： `go get gitee.com/rocket049/rpc2d` 或者`go get github.com/rocket049/rpc2d`

*`NewRpcNodeByConn` 函数兼容 `gitee.com/rocket049/pipeconn`*

###提供下列 API 和类型
```
type ProviderType struct {
	Client *rpc.Client
	Data   interface{}
}
type RpcNode struct {
	Server         *rpc.Server
	Client         *rpc.Client
	//private field
}
func Accept(l net.Listener, provider interface{}) (*RpcNode, error)
func NewRpcNode(provider interface{}) *RpcNode
func NewRpcNodeByConn(provider interface{}, conn io.ReadWriteCloser) *RpcNode
func (self *RpcNode) Close()
func (self *RpcNode) Dial(addr string) error
```
###示例
```
//server.go
package main

import (
	"fmt"
	"log"
	"net"

	"gitee.com/rocket049/rpc2d"
)

type Server rpc2d.ProviderType

var count = 0

func (self *Server) Show(arg string, reply *int) error {
	fmt.Printf("Recv: %s, count: %d\n", arg, count)
	*reply = count
	count++
	var ret int
	self.Client.Call("Client.Show", fmt.Sprintf("callback:%s.", arg), &ret)
	return nil
}

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:5678")
	if err != nil {
		log.Fatal("Listen:", err)
	}
	defer l.Close()
	p := new(Server)
	node1, err := rpc2d.Accept(l, p)
	if err != nil {
		log.Fatal("Accept:", err)
	}
	defer node1.Close()
	p.Client = node1.Client
	var s string
	var ret int
	for i := 0; i < 5; i++ {
		s = fmt.Sprintf("server message %d\n", i)
		node1.Client.Call("Client.Show", s, &ret)
		fmt.Printf("Return:%d\n", ret)
	}

	select {}
}


//client.go
package main

import (
	"fmt"
	"log"

	"gitee.com/rocket049/rpc2d"
)

type Client int

var count = 10

func (self *Client) Show(arg string, reply *int) error {
	fmt.Printf("Recv: %s\n", arg)
	*reply = count
	count++
	return nil
}

func main() {
	p := new(Client)
	node1 := rpc2d.NewRpcNode(p)
	err := node1.Dial("127.0.0.1:5678")
	if err != nil {
		log.Fatal("Dial:", err)
	}
	//p.Client = node1.Client
	defer node1.Close()
	var s string
	var ret int
	for i := 0; i < 5; i++ {
		s = fmt.Sprintf("client message %d\n", i)
		node1.Client.Call("Server.Show", s, &ret)
		fmt.Printf("Return: %d\n", ret)
	}
	select {}
}

```