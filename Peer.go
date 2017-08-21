package main

import "net/http"
import "strconv"

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


func(peer *Peer) request(url string) (*http.Response, error) {
	resp, err := http.Get(peer.Host + ":" + strconv.Itoa(peer.Port) + url);
	return resp, err;
}
