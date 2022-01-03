package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"strconv"
	"time"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

//var Blockchain []Block

func CalculateHash(block Block) string {
	blockReflection := reflect.ValueOf(block)
	typeOfBlock := blockReflection.Type()
	var record string
	var nameOfField string
	var typeOfField string
	for i := 0; i < blockReflection.NumField(); i++ {
		nameOfField = typeOfBlock.Field(i).Name
		typeOfField = reflect.TypeOf(blockReflection.Field(i).Interface()).String()

		if nameOfField != "Hash" && nameOfField != "PrevHash" {
			if typeOfField == "string" {
				record += blockReflection.Field(i).Interface().(string)
			} else if typeOfField == "int" {
				record += strconv.Itoa(blockReflection.Field(i).Interface().(int))
			}
		}

	}
	encodedBytes := sha256.New()
	encodedBytes.Write([]byte(record))
	hashed := encodedBytes.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateBlock(oldBlock Block, BPM int) (Block, error) {
	currTime := time.Now().String()
	newBlock := Block{Index: oldBlock.Index + 1, Timestamp: currTime, BPM: BPM, Hash: "", PrevHash: oldBlock.Hash}
	newBlock.Hash = CalculateHash(newBlock)
	return newBlock, nil
}

func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func ReplaceChain(newBlocks, oldBlocks []Block) []Block {
	var resultChain []Block
	if len(newBlocks) > len(oldBlocks) {
		resultChain = newBlocks
	} else {
		resultChain = oldBlocks
	}
	return resultChain
}
