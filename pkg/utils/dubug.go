package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func Debug(data any) {
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Printf("Error %v ", err)
	}
	fmt.Printf("%v\n", string(bytes))
}

func Output(data any) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error %v", err)
	}
	return bytes
}
