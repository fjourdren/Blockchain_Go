package main

import "strconv"
import "strings"
import "encoding/json"


type Block struct {
	Index int
    Timestamp int
    Difficulty int
    NextBlockDifficulty int
    Data string
    Hash string
    PreviousHash string
    Nonce int
};




func Construct_block(indexVal int,
    timestampVal int,
    difficultyVal int,
    nextBlockDifficultyVal int,
    dataVal string,
    hashVal string,
    previousHashVal string,
    nonceVal int) Block {
	block := Block{Index: indexVal,
    				Timestamp: timestampVal,
    				Difficulty: difficultyVal,
    				NextBlockDifficulty: nextBlockDifficultyVal,
    				Data: dataVal,
    				Hash: hashVal,
    				PreviousHash: previousHashVal,
    				Nonce: nonceVal};
	return block;
}


func(block Block) calculate_hash() string {
    total := strconv.Itoa(block.Index) + strconv.Itoa(block.Nonce) + block.PreviousHash + strconv.Itoa(block.Difficulty) + strconv.Itoa(block.NextBlockDifficulty) + strconv.Itoa(block.Timestamp) + block.Data;
    return hash(total);
}


func(block Block) hash_is_valid(hash string) bool {

    if block.Difficulty == 0 {
        return true;
    }

    return strings.HasPrefix(hash, get_difficulty_string_patern(block.Difficulty));
}


func(block *Block) mine_hash() {
    for !block.hash_is_valid(block.Hash) {
        block.Timestamp = now();
        block.Nonce = random(4294967296);
        block.Hash = block.calculate_hash();
    }
}


func(block Block) is_valid() bool {
    if(block.calculate_hash() != block.Hash) {
        return false;
    }

    if !block.hash_is_valid(block.Hash) {
        return false;
    }

    if len(block.to_json()) > 2097152 { //2MB
        return false;
    }

    return true;
}


func(block *Block) to_json() []byte {
    payload, _ := json.Marshal(block);
    return payload;
}

func block_json_to_object(content []byte) *Block {
    block := new(Block);
    json.Unmarshal(content, block);

    return block;
}