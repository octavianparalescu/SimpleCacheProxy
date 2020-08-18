package HTTP

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

// go binary encoder
func EncodeResponse(httpResponse Response) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(httpResponse)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// go binary decoder
func DecodeResponse(str string) Response {
	m := Response{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return m
}

func EncodePath(path string) string {
	pathMD5sum := md5.Sum([]byte(path))
	pathMD5 := hex.EncodeToString(pathMD5sum[:])

	return pathMD5
}
