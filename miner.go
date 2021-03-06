package main

import "fmt"
//import "time"
import "strconv"
import "math/rand"



func get_difficulty_string_patern(difficulty int) string {
	var out string = "";

	for i := 0; i < difficulty; i++ {
		out += "0";
	}

	return out;
}


func miner(networkManager *NetworkManager) {
	nbMine := 0;
	startTime := now();
	restart := false
	for {
		restart = false;
		//time_mine_start := time.Now();

		latest_block := networkManager.Blockchain.get_latest_block();
		next_block_difficulty := latest_block.NextBlockDifficulty;

		//calcule next difficulty if needed
		if (latest_block.Index + 1) % (networkManager.Blockchain.NbBlocksForDifficultyCalculation - 1) == 0 {
			next_block_difficulty = networkManager.Blockchain.calculate_new_difficulty();
			fmt.Println("New difficulty: " + strconv.Itoa(next_block_difficulty));
		}


		difficulty_for_this_block := latest_block.NextBlockDifficulty;
		if difficulty_for_this_block < 1 {
			difficulty_for_this_block = 1;
		}

		block := Construct_block(latest_block.Index + 1,
							    now(),
							    difficulty_for_this_block,
							    next_block_difficulty,
							    "data",
							    "",
							    networkManager.Blockchain.get_block(networkManager.Blockchain.get_chain_length() - 1).Hash,
							    0);


		//mining
		for !block.hash_is_valid(block.Hash) && !restart {
	        block.Timestamp = now();
	        block.Nonce = rand.Intn(4294967296);
	        block.Hash = block.calculate_hash();

	        //mineHash calculation & network sync
	        nbMine++;
			if now() > startTime + timeBeetweenNetworkSync {

				/*if len(networkManager.Peers) > 0 {
					//ask to random peer the network if he has better block
					response, err := networkManager.randomPeer("/synch?indexBlock=" + strconv.Itoa(networkManager.Blockchain.get_latest_block().Index) + "&index=" + networkManager.Me.Index + "&popularity=" + strconv.Itoa(networkManager.Me.Popularity) + "&host=" + networkManager.Me.Host + "&port=" + strconv.Itoa(networkManager.Me.Port));
					if err == nil {

						//response to []byte
						body, err := ioutil.ReadAll(response.Body);
						if err != nil {
						    panic(err.Error());
						}

						//extract data
						networkManagerDist, _ := NetworkManagerFromJSON(body);

						//update Peer
						//networkManager.update_Peer(networkManager.Peers[networkManager.get_peer_index(networkManagerDist.Me)], networkManagerDist.Me);
						networkManager.Peers[networkManager.get_peer_index(networkManagerDist.Me)] = networkManagerDist.Me;

						if networkManagerDist.LastBlockIndex > networkManager.Blockchain.get_latest_block().Index {
							networkManager.syncChain(networkManagerDist.Me, networkManagerDist.LastBlockIndex);
							fmt.Println("Random peer sync");
						}
					}

				}*/



				//1 second has pass
				restart = true;

				timeSinceLastMiningRateCalc := (now() - (startTime + timeBeetweenNetworkSync));
				fmt.Println(strconv.Itoa((nbMine / timeSinceLastMiningRateCalc) / 1000) + "Kh/s");
				nbMine = 0;
				startTime = now();
			}
	    }

	    //time_to_mine := time.Since(time_mine_start).String();


	    //add mined block
		if networkManager.Blockchain.add_block_with_verification(block) && !restart {
			//fmt.Println("Block #" + strconv.Itoa(block.Index) + " difficulty(" + strconv.Itoa(block.Difficulty) + ") mined in " + time_to_mine + " with nonce " + strconv.Itoa(block.Nonce) + " (" + block.Hash + ")");
			//networkManager.broadcast_udp("/foundBlock?indexBlock=" + strconv.Itoa(block.Index) + "&index=" + networkManager.Me.Index + "&popularity=" + strconv.Itoa(networkManager.Me.Popularity) + "&host=" + networkManager.Me.Host + "&port=" + strconv.Itoa(networkManager.Me.Port));
			packet := Construct_Packet(networkManager.Me, "BROADCAST", "FOUNDBLOCK", []byte(strconv.Itoa(block.Index)));
  			packetJson := packet.to_json();
			networkManager.broadcast(packetJson);
		}

	}
}
