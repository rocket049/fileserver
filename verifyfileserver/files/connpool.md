# connpool:go语言TCP连接池

### 安装：
```
	go get -v -u github.com/rocket049/connpool
	go get -v -u gitee.com/rocket049/connpool
```

### 说明
`rocket049/connpool`包是本人用go语言开发的，提供一个通用的TCP连接池，初始化参数包括最高连接数、超时秒数、连接函数，放回连接池的连接被重新取出时，如果已经超时，将会自动重新连接；如果没有超时，连接将被复用。

### 可调用的函数
```
type Conn
    func (s *Conn) Read(p []byte) (int, error)
    func (s *Conn) Write(p []byte) (int, error)
    func (s *Conn) Close() error
    func (s *Conn) Timeout() bool
type Pool
    func NewPool(max, timeout int, factory func() (net.Conn, error)) *Pool
    func (s *Pool) Get() (*Conn, error)
    func (s *Pool) Put(conn1 *Conn)
    func (s *Pool) Close()
```

### 调用示例
```
import "gitee.com/rocket049/connpool"

func factory() (net.Conn,error) {
	return net.Dial("tcp","127.0.0.1:7060")
}

func UsePool() {
	// 设置连接池大小为10,超时秒数为30,连接函数为 factory
	pool1 := connpool.NewPool(10, 30 ,factory)
	defer pool1.Close()
	var wg sync.WaitGroup
	for i:=0;i<50;i++ {
		wg.Add(1)
		go func(n int){
		    // 取出连接
			conn ,err := pool1.Get()
			if err!=nil {
				...
			}
			//发送数据
			_,err = conn.Write( msg )
			if err!=nil{
				...
			}
			//读取数据
			n1,err := conn.Read( buf )
			if err!=nil{
				...
			}
			//超时检测和重连
			if conn.Timeout() {
				pool1.Put(conn)
				conn ,err := pool1.Get()
				...
			}
			//放回连接池
			pool1.Put(conn)
			wg.Done()
		}(i)
	}
	wg.Wait()

}
```