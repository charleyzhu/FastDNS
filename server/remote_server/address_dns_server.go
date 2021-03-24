/*
@Time : 2021/3/10 11:51 PM
@Author : charley
@File : address_dns_server
*/
package remote_server

import (
	"fmt"
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
	"net"
	"time"
)

type AddressDnsServer struct {
	address string
}

func (ads *AddressDnsServer) DnsServerType() constant.DnsServerType {
	return constant.AddressServer
}

func (ads *AddressDnsServer) Query(m *dns.Msg) (*dns.Msg, time.Duration, error) {
	msg := &dns.Msg{}
	msg.SetReply(m)
	switch m.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600},
			A:   net.ParseIP(ads.address),
		})
	default:
		return nil, 0, fmt.Errorf("unsupported Question type %d", m.Question[0].Qtype)
	}
	return msg, 3600, nil
}

func (ads *AddressDnsServer) String() string {
	return fmt.Sprintf("[Server Type:%s Address:%s]", ads.DnsServerType().String(), ads.address)
}

func NewAddressDnsServer(address string) *AddressDnsServer {
	return &AddressDnsServer{
		address: address,
	}
}
