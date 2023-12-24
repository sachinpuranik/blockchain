package block

// Importing fmt and time
import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Transaction struct {
	senderAddress    string
	recipientAddress string
	amount           float32
}

func NewTransaction(sender string, receiver string, amount float32) *Transaction {
	return &Transaction{sender, receiver, amount}
}

func (t *Transaction) Print() {
	fmt.Printf("\n		SenderAddress		:%v", t.senderAddress)
	fmt.Printf("\n		RecipientAddress 	:%v", t.recipientAddress)
	fmt.Printf("\n		Amount 				:%f", t.amount)
}

func (t Transaction) MarshaJSONl() []byte {
	jsonBytes, _ := json.Marshal(struct {
		SenderAddress    string  `json:"sender_address"`
		RecipientAddress string  `json:"recipient_address"`
		Amount           float32 `json:"amount"`
	}{
		SenderAddress:    t.senderAddress,
		RecipientAddress: t.recipientAddress,
		Amount:           t.amount,
	})
	return jsonBytes
}

// ***********

type Block struct {
	perviousHash []byte
	transactions []*Transaction
	timestamp    int64
	nonce        int
}

func NewBlock(nonce int, previousHash []byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.transactions = transactions
	b.nonce = nonce
	b.perviousHash = previousHash
	return b
}

func (blk *Block) Print() {
	fmt.Printf("\nperviousHash 	:%v", blk.perviousHash)
	fmt.Printf("\ntimestamp 	:%v", blk.timestamp)
	fmt.Printf("\nnonce 		:%v", blk.nonce)
	for i, t := range blk.transactions {
		fmt.Printf("\ntransaction[%d]", i)
		t.Print()
	}
}

func (blk *Block) Hash() []byte {
	blkJson := blk.MarshalJSON()
	hash := sha256.Sum256(blkJson)
	return hash[:]
}

func (blk Block) MarshalJSON() []byte {
	out, _ := json.Marshal(struct {
		PerviousHash []byte         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
	}{
		PerviousHash: blk.perviousHash,
		Transactions: blk.transactions,
		Timestamp:    blk.timestamp,
		Nonce:        blk.nonce,
	})
	return out
}

// ***********

func NewGenesisBlock() *Block {
	return NewBlock(0, []byte("initial hash"), []*Transaction{{}})
}

type BlockChain struct {
	chain             []*Block
	transactionPool   []*Transaction
	blockchainAddress string
}

func NewBlockChain(blockchainAddress string) *BlockChain {
	blk := NewGenesisBlock()
	bc := &BlockChain{}
	bc.blockchainAddress = blockchainAddress
	bc.chain = append(bc.chain, blk)
	return bc
}

func (bc *BlockChain) Print() {
	fmt.Printf("\n%s", strings.Repeat("=", 25))
	for _, blk := range bc.chain {
		blk.Print()
	}
}

func (bc *BlockChain) CreateBlock(nonce int, previousHash []byte) {
	blk := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, blk)
	bc.transactionPool = []*Transaction{}
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) AddTransaction(t *Transaction) {
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *BlockChain) CopyTransactionPool() []*Transaction {
	copyPool := make([]*Transaction, len(bc.transactionPool))
	for i, t := range bc.transactionPool {
		copyPool[i] = NewTransaction(t.senderAddress, t.recipientAddress, t.amount)
	}
	return copyPool
}

func (bc *BlockChain) ValidateProof(nonce int, previousHash []byte, transactions []*Transaction, difficultyLevel int) bool {
	zeros := strings.Repeat("0", difficultyLevel)
	guessBlk := NewBlock(nonce, previousHash, transactions)
	guessHash := fmt.Sprintf("%x", guessBlk.Hash())
	return guessHash[:difficultyLevel] == zeros
}

func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidateProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce++
	}
	return nonce
}

func (bc *BlockChain) Mining() bool {
	t := NewTransaction(MINING_SENDER, bc.blockchainAddress, float32(MINING_REWARD))
	bc.AddTransaction(t)
	nonce := bc.ProofOfWork()
	bc.CreateBlock(nonce, bc.LastBlock().Hash())
	return true
}
