/*
@Time : 2021/3/9 9:30 AM
@Author : charley
@File : remote_server
*/
package constant

import (
	"github.com/miekg/dns"
	"time"
)

const (
	UDPServer DnsServerType = iota
	TCPServer
	TSLServer

	AddressServer
	GroupServer
)

type DnsServerType int

func (st DnsServerType) String() string {
	switch st {
	case UDPServer:
		return "udp"
	case TCPServer:
		return "tcp"
	case TSLServer:
		return "tls"
	case AddressServer:
		return "AddressServer"
	case GroupServer:
		return "Group"

	default:
		return "Unknown"
	}
}

type RemoteDnsServer interface {
	DnsServerType() DnsServerType
	Query(m *dns.Msg) (*dns.Msg, time.Duration, error)
	String() string
}
