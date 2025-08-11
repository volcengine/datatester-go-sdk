package utils

import "github.com/bytedance/sonic"

var NumberJsonApi = sonic.Config{
	CopyString:  true,
	SortMapKeys: true,
	EscapeHTML:  true,
	UseNumber:   true,
}.Froze()
