package feeder

import "dotaWorstPlayerChacker/internal/openDota"

func CheckMostDeaths(players []openDota.Player, acc *Acc) {
	if len(players) == 0 {
		return
	}
	maxDeaths, idx := -1, -1
	all := make([]float64, 0, len(players))
	for i, p := range players {
		all = append(all, float64(p.Deaths))
		if p.Deaths > maxDeaths {
			maxDeaths, idx = p.Deaths, i
		}
	}
	who := players[idx]

	acc.Bump(who, "Найбільше смертей", zscore(float64(who.Deaths), all))
}

func CheckDPM(players []openDota.Player, match openDota.Match, acc *Acc) {
	if len(players) == 0 {
		return
	}
	mins := max(1, match.MatchDurSec/60)
	allDPM := make([]float64, len(players))
	idx, best := 0, -1.0
	for i, p := range players {
		allDPM[i] = fdiv(p.Deaths, mins)
		if allDPM[i] > best {
			idx, best = i, allDPM[i]
		}
	}
	who := players[idx]
	acc.Bump(who, "Найбільше смертей за хвилину", zscore(float64(best), allDPM))
}

func CheckEarlyDeaths(players []openDota.Player, acc *Acc) {
	if len(players) == 0 {
		return
	}
	maxEarly, idx := -1, -1
	all := make([]float64, 0, len(players))
	for i, p := range players {
		ed := 0
		for _, d := range p.DeathsLog {
			if d.Time <= 10*60 {
				ed++
			}
		}
		all = append(all, float64(ed))
		if ed > maxEarly {
			maxEarly, idx = ed, i
		}
	}
	if idx < 0 {
		return
	}
	who := players[idx]
	acc.Bump(who, "Нафідив на лайн стадії", zscore(float64(maxEarly), all))

}

func CheckPostBuybackDeaths(players []openDota.Player, acc *Acc) {
	if len(players) == 0 {
		return
	}
	const window = 120
	type count struct {
		i, n int
	}
	best := count{-1, -1}
	all := make([]float64, 0, len(players))
	for i, p := range players {
		if len(p.BuybackLog) == 0 || len(p.DeathsLog) == 0 {
			all = append(all, 0)
			continue
		}
		n := 0
		for _, bb := range p.BuybackLog {
			for _, d := range p.DeathsLog {
				if d.Time > bb.Time && d.Time < bb.Time+window {
					n++
				}
			}
		}
		all = append(all, float64(n))
		if n > best.n {
			best = count{i, n}
		}
	}

	if best.i < 0 || best.n <= 0 {
		return
	}
	who := players[best.i]
	acc.Bump(who, "Смерті після байбека", zscore(float64(best.n), all))
}

func CheckLowImpactPerMin(players []openDota.Player, match openDota.Match, acc *Acc) {
	if len(players) == 0 {
		return
	}
	mins := max(1, match.MatchDurSec/60)
	all := make([]float64, len(players))
	idx := 0
	for i, p := range players {
		all[i] = (float64(p.HeroDamage) + 0.5*float64(p.TowerDamage)) / float64(mins)
		if all[i] < all[idx] {
			idx = i
		}
	}
	who := players[idx]
	acc.Bump(who, "найнижчий імпакт(хіро/тавер демедж)", -zscore(all[idx], all))

}

func teamKills(players []openDota.Player) int {
	sum := 0
	for _, p := range players {
		sum += p.Kills
	}
	return sum
}

func CheckLowKP(players []openDota.Player, acc *Acc) {
	if len(players) == 0 {
		return
	}
	type count struct {
		i  int
		kp float64
	}
	var allKP []float64
	best := count{-1, 1e9}
	tk := teamKills(players)
	if tk == 0 {
		return
	}
	for i, p := range players {
		kp := float64(p.Kills+p.Assists) / float64(tk)
		allKP = append(allKP, kp)
		if kp < best.kp {
			best = count{i, kp}
		}
	}
	if best.i < 0 {
		return
	}
	who := players[best.i]
	acc.Bump(who, "Найменьше вбивств та асістів порівняно з командою", -zscore(best.kp, allKP))

}
