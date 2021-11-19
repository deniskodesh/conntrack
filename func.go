package main

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"time"

	ct "github.com/florianl/go-conntrack"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// See https://golang.org/pkg/container/heap/
func getHeap(m map[string]int) *KVHeap {
	h := &KVHeap{}
	heap.Init(h)

	for k, v := range m {
		heap.Push(h, kv{k, v})
	}
	return h
}

// See https://golang.org/pkg/container/heap/
func (h KVHeap) Less(i, j int) bool { return h[i].Value > h[j].Value }
func (h KVHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h KVHeap) Len() int           { return len(h) }

// See https://golang.org/pkg/container/heap/
func (h *KVHeap) Push(x interface{}) {
	*h = append(*h, x.(kv))
}

// See https://golang.org/pkg/container/heap/
func (h *KVHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// This func read from file and expose byte slice
func readFromFile(path string) []byte {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	return content
}

//This func uses Atoi for return float64
func Float64frombytes(bytesSlice []byte) float64 {
	if len(bytesSlice) > 1 {
		bytesSlice = bytesSlice[:len(bytesSlice)-1]

	}
	intNumber, _ := strconv.Atoi(string(bytesSlice))

	return float64(intNumber)

}

//Gets dump of record from conntrack table
func GetRecordsFromTable() []string {

	temp := []string{}
	nfct, err := ct.Open(&ct.Config{})
	if err != nil {
		fmt.Println("Could not create nfct:", err)

	}
	defer nfct.Close()
	sessions, err := nfct.Dump(ct.Conntrack, ct.IPv4)
	if err != nil {
		fmt.Println("Could not dump sessions:", err)

	}
	for _, session := range sessions {
		//fmt.Printf("[%2d] %s - %s\n", session.Origin.Proto.Number, session.Origin.Src, session.Origin.Dst)

		temp = append(temp, session.Origin.Src.String())

	}

	return temp
}

//Calculate how many one IP match in string slice based on this info creates map like ip - session count
func HowMatches(IPs []string) map[string]int {

	dict := make(map[string]int)

	for _, ip := range IPs {
		dict[ip] = dict[ip] + 1
	}

	//for k, v := range dict {
	//fmt.Printf("%s -> %s\n", k, v)
	//}

	return dict
}

//This func get top values
func getTopValues(count int, sessions []string) []kv {
	var results []kv
	h := *getHeap(HowMatches(sessions))

	sort.SliceStable(h, func(i, j int) bool { return h[i].Value > h[j].Value })

	if count > len(h) {

		count = len(h) - 1
	}

	for i := 0; i <= count; i++ {

		results = append(results, h[i])
	}
	return results

}

//Gets metrics conntract
func recordMetrics() {
	go func() {
		for {

			fileBytes := readFromFile(settings.PathToConntrackCount)
			conntrack_Total.Set(Float64frombytes(fileBytes))
			log.WithFields(log.Fields{
				"value": Float64frombytes(fileBytes),
			}).Info("Conntrack count")

			time.Sleep(time.Duration(settings.ConntrackCountCheckInterval) * time.Second)
		}
	}()

	go func() {
		for {

			fileBytes := readFromFile(settings.PathToConntrackMax)
			conntrack_Max.Set(Float64frombytes(fileBytes))
			log.WithFields(log.Fields{
				"value": Float64frombytes(fileBytes),
			}).Info("Conntrack max")

			time.Sleep(time.Duration(settings.ConntrackMaxCheckInterval) * time.Second)
		}
	}()

	go func() {
		for {
			sessions := GetRecordsFromTable()
			results := getTopValues(settings.TopRecordsCount-1, sessions)
			log.WithFields(log.Fields{
				"Top": settings.TopRecordsCount - 1,
			}).Info("Top count")
			for _, el := range results {
				Top.With(prometheus.Labels{"ip": el.Key}).Set(float64(el.Value))
				log.WithFields(log.Fields{
					"ip":    el.Key,
					"value": el.Value,
				}).Info("Top values")
			}

			time.Sleep(time.Duration(settings.ConntrackTopCheckInterval) * time.Second)
		}
	}()
}
