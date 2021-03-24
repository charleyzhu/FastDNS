/*
@Time : 2021/3/11 10:42 AM
@Author : charley
@File : query_model
*/
package remote_server

import (
	"github.com/miekg/dns"
	"time"
)

type QueryModel struct {
	msg *dns.Msg
	rtt time.Duration

	err error
}

func NewQueryModel(msg *dns.Msg, rtt time.Duration, err error) QueryModel {
	return QueryModel{
		msg: msg,
		rtt: rtt,
		err: err,
	}
}

type AvgRttQueryModel struct {
	msg  *dns.Msg
	rtt  time.Duration
	ping time.Duration

	err error
}

func NewAvgRttQueryModel(msg *dns.Msg, rtt, ping time.Duration, err error) AvgRttQueryModel {
	return AvgRttQueryModel{
		msg:  msg,
		rtt:  rtt,
		ping: ping,
		err:  err,
	}
}
