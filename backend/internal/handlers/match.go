package handlers

import (
	"dotaWorstPlayerChacker/internal/feeder"
	"dotaWorstPlayerChacker/internal/openDota"
	"encoding/json"
	"net/http"
	"strings"
)

func FeederByMatch(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/match/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "feeder" || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	matchID := parts[0]

	match, err := deps.OD.GetMatches(matchID)
	if err != nil {
		http.Error(w, "failed to fetch match from OpenDota: "+err.Error(), http.StatusBadGateway)
		return
	}

	selectedTeam := r.FormValue("team")
	if selectedTeam != "radiant" && selectedTeam != "dire" {
		selectedTeam = "radiant"
	}
	var filteredPlayers []openDota.Player
	heroNameCache := make(map[int]string)
	for _, player := range match.Players {
		includeRadiant := selectedTeam == "radiant" && player.IsRadiant()
		includeDire := selectedTeam == "dire" && player.IsDire()
		if !includeRadiant && !includeDire {
			continue
		}

		if player.HeroID != 0 {
			name, ok := heroNameCache[player.HeroID]
			if !ok {
				if heroName, found, err := deps.Store.GetHeroName(r.Context(), player.HeroID); err == nil && found {
					name = heroName
				}
				heroNameCache[player.HeroID] = name
			}
			if name != "" {
				player.HeroName = name
			}
		}
		filteredPlayers = append(filteredPlayers, player)
	}

	worst, err := feeder.EvaluateWorstPlayer(*match, filteredPlayers)
	if err != nil {
		http.Error(w, "failed to evaluate players: "+err.Error(), http.StatusBadRequest)
		return
	}

	feederHeroName := worst.Player.HeroName
	resp := map[string]any{
		"match_id":    matchID,
		"radiant_win": match.RadiantWin,
		"feeder":      worst.Player,
		"feeder_hero": feederHeroName,
		"score":       worst.Score,
		"tie":         worst.Tie,
		"reasons":     worst.Reasons,
		"explain":     strings.Join(worst.Reasons, "; "),
		"players":     filteredPlayers,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
