package main

import (
	"container/heap"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	ct "github.com/florianl/go-conntrack"
	"github.com/prometheus/procfs"
)

type KVHeap []kv

type kv struct {
	Key   string
	Value int
}

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

	data := binary.BigEndian.Uint64(bytesSlice)
	fmt.Println(data)
	data1 := float64(data)

	return data1
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

func GetRecordsFromTable() ([]string, float64) {

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

	return temp, float64(len(sessions))
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

func getTopValues(count int, sessions []string) {

	//topValues := []KVHeap{}

	h := *getHeap(HowMatches(sessions))

	sort.SliceStable(h, func(i, j int) bool { return h[i].Value > h[j].Value })
	println(len(h))
	for i := 0; i < count; i++ {
		println(count, "ip", h[i].Value, "sessions-", h[i].Value)
	}
}
