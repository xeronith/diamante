package model

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/xeronith/diamante/protobuf"
)

func NewCompoundObject() *protobuf.CompoundObject {
	return &protobuf.CompoundObject{
		Sample: &protobuf.Sample{},
	}
}

func CreateRandomCompoundObject() *protobuf.CompoundObject {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(89999) + 10000
	sample := CreateRandomSample()
	return &protobuf.CompoundObject{
		Title:  fmt.Sprintf("Compound Object %d", id),
		Sample: sample,
	}
}
