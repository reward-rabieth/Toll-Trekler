package main

import (
	"github.com/gorilla/websocket"
	"github.con/reward-rabieth/Troll-Trekler/types"
	"log"
	"math/rand"
	"time"
)

const wsEndpoint = "ws://127.0.0.1:7000/ws"

var sendInterval = time.Second * 5

func main() {
	obuIDS := generateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, lang := genLatLong()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat:   lat,
				Long:  lang,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(sendInterval)
	}
}

func genLatLong() (float64, float64) {
	return genCord(), genCord()
}

func genCord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()

	return n + f
}

func generateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(100000)
	}
	return ids
}

func init() {

	rand.NewSource(time.Now().UnixNano())
}
