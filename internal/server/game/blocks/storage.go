package blocks

import (
	"encoding/json"
	"log"
	"os"
)

type Page struct {
	Number int       `json:"number"`
	Blocks [6]string `json:"blocks"`
}

type MaterialBlocks struct {
	Pages []Page `json:"pages"`
}

type Materials struct {
	Wood  MaterialBlocks `json:"wood"`
	Stone MaterialBlocks `json:"stone"`
	Metal MaterialBlocks `json:"metal"`
}

type Storage struct {
	Blocks Materials
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) LoadBlocks() {
	jsonData, err := os.ReadFile("blocks/pages.json")
	if err != nil {
		log.Printf("Error! Cant read pages json file: %v\n", err)
	}

	if err = json.Unmarshal(jsonData, &s.Blocks); err != nil {
		log.Printf("Error! Cant unmarshal binary data from pages json file: %v\n", err)
	}
}
