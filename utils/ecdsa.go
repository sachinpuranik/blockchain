package utils

import (
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) Print() {
	fmt.Printf("%x%x", s.R, s.S)
}
