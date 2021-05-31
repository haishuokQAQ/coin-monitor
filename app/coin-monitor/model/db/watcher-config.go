package db

type WatcherConfig struct {
	Symbol    string `json:"symbol"`
	Types     string `json:"types"`
	Intervals string `json:"intervals"`
}
