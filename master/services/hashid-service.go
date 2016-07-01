package services

import (
	"time"

	"github.com/speps/go-hashids"
)

type HashIDService struct {
	hashGen *hashids.HashID
}

func (h *HashIDService) Init() {
	hashData := hashids.NewData()
	hashData.Salt = "privateSalt"
	h.hashGen = hashids.NewWithData(hashData)
}

func (h HashIDService) GenerateHash(v int) string {

	time := time.Now()

	hash, _ := h.hashGen.Encode([]int{v, time.Year(), time.YearDay(), time.Hour(), time.Minute(), time.Second(), time.Nanosecond()})

	return hash
}

func (h HashIDService) DecodeHash(hash string) int {
	numbers, _ := h.hashGen.DecodeWithError(hash)

	return numbers[0]
}
