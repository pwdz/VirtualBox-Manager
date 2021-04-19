package pkg

import (
	"strings"
)
func ConcatJsons(in ...[]byte)[]byte{
	if len(in) < 2{
		return nil
	}
	dst := string(in[0])

	for _, inputJson := range in[1: len(in)]{
		dst = strings.TrimRight(dst, "}")	
		dst = dst + "," + strings.TrimLeft(string(inputJson), "{")	
	}
	return []byte(dst)
}