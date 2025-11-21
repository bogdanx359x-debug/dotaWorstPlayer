package openDota

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	http *http.Client
	Base string
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{Timeout: 3 * time.Second},
		Base: "https://api.opendota.com/api",
	}
}

type Match struct {
	RadiantWin bool     `json:"radiant_win"`
	MatchDurSec int      `json:"duration"`
	Players    []Player `json:"players"`
}

type Player struct {
	AccountID   *int   `json:"account_id"`
	PlayerSlot  int    `json:"player_slot"`
	HeroID      int    `json:"hero_id"`
	Kills       int    `json:"kills"`
	Deaths      int    `json:"deaths"`
	Assists     int    `json:"assists"`
	GoldPerMin  int    `json:"gold_per_min"`
	XpPerMin    int    `json:"xp_per_min"`
	ObsPlaced   int    `json:"obs_placed"`
	SenPlaced   int    `json:"sen_placed"`
	TowerDamage int    `json:"tower_damage"`
	HeroDamage  int    `json:"hero_damage"`
	PersonaName string `json:"personaname"`
	HeroName    string `json:"hero_name,omitempty"`
	LaneRole     int     `json:"lane_role"`
	LastHits    int    `json:"last_hits"` 
	    DeathsLog  []struct {
        Time int    `json:"time"`
        Key  string `json:"key"`
    } `json:"deaths_log"`
	  BuybackLog []struct {
        Time int `json:"time"`
    } `json:"buyback_log"`
}

type Hero struct {
	HeroID        int    `json:"id"`
	Name          string `json:"name"`
	LocalizedName string `json:"localized_name"`
}
type Heroes []Hero

func (p *Player) IsRadiant() bool {
	return p.PlayerSlot < 128
}

func (p *Player) IsDire() bool {
	return p.PlayerSlot >= 128
}

func (p *Player) GetTeam() string {
	if p.IsRadiant() {
		return "radiant"
	}
	return "dire"
}
func (c *Client) GetHeroes() (Heroes, error) {
	url := strings.TrimRight(c.Base, "/") + "/heroes"
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request failed: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenDota request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10)) // 4KB
		return nil, fmt.Errorf("OpenDota status %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	var h Heroes
	if err := json.NewDecoder(resp.Body).Decode(&h); err != nil {
		return nil, fmt.Errorf("Decode match failed: %w", err)
	}
	return h, nil
}

func (c *Client) GetMatches(matchID string) (*Match, error) {
	url := fmt.Sprintf("%s/matches/%s", c.Base, matchID)
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenDota request filed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("OpenDota StatusCode: %s", resp.Status)
	}

	var m Match
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("Decode match filed: %w", err)
	}
	return &m, nil
}
