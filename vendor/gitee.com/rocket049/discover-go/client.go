package discover

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/url"
	"time"
)

type DiscoverClient struct {
	Conn  *net.UDPConn
	RAddr *net.UDPAddr
}

func NewClient() *DiscoverClient {
	defer func() {
		e := recover()
		if e != nil {
			log.Println(e)
		}
	}()
	var c *net.UDPConn
	var err error
	raddr, err := net.ResolveUDPAddr("udp4", "239.9.0.99:9000")
	chk(err)

	c, err = net.ListenUDP("udp4", nil)
	chk(err)
	return &DiscoverClient{Conn: c, RAddr: raddr}
}

func (s *DiscoverClient) Query() (res []ServeNode) {
	defer func() {
		e := recover()
		if e != nil {
			log.Println(e)
		}
	}()
	_, err := s.Conn.WriteToUDP([]byte("<query />\n\r"), s.RAddr)
	chk(err)
	res = []ServeNode{}
	var buf [1024]byte
	for {
		s.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
		n, _, err := s.Conn.ReadFromUDP(buf[:])
		if err != nil {
			break
		}
		if n > 0 {
			//log.Println(string(buf[:n]))
			var msg serveData
			xml.Unmarshal(buf[:n], &msg)
			res = append(res, ServeNode{Href: msg.Href, Title: msg.Title, Name: msg.Name})
		}
	}
	return
}

func (s *DiscoverClient) Append(scheme string, port int, uri string, name, title string) bool {
	var msg string

	msg = fmt.Sprintf("<append scheme=\"%s\" port=\"%d\" uri=\"%s\" title=\"%s\" name=\"%s\" />\n\r",
		scheme, port, url.PathEscape(uri), xmlEscape(title), xmlEscape(name))

	s.Conn.WriteToUDP([]byte(msg), s.RAddr)
	s.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
	var buf [64]byte
	n, _, err := s.Conn.ReadFromUDP(buf[:])
	if err != nil || n < 1 {
		return false
	}
	//log.Println(from.String(), string(buf[:n]))
	return true
}

func (s *DiscoverClient) Remove(scheme string, port int, uri string) bool {
	msg := fmt.Sprintf("<remove scheme=\"%s\" port=\"%d\" uri=\"%s\" />\n\r",
		scheme, port, url.PathEscape(uri))
	s.Conn.WriteToUDP([]byte(msg), s.RAddr)
	s.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
	var buf [64]byte
	n, _, err := s.Conn.ReadFromUDP(buf[:])
	if err != nil || n < 1 {
		return false
	}
	//log.Println(from.String(), string(buf[:n]))
	return true
}
