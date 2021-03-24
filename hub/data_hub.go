/*
@Time : 2021/3/10 2:50 PM
@Author : charley
@File : data_hub
*/
package hub

import (
	"github.com/charleyzhu/FastDNS/config"
	"github.com/charleyzhu/FastDNS/log"
	"github.com/miekg/dns"
)

type DataHub struct {
	config *config.Config
	listen []*dns.Server
}

func (dh *DataHub) ListenAndServe() {
	//var listen []*dns.Server
	//
	//for _, localServer := range dh.config.Listen {
	//	log.Logger.Infof("listen server %s port:%s ", localServer.ServerType.String(), localServer.port)
	//	srv := &dns.Server{Addr: ":" + localServer.port, Net: localServer.ServerType.String()}
	//	listen = append(listen, srv)
	//}
	//log.Logger.Infof("listen server %s port:%s ", ls.ServerType.String(), ls.port)
	//srv := &dns.Server{Addr: ":" + ls.port, Net: ls.ServerType.String()}
	//srv.Handler = &ls
	//if err := srv.ListenAndServe(); err != nil {
	//	log.Logger.Errorf("Failed to set %s listener %s\n", ls.ServerType.String(), err.Error())
	//}

}

func (dh *DataHub) Restart() {

}

func (dh *DataHub) Shutdown() {
	for _, srv := range dh.listen {
		err := srv.Shutdown()
		if err != nil {
			log.Logger.Error(err)
		}
	}
}
