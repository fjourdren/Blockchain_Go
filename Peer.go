package main

import "strconv"
import "net"

type Peer struct {
    Index string
	Popularity int
	Host string
	Port int
};




func Construct_peer(index string, popularity int, host string, port int) Peer {
	peer := Peer{Index: index,
				Popularity: popularity,
				Host: host,
				Port: port};

	return peer;
}


func(peer *Peer) get_address() string {
	return peer.Host + ":" + strconv.Itoa(peer.Port);
}


func(peer *Peer) tcp(content []byte) []byte {
	conn, err := net.Dial("tcp", peer.Host + ":" + strconv.Itoa(peer.Port))

	check_error(err);

	//send
	conn.Write([]byte(content));

	//receive
	buff := make([]byte, 800000);
	n, _ := conn.Read(buff);

	conn.Close();
	
	return buff[:n];
}
