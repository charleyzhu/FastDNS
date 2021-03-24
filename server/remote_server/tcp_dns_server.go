/*
@Time : 2021/3/9 10:22 AM
@Author : charley
@File : tcp_dns_server
*/
package remote_server

import (
	"fmt"
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
	"net"
	"strconv"
	"strings"
	"time"
)

type TCPDnsServer struct {
	name      string
	ipAddress string
	port      int
}

func (server *TCPDnsServer) DnsServerType() constant.DnsServerType {
	return constant.TCPServer
}

func (server *TCPDnsServer) Query(m *dns.Msg) (*dns.Msg, time.Duration, error) {
	c := new(dns.Client)
	c.Net = "tcp"
	r, rtt, err := c.Exchange(m, net.JoinHostPort(server.ipAddress, strconv.Itoa(server.port)))
	if err != nil {
		return nil, 0, err
	}
	return r, rtt, err
}

func (server *TCPDnsServer) String() string {
	return fmt.Sprintf("[ServerName:%s Server Type:%s Address:%s Port:%d]", server.DnsServerType().String(), server.name, server.ipAddress, server.port)
}

func NewTCPDnsServer(name, ipAddress string, port int) *TCPDnsServer {
	return &TCPDnsServer{
		name:      strings.ToLower(name),
		ipAddress: ipAddress,
		port:      port,
	}
}
