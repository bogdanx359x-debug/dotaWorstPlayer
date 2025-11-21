package handlers

import (
	"dotaWorstPlayerChacker/internal/openDota"
	"dotaWorstPlayerChacker/internal/redisStore"
)

type Deps struct {
	OD    *openDota.Client
	Store *redisstore.Store
}

var deps Deps

func SetDeps(d Deps) {
	deps = d
}
