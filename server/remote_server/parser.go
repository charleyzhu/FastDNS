/*
@Time : 2021/3/9 11:52 AM
@Author : charley
@File : parser
*/
package remote_server

import "C"
import (
	"fmt"
	"github.com/charleyzhu/FastDNS/constant"
)

func ParseRemoteServer(serverType, name, address string, port int) (constant.RemoteDnsServer, error) {
	var (
		parseErr error
		parsed   constant.RemoteDnsServer
	)

	switch serverType {
	case constant.UDPServer.String():
		parsed = NewUDPDnsServer(name, address, port)

	case constant.TCPServer.String():
		parsed = NewTCPDnsServer(name, address, port)

	case constant.TSLServer.String():
		parsed = NewTSLDnsServer(name, address, port)

	default:
		parseErr = fmt.Errorf("unsupported rule type %s", serverType)
	}

	return parsed, parseErr
}
