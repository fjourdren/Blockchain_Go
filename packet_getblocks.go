package main

import "encoding/json"

type packet_getblocks struct {
    Start int
    Stop int
};




func Construct_Packet_getBlocks(start int, stop int) packet_getblocks {
	packet := packet_getblocks{Start: start,
								Stop: stop};

	return packet;
}


func(packet *packet_getblocks) to_json() []byte {
	payload, _ := json.Marshal(packet);
	return payload;
}


func packet_getBlocks_json_to_object(content []byte) *packet_getblocks {
	packet := new(packet_getblocks);
    json.Unmarshal(content, packet);

    return packet;
}
