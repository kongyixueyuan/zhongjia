package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutputs struct {
	ZjOutputs []*TXOutput
}

/**
	序列化对象
 */
func (outputs *TXOutputs) ZjSerialize() []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(outputs)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

/**
	反序列化，将[]byte反序列化成对象
 */
func ZjDeserializeOutputs(outputsBytes []byte) *TXOutputs {
	var outputs TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(outputsBytes))
	err := decoder.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	return &outputs
}