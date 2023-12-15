package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	buf, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	buf = bytes.TrimSpace(buf)
	words := bytes.Split(buf, []byte{','})

	var boxes [256]Box
	sum := 0
	for _, i := range words {
		sum += hashWord(i)
		processInstr(i, boxes[:])
	}
	fmt.Println("Part 1:", sum)

	sum = 0
	for n, i := range boxes {
		for k, j := range i.values {
			sum += (n + 1) * (k + 1) * j.v
		}
	}
	fmt.Println("Part 2:", sum)
}

type (
	Box struct {
		values []BoxValue
	}

	BoxValue struct {
		id string
		v  int
	}
)

func processInstr(v []byte, m []Box) {
	if r, ok := bytes.CutSuffix(v, []byte{'-'}); ok {
		h := hashWord(r)
		label := string(r)
		for n, i := range m[h].values {
			if i.id == label {
				copy(m[h].values[n:], m[h].values[n+1:])
				m[h].values = m[h].values[:len(m[h].values)-1]
				return
			}
		}
		return
	}
	if r, b, ok := bytes.Cut(v, []byte{'='}); ok {
		num, err := strconv.Atoi(string(b))
		if err != nil {
			log.Fatalln(err)
		}
		h := hashWord(r)
		label := string(r)
		for n, i := range m[h].values {
			if i.id == label {
				m[h].values[n] = BoxValue{
					id: label,
					v:  num,
				}
				return
			}
		}
		m[h].values = append(m[h].values, BoxValue{
			id: label,
			v:  num,
		})
		return
	}
	log.Fatalln("Malformed op")
}

func hashWord(v []byte) int {
	currentValue := 0
	for _, i := range v {
		currentValue = hashStep(currentValue, int(i))
	}
	return currentValue
}

func hashStep(v, b int) int {
	return ((v + b) * 17) % 256
}
