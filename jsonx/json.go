package jsonx

import "github.com/bytedance/sonic"

var (
	json          = sonic.ConfigFastest
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)
