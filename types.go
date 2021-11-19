package main

type Config struct {
	Port                        string `json:"port"`
	TopRecordsCount             int    `json:"top_records_count"`
	PathToConntrackCount        string `json:"path_to_conntrack_count"`
	PathToConntrackMax          string `json:"path_to_conntrack_max"`
	MetricsRoutePath            string `json:"metrics_route_path"`
	ConntrackCountCheckInterval int    `json:"conntrack_count_check_interval"`
	ConntrackMaxCheckInterval   int    `json:"conntrack_max_check_interval"`
	ConntrackTopCheckInterval   int    `json:"conntrack_top_check_interval"`
	LogDebug                    bool   `json:"log_debug"`
}

type KVHeap []kv

type kv struct {
	Key   string
	Value int
}
