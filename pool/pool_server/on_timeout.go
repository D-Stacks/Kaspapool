package pool_server

import (
	"KPool/miner"
	"fmt"
	"strings"
)

func (p *Pool) onTimeOut(timeOut miner.TimeOutEvent) {
	fmt.Println(timeOut.Reason.Error())
	stringAddrs := timeOut.Conn.RemoteAddr().String()
	splitIndex := strings.LastIndex(stringAddrs, ":")
	addrs, _ := stringAddrs[:splitIndex], stringAddrs[splitIndex:]

	connStatus, found := p.ConnectionStatus[addrs]
	if found {
		connStatus.TimedOut = timeOut
		p.ConnectionStatus[timeOut.Conn.RemoteAddr().String()] = connStatus
		timeOut.Conn.Close()
	} 
	
	fmt.Println(timeOut.Reason.Error())
}