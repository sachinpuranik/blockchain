package block

// Importing fmt and time
import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sachinpuranik/blockchain/utils"
)

var (
	MINING_DIFFICULTY         = 3
	MINING_SENDER             = "THE BLOCKCHAIN"
	MINING_REWARD     float32 = 1.0
)

type Transaction struct {
	senderBlockChainAddress    string
	recipientBlockChainAddress string
	amount                     float32
}

func NewTransaction(senderBlockChainAddress string, recipientBlockChainAddress string, amount float32) *Transaction {
	t := new(Transaction)
	t.senderBlockChainAddress = senderBlockChainAddress
	t.recipientBlockChainAddress = recipientBlockChainAddress
	t.amount = amount
	return t
}

func (t *Transaction) Print() {
	fmt.Printf("\n	SenderAddress 	:%s", t.senderBlockChainAddress)
	fmt.Printf("\n	RecipientAddress 	:%s", t.recipientBlockChainAddress)
	fmt.Printf("\n	Amount 				:%f", t.amount)
}

func (t Transaction) MarshaJSON() []byte {
	jsonBytes, _ := json.Marshal(struct {
		SenderBlockChainAddress    string  `json:"sender_blockchain_address"`
		RecipientBlockChainAddress string  `json:"recipient_blockchain_address"`
		Amount                     float32 `json:"amount"`
	}{
		SenderBlockChainAddress:    t.senderBlockChainAddress,
		RecipientBlockChainAddress: t.recipientBlockChainAddress,
		Amount:                     t.amount,
	})
	return jsonBytes
}

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

func (bc *BlockChain) AddTransaction(sender string, recipient string, amount float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	if sender == recipient {
		fmt.Println("sender and receiver can not be same")
		return false
	}
	t := NewTransaction(sender, recipient, amount)
	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		/*if bc.CalculateTotalAmount(sender) < amount {
			return false
		}*/
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}
	return false
}

func (bc *BlockChain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *BlockChain) CopyTransactionPool() []*Transaction {
	copyPool := make([]*Transaction, len(bc.transactionPool))
	for i, t := range bc.transactionPool {
		copyPool[i] = NewTransaction(t.senderBlockChainAddress, t.recipientBlockChainAddress, t.amount)
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
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	bc.CreateBlock(nonce, bc.LastBlock().Hash())
	return true
}

func (bc *BlockChain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			amount := t.amount
			if blockchainAddress == t.recipientBlockChainAddress {
				totalAmount += amount
			}

			if blockchainAddress == t.senderBlockChainAddress {
				totalAmount -= amount
			}
		}
	}
	return totalAmount
}
