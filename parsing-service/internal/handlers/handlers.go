package handlers

import "parsing-service/internal/db"

type Parsing interface {
	Parse(r db.Request)
}

func Parse(p Parsing, r db.Request) {
	p.Parse(r)
}
