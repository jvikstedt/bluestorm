package network

import (
	"github.com/satori/go.uuid"
)

func GenerateID() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return id.String()
}
