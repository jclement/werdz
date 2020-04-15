package main

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"gitlab.adipose.net/jeff/werdz/models/game"
)

type gameStateMessage struct {
	State       game.State          `json:"state"`
	Mode        game.Mode           `json:"mode"`
	RoundID     string              `json:"mode"`
	Word        string              `json:"word"`
	Remaining   int                 `json:"remaining"`
	Players     []playerMessage     `json:"players"`
	Definitions []definitionMessage `json:"definitions"`
}

type definitionMessage struct {
	ID         string `json:"id"`
	Definition string `json:"definition`
}

type playerMessage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func newGameStateMessage(g *game.Game, targetPlayerID game.PlayerID) gameStateMessage {
	m := gameStateMessage{
		State: g.State,
		Mode:  g.Mode,
	}

	if g.State != game.StateNew {
		r := g.CurrentRound()
		m.RoundID = string(r.ID)
		m.Word = r.Word
		if r.State == game.RoundStateOpen {
			m.Remaining = int(time.Since(r.RoundStartTime).Seconds())
		}
		if r.State == game.RoundStateVoting {
			m.Remaining = int(time.Since(r.VotingStartTime).Seconds())
			for _, d := range r.Definitions {
				// don't give the player their own definition
				if d.Player != targetPlayerID {
					m.Definitions = append(m.Definitions, definitionMessage{
						ID:         string(d.ID),
						Definition: d.Definition,
					})
				}
			}
			// randomize the order of definitions 
			for i := range m.Definitions {
				j := rand.Intn(i + 1)
				m.Definitions[i], m.Definitions[j] = m.Definitions[j], m.Definitions[i]
			}
		}
	}

	for _, p := range g.Players {
		m.Players = append(m.Players, playerMessage{
			ID:    string(p.ID),
			Name:  p.Name,
			Score: p.Score,
		})

	}

	// sort the players by score and then by name
	sort.Slice(g.Players, func(i, j int) bool {
		if g.Players[i].Score == g.Players[j].Score {
			return strings.ToUpper(g.Players[i].Name) < strings.ToUpper(g.Players[j].Name)
		}
		return g.Players[i].Score < g.Players[j].Score
	})
	return m
}

func (a *App) loop() {
	for {
		for _, g := range a.games {
			g.Dirty = g.Dirty || g.Game.Tick()
			// if something has happened in this game, push an update
			if g.Dirty {
				for c, p := range g.Clients {
					m := newGameStateMessage(g.Game, p)
					err := c.WriteJSON(m)
					if err != nil {
						g.Lock.Lock()
						g.Game.SetPlayerInactive(p, true)
						g.Lock.Unlock()
						c.Close()
						delete(g.Clients, c)
					}
				}
				g.Dirty = false
			}
		}
		time.Sleep(time.Second)
	}
}
