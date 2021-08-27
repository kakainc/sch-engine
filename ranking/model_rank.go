package ranking

import (
	"Sch_engine/init_sch"
	"Sch_engine/rely"
	"bytes"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func buildRanker(querry string, pids []rely.PidScore, mg *mongo.Collection) map[string]string {
	ranker := map[string]string{
		"querry": querry,
	}
	for _,pid := range pids {
		content:= init_sch.FindOne(mg, strconv.FormatInt(pid.Pid, 10)).Content
		ranker[strconv.FormatInt(pid.Pid, 10)] = content
	}
	return ranker
}

func post(querry string, pids []rely.PidScore, mg *mongo.Collection, rkChan chan map[string]float64) {
	data := buildRanker(querry, pids, mg)
	bytesData, _ := json.Marshal(data)
	var rkscore map[string]float64
	init_sch.Info.Printf("go ranking......")
	//defer wg.Done()

	client := &http.Client{Timeout: 15 * time.Millisecond}  // 15ms超时
	resp, err := client.Post("http://127.0.0.1:8001","application/json", bytes.NewReader(bytesData))
	if err != nil {
		init_sch.Info.Printf("ranking client Post, err!!!")
		rkChan  <- rkscore
		close(rkChan)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &rkscore)

	for pid, score := range rkscore {
		init_sch.Info.Printf("%s bm25rank score: %f", pid, score)
	}

	rkChan  <- rkscore
	close(rkChan)
	return
}
