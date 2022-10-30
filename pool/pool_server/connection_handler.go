package pool_server

import "fmt"

func (p *Pool) connectListen() error {
	for {
		conn, err := p.Listener.Accept()
		fmt.Println(conn.RemoteAddr())
		if err != nil {
			return err
		}
		go func() { p.IncommingConnectionChan <- conn }()
	}
}
