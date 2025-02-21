package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
)

type Person struct {
	Name        string
	DateOfBirth time.Time
	MassKg      float32
}

func main() {
	me := Person{Name: "SiuYin",
		DateOfBirth: time.Date(1962, 6, 28, 0, 0, 0, 0, time.UTC), MassKg: 103}
	writeGOB(&me, "/tmp/mydat.gob")

	me2 := readGOB("/tmp/mydat.gob")
	fmt.Printf("me2: %v\n", me2)
	fmt.Printf("Date of Birth: %s, Mass: %.1fkg\n",
		me2.DateOfBirth.Format("2 Jan 2006"), me2.MassKg)
}

func writeGOB(p *Person, fn string) {
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}

	enc := gob.NewEncoder(f)
	if err := enc.Encode(p); err != nil {
		log.Fatal(err)
	}
}

func readGOB(fn string) Person {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}

	dec := gob.NewDecoder(f)
	var p Person
	if err := dec.Decode(&p); err != nil {
		log.Fatal(err)
	}
	return p
}
