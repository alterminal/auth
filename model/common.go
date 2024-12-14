package model

type Pair[L, R any] struct {
	Left  L `json:"left"`
	Right R `json:"right"`
}
