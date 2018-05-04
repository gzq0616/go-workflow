package go_workflow

import (
	"encoding/hex"
	"crypto/md5"
)

func md5Sum(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	cipherStr := hash.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func listRemove(list []string, item string) []string {
	var ret []string
	for index, value := range list {
		if value == item {
			end := index + 1
			ret = append(ret, list[:index]...)
			ret = append(ret, list[end:]...)
		}
	}
	return ret
}
