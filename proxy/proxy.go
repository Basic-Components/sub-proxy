//
// 发布订阅模式的代理组件,沟通发布者和订阅者,降低发布者负载
//
package proxy

import (
	log "github.com/sirupsen/logrus"

	"github.com/Basic-Components/sub-proxy/consts"
	loadconfig "github.com/Basic-Components/sub-proxy/loadconfig"

	zmq "github.com/pebbe/zmq4"
)

// 代理本体
func Run(config loadconfig.Config) {
	//  Prepare our sockets
	frontend, _ := zmq.NewSocket(zmq.XSUB)
	defer frontend.Close()

	backend, _ := zmq.NewSocket(zmq.XPUB)
	defer backend.Close()

	if config.Conflate {
		frontend.SetConflate(true)
		backend.SetConflate(true)
	} else {
		if config.RCVHWM >= 0 {
			frontend.SetRcvhwm(config.RCVHWM)
		}
		if config.SNDHWM >= 0 {
			backend.SetSndhwm(config.SNDHWM)
		}
	}
	frontend.Connect(config.FrontendURL)
	backend.Bind(config.BackendURL)

	err := zmq.Proxy(frontend, backend, nil)
	log.WithFields(log.Fields{
		consts.TYPE: consts.NAME,
	}).Fatalln("Proxy interrupted:", err)

}
