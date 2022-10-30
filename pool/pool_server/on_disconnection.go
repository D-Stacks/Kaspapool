package pool_server

import "net"

func (p *Pool) onDisconnection(disconnection net.Conn) {
	p.RemoveMiner(disconnection)
}

func (p *Pool) RemoveMiner(disconnection net.Conn) {
	miner := p.MinersByConnection[disconnection]
	p.ExtraNonces.RemoveExtraNonce(miner.State.ExtraNonce)
	delete(p.MinersByConnection, miner.Connection.Conn)
	miner.Connection.Conn.Close()
}