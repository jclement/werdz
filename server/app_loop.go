package main

import (
	"sort"
	"strings"
	"time"

	"gitlab.adipose.net/jeff/werdz/models/game"
)

type gameStateMessage struct {
	State       game.State           `json:"state"`
	Mode        game.Mode            `json:"mode"`
	RoundID     string               `json:"roundId"`
	Round       int                  `json:"round"`
	RoundState  game.RoundState      `json:"roundState"`
	Word        string               `json:"word"`
	Remaining   int                  `json:"remaining"`
	Players     []*playerMessage     `json:"players"`
	Definitions []*definitionMessage `json:"definitions"`
	CanSubmit   bool                 `json:"canSubmit"`
	CanVote     bool                 `json:"canVote"`
	CanStart    bool                 `json:"canStart"`
}

type definitionMessage struct {
	ID            string `json:"id"`
	Definition    string `json:"definition"`
	OwnDefinition bool   `json:"ownDefinition"`
}

type playerMessage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Score     int    `json:"score"`
	Voted     bool   `json:"voted"`
	Submitted bool   `json:"submitted"`
}

func newGameStateMessage(g *game.Game, targetPlayerID game.PlayerID) gameStateMessage {
	m := gameStateMessage{
		State:    g.State,
		Mode:     g.Mode,
		CanStart: g.CanStartGame(),
	}

	if g.State != game.StateNew {
		r := g.CurrentRound()
		m.Round = len(g.Rounds)
		m.RoundID = string(r.ID)
		m.RoundState = r.State
		m.Word = r.Word
		if r.State == game.RoundStateOpen {
			m.Remaining = g.SubmissionDuration - int(time.Since(r.RoundStartTime).Seconds())
			m.CanSubmit = true
			for _, d := range r.Definitions {
				if d.Player == targetPlayerID {
					m.CanSubmit = false
				}
			}
		}
		if r.State == game.RoundStateVoting {
			m.Remaining = g.VotingDuration - int(time.Since(r.VotingStartTime).Seconds())
			m.CanVote = true
			for _, d := range r.Definitions {
				dm := &definitionMessage{
					ID:         string(d.ID),
					Definition: d.Definition,
				}
				if d.Player != targetPlayerID {
					for _, v := range d.Votes {
						if v == targetPlayerID {
							m.CanVote = false
						}
					}
				} else {
					dm.OwnDefinition = true
				}
				m.Definitions = append(m.Definitions, dm)
			}
		}
		if r.State == game.RoundStateVotingComplete {
			m.Remaining = g.VotingCompleteDuration - int(time.Since(r.VotingCompleteStartTime).Seconds())
			for _, d := range r.Definitions {
				dm := &definitionMessage{
					ID:         string(d.ID),
					Definition: d.Definition,
				}
				if d.Player == targetPlayerID {
					dm.OwnDefinition = true
				}
				m.Definitions = append(m.Definitions, dm)
			}
		}
	}

	voted := make(map[game.PlayerID]bool)
	submitted := make(map[game.PlayerID]bool)

	if g.State != game.StateNew {
		for _, d := range g.CurrentRound().Definitions {
			submitted[d.Player] = true
			for _, v := range d.Votes {
				voted[v] = true
			}
		}
	}

	for _, p := range g.Players {
		msg := playerMessage{
			ID:    string(p.ID),
			Name:  p.Name,
			Score: p.Score,
		}
		if _, ok := voted[p.ID]; ok {
			msg.Voted = true
		}
		if _, ok := submitted[p.ID]; ok {
			msg.Submitted = true
		}
		m.Players = append(m.Players, &msg)
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

func (a *App) gameLoop(g *gameState) {
	for {

		// Wait for something to happen
		<-g.broadcastChan

		// TODO: probably need some smarter locking around the reads to the message without
		// locking on sends
		for c, p := range g.Clients {
			m := newGameStateMessage(g.Game, p)
			err := c.WriteJSON(m)
			if err != nil {
				g.lock.Lock()
				g.Game.SetPlayerInactive(p, true)
				g.lock.Unlock()
				c.Close()
				delete(g.Clients, c)
			}
		}

		// kill the loop for this game if it's complete
		if g.Game.State == game.StateComplete {
			//return
		}
	}
}

func (a *App) loop() {
	for {
		for _, g := range a.games {
			if g.Game.Tick() {
				g.PushUpdate()
			}
		}
		time.Sleep(time.Second)
	}
}
