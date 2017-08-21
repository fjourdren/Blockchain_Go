package main

import "fmt"
import "net/http"
import "encoding/json"
import "strconv"
import "io/ioutil"
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

	response, err := networkManager.Peers[0].request("/join?index=" + networkManager.Me.Index + "&popularity=" + strconv.Itoa(networkManager.Me.Popularity) + "&host=" + networkManager.Me.Host + "&port=" + strconv.Itoa(networkManager.Me.Port));
	if err != nil {
	    panic("Unknown init host.");
	}


	//response to []byte
	body, err := ioutil.ReadAll(response.Body);
	if err != nil {
	    panic(err.Error());
	}

	//extract data
	networkManagerDist, _ := NetworkManagerFromJSON(body);

	//update init Peer
	//networkManager.update_Peer(networkManager.Peers[0], networkManagerDist.Me);
	networkManager.Peers[0] = networkManagerDist.Me;

	//sync chain
	networkManager.syncChain(networkManagerDist.Me, networkManagerDist.LastBlockIndex);

	for _, peerElement := range networkManagerDist.Peers {
		if peerElement != networkManager.Me {
			networkManager.add_Peer(peerElement);
		}
	}

	return networkManager;
}






func(networkManager *NetworkManager) server() {

	//routes
	http.HandleFunc("/join", func(writer http.ResponseWriter, request *http.Request) {

		popularity, err := strconv.Atoi(request.URL.Query().Get("popularity"));
		port, err1 := strconv.Atoi(request.URL.Query().Get("port"));

		if err != nil && err1 != nil {
  			if err != nil {
  				fmt.Print(err);
  			}

  			if err1 != nil {
  				fmt.Print(err1);
  			}
  		} else {

  			networkManager.LastBlockIndex = networkManager.Blockchain.get_latest_block().Index;
  			payload, _ := json.Marshal(networkManager);

			writer.Header().Add("Content-Type", "application/json");
			writer.Write(payload);

			peer := Construct_peer(request.URL.Query().Get(("index")), popularity, request.URL.Query().Get("host"), port);

			if !networkManager.has_peer(peer) {
				networkManager.add_Peer(peer);
			}

		}
    });


    http.HandleFunc("/getblock", func(writer http.ResponseWriter, request *http.Request) {

    	startIndex, err := strconv.Atoi(request.URL.Query().Get(("startIndex")));
    	stopIndex, err1 := strconv.Atoi(request.URL.Query().Get(("stopIndex")));

    	if err != nil || err1 != nil {
  			if err != nil {
  				fmt.Print(err);
  			}

  			if err1 != nil {
  				fmt.Print(err1);
  			}
  		} else {

  			if startIndex > stopIndex {
  				temp := startIndex;
  				startIndex = stopIndex;
  				stopIndex = temp;
  			}

  			if startIndex < 0 {
  				startIndex = 0;
  			}

  			var outBlocks []Block;

  			index := startIndex;
  			for index <= stopIndex {
	        	block := networkManager.Blockchain.get_block(index);
  				outBlocks = append(outBlocks, block);
  				index++;
  			}

	        payload, _ := json.Marshal(outBlocks);

  			writer.Header().Add("Content-Type", "application/json");
  			writer.Write(payload);
  		}

    });



	http.HandleFunc("/foundBlock", func(writer http.ResponseWriter, request *http.Request) {

    	indexBlock, err := strconv.Atoi(request.URL.Query().Get(("indexBlock")));
    	popularity, err1 := strconv.Atoi(request.URL.Query().Get("popularity"));
		port, err2 := strconv.Atoi(request.URL.Query().Get("port"));

    	if err != nil || err1 != nil || err2 != nil {
    		if err != nil {
				fmt.Print(err);
    		}

    		if err1 != nil {
				fmt.Print(err1);
    		}

    		if err2 != nil {
				fmt.Print(err2);
    		}
		} else {
			peer := Construct_peer(request.URL.Query().Get(("index")), popularity, request.URL.Query().Get("host"), port);

			//update Peer
			networkManager.update_Peer(networkManager.Peers[networkManager.get_peer_index(peer)], peer);

			if indexBlock > networkManager.Blockchain.get_latest_block().Index {
				networkManager.syncChain(peer, indexBlock);
			}



			writer.Header().Add("Content-Type", "application/json");
			writer.Write([]byte("{ \"success\": true }"));
		}

    });


    http.HandleFunc("/synch", func(writer http.ResponseWriter, request *http.Request) {

    	indexBlock, err := strconv.Atoi(request.URL.Query().Get(("indexBlock")));
    	popularity, err1 := strconv.Atoi(request.URL.Query().Get("popularity"));
		port, err2 := strconv.Atoi(request.URL.Query().Get("port"));

    	if err != nil || err1 != nil || err2 != nil {
    		if err != nil {
				fmt.Print(err);
    		}

    		if err1 != nil {
				fmt.Print(err1);
    		}

    		if err2 != nil {
				fmt.Print(err2);
    		}
		} else {
			peer := Construct_peer(request.URL.Query().Get(("index")), popularity, request.URL.Query().Get("host"), port);

			//update Peer
			networkManager.update_Peer(networkManager.Peers[networkManager.get_peer_index(peer)], peer);

			networkManager.LastBlockIndex = networkManager.Blockchain.get_latest_block().Index;

			payload, _ := json.Marshal(networkManager);

			writer.Header().Add("Content-Type", "application/json");

			writer.Write(payload);
			if indexBlock > networkManager.Blockchain.get_latest_block().Index {
				networkManager.syncChain(peer, indexBlock);
			}
		}

    });



	//start server
    http.ListenAndServe(":" +  strconv.Itoa(networkManager.Me.Port), nil);
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


func(networkManager *NetworkManager) broadcast(url string) {
	for _, peer := range networkManager.Peers {
		_, err := peer.request(url);
		if err != nil {
		    networkManager.remove_Peer(peer);
		}
	}
}


func(networkManager *NetworkManager) download_block(peer Peer, startIndex int, stopIndex int) []Block {
	startIndexString := strconv.Itoa(startIndex);
	stopIndexString := strconv.Itoa(stopIndex);

	content, err := peer.request("/getblock?startIndex=" + startIndexString + "&stopIndex=" + stopIndexString);
	if err != nil {
		networkManager.remove_Peer(peer);
	}

	body, err := ioutil.ReadAll(content.Body)
	if err != nil {
	    panic(err.Error())
	}

	out := make([]Block, 0);
	json.Unmarshal(body, &out);

	return out;
}


func block_is_in_chain(a Block, list []Block) int {
    for i, b := range list {
        if b.Index == a.Index {
            return i;
        }
    }

    return -1;
}


func(networkManager *NetworkManager) syncChain(peer Peer, index int) {

	blockchain_test := networkManager.Blockchain;

	blocks := networkManager.download_block(peer, 0, index);
	blockchain_test.Chain = blocks;

	if blockchain_test.is_valid() {
		networkManager.Blockchain = blockchain_test;
		fmt.Println("Blockchain downloaded.")
	}

	/*nb_block_dl_same_time := 10;

	blockchain_to_validate := networkManager.Blockchain;

	added := 0;
	run := 0;

	for blockchain_to_validate.is_valid() {

		startIndex := (index - ((run + 1) * nb_block_dl_same_time)) + 1;
		stopIndex := (index - (run * nb_block_dl_same_time));

		if startIndex < 0 {
			startIndex = 0;
		}

		if stopIndex < 0 {
			panic("Can't download a < 0 block");
		}

		blocks := networkManager.download_block(peer, startIndex, stopIndex);

		for _, block_to_add := range blocks {
			index_in_chain := block_is_in_chain(block_to_add, blockchain_to_validate.Chain);
			if index_in_chain > -1 {
				blockchain_to_validate.Chain[index_in_chain] = block_to_add;
			} else {
				blockchain_to_validate.Chain = append(blockchain_to_validate.Chain, block_to_add);
			}
		}


		for i, iElement := range blockchain_to_validate.Chain {
			for j, jElement := range blockchain_to_validate.Chain {
				if iElement.Index < jElement.Index {
					blockchain_to_validate.Chain[i], blockchain_to_validate.Chain[j] = blockchain_to_validate.Chain[j], blockchain_to_validate.Chain[i];
				}
			}
		}

		added += len(blocks);
		run++;
	}

	networkManager.Blockchain = blockchain_to_validate;
	fmt.Println(strconv.Itoa(added) + " blocks from blockchain downloaded in " + strconv.Itoa(run) + " steps.");*/





	/*nbBlockDownloadInSameTime := 3;
	blockchainTest := networkManager.Blockchain;

	//init download
	if len(blockchainTest.Chain) == 0 {
		blocks := networkManager.download_block(peer, 0, (nbBlockDownloadInSameTime * 2) - 1);
		blockchainTest.Chain = blocks;

		networkManager.Blockchain = blockchainTest;
		added := len(blockchainTest.Chain);
		fmt.Println(strconv.Itoa(added) + " init blocks from blockchain downloaded.");
	} else { //update download

		run := 0;
		added := 0;

		for blockchainTest.is_valid() {

			startIndex := index - added - nbBlockDownloadInSameTime;
			stopIndex := index - added - 1;

			if startIndex < 0 {
				startIndex = 0;
			}

			blocks := networkManager.download_block(peer, startIndex, stopIndex);

			for _, block_to_add := range blocks {
				index_in_chain := block_is_in_chain(block_to_add, blockchainTest.Chain);
				fmt.Println(index_in_chain)
				if index_in_chain > -1 {
					blockchainTest.Chain[index_in_chain] = block_to_add;
				} else {
					blockchainTest.Chain = append(blockchainTest.Chain, block_to_add);
				}
			}



			for i, iElement := range blockchainTest.Chain {
				for j, jElement := range blockchainTest.Chain {
					if iElement.Index < jElement.Index {
						blockchainTest.Chain[i], blockchainTest.Chain[j] = blockchainTest.Chain[j], blockchainTest.Chain[i];
					}
				}
			}
		}

		networkManager.Blockchain = blockchainTest;
		fmt.Println(strconv.Itoa(added) + " blocks from blockchain downloaded in " + strconv.Itoa(run) + " steps.");

	}*/








	/*else {

	 	run := 0;
	 	added := 0;

	 	saveBlockchainTest := blockchainTest;

	 	for run == 0 || !blockchainTest.is_valid() {

	 		var added_blocks []Block;

	 		low_index := index - run - nbBlockDownloadInSameTime;
	 		if low_index < 0 {
				low_index = 0;
	 		}

	 		fmt.Println("download : " + strconv.Itoa(low_index) + " " + strconv.Itoa(index - run));

	 		//download blocks
		  	blocks := networkManager.download_block(peer, low_index, index - run);

			for _, block_to_add := range blocks {
				index_in_chain := block_is_in_chain(block_to_add, added_blocks);
				if index_in_chain > -1 {
					added_blocks[index_in_chain] = block_to_add;
				} else {
					added_blocks = append(added_blocks, block_to_add);
				}
			}

	 		added += len(blocks);
	 		run += added - 1;

	 		blockchainTest.Chain = append(saveBlockchainTest.Chain, blocks...);

	 		if run % nbBlockDownloadInSameTime >= index % nbBlockDownloadInSameTime {
	 			return;
	 		}
	 	}

	 	networkManager.Blockchain = blockchainTest;
		fmt.Println(strconv.Itoa(added) + " blocks from blockchain downloaded.");
	}*/








  	//download blocks
  	/*blocks := networkManager.download_block(peer, 0, index);
	if blocks[0].Index == 0 {
		blockchainTest.Chain = blocks;
	} else {
		blockchainTest.Chain = append(blockchainTest.Chain, blocks...);
	}*/

	/*for i := 0; i < len(blockchainTest.Chain); i++ {
		fmt.Println(blockchainTest.get_block(i).Index);
	}*/

	/*if blockchainTest.is_valid() {
		networkManager.Blockchain = blockchainTest;
		fmt.Println("Blockchain downloaded.")
	}*/

}

func(networkManager *NetworkManager) randomPeer(url string) (*http.Response, error) {

	n := rand.Int() % len(networkManager.Peers);
	peer := networkManager.Peers[n];
	response, err := peer.request(url);
	if err != nil {
		networkManager.remove_Peer(peer);
	}

	return response, err;

}
