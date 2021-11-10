package main

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"

	ct "github.com/florianl/go-conntrack"
	"github.com/prometheus/procfs"
)

func getHeap(m map[string]int) *KVHeap {
	h := &KVHeap{}
	heap.Init(h)
	for k, v := range m {
		heap.Push(h, kv{k, v})
	}
	return h
}

// See https://golang.org/pkg/container/heap/

// Note that "Less" is greater-than here so we can pop *larger* items.
func (h KVHeap) Less(i, j int) bool { return h[i].Value > h[j].Value }
func (h KVHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h KVHeap) Len() int           { return len(h) }

func (h *KVHeap) Push(x interface{}) {
	*h = append(*h, x.(kv))
}

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

func Float64frombytes(bytesSlice []byte) float64 {
	if len(bytesSlice) > 1 {
		bytesSlice = bytesSlice[:len(bytesSlice)-1]

	}
	intNumber, _ := strconv.Atoi(string(bytesSlice))

	fmt.Println(float64(intNumber))
	return float64(intNumber)

}

func printslice(slice []string) {
	//fmt.Println("slice = ", slice)
}

func dup_count(list []string) map[string]int {

	duplicate_frequency := make(map[string]int)

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map

		_, exist := duplicate_frequency[item]

		if exist {
			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item] = 1 // else start counting from 1
		}
	}
	return duplicate_frequency
}

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

func GetTableEntriesNumber() float64 {

	fs, err := procfs.NewFS("/")
	stats, err := fs.ConntrackStat()
	if err != nil {
		fmt.Println("No file", err)

	}
	//count := 10

	println(len(stats))
	// for _, el := range stats {
	// 	println("****************")
	// 	println(el.Entries)
	// 	println("****************")

	// 	count = count + int(el.Entries)

	// }

	return float64(len(stats))
}

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
