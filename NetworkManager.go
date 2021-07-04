package main

import "fmt"
import "net"
import "encoding/json"
import "strconv"
import "math/rand"

var nb_peer_to_have = 10;

type NetworkManager struct {
    Me Peer
	Blockchain *Blockchain `json:"-"`
	Peers []Peer
	LastBlockIndex int
};




func NetworkManagerFromJSON(content []byte) (*NetworkManager, error) {
    NetworkManagerDist := new(NetworkManager);
    err := json.Unmarshal(content, NetworkManagerDist);
    return NetworkManagerDist, err;
}



func Construct_NetworkManager(index string, popularity int, host string, port int, blockchain *Blockchain) NetworkManager {
	me := Construct_peer(index, popularity, host, port);

	networkManager := NetworkManager{Me: me,
									Blockchain: blockchain};

	return networkManager;
}


func Construct_Init_NetworkManager(index string, popularity int, host string, port int, blockchain *Blockchain, initHost string, initPort int) NetworkManager {
	networkManager := Construct_NetworkManager(index, popularity, host, port, blockchain);

	initpeer := Construct_peer("0", 0, initHost, initPort)
	networkManager.add_Peer(initpeer);

	packet := Construct_Packet(networkManager.Me, "SINGLE", "JOIN", nil);
	return_value := networkManager.Peers[0].tcp(packet.to_json()); //send packet
	networkManager.process_packet(return_value); //process return

	packetLatestBlock := Construct_Packet(networkManager.Me, "SINGLE", "GETLATESTBLOCK", nil);
	return_valuepacketLatestBlock := networkManager.Peers[0].tcp(packetLatestBlock.to_json()); //send packet
	_, blockJSON := networkManager.process_packet(return_valuepacketLatestBlock); //process return
	block := block_json_to_object(blockJSON);

	networkManager.syncChain(networkManager.Peers[0], block.Index);

	return networkManager;
}


func(networkManager *NetworkManager) server_tcp() {
    ln, err := net.Listen("tcp", ":" + strconv.Itoa(networkManager.Me.Port));
    check_error(err);

    defer ln.Close();

    for {
    	conn, err := ln.Accept();
    	check_error(err);

		buffer := make([]byte, 4096);
 		n, err := conn.Read(buffer);
 		check_error(err);

		reply, _ := networkManager.process_packet(buffer[:n]);
		conn.Write(reply);

    	conn.Close();
    }
}

								
																	//response, return value
func(networkManager *NetworkManager) process_packet(content []byte) ([]byte, []byte) {
	packet := packet_json_to_object(content);


	switch packet.Type {
		case "SINGLE":
			switch packet.Name {
				case "JOIN":
					if !networkManager.has_peer(packet.Sender) {
						networkManager.add_Peer(packet.Sender);
					}

					//answer the last block id, popularity and peer list
					networkManager.LastBlockIndex = networkManager.Blockchain.get_latest_block().Index; //update last block index
  					networkManagerJSON := networkManager.to_json();
  					packetJoinAnswer := Construct_Packet(networkManager.Me, "SINGLE", "JOINANSWER", networkManagerJSON);
  					payload := packetJoinAnswer.to_json();

  					return payload, nil;

				case "JOINANSWER":
					networkManagerDist, _ := NetworkManagerFromJSON(packet.Content);
					networkManager.update_Peer(networkManager.Peers[0], networkManagerDist.Me);

					return nil, packet.Content;

				case "GETLATESTBLOCK":
					block := networkManager.Blockchain.get_latest_block();
					blockJSON := block.to_json();
					answer := Construct_Packet(networkManager.Me, "SINGLE", "GETLATESTBLOCKANSWER", blockJSON);
  					payload := answer.to_json();

  					return payload, nil;

				case "GETLATESTBLOCKANSWER":
					return nil, packet.Content;

				case "DOWNLOADBLOCK":
					startandstop := packet_getBlocks_json_to_object(packet.Content);

					if startandstop.Start > startandstop.Stop {
		  				temp := startandstop.Start;
		  				startandstop.Start = startandstop.Stop;
		  				startandstop.Stop = temp;
		  			}

		  			if startandstop.Start < 0 {
		  				startandstop.Start = 0;
		  			}

		  			var outBlocks []Block;

		  			index := startandstop.Start;
		  			for index <= startandstop.Stop {
			        	block := networkManager.Blockchain.get_block(index);
		  				outBlocks = append(outBlocks, block);
		  				index++;
		  			}

			        payload, _ := json.Marshal(outBlocks);



			      	answer := Construct_Packet(networkManager.Me, "SINGLE", "DOWNLOADBLOCKANSWER", payload);
  					answerJson := answer.to_json();

  					return answerJson, nil;

				case "DOWNLOADBLOCKANSWER":
					return nil, packet.Content;

				default:
					fmt.Println("Unknow packet name.")
			}

		case "BROADCAST":
			if !broadcastManager_has_packet(packet.Index) {
				switch packet.Name {
					case "FOUNDBLOCK":
						fmt.Println(packet.Index)

						index, _ := strconv.Atoi(string(packet.Content));
						if networkManager.syncChain(packet.Sender, index) {
							packet.Sender = networkManager.Me;
							networkManager.broadcast(packet.to_json());
						}

					default:
						fmt.Println("Unknow packet name.")
				}
			}

		default:
			fmt.Println("Unknow packet type.")
	}

	return nil, nil;
	
}


func(networkManager *NetworkManager) server() {

	go networkManager.server_tcp();

}


func(networkManager *NetworkManager) add_Peer(peer Peer) {
	networkManager.Peers = append(networkManager.Peers, peer);
	networkManager.Me.Popularity = len(networkManager.Peers);
}


func(networkManager *NetworkManager) remove_Peer(peer Peer) {
	index := networkManager.get_peer_index(peer);

	networkManager.Peers = append(networkManager.Peers[:index], networkManager.Peers[index+1:]...);
	networkManager.Me.Popularity = len(networkManager.Peers);
}


func(networkManager *NetworkManager) update_Peer(oldPeer Peer, newPeer Peer) {

	networkManager.remove_Peer(oldPeer);
	networkManager.add_Peer(newPeer);

}


func(networkManager *NetworkManager) get_peer_index(peer Peer) int {

	for index, peerElement := range networkManager.Peers {
		if peerElement.Index == peer.Index {
			return index;
		}
	}

	return -1;

}


func(networkManager *NetworkManager) has_peer(peer Peer) bool {

	index := networkManager.get_peer_index(peer);

	if index == -1 {
		return false;
	} else {
		return true;
	}

}


func(networkManager *NetworkManager) broadcast(content []byte) {
	packet := packet_json_to_object(content);

	if !broadcastManager_has_packet(packet.Index) {
		for _, peer := range networkManager.Peers {
			err := peer.tcp(content);
			if err != nil {
			    networkManager.remove_Peer(peer);
			}
		}

		broadcastManager_add_packet(packet.Index);
	}

}


func(networkManager *NetworkManager) download_block(peer Peer, startIndex int, stopIndex int) []Block {

	packet_startstop := Construct_Packet_getBlocks(startIndex, stopIndex);
	packet := Construct_Packet(networkManager.Me, "SINGLE", "DOWNLOADBLOCK", packet_startstop.to_json());
  	payload := packet.to_json();

  	return_valuepacket := networkManager.Peers[0].tcp(payload); //send packet
	_, blocks := networkManager.process_packet(return_valuepacket); //process return

	out := make([]Block, 0);
	json.Unmarshal(blocks, &out);

	return out;
}


func(networkManager *NetworkManager) syncChain(peer Peer, index int) bool {

	blockchain_test := networkManager.Blockchain;

	blocks := networkManager.download_block(peer, 0, index);
	blockchain_test.Chain = blocks;

	if blockchain_test.is_valid() {
		networkManager.Blockchain = blockchain_test;
		fmt.Println("Blockchain downloaded.")
		return true;
	}

	return false;
}


func(networkManager *NetworkManager) randomPeer(content []byte) {

	n := rand.Int() % len(networkManager.Peers);
	peer := networkManager.Peers[n];
	err := peer.tcp(content);
	if err != nil {
		networkManager.remove_Peer(peer);
	}

}


func(networkManager *NetworkManager) to_json() []byte {
	payload, _ := json.Marshal(networkManager);
	return payload;
}