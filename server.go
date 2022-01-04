package main

import (
	. "blockchain-example/blockchain"
	"bufio"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var Blockchain []Block

var BlockchainServer chan []Block

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	BlockchainServer = make(chan []Block)
	currTime := time.Now()
	genesisBlock := Block{Index: 0, Timestamp: currTime.String(), BPM: 0, Hash: "", PrevHash: ""}
	genesisBlock.Hash = CalculateHash(genesisBlock)
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)
	// start and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("PORT"))

	if err != nil {
		log.Fatal(err)
	}

	defer func(server net.Listener) {
		err := server.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(server)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	_, err := io.WriteString(conn, "Enter a new BPM: ")
	if err != nil {
		log.Println("Error responding to client with prompt")
		return
	}
	scanner := bufio.NewScanner(conn)
	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v is NaN: %v", scanner.Text(), err)
				continue
			}
			var oldBlock = Blockchain[len(Blockchain)-1]
			newBlock, err := GenerateBlock(oldBlock, bpm)
			if IsBlockValid(newBlock, oldBlock) {
				newBlockchain := append(Blockchain, newBlock)
				Blockchain = ReplaceChain(Blockchain, newBlockchain)
			}
			BlockchainServer <- Blockchain
			io.WriteString(conn, "\nEnter a new BPM: ")

		}
	}()
	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()

	for _ = range BlockchainServer {
		spew.Dump(Blockchain)
	}
}
