package main

import "encoding/json"

type Packet struct {
    Index string
    Sender Peer
    Type string
    Name string
    Content []byte
};




func Construct_Packet(sender_var Peer, type_var string, name_var string, content_var []byte) Packet {
	packet := Packet{Sender: sender_var,
					Type: type_var,
					Name: name_var,
					Content: content_var};

	packet.Index = generate_UUID();

	return packet;
}


func(packet *Packet) to_json() []byte {
	payload, _ := json.Marshal(packet);
	return payload;
}


func packet_json_to_object(content []byte) *Packet {
	packet := new(Packet);
    json.Unmarshal(content, packet);

    return packet;
}
