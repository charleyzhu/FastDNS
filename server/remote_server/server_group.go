/*
@Time : 2021/3/9 10:51 AM
@Author : charley
@File : server_group
*/
package remote_server

import (
	"fmt"
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/go-ping/ping"
	"github.com/miekg/dns"
	"math"
	"sort"
	"sync"
	"time"
)

const (
	Balancing DnsServerGroupType = iota
	Parallel
	FastTest
)

type DnsServerGroupType int

func (d DnsServerGroupType) String() string {
	switch d {
	case Balancing:
		return "Balancing"
	case Parallel:
		return "Parallel"
	case FastTest:
		return "FastTest"
	default:
		return "Unknown"
	}
}

type DnsServerGroup struct {
	name      string
	groupType DnsServerGroupType
	servers   []constant.RemoteDnsServer
}

func (server *DnsServerGroup) DnsServerType() constant.DnsServerType {
	return constant.GroupServer
}

func (server *DnsServerGroup) Query(r *dns.Msg) (*dns.Msg, time.Duration, error) {
	switch server.groupType {
	case Balancing:
		return server.BalancingQuery(r)
	case Parallel:
		return server.ParallelQuery(r)
	case FastTest:
		return server.FastTestQuery(r)
	default:
		return server.BalancingQuery(r)
	}
}

func (server *DnsServerGroup) BalancingQuery(r *dns.Msg) (*dns.Msg, time.Duration, error) {
	remoteServers := server.servers
	for _, remoteServer := range remoteServers {
		r, rtt, err := remoteServer.Query(r)
		if err == nil {
			return r, rtt, err
		}
	}
	return nil, 0, fmt.Errorf("all servers failed to resolve")
}

func (server *DnsServerGroup) ParallelQuery(r *dns.Msg) (*dns.Msg, time.Duration, error) {

	qmChan := make(chan QueryModel)

	remoteServers := server.servers

	for _, remoteServer := range remoteServers {
		go func(rServer constant.RemoteDnsServer) {
			r, rtt, err := rServer.Query(r)
			rqm := NewQueryModel(r, rtt, err)
			qmChan <- rqm

		}(remoteServer)
	}

	var queryModel *QueryModel
	for i := 0; i <= len(remoteServers); i++ {
		qm := <-qmChan
		if qm.err == nil {
			queryModel = &qm
			break
		}
	}

	if queryModel == nil {
		return nil, 0, fmt.Errorf("all servers failed to resolve")
	} else {
		return queryModel.msg, queryModel.rtt, queryModel.err
	}

}

func (server *DnsServerGroup) FastTestQuery(r *dns.Msg) (*dns.Msg, time.Duration, error) {
	var wg sync.WaitGroup //定义一个同步等待的组

	remoteServers := server.servers
	wg.Add(len(remoteServers))
	qmChan := make(chan AvgRttQueryModel, len(remoteServers))
	for _, remoteServer := range remoteServers {

		go func(rServer constant.RemoteDnsServer, qc chan AvgRttQueryModel) {
			r, rtt, err := rServer.Query(r)
			if err != nil {
				wg.Done()
				return
			}
			if len(r.Answer) <= 0 {
				wg.Done()
				return
			}
			addr := r.Answer[0].String()
			avgRtt := server.TestPing(addr)
			rqm := NewAvgRttQueryModel(r, rtt, avgRtt, err)
			qc <- rqm
			wg.Done()

		}(remoteServer, qmChan)
	}

	var qmArray []AvgRttQueryModel

	wg.Wait()
	close(qmChan)

	for qm := range qmChan {
		qmArray = append(qmArray, qm)
	}

	if len(qmArray) <= 0 {
		return nil, 0, fmt.Errorf("all servers failed to resolve")
	}
	sort.Slice(qmArray, func(i, j int) bool {
		return qmArray[i].ping > qmArray[j].ping
	})
	qm := qmArray[0]
	return qm.msg, qm.rtt, qm.err
}

func (server *DnsServerGroup) TestPing(address string) time.Duration {
	pinger, err := ping.NewPinger(address)
	if err != nil {
		return math.MaxInt32
	}
	pinger.Count = 5
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return math.MaxInt32
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	return stats.AvgRtt
}

func (server *DnsServerGroup) String() string {
	return fmt.Sprintf("[Group name:%s  Type:%s ]", server.name, server.groupType.String())
}

func NewDnsServerGroup(name string, groupType DnsServerGroupType, servers []constant.RemoteDnsServer) *DnsServerGroup {
	return &DnsServerGroup{
		name:      name,
		groupType: groupType,
		servers:   servers,
	}
}
