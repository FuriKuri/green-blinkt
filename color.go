package main

var Colors map[string]LedColor

type LedColor struct {
	Red   int32 `json:"red"`
	Green int32 `json:"green"`
	Blue  int32 `json:"blue"`
	Led   int32 `json:"led"`
}
