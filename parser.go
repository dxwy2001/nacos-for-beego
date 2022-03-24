package nacos

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

var (
	bDQuote = []byte{'"'}
)

type Parser interface {
	Parse(content string) (map[string]string, error)
}

type IniParser struct {
}

func (i IniParser) Parse(content string) (map[string]string, error) {
	data := make(map[string]string)
	tmpBuf := bytes.NewBuffer(nil)
	buf := bufio.NewReader(bytes.NewBuffer([]byte(content)))
	for {
		tmpBuf.Reset()
		//line, _, err := buf.ReadLine()
		line, _, err := buf.ReadLine()
		fmt.Println(string(bytes.Trim(line, string('\n'))))
		if err != nil {
			if err == io.EOF {
				break
			}
			//fmt.Println(err.Error())
		}
		tmpBuf.Write(line)
		row := tmpBuf.Bytes()

		if bytes.Equal(row, []byte{}) || bytes.HasPrefix(row, []byte("#")) || bytes.HasPrefix(row, []byte(";")) {
			continue
		}
		keyValue := bytes.SplitN(row, bDQuote, 2)
		key := string(bytes.TrimSpace(keyValue[0]))
		key = strings.ToLower(key)

		if len(keyValue) == 1 {
			continue
		}
		val := bytes.TrimSpace(keyValue[1])

		if bytes.HasPrefix(val, []byte{'"'}) {
			val = bytes.Trim(val, `=`)
		}
		data[key] = string(val)
	}
	return data, nil
}
