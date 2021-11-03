package main

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	ct "github.com/florianl/go-conntrack"
)

// This func read from file and expose byte slice
func readFromFile(path string) []byte {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(string(content))

	return content
}

func StringToFloat(content string) float64 {

	s, err := strconv.ParseFloat(content, 64)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func printslice(slice []string) {
	fmt.Println("slice = ", slice)
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

type KVHeap []kv

func (h KVHeap) Len() int           { return len(h) }
func (h KVHeap) Less(i, j int) bool { return h[i].Value > h[j].Value }
func (h KVHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

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

	// dup_map := dup_count(temp)

	// // for k, v := range dup_map {
	// // 	fmt.Printf("Item : %s , Count : %d\n", k, v)
	// // }

	// h := getHeap(dup_map)
	// n := 10
	// for i := 0; i < n; i++ {
	// 	//fmt.Printf("%d) %#v\n", i+1, heap.Pop(h))

	// 	fmt.Print(heap.Pop(h))
	return temp
}

func HowMatches(IPs []string) map[string]int {

	dict := make(map[string]int)

	for _, ip := range IPs {
		dict[ip] = dict[ip] + 1
	}
	return dict
}
