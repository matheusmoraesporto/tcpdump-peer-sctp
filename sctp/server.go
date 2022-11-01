package sctp

import (
	"fmt"
	"syscall"
	"unisinos/redes-i/tgb/address"

	sctp "github.com/thebagchi/sctp-go"
)

func (_ ConnectionSCTP) RunServer(ip string, port int, responseAddresses []address.Address) {
	staddr := fmt.Sprintf("%s:%d", ip, port)
	addr, err := sctp.MakeSCTPAddr(SCTPNetowrk, staddr)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	initMsg := NewSCTPInitMessage()
	server, err := sctp.ListenSCTP(SCTPNetowrk, syscall.SOCK_STREAM, addr, &initMsg)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	defer server.Close()

	for {
		fmt.Printf("Servidor executando no endereço %s\n", addr.String())
		// Aguarda um conexão
		conn, err := server.AcceptSCTP()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		// Conexão encontrada
		if remote := conn.RemoteAddr(); nil != remote {
			fmt.Println("Conexão estabelecida com o endereço:", remote)
		}
		// obtém os dados recebidos do client
		go handleClient(conn)
	}
}

func handleClient(conn *sctp.SCTPConn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	data := make([]byte, 8192)
	flag := 0

	for {
		info := &sctp.SCTPSndRcvInfo{}
		len, err := conn.RecvMsg(data, info, &flag)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		if len == 0 {
			fmt.Printf("Conexão com o endereço %s foi encerrada!\n", remoteAddr)
			break
		}
		buffer := string(data[:len])

		fmt.Println("=============================================================")
		fmt.Printf("Pacote sniffado e recebido pelo endereço: %s\n\n%s\n", remoteAddr, buffer)
		fmt.Println("=============================================================")
	}
}
