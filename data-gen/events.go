package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// Event represents user behavior events
type Event struct {
	EventType     string `json:"event_type"`
	VideoID       string `json:"video_id"`
	VideoCategory string `json:"video_category"`
	ChannelID     string `json:"channel_id"`
	IsSubscribed  bool   `json:"is_subscirbed"`
	SearchTerm    string `json:"search_term"`
}

type categoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// Channel represents Youtube channel category views
type Channel struct {
	ChannelID      string          `json:"_id"`
	CategoryCounts []categoryCount `json:"categories"`
}

func genEvents(ecnt int, ef *os.File, chf *os.File) {
	catf, err := os.OpenFile("categories.csv", os.O_RDONLY, 0622)
	if err != nil {
		log.Fatal(err)
	}
	cs, err := bufio.NewReader(catf).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	catf.Close()
	ew := bufio.NewWriter(ef)
	chw := bufio.NewWriter(chf)

	chcnt := ecnt / 1000
	vcnt := ecnt / 10
	cats := strings.Split(cs, ", ")
	etype := []string{"ad_click", "ad_watch"}
	ttype := []string{"video_watch", "video_exit"}
	terms := []string{"go", "c++", "javascript", "perl", "python", "java"}
	cnt := 0

	chs := make(map[string]map[string]int)

	for cnt < ecnt {
		cid := strconv.Itoa(rand.Intn(chcnt))
		vid := strconv.Itoa(rand.Intn(vcnt))
		cat := cats[rand.Intn(len(cats))]
		isSub := rand.Intn(2) == 1
		term := terms[rand.Intn(len(terms))]

		e := Event{"video_click", vid, cat, cid, isSub, term}
		b, err := json.Marshal(e)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := ew.WriteString(string(b) + "\n"); err != nil {
			log.Fatal(err)
		}
		cnt++
		// fmt.Printf("%#v\n", e)
		if chs[cid] == nil {
			chs[cid] = make(map[string]int)
		}
		chs[cid][cat]++

		bias := len(etype)
		if isSub {
			bias = len(etype)/2 + 1
		}
		for eroll := rand.Intn(len(etype) + bias); eroll < len(etype); {
			e = Event{etype[eroll], vid, cat, cid, isSub, term}
			b, err = json.Marshal(e)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := ew.WriteString(string(b) + "\n"); err != nil {
				log.Fatal(err)
			}
			// fmt.Printf("%#v\n", e)
			cnt++
			if cnt >= ecnt-1 {
				break
			}
			eroll = rand.Intn(len(etype) + 1)
		}

		e = Event{ttype[rand.Intn(len(ttype))], vid, cat, cid, isSub, term}
		b, err = json.Marshal(e)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := ew.WriteString(string(b) + "\n"); err != nil {
			log.Fatal(err)
		}
		if err := ew.Flush(); err != nil {
			log.Fatal(err)
		}
		cnt++
		// fmt.Printf("%#v\n", e)
	}

	fmt.Println("Done writing events")

	for cid, ccnt := range chs {
		ch := Channel{}
		ch.ChannelID = cid
		for cat, cnt := range ccnt {
			ch.CategoryCounts = append(ch.CategoryCounts, categoryCount{cat, cnt})
		}
		b, err := json.Marshal(ch)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := chw.WriteString(string(b) + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	if err := chw.Flush(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done writing channels")
}
