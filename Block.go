package main

import "strconv"
import "strings"
import "math/rand"
import "encoding/json"

type Block struct {
	Index int
    Timestamp int
    Difficulty int
    Data string
    Hash string
    PreviousHash string
    Nonce int
};




func Construct_block(indexVal int,
    timestampVal int,
    difficultyVal int,
    dataVal string,
    hashVal string,
    previousHashVal string,
    nonceVal int) Block {
	block := Block{Index: indexVal, 
						Timestamp: timestampVal, 
						Difficulty: difficultyVal,
						Data: dataVal,
						Hash: hashVal,
						PreviousHash: previousHashVal,
						Nonce: nonceVal};
	return block;
}


func(block Block) calculate_hash() string {
    total := strconv.Itoa(block.Index) + strconv.Itoa(block.Nonce) + block.PreviousHash + strconv.Itoa(block.Difficulty) + strconv.Itoa(block.Timestamp) + block.Data;
    return hash(total);
}


func(block Block) hash_is_valid(hash string) bool {

    if(block.Difficulty == 0) {
        return true;
    }

    return strings.HasPrefix(hash, get_difficulty_string_patern(block.Difficulty));
}


func(block Block) is_valid() bool {
    if(block.calculate_hash() != block.Hash) {
        return false;
    }

    if(!block.hash_is_valid(block.Hash)) {
        return false;
    }

    if(len(block.to_string()) > 2097152) { //2MB
        return false;
    }

    return true;
}


func(block *Block) mine_hash() {
    for !block.hash_is_valid(block.Hash) {
        block.Timestamp = now();
        block.Nonce = rand.Intn(4294967296);
        block.Hash = block.calculate_hash();
    }
}

func(block *Block) to_string() []byte {
    payload, _ := json.Marshal(block);
    return payload;
}
