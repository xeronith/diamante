package model

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/xeronith/diamante/protobuf"
)

func NewSample() *protobuf.Sample {
	return &protobuf.Sample{}
}

func CreateEmptySample() *protobuf.Sample {
	return &protobuf.Sample{}
}

func CreateRandomSample() *protobuf.Sample {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(89999) + 10000
	return &protobuf.Sample{
		StringField: fmt.Sprintf("Sample %d", id),
		Int32Field:  int32(id),
		BoolField:   id >= 50000,
	}
}
