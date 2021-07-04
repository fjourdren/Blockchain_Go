package main

/* TODO */
	//don't use the array length
	//data processing
	//save data
	//storage and network without miner
	//miner without real storage
	//test unitaire
	//rewrite miner
	//config file
	//better network (full packet system, multi thread server)



	//sync chain
	//tcp / http for transactions
	//better broadcast network
	//popularity repartition


import "os"
import "fmt"
import "strconv"



var timeBeetweenNetworkSync int = 1;


func main() {

	blockchain := Construct_blockchain(50, 20);
	blockchainPointer := &blockchain;

	var networkManager NetworkManager;

	index := hash(strconv.Itoa(random(4294967296)));

	if len(os.Args) > 1 {
		port, _ := strconv.Atoi(os.Args[2]);
		networkManager = Construct_Init_NetworkManager(index, 0, "127.0.0.1", 8081, blockchainPointer, os.Args[1], port);
	} else {
		blockchainPointer.create_genesis_block();
		networkManager = Construct_NetworkManager(index, 0, "127.0.0.1", 8080, blockchainPointer);
	}

	go networkManager.server();
	fmt.Println("P2P Network started !");

	miner(&networkManager);
}
