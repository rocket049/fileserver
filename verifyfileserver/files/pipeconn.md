# 用pipeconn帮助go语言程序编写不依靠网络连接的 rpc server/client
包路径： `gitee.com/rocket049/pipeconn`

**[代码托管页面链接](https://gitee.com/rocket049/pipeconn)**

`pipeconn`用标准输入输出和管道模拟 `io.ReadWriteCloser`，可以用于编写不依靠网络连接的 `rpc` 服务器和客户端，通过它程序可以不通过网络也使用`RPC`方式调用另一个本机程序。

服务器方建立连接时调用`pipeconn.NewServerPipeConn()`

客户端建立连接时调用`NewClientPipeConn(progPath , args...)`

下面是示例程序：
#### 服务端程序：
```
//pipe-server.go
package main
import (
	"net/rpc"
	"gitee.com/rocket049/pipeconn"
)

//定义服务类型 Arith
......

func main() {
	arith := new(Arith)
	server := rpc.NewServer()
	server.Register(arith)
	conn := pipeconn.NewServerPipeConn()
	server.ServeConn(conn)
}
```

#### 客户端程序：
```
//pipe-client.go
package main
import (
	"net/rpc"
	"gitee.com/rocket049/pipeconn"
)

func callRpc()  error{
	conn, err := pipeconn.NewClientPipeConn("./pipe-server")
	if err != nil {
		return error
	}
	client := rpc.NewClient(conn)
	defer client.Close()
	//使用 rpc 调用
	......
}
```

**[完整示例程序](https://gitee.com/rocket049/pipeconn/tree/master/rpc)**

**`rpc` 目录中的是一个示例程序。**