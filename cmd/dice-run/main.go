package main

import (
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Square struct {
	Apply func(s *Status)
}

type Status struct {
	SquareIndex       int
	Laps              int
	Pieces            int
	HasumiPieces      int
	EnhancementStones int
	Reports           int
	CreditsK          int
}

var Squares = []*Square{
	{
		Apply: func(s *Status) {
			// do nothing
		},
	},
	{
		Apply: func(s *Status) {
			s.Pieces++
		},
	},
	{
		Apply: func(s *Status) {
			s.EnhancementStones += 6
		},
	},
	{
		Apply: func(s *Status) {
			s.Reports += 6
		},
	},
	{
		Apply: func(s *Status) {
			s.SquareIndex += 3
		},
	},
	{
		Apply: func(s *Status) {
			s.CreditsK += 800
		},
	},
	{
		Apply: func(s *Status) {
			s.Pieces++
		},
	},
	{
		Apply: func(s *Status) {
			s.CreditsK += 500
		},
	},
	{
		Apply: func(s *Status) {
			s.EnhancementStones += 6
		},
	},
	{
		Apply: func(s *Status) {
			s.CreditsK += 1200
		},
	},
	{
		Apply: func(s *Status) {
			s.SquareIndex++
		},
	},
	{
		Apply: func(s *Status) {
			s.Reports += 4
		},
	},
	{
		Apply: func(s *Status) {
			s.Pieces++
		},
	},
	{
		Apply: func(s *Status) {
			s.Reports += 2
		},
	},
	{
		Apply: func(s *Status) {
			s.EnhancementStones += 2
		},
	},
	{
		Apply: func(s *Status) {
			s.HasumiPieces += 3
		},
	},
	{
		Apply: func(s *Status) {
			s.SquareIndex += 2
		},
	},
	{
		Apply: func(s *Status) {
			s.HasumiPieces += 2
		},
	},
}

var r = rand.New(rand.NewSource(100))

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: [num of iterate]")
	}
	n, err := strconv.Atoi(strings.ReplaceAll(os.Args[1], ",", ""))
	if err != nil {
		log.Fatal("usage: [num of iterate]")
	}

	log.Printf("Num of iterate: %d", n)
	if err := run(n); err != nil {
		log.Fatal(err)
	}
}

func run(n int) error {
	s := &Status{}

	for i := 0; i < n; i++ {
		dice := r.Intn(6) + 1
		prevIndex := s.SquareIndex
		nextIndex := (s.SquareIndex + dice) % len(Squares)

		s.SquareIndex = nextIndex
		Squares[nextIndex].Apply(s)

		s.SquareIndex %= len(Squares)
		if s.SquareIndex < prevIndex {
			s.Laps++
		}
	}
	s.HasumiPieces += s.Laps

	log.Printf("Status: %+v", s)

	per := float64(10)
	p := message.NewPrinter(language.Japanese)

	rv := reflect.ValueOf(*s)
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Name == "SquareIndex" {
			continue
		}

		val := rv.FieldByName(field.Name)
		if val.CanInt() {
			num := val.Int()
			p.Printf("%s: %v per %d, (%d)\n",
				field.Name,
				strconv.FormatFloat(float64(num)/float64(n)*per, 'f', 1, 64),
				int(per),
				num,
			)
		}
	}

	return nil
}
