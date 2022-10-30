package pool_server

import (
	"KPool/miner"
	//"fmt"
	//"errors"
	"github.com/pbnjay/memory"
	"net"
	"strings"
	"time"
)

func (p *Pool) onConnection(conn net.Conn) {
	passed := p.checkConnection(conn)
	if passed {
		p.addMiner(conn)
	}

}

func (p *Pool) checkConnection(conn net.Conn) bool {
	if float32(memory.FreeMemory() / memory.TotalMemory()) > 0.9 {
		conn.Close()
		return false
	}
	stringAddrs := conn.RemoteAddr().String()
	splitIndex := strings.LastIndex(conn.RemoteAddr().String(), ":")
	addrs, port := stringAddrs[:splitIndex], stringAddrs[splitIndex:]
	status, found := p.ConnectionStatus[addrs]; if found {
		if len(status.portsConnected) > 20{
			conn.Close()
			return false
		} else if found {
			portFound := status.portsConnected[port]; if portFound {
			conn.Close()
			return false
			}
		} else if time.Now().Before(status.TimedOut.End) {
			conn.Close()
			return false
		}
		status.portsConnected[port] = true
	}
	if !found {
		p.ConnectionStatus[addrs] = connectionStatus{
			Strikes: 0,
			portsConnected: make(map[string]bool),
		}
		p.ConnectionStatus[addrs].portsConnected[port] = true
	}
	return true

}

func (p *Pool) addMiner(conn net.Conn) {

	extraNonce1 := p.ExtraNonces.GetExtraNonce()

	miner := miner.NewMiner(conn,extraNonce1, p.Jobs.Cache, p.IncommingTimeOutChan, p.Jobs.SubmitBlock)
	
	p.Miners = append(p.Miners, miner)
	miner.Run()
}