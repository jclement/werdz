package main

import (
	"sort"
	"strings"
	"time"

	"gitlab.adipose.net/jeff/werdz/models/game"
)

type gameStateMessage struct {
	State         game.State       `json:"state"`
	Mode          game.Mode        `json:"mode"`
	RemainingTime int              `json:"remainingTime"`
	TotalTime     int              `json:"totalTime"`
	TotalRounds   int              `json:"totalRounds"`
	Players       []*playerMessage `json:"players"`
	Rounds        []*roundMessage  `json:"rounds"`
	CurrentRound  *roundMessage    `json:"currentRound"`
	CanSubmit     bool             `json:"canSubmit"`
	CanVote       bool             `json:"canVote"`
	CanStart      bool             `json:"canStart"`
}

type roundMessage struct {
	ID          string               `json:"id"`
	Num         int                  `json:"num"`
	State       game.RoundState      `json:"state"`
	Word        string               `json:"word"`
	Definitions []*definitionMessage `json:"definitions"`
}

type definitionMessage struct {
	ID            string   `json:"id"`
	Definition    string   `json:"definition"`
	OwnDefinition bool     `json:"ownDefinition"`
	Player        string   `json:"player"`
	Votes         []string `json:"votes"`
}

type playerMessage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Score     int    `json:"score"`
	Voted     bool   `json:"voted"`
	Submitted bool   `json:"submitted"`
}

func generateRoundSummary(targetPlayerID game.PlayerID, g *game.Game, r *game.RoundData, roundNumber int) *roundMessage {
	rm := roundMessage{
		Num:   roundNumber,
		ID:    string(r.ID),
		State: r.State,
		Word:  r.Word,
	}
	for _, d := range r.Definitions {
		dm := &definitionMessage{
			ID:            string(d.ID),
			Definition:    strings.ToLower(d.Definition),
			OwnDefinition: targetPlayerID == d.Player,
		}
		rm.Definitions = append(rm.Definitions, dm)
		if g.State == game.StateComplete || r.State == game.RoundStateVotingComplete {
			for _, v := range d.Votes {
				dm.Votes = append(dm.Votes, string(v))
			}
			dm.Player = string(d.Player)
		}
	}

	return &rm
}

func newGameStateMessage(g *game.Game, targetPlayerID game.PlayerID) gameStateMessage {
	m := gameStateMessage{
		State:       g.State,
		Mode:        g.Mode,
		CanStart:    g.CanStartGame(),
		TotalRounds: g.NumRounds,
	}

	if g.State == game.StateActive {
		r := g.CurrentRound()
		if r.State == game.RoundStateOpen {
			m.RemainingTime = g.SubmissionDuration - int(time.Since(r.RoundStartTime).Seconds())
			m.TotalTime = g.SubmissionDuration
			m.CanSubmit = true
			for _, d := range r.Definitions {
				if d.Player == targetPlayerID {
					m.CanSubmit = false
				}
			}
		}
		if r.State == game.RoundStateVoting {
			m.RemainingTime = g.VotingDuration - int(time.Since(r.VotingStartTime).Seconds())
			m.TotalTime = g.VotingDuration
			m.CanVote = true
			for _, d := range r.Definitions {
				for _, v := range d.Votes {
					if v == targetPlayerID {
						m.CanVote = false
					}
				}
			}
		}
		if r.State == game.RoundStateVotingComplete {
			m.RemainingTime = g.VotingCompleteDuration - int(time.Since(r.VotingCompleteStartTime).Seconds())
			m.TotalTime = g.VotingCompleteDuration
		}
		m.CurrentRound = generateRoundSummary(targetPlayerID, g, r, len(g.Rounds))
	}

	if g.State == game.StateComplete {
		for i, t := range g.Rounds {
			m.Rounds = append(m.Rounds, generateRoundSummary(targetPlayerID, g, t, i+1))
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
		return g.Players[i].Score > g.Players[j].Score
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
