package main

/* TODO */
	//verify difficulty next block
	//don't use the array length
	//data processing
	//network
	//popularity
	//save data
	//storage and network without miner
	//miner without real storage
	//test unitaire
	//rewrite miner
	//config file
	//synch chain


	//sync chain
	//repartition des noeud en fonction de la popularité
	//next block difficulty

import "os"
import "fmt"
import "time"
import "strconv"
import (
    "crypto/sha512"
    "encoding/base64"
)


var timeBeetweenNetworkSync int = 1;


func now() int {
	return int(time.Now().Unix());
}


func hash(input string) string {
	hasher := sha512.New();

    hasher.Write([]byte(input));
    return base64.URLEncoding.EncodeToString(hasher.Sum(nil));
}


func main() {

	blockchain := Construct_blockchain(50, 20);
	blockchainPointer := &blockchain;

	var networkManager NetworkManager;

	if len(os.Args) > 1 {
		port, _ := strconv.Atoi(os.Args[2]);
		networkManager = Construct_Init_NetworkManager(0, 0, "http://127.0.0.1", 8081, blockchainPointer, os.Args[1], port);
	} else {
		blockchainPointer.create_genesis_block();
		networkManager = Construct_NetworkManager(0, 0, "http://127.0.0.1", 8080, blockchainPointer);
	}

	go networkManager.server();
	fmt.Println("P2P Network started !");

	miner(&networkManager);
}
