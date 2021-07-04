package main

import "fmt"

var packets_already_broadcast []string;


func broadcastManager_add_packet(packet_id string) {
    packets_already_broadcast = append(packets_already_broadcast, packet_id);
}


func broadcastManager_remove_packet(packet_id string) {
    index := broadcastManager_get_packet_index(packet_id);
    packets_already_broadcast = append(packets_already_broadcast[:index], packets_already_broadcast[index+1:]...);
}


func broadcastManager_get_packet_index(packet_id string) int {
    for index, packet_sign := range packets_already_broadcast {
        if packet_sign == packet_id {
            return index;
        }
    }

    return -1;
}


func broadcastManager_has_packet(packet_id string) bool {
    index := broadcastManager_get_packet_index(packet_id);

    fmt.Println(index)

    if index == -1 {
        return false;
    } else {
        return true;
    }
}