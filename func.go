package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"

	ct "github.com/florianl/go-conntrack"
	"github.com/prometheus/procfs"
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

	// make an array of type string to store our keys
	keys := []string{}

	// iterate over the map and append all keys to our
	// string array of keys
	for key := range dict {
		keys = append(keys, key)
	}

	// use the sort method to sort our keys array
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Println(key, dict[key])
	}
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
