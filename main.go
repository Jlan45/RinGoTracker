package main

import "sync"

type Peer struct {
	PeerID   string `bencode:"peer_id"`
	InfoHash string `bencode:"info_hash"`
	IP       string `bencode:"ip"`
	Port     int    `bencode:"port"`
	//Uploaded   int    `bencode:"uploaded"`
	//Downloaded int    `bencode:"downloaded"`
	//Left       int    `bencode:"left"`
	Event string `bencode:"event"`
}
type TrackerPeer struct {
	PeerID string `bencode:"peer_id"`
	IP     string `bencode:"ip"`
	Port   int    `bencode:"port"`
}

var peerStoreLock sync.RWMutex
var peerStore = make(map[string]Peer)

// 存储每个infohash下必要信息，无需PeerID
// key是infohash
var peerShare = make(map[string][]TrackerPeer)

func main() {
	startHttpServer(8080)
	//startUDPServer(8081)
	//make(map[string]string)
}
