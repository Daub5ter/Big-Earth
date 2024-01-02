package random

import (
	"crypto/rand"
	"math/big"
	"time"
)

func CreateNumber(min int64, max int64) (time.Duration, error) {
	number := big.NewInt(max - min + 1)
	randInt, err := rand.Int(rand.Reader, number)

	if err != nil {
		return time.Duration(0), err
	}

	return time.Duration(randInt.Int64() + min), nil
}
