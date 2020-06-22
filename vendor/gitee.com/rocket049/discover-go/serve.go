package discover

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"sync"
)

/*
224.0.0.0～224.0.0.255为预留的组播地址（永久组地址），地址224.0.0.0保留不做分配，其它地址供路由协议使用。
224.0.1.0～238.255.255.255为用户可用的组播地址（临时组地址），全网范围内有效。
239.0.0.0～239.255.255.255为本地管理组播地址，仅在特定的本地范围内有效。
本协议地址:
239.9.0.99:9000
*/

type DiscoverServer struct {
	Conn     *net.UDPConn
	Services []ServeNode
	Lock     sync.RWMutex
	ReadOnly bool
}

func newServerConn(c *net.UDPConn) *DiscoverServer {
	return &DiscoverServer{Conn: c, Services: []ServeNode{}}
}

func NewServer() *DiscoverServer {
	return &DiscoverServer{Services: []ServeNode{}}
}

func Serve() (err error) {
	defer func() {
		recover()
	}()
	addr, err := net.ResolveUDPAddr("udp4", "239.9.0.99:9000")
	chk(err)
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	chk(err)
	server := newServerConn(conn)

	buf := make([]byte, 1024)
	var n int
	var from *net.UDPAddr
	for {
		n, from, err = conn.ReadFromUDP(buf)
		chk(err)
		if n > 0 {
			server.ParseMsg(buf[:n], from)
		}
	}
}

func (s *DiscoverServer) ParseMsg(msg []byte, from *net.UDPAddr) {
	defer func() {
		e := recover()
		if e != nil {
			log.Println(e)
		}
	}()

	n := bytes.Index(msg, []byte{' '})
	//log.Println(string(msg[1:n]))

	switch string(msg[1:n]) {
	case "query":
		s.responseQuery(msg, from)
	case "append":
		s.responseAppend(msg, from)
	case "remove":
		s.responseRemove(msg, from)
	}

}

func (s *DiscoverServer) responseQuery(data []byte, from *net.UDPAddr) {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	for _, serve := range s.Services {
		dgam := fmt.Sprintf("<serve href=\"%s\" title=\"%s\"  name=\"%s\" />\n\r", serve.Href, xmlEscape(serve.Title), xmlEscape(serve.Name))
		//log.Println("resp:", dgam)
		s.Conn.WriteToUDP([]byte(dgam), from)
	}
}

func (s *DiscoverServer) responseAppend(data []byte, from *net.UDPAddr) {
	if s.ReadOnly {
		return
	}

	var msg appendData
	err := xml.Unmarshal(data, &msg)
	chk(err)
	url := createUrl(msg.Scheme, msg.Port, msg.Uri, from.IP.String())

	var isExist bool = false

	for i := range s.Services {
		if s.Services[i].Href == url {
			isExist = true
			s.Lock.Lock()
			s.Services[i].Name = msg.Name
			s.Services[i].Title = msg.Title
			s.Lock.Unlock()
			break
		}
	}
	if isExist == false {
		s.Lock.Lock()
		s.Services = append(s.Services, ServeNode{Href: url, Title: msg.Title, Name: msg.Name})
		s.Lock.Unlock()
	}

	//log.Println("response ok")
	s.Conn.WriteToUDP([]byte("<response name=\"ok\" />\n\r"), from)

}

func (s *DiscoverServer) responseRemove(data []byte, from *net.UDPAddr) {
	if s.ReadOnly {
		return
	}

	var msg removeData
	err := xml.Unmarshal(data, &msg)
	chk(err)
	url := createUrl(msg.Scheme, msg.Port, msg.Uri, from.IP.String())

	s.Lock.RLock()
	services := make([]ServeNode, 0, len(s.Services))

	for i := range s.Services {
		if s.Services[i].Href != url {
			services = append(services, s.Services[i])
		}
	}
	s.Lock.RUnlock()
	s.Lock.Lock()
	s.Services = services
	s.Lock.Unlock()
	s.Conn.WriteToUDP([]byte("<response name=\"ok\" />\n\r"), from)
}

func (s *DiscoverServer) Append(scheme, ip, port, uri, name, title string) {
	url := createUrl(scheme, port, uri, ip)

	var isExist bool = false

	for i := range s.Services {
		if s.Services[i].Href == url {
			isExist = true
			s.Lock.Lock()
			s.Services[i].Name = name
			s.Services[i].Title = title
			s.Lock.Unlock()
			break
		}
	}
	if isExist == false {
		s.Lock.Lock()
		s.Services = append(s.Services, ServeNode{Href: url, Title: title, Name: name})
		s.Lock.Unlock()
	}
}

func (s *DiscoverServer) Remove(scheme, ip, port, uri string) {
	url := createUrl(scheme, port, uri, ip)

	s.Lock.RLock()
	services := make([]ServeNode, 0, len(s.Services))

	for i := range s.Services {
		if s.Services[i].Href != url {
			services = append(services, s.Services[i])
		}
	}
	s.Lock.RUnlock()
	s.Lock.Lock()
	s.Services = services
	s.Lock.Unlock()
}

func (s *DiscoverServer) Serve(readonly bool) {
	defer func() {
		recover()
	}()
	s.ReadOnly = readonly

	addr, err := net.ResolveUDPAddr("udp4", "239.9.0.99:9000")
	chk(err)
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	chk(err)

	s.Conn = conn

	buf := make([]byte, 1024)
	var n int
	var from *net.UDPAddr
	for {
		n, from, err = conn.ReadFromUDP(buf)
		chk(err)
		if n > 0 {
			s.ParseMsg(buf[:n], from)
		}
	}
}
