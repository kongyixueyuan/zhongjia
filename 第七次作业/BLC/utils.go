package BLC

import (
	"encoding/binary"
	"bytes"
	"log"
	"encoding/json"
	"fmt"
	"encoding/gob"
)

func ZjJSONToArray(jsonString string) []string {

	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}

func ZjIntToHex(num int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))
	return buf
}

func ZjIntToHex2(num int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

// 字节数组反转
func ZjReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}


//version 转字节数组
func zjCommandToBytes(command string) []byte {
	var bytes [COMMANDLENGTH]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

//字节数组转version
func zjBytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

// 将结构体序列化成字节数组
func zjGobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
