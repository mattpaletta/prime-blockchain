package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
	"google.golang.org/grpc"
	"github.com/davecgh/go-spew/spew"
	"net"
	"github.com/mattpaletta/prime-blockchain/blockchain"
	"context"
)

type bookServiceServer struct {}

func (s *bookServiceServer) GetBlock(context context.Context, block *blockchain.BlockRequest) (*blockchain.Block, error) {
	return *Blockchain[block.Index], nil
}

const difficulty = 1

var Blockchain []*blockchain.Block

var mutex = &sync.Mutex{}

func run() error {

	port := 6000

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	blockchain.RegisterBookServiceServer(grpcServer, bookServiceServer{})

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func getNextPrimeNumber(start_from int64) int64 {
	c1 := make(chan int64, 1000)
	c2 := make(chan bool, 1000)

	numToDivide := int64(0)
	for j := start_from + 1; true; j++ {
		fmt.Println("Finding next prime number.")
		for i := int64(1); i <= j; i++ {
			c1 <- i

			go func() {
				num := <-c1
				if j % num == 0 {
					c2 <- true
				} else {
					c2 <- false
				}
			}()
		}

		num_divided := 0

		for i := int64(1); i <= j; i++ {
			didDivide := <- c2
			if didDivide {
				num_divided++
			}
		}

		if num_divided == 2 {
			numToDivide = j
			break
		}
	}
	fmt.Println("Found next number: %d", numToDivide)
	return numToDivide
}

func generateBlock(oldBlock blockchain.Block) blockchain.Block {
	var newBlock blockchain.Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.VAL = getNextPrimeNumber(oldBlock.VAL)
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {
		hexValue := fmt.Sprintf("%x", i)
		newBlock.Nonce = hexValue
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " work done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}
	return newBlock
}

func isHashValid(hash string, difficulty int64) bool {
	prefix := strings.Repeat("0", int(difficulty))
	return strings.HasPrefix(hash, prefix)
}

func isBlockValid(newBlock, oldBlock blockchain.Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func calculateHash(block blockchain.Block) string {
	record := strconv.Itoa(int(block.Index)) + block.Timestamp + strconv.Itoa(int(block.VAL)) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal(err)
	//}

	go func() {
		t := time.Now()

		genesisBlock := blockchain.Block{
											Index: 0,
											Timestamp: t.String(),
											VAL: 2,
											Hash: "",
											PrevHash: "",
											Difficulty: difficulty,
											Nonce: "",
		}

		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()

	go func() {
		// How to ensure the first go routine ran already?
		time.Sleep(time.Second * 10)
		for {
			//ensure atomicity when creating new block
			mutex.Lock()
			newBlock := generateBlock(Blockchain[len(Blockchain)-1])
			mutex.Unlock()

			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				Blockchain = append(Blockchain, newBlock)
				spew.Dump(Blockchain)
			}
		}
	}()

	log.Fatal(run())

}
