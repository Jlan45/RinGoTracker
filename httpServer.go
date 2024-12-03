package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeebo/bencode"
	"strconv"
)

type FailedResponse struct {
}
type SuccessResponse struct {
	Interval    int           `bencode:"interval"`
	MinInterval int           `bencode:"min interval"`
	Peers       []TrackerPeer `bencode:"peers"`
}

func startHttpServer(port int) {
	httpServer := gin.New()
	httpServer.GET("/announce", announceHandler)
	httpServer.Run(fmt.Sprintf(":%d", port))
}
func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
func announceHandler(c *gin.Context) {
	//参数处理
	peer := Peer{
		PeerID:   c.Query("peer_id"),
		InfoHash: c.Query("info_hash"),
		IP:       c.Query("ip"),
		Port:     atoi(c.Query("port")),
		//Uploaded:   atoi(c.Query("uploaded")),
		//Downloaded: atoi(c.Query("downloaded")),
		//Left:       atoi(c.Query("left")),
		Event: c.Query("event"),
	}
	if peer.PeerID == "" || peer.InfoHash == "" {
		c.JSON(400, gin.H{"failure reason": "peer_id is required"})
		return
	}
	if peer.IP == "" {
		peer.IP = c.ClientIP()
	}
	peerStoreLock.Lock()
	defer peerStoreLock.Unlock()
	switch peer.Event {
	case "started":
		//添加peer到peerStore
		peerStore[peer.PeerID] = peer
		//判断peerShare是否存在infohash
		if _, ok := peerShare[peer.InfoHash]; !ok {
			peerShare[peer.InfoHash] = make([]TrackerPeer, 0)
		}
		peerShare[peer.InfoHash] = append(peerShare[peer.InfoHash], TrackerPeer{IP: peer.IP, Port: peer.Port})

		//判断peer是否已经存在，存在就删除
		for i, p := range peerShare[peer.InfoHash] {
			if p.IP == peer.IP && p.Port == peer.Port {
				peerShare[peer.InfoHash] = append(peerShare[peer.InfoHash][:i], peerShare[peer.InfoHash][i+1:]...)
				break
			}
		}
		peerShare[peer.InfoHash] = append(peerShare[peer.InfoHash], TrackerPeer{IP: peer.IP, Port: peer.Port})
		res := SuccessResponse{
			Interval:    3600,
			MinInterval: 1200,
			Peers:       peerShare[peer.InfoHash],
		}
		resText, _ := bencode.EncodeString(res)
		c.String(200, resText)
		break
	case "stopped":
		//从peerStore删除peer
		delete(peerStore, peer.PeerID)
		//从peerShare删除peer
		for i, p := range peerShare[peer.InfoHash] {
			if p.IP == peer.IP && p.Port == peer.Port {
				peerShare[peer.InfoHash] = append(peerShare[peer.InfoHash][:i], peerShare[peer.InfoHash][i+1:]...)
				break
			}
		}
		break
	case "completed":
		//标记peer下载完成
		peerStore[peer.PeerID] = peer
		break
	}
}
