package tcp

import (
	"bufio"
	"fmt"
	"net"
)

func (_ ConnectionTCP) RunClient(ipLocal, ipRemote string, port int) []string {
	localAddr := HandleTCPAddress(ipLocal, port)
	remoteAddr := HandleTCPAddress(ipRemote, port)

	connection, err := net.DialTCP(TCPProtocol, localAddr, remoteAddr)
	if err != nil {
		fmt.Printf("Client side: Errro -> %s\n", err)
		return nil
	}

	defer func() {
		if err := connection.Close(); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Client side: conexão encerrada com o endereço %s\n", remoteAddr.IP)
		}
	}()

	return waitPackets(connection)
}

func waitPackets(connection *net.TCPConn) []string {
	packets := make([]string, 10)
	for {
		data, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Printf("Client side: Erro -> %s\n", err)
			return nil
		}
		packets = append(packets, data)

		if len(packets) == 10 {
			break
		}
	}
	return packets
}
