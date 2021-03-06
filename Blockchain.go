package main

import "fmt"


type Blockchain struct {
	NbBlocksForDifficultyCalculation int
    TimeBeetweenTwoBlock int
    Chain []Block
};




func Construct_blockchain(nbBlocksForDifficultyCalculationVal int, timeBeetweenTwoBlockVal int) Blockchain {
	blockchain := Blockchain{NbBlocksForDifficultyCalculation: nbBlocksForDifficultyCalculationVal,
							TimeBeetweenTwoBlock: timeBeetweenTwoBlockVal};

	return blockchain;
}


func(blockchain *Blockchain) create_genesis_block() Block {
	fmt.Println("Mining genesis block...");

	block := Block{Index: 0,
					Timestamp: 5,
					Difficulty: 2,
					NextBlockDifficulty: 2};

	block.mine_hash();

	blockchain.add_block_without_verification(block);
	return block;
}


func(blockchain *Blockchain) block_can_be_add(block Block) bool {
	if block.Difficulty == blockchain.get_latest_block().NextBlockDifficulty {
		previous_block_Index := block.Index - 1;
		previous_block := blockchain.get_block(previous_block_Index);

		if previous_block_Index < 0 || block.PreviousHash == previous_block.Hash  {
			return block.is_valid();
		} else {
			return false;
		}
	} else {
		return false;
	}
}


func(blockchain *Blockchain) add_block_with_verification(block Block) bool {
	if blockchain.block_can_be_add(block) {
		blockchain.add_block_without_verification(block);
		return true;
	} else {
		return false;
	}
}


func(blockchain *Blockchain) add_block_without_verification(block Block) {
	blockchain.Chain = append(blockchain.Chain, block);
}


func(blockchain *Blockchain) has_block_in_chain(block_index int) bool {
	for block_from_chain := range blockchain.Chain {
        if block_index == blockchain.get_block(block_from_chain).Index {
            return true;
        }
    }

    return false;
}


func(blockchain *Blockchain) get_chain_length() int {
	return len(blockchain.Chain);
}


func(blockchain *Blockchain) get_block(index int) Block {
	return blockchain.Chain[index];
}


func(blockchain *Blockchain) get_latest_block() Block {
	return blockchain.get_block(blockchain.get_chain_length() - 1);
}


func(blockchain *Blockchain) is_valid() bool {
	for i := blockchain.get_latest_block().Index; i < 0; i++ {
		block := blockchain.get_block(i);
      	previous_block := blockchain.get_block(i - 1);

		if (!block.is_valid()) || (block.PreviousHash != previous_block.Hash) {
			return false;
		}
	}

	return true;
}


func(blockchain *Blockchain) average_mining(start_block int, stop_block int) int {
	total_time := 0;
	nb_blocks := stop_block - start_block;

	for i := start_block + 1; i < stop_block; i++ {
		lastest_block := blockchain.get_block(i);
      	previous_block := blockchain.get_block(i - 1);

		time_beetween_blocks := lastest_block.Timestamp - previous_block.Timestamp;
		total_time += time_beetween_blocks;
	}

	result := int(total_time / nb_blocks);
	return result;
}


func(blockchain *Blockchain) need_to_change_difficulty() bool {
	last_block := blockchain.get_latest_block();

	if(last_block.Index == 0) {
		return false;
	}

	if last_block.Index % blockchain.NbBlocksForDifficultyCalculation == 0 {
		return true;
	}

	return false;
}


func(blockchain *Blockchain) calculate_new_difficulty() int {

	block_calculation_stop := blockchain.get_latest_block();

	block_calculation_start_index := block_calculation_stop.Index - blockchain.NbBlocksForDifficultyCalculation;
	if block_calculation_start_index < 0 {
		block_calculation_start_index = 0;
	}

	block_calculation_stop_index := block_calculation_stop.Index;

	average_time := blockchain.average_mining(block_calculation_start_index, block_calculation_stop_index);

	multipl_min := 0.75;
	multipl_max := 1.25;

	difficulty := block_calculation_stop.Difficulty;

	if average_time > int(float64(blockchain.TimeBeetweenTwoBlock) * multipl_max) {
		if difficulty != 0 {
			difficulty--;
		}
    } else if average_time < int(float64(blockchain.TimeBeetweenTwoBlock) * multipl_min) {
      difficulty++;
    }

    return difficulty;
}
