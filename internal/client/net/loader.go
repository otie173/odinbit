package net

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) LoadTextures(addr string) ([]byte, error) {
	_, err := http.Get(fmt.Sprintf("%s/ping", addr))
	if err != nil {
		log.Printf("Error! Cant ping server: %v\n", err)
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/textures", addr))
	if err != nil {
		log.Printf("Error! Cant load textures from server: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error! Cant read response body: %v\n", err)
		return nil, err
	}

	return body, nil
}

func (l *Loader) LoadWorld(addr string) ([]byte, error) {
	_, err := http.Get(fmt.Sprintf("%s/ping", addr))
	if err != nil {
		log.Printf("Error! Cant ping server: %v\n", err)
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%s/world", addr))
	if err != nil {
		log.Printf("Error! Cant load world from server: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error! Cant read response body: %v\n", err)
		return nil, err
	}

	return body, nil
}
