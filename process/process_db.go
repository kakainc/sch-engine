package process

import (
	"Sch_engine/init_sch"
	"Sch_engine/rely"
	"fmt"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"sync"
)

func DoDbProcess(pids []rely.PidScore, rd *redis.Client, topN int) rely.RetPids {
	if len(pids)<topN {
		dePids := init_sch.DoErrDegree(rd, int64(topN-len(pids)))
		for _, pid := range dePids {
			pids = append(pids, pid)
		}
	}
	init_sch.Info.Printf("process get %d posts!", len(pids))

	var finalPids rely.RetPids
	for i := 0; i < topN; i++ {
		pid:=pids[i]
		finalPids.Pids = append(finalPids.Pids, pid.Pid)
	}
	return finalPids
}

func DoDbaProcess(pids []rely.PidScore, rd *redis.Client, mg *mongo.Collection, topN int,
	cacheContent map[int64]string, cacheComment map[int64]string, cLock *sync.Mutex, mLock *sync.Mutex) string {
	if len(pids)<topN {
		dePids := init_sch.DoErrDegree(rd, int64(topN-len(pids)))
		for _, pid := range dePids {
			pids = append(pids, pid)
		}
	}
	init_sch.Info.Printf("process get %d posts!", len(pids))
	var finalcontent string
	for i := 0; i < topN; i++ {
		pid:=pids[i].Pid
		var content string
		if cn, ok := cacheContent[pid]; ok {
			content = cn
		} else{
			content= init_sch.FindOne(mg, strconv.FormatInt(pid, 10)).Content
			//cLock.Lock()
			cacheContent[pid] = content
			//cLock.Unlock()
		}

		var comments string
		if cn, ok := cacheComment[pid]; ok {
			comments = cn
		} else{
			comments= init_sch.FindOne(mg, strconv.FormatInt(pid, 10)).Comment
			//mLock.Lock()
			cacheComment[pid] = comments
			//mLock.Unlock()
		}

		s := fmt.Sprintf("\033[32m[content%d:]\033[0m", i+1)
		finalcontent= finalcontent + s + content + "\n"

		if comments == "" {
			continue
		} else {
			commentSlice := strings.Split(comments, "\t")
			init_sch.Info.Printf("post %d get %d comments!", pid, len(pids))
			mh := fmt.Sprintf("\033[33m[comments:]\033[0m")
			finalcontent= finalcontent + mh
			for j, comment := range commentSlice {
				var m string
				if j==0 {
					m = fmt.Sprintf("%d %s \n", j+1, comment)
				} else {
					m = fmt.Sprintf("           "+"%d %s \n", j+1, comment)   // 11*""
				}

				finalcontent= finalcontent + m
				if j >= 3 {
					break
				}
			}
		}
	}
	//fmt.Println(finalcontent)
	return finalcontent
}

func DoDbcProcess(pids []rely.PidScore, rd *redis.Client, mg *mongo.Collection, topN int,
	cacheContent map[int64]string, cacheComment map[int64]string, cLock *sync.Mutex, mLock *sync.Mutex) string {
	if len(pids)<topN {
		dePids := init_sch.DoErrDegree(rd, int64(topN-len(pids)))
		for _, pid := range dePids {
			pids = append(pids, pid)
		}
	}
	init_sch.Info.Printf("process get %d posts!", len(pids))
	var finalcontent string
	for i := 0; i < topN; i++ {
		pid:=pids[i].Pid
		var content string
		if cn, ok := cacheContent[pid]; ok {
			content = cn
		} else{
			content= init_sch.FindOne(mg, strconv.FormatInt(pid, 10)).Content
			//cLock.Lock()
			cacheContent[pid] = content
			//cLock.Unlock()
		}

		var comments string
		if cn, ok := cacheComment[pid]; ok {
			comments = cn
		} else{
			comments= init_sch.FindOne(mg, strconv.FormatInt(pid, 10)).Comment
			//mLock.Lock()
			cacheComment[pid] = comments
			//mLock.Unlock()
		}

		s := fmt.Sprintf("content %d :\n", i+1)
		finalcontent= finalcontent + s + content + "\n"

		if comments == "" {
			continue
		} else {
			commentSlice := strings.Split(comments, "\t")
			init_sch.Info.Printf("post %d get %d comments!", pid, len(pids))
			mh := fmt.Sprintf("\t comment :\n")
			finalcontent= finalcontent + mh
			for j, comment := range commentSlice {
				m := fmt.Sprintf("\t %d %s \n", j+1, comment)
				finalcontent= finalcontent + m
				if j >= 3 {
					break
				}
			}
		}
	}
	fmt.Println(finalcontent)
	return finalcontent
}