package game

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"time"
)

const (
	// MaximumRounds is the upper limit on the number of rounds in a game
	MaximumRounds = 10
	// MinimumRounds is the lower limit on the number of rounds in a game
	MinimumRounds = 1
	// MaximumSubmissionDuration is the maximum duration (seconds) for the submission phase of a round
	MaximumSubmissionDuration = 600
	// MinimumSubmissionDuration is the minimum duration (seconds) for the submission phase of a round
	MinimumSubmissionDuration = 60
	// MinimumVotingDuration is the minimum duration (seconds) for the voting phase of a round
	MinimumVotingDuration = 30
	// MaximumVotingDuration is the maximum duration (seconds) for the voting phase of a round
	MaximumVotingDuration = 300
)

// Mode represents the mode of a game
type Mode int

const (
	// ModeNormal is for a standard game
	ModeNormal Mode = iota
	// ModeFun for a game with no correct answer
	ModeFun
)

// State represents the status of a game
type State int

const (
	// StateNew is for a new non-started game
	StateNew State = iota
	// StateActive for an active game
	StateActive
	// StateComplete for a completed game
	StateComplete
)

// GID is the unique ID for a game
type GID string

// PlayerID is the unique ID for a player
type PlayerID string

const rightAnswerPlayerID PlayerID = ""

// RoundID is the unique ID for a round
type RoundID string

// DefinitionID is the unique ID for a definition
type DefinitionID string

// RoundState represents the status of a round
type RoundState int

const (
	// RoundStateOpen is for a round that is open for submissions
	RoundStateOpen RoundState = iota
	// RoundStateVoting for a round that is open for voting
	RoundStateVoting
	// RoundStateComplete for a completed round
	RoundStateComplete
)

// PlayerState represents the state of a player in the game
type PlayerState struct {
	ID       PlayerID
	Name     string
	Score    int
	Inactive bool
	Deleted  bool
}

// Definition represents a possible definition
type Definition struct {
	ID         DefinitionID
	Player     PlayerID
	Definition string
	Votes      []PlayerID
}

// RoundData represents data collected during a round of the game
type RoundData struct {
	ID              RoundID
	State           RoundState
	Word            string
	RoundStartTime  time.Time
	VotingStartTime time.Time
	Definitions     []*Definition
}

// Game represents the state of an instance of a game
type Game struct {
	ID                 GID
	State              State
	Mode               Mode
	SubmissionDuration int
	VotingDuration     int
	Players            []*PlayerState
	NumRounds          int
	Rounds             []*RoundData
	StartTime          time.Time
	wordSource         func() (word, definition string)
}

func generateID(unambiguous bool, length int) string {
	var letters []rune

	if unambiguous {
		letters = []rune("23456789ABCDEFGHJKMNPQRSTUVWXYZ")
	} else {
		letters = []rune("abcdefghijklmnopqrstuvwxyz123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}

	id := make([]rune, length)

	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}

	return string(id)

}

func generateGameID() GID {
	return GID(generateID(true, 5))
}

func generateDefinitionID() DefinitionID {
	return DefinitionID(generateID(false, 10))
}

func generateRoundID() RoundID {
	return RoundID(generateID(false, 20))
}

// GeneratePlayerID generates a unique ID for a player
func GeneratePlayerID() PlayerID {
	return PlayerID(generateID(false, 20))
}

// NewGame returns a new game instance
func NewGame(wordSource func() (word, definition string), mode Mode, numRounds int, submissionDuration int, votingDuration int) (*Game, error) {
	if numRounds < MinimumRounds || numRounds > MaximumRounds {
		return nil, fmt.Errorf("number of rounds must be between %d and %d", MinimumRounds, MaximumRounds)
	}
	if submissionDuration < MinimumSubmissionDuration || submissionDuration > MaximumSubmissionDuration {
		return nil, fmt.Errorf("submission duration must be between %d and %d", MinimumSubmissionDuration, MaximumSubmissionDuration)
	}
	if votingDuration < MinimumVotingDuration || votingDuration > MaximumVotingDuration {
		return nil, fmt.Errorf("voting duration must be between %d and %d", MinimumVotingDuration, MaximumVotingDuration)
	}
	return &Game{
		ID:                 generateGameID(),
		Mode:               mode,
		NumRounds:          numRounds,
		SubmissionDuration: submissionDuration,
		VotingDuration:     votingDuration,
		State:              StateNew,
		wordSource:         wordSource,
	}, nil
}

func (g *Game) findPlayer(id PlayerID) (index int, player *PlayerState, err error) {

	for i, p := range g.Players {
		if p.ID == id {
			return i, p, nil
		}
	}
	return 0, nil, fmt.Errorf("player not found")
}

// AddPlayer adds a player to the game
func (g *Game) AddPlayer(id PlayerID, name string) error {
	for _, p := range g.Players {
		if p.ID == id {
			return fmt.Errorf("player with this ID already part of this game")
		}
		if p.Name == name {
			return fmt.Errorf("player with this name already part of this game")
		}
	}
	g.Players = append(g.Players, &PlayerState{ID: id, Name: name})
	return nil
}

// RemovePlayer removes a player to the game
// Deleted players are considerd gone (but left for historical scorting purposes).
// Deleted players are not able to take actions in a game
func (g *Game) RemovePlayer(id PlayerID) error {
	_, p, err := g.findPlayer(id)
	if err != nil || p.Deleted {
		return fmt.Errorf("player does not exist")
	}
	p.Deleted = true
	return nil
}

// SetPlayerInactive gives a way to mark a player is inactive in the game.
// This is mostly just a display thing and a way of shortening the voting / submission
// process since we are waiting for the timer / all active players to vote.
func (g *Game) SetPlayerInactive(id PlayerID, inactive bool) error {
	_, p, err := g.findPlayer(id)
	if err != nil || p.Deleted {
		return fmt.Errorf("player does not exist")
	}
	p.Inactive = inactive
	return nil
}

func (g *Game) createNewRound() *RoundData {
	word, def := g.wordSource()

	r := RoundData{
		ID:             generateRoundID(),
		State:          RoundStateOpen,
		Word:           word,
		RoundStartTime: time.Now(),
	}

	if g.Mode == ModeNormal {
		// in normal play mode, we append the real definition to the list of options
		r.Definitions = append(r.Definitions, &Definition{
			ID:         generateDefinitionID(),
			Definition: def,
			Player:     rightAnswerPlayerID,
		})
	}

	return &r
}

// StartGame begins the game
func (g *Game) StartGame() error {
	if g.State != StateNew {
		return fmt.Errorf("game has alread been started")
	}
	if len(g.Players) == 0 {
		return fmt.Errorf("starting a game requires players")
	}
	g.State = StateActive
	g.StartTime = time.Now()

	g.Rounds = append(g.Rounds, g.createNewRound())

	return nil
}

// Tick is the function expected to be called by the outside appliaction to handle the automatic
// closing / scoring of rounds.
func (g *Game) Tick() (round RoundID, state RoundState, secondsRemaining int, err error) {

	// if the game hasn't started... nothing interesting going on here
	if g.State == StateNew {
		return RoundID(""), RoundStateComplete, 0, fmt.Errorf("game not started")
	}

	r := g.CurrentRound()

	var rem int

	if r.State == RoundStateOpen {
		rem = g.SubmissionDuration - int(time.Since(r.RoundStartTime).Seconds())
		if rem < 0 {
			g.closeSubmissionsForCurrentRound()
		}
	}

	if r.State == RoundStateVoting {
		rem = g.VotingDuration - int(time.Since(r.VotingStartTime).Seconds())
		if rem < 0 {
			g.completeCurrentRound()
			if g.State != StateComplete {
				rem = g.SubmissionDuration
				r = g.CurrentRound()
			}
		}
	}

	return r.ID, r.State, rem, nil
}

// CurrentRound returns the current round
func (g *Game) CurrentRound() *RoundData {
	return g.Rounds[len(g.Rounds)-1]
}

// CloseSubmissionsForCurrentRound closes the round for new definions from players and starts voting
func (g *Game) closeSubmissionsForCurrentRound() error {
	if g.State != StateActive {
		return fmt.Errorf("game is not active")
	}
	r := g.CurrentRound()
	if r.State != RoundStateOpen {
		return fmt.Errorf("round is not open")
	}
	r.State = RoundStateVoting
	r.VotingStartTime = time.Now()
	return nil
}

func (g *Game) scoreRound() {
	r := g.CurrentRound()
	for _, def := range r.Definitions {
		if def.Player == rightAnswerPlayerID {
			// this is the right answer.  give each of these people 3 pts
			for _, id := range def.Votes {
				if _, p, err := g.findPlayer(id); err == nil {
					p.Score += 3
				}
			}
		} else {
			// this is a note.  +1 to the owner
			if _, p, err := g.findPlayer(def.Player); err == nil {
				p.Score += len(def.Votes)
			}
		}
	}
}

// CompleteCurrentRound closes voting on the current round and tallies up the scores
func (g *Game) completeCurrentRound() error {
	if g.State != StateActive {
		return fmt.Errorf("game is not active")
	}

	r := g.CurrentRound()
	r.State = RoundStateComplete
	g.scoreRound()

	if len(g.Rounds) < g.NumRounds {
		// if we aren't at maxRounds, spin up a new game
		newRound := g.createNewRound()
		g.Rounds = append(g.Rounds, newRound)
	} else {
		g.EndGame()
	}

	return nil
}

// SubmitWord submits a definition
func (g *Game) SubmitWord(player PlayerID, round RoundID, definition string) error {
	if g.State != StateActive {
		return fmt.Errorf("game is not active")
	}
	r := g.CurrentRound()
	if r.State != RoundStateOpen {
		return fmt.Errorf("round is not open")
	}
	if r.ID != round {
		return fmt.Errorf("not the correct round")
	}
	if _, p, err := g.findPlayer(player); err != nil || p.Deleted {
		return fmt.Errorf("player not found")
	}
	for _, def := range r.Definitions {
		if def.Player == player {
			return fmt.Errorf("player already submitted a definition for this round")
		}
	}
	r.Definitions = append(r.Definitions, &Definition{
		ID:         generateDefinitionID(),
		Player:     player,
		Definition: definition,
	})
	return nil
}

// Vote casts a vote for a player
func (g *Game) Vote(player PlayerID, round RoundID, definition DefinitionID) error {
	if g.State != StateActive {
		return fmt.Errorf("game is not active")
	}
	r := g.CurrentRound()
	if r.State != RoundStateVoting {
		return fmt.Errorf("round is not voting")
	}
	if r.ID != round {
		return fmt.Errorf("not the correct round")
	}
	if _, p, err := g.findPlayer(player); err != nil || p.Deleted {
		return fmt.Errorf("player not found")
	}
	for _, def := range r.Definitions {
		for _, vote := range def.Votes {
			if vote == player {
				return fmt.Errorf("player already voted")
			}
		}
	}
	for _, def := range r.Definitions {
		if def.ID == definition {
			if def.Player == player {
				return fmt.Errorf("can't vote for self")
			}
			def.Votes = append(def.Votes, player)
		}
	}
	return nil
}

// EndGame sets a game status to complete
func (g *Game) EndGame() error {
	if g.State == StateComplete {
		return fmt.Errorf("game is already complete")
	}

	if g.State == StateActive {
		r := g.CurrentRound()
		if r.State != RoundStateComplete {
			r.State = RoundStateComplete
			g.scoreRound()
		}
	}

	g.State = StateComplete
	return nil
}

// Serialize the game to a writer
func (g *Game) Serialize(writer io.Writer) {
	enc := json.NewEncoder(writer)
	enc.Encode(g)
	writer.Write([]byte("test"))
}

// Unserialize a game from a reader
func Unserialize(reader io.Reader, g *Game) error {
	dec := json.NewDecoder(reader)
	return dec.Decode(&g)
}
