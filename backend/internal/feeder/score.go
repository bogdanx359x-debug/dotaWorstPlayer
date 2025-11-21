package feeder

import (
	"dotaWorstPlayerChacker/internal/openDota"
	"errors"
	"math"
	"sort"
)

type Acc struct {
	score   map[int64]int
	reasons map[int64][]string
	tie     map[int64]float64

	idx map[int64]int
}

func NewAcc() *Acc {
	return &Acc{
		score:   map[int64]int{},
		reasons: map[int64][]string{},
		tie:     map[int64]float64{},
		idx:     map[int64]int{},
	}
}

func (a *Acc) ensureIdx(players []openDota.Player) {
	for i, p := range players {
		if p.AccountID == nil {
			continue
		}
		id := int64(*p.AccountID)

		if _, ok := a.idx[id]; !ok {
			a.idx[id] = i
		}
	}
}

func (a *Acc) Bump(p openDota.Player, reason string, tieDelta float64) {
	if p.AccountID == nil {
		return
	}
	id := int64(*p.AccountID)
	a.score[id]++
	if reason != "" {
		a.reasons[id] = append(a.reasons[id], reason)
	}

	a.tie[id] += tieDelta

}

type Winner struct {
	Player  openDota.Player
	Score   int
	Tie     float64
	Reasons []string
}

func (a *Acc) PickWinner(players []openDota.Player) Winner {
	a.ensureIdx(players)
	byID := make(map[int64]openDota.Player, len(players))
	for _, p := range players {
		if p.AccountID == nil {
			continue
		}
		accId := int64(*p.AccountID)
		byID[accId] = p
	}
	type row struct {
		id    int64
		score int
		tie   float64
		order int
	}
	list := make([]row, 0, len(players))
	for _, p := range players {
		if p.AccountID == nil {
			continue
		}
		accId := int64(*p.AccountID)
		list = append(list, row{
			id:    accId,
			score: a.score[accId],
			tie:   a.tie[accId],
			order: a.idx[int64(*p.AccountID)],
		})
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].score != list[j].score {
			return list[i].score > list[j].score
		}
		if list[i].tie != list[j].tie {
			return list[i].tie > list[j].tie
		}
		return list[i].order < list[j].order
	})
	if len(list) == 0 {
		return Winner{}
	}
	top := list[0]
	return Winner{
		Player:  byID[top.id],
		Score:   top.score,
		Tie:     top.tie,
		Reasons: a.reasons[top.id],
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func fdiv(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

func zscore(x float64, xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	var m float64
	for _, v := range xs {
		m += v
	}
	m /= float64(len(xs))
	var v float64
	for _, u := range xs {
		d := u - m
		v += d * d
	}
	std := math.Sqrt(v / float64(len(xs)))
	if std == 0 {
		return 0
	}
	return (x - m) / std
}

func EvaluateWorstPlayer(match openDota.Match, players []openDota.Player) (Winner, error) {
	if len(players) == 0 {
		return Winner{}, errors.New("немає гравців для оцінки")
	}

	acc := NewAcc()
	acc.ensureIdx(players)

	CheckMostDeaths(players, acc)
	CheckDPM(players, match, acc)
	CheckEarlyDeaths(players, acc)
	CheckPostBuybackDeaths(players, acc)
	CheckLowImpactPerMin(players, match, acc)
	CheckLowKP(players, acc)

	w := acc.PickWinner(players)
	if w.Player.AccountID == nil {
		return Winner{}, errors.New("не вдалося визначити гравця (немає account_id)")
	}
	return w, nil
}
