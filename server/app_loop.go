package main

import (
	"time"

	"gitlab.adipose.net/jeff/werdz/models/game"
)

type gameStateMessage struct {
	State     game.State `json:"state"`
	Mode      game.Mode  `json:"mode"`
	RoundID   string     `json:"mode"`
	Remaining int        `json:"remaining"`
}

func newGameStateMessage(g *game.Game) gameStateMessage {
	m := gameStateMessage{
		State: g.State,
		Mode:  g.Mode,
	}

	if g.State != game.StateNew {
		r := g.CurrentRound()
		m.RoundID = string(r.ID)
		if r.State == game.RoundStateOpen {
			m.Remaining = int(time.Since(r.RoundStartTime).Seconds())
		}
		if r.State == game.RoundStateVoting {
			m.Remaining = int(time.Since(r.VotingStartTime).Seconds())
		}
	}

	return m
}

func (a *App) loop() {
	for {
		for _, g := range a.games {
			g.Game.Tick()
			if g.Dirty {
				m := newGameStateMessage(g.Game)
				for c, p := range g.Clients {
					err := c.WriteJSON(m)
					if err != nil {
						g.Lock.Lock()
						g.Game.SetPlayerInactive(p, true)
						g.Lock.Unlock()
						c.Close()
						delete(g.Clients, c)
					}
				}
			}
		}
		time.Sleep(time.Second)
	}
}
