/*
@Time : 2021/3/9 10:05 AM
@Author : charley
@File : remote_dns_server
*/
package remote_server

import (
	"fmt"
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
	"net"
	"strconv"
	"time"

	"strings"
)

type UDPDnsServer struct {
	name      string
	ipAddress string
	port      int
}

func (server *UDPDnsServer) DnsServerType() constant.DnsServerType {
	return constant.UDPServer
}

func (server *UDPDnsServer) Query(m *dns.Msg) (*dns.Msg, time.Duration, error) {
	c := new(dns.Client)

	r, rtt, err := c.Exchange(m, net.JoinHostPort(server.ipAddress, strconv.Itoa(server.port)))
	if err != nil {
		return nil, 0, err
	}
	return r, rtt, err
}

func (server *UDPDnsServer) String() string {
	return fmt.Sprintf("[ServerName:%s Server Type:%s ip Address:%s port:%d]", server.DnsServerType().String(), server.name, server.ipAddress, server.port)
}

func NewUDPDnsServer(name, ipAddress string, port int) *UDPDnsServer {
	return &UDPDnsServer{
		name:      strings.ToLower(name),
		ipAddress: ipAddress,
		port:      port,
	}
}
