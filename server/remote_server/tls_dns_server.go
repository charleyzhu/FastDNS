/*
@Time : 2021/3/9 10:23 AM
@Author : charley
@File : tls_dns_server
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

type TLSDnsServer struct {
	name    string
	address string
	port    int
}

func (server *TLSDnsServer) DnsServerType() constant.DnsServerType {
	return constant.TSLServer
}

func (server *TLSDnsServer) Query(m *dns.Msg) (*dns.Msg, time.Duration, error) {
	c := new(dns.Client)
	c.Net = "tcp-tls"

	r, rtt, err := c.Exchange(m, net.JoinHostPort(server.address, strconv.Itoa(server.port)))
	if err != nil {
		return nil, 0, err
	}
	return r, rtt, err
}

func (server *TLSDnsServer) String() string {
	return fmt.Sprintf("[ServerName:%s Server Type:%s Address:%s Port:%d]", server.DnsServerType().String(), server.name, server.address, server.port)
}

func NewTSLDnsServer(name, address string, port int) *TLSDnsServer {
	return &TLSDnsServer{
		name:    strings.ToLower(name),
		address: address,
		port:    port,
	}
}
