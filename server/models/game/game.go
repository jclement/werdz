package game

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"time"
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
	ID      PlayerID
	Name    string
	Score   int
	Deleted bool
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
	ID          RoundID
	State       RoundState
	Word        string
	StartTime   time.Time
	Definitions []*Definition
}

// Game represents the state of an instance of a game
type Game struct {
	ID            GID
	State         State
	Mode          Mode
	RoundDuration int
	Players       []*PlayerState
	MaxRounds     int
	Rounds        []*RoundData
	StartTime     time.Time
	wordSource    func() (word, definition string)
}

func generateGameID() GID {
	letters := []rune("23456789ABCDEFGHJKMNPQRSTUVWXYZ")
	id := make([]rune, 5)

	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}

	return GID(id)
}

func generateDefinitionID() DefinitionID {
	letters := []rune("abcdefghijklmnopqrstuvwxyz123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	id := make([]rune, 10)

	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}

	return DefinitionID(id)
}

func generateRoundID() RoundID {
	letters := []rune("abcdefghijklmnopqrstuvwxyz123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	id := make([]rune, 10)

	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}

	return RoundID(id)
}

// GeneratePlayerID generates a unique ID for a player
func GeneratePlayerID() PlayerID {
	letters := []rune("abcdefghijklmnopqrstuvwxyz123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	id := make([]rune, 20)

	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}

	return PlayerID(id)
}

// NewGame returns a new game instance
func NewGame(wordSource func() (word, definition string), mode Mode, maxRounds int, roundDuration int) Game {
	// bail on unreasonable maxRounds or roundDuration
	return Game{
		ID:            generateGameID(),
		Mode:          mode,
		MaxRounds:     maxRounds,
		RoundDuration: roundDuration,
		State:         StateNew,
		wordSource:    wordSource,
	}
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

// RemovePlayer adds a player to the game
func (g *Game) RemovePlayer(id PlayerID) error {
	_, p, err := g.findPlayer(id)
	if err != nil || p.Deleted {
		return fmt.Errorf("player does not exist")
	}
	p.Deleted = true
	return nil
}

func (g *Game) createNewRound() *RoundData {
	word, def := g.wordSource()

	r := RoundData{
		ID:        generateRoundID(),
		State:     RoundStateOpen,
		Word:      word,
		StartTime: time.Now(),
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

// CurrentRound returns the current round
func (g *Game) CurrentRound() *RoundData {
	return g.Rounds[len(g.Rounds)-1]
}

// CloseSubmissionsForCurrentRound closes the round for new definions from players and starts voting
func (g *Game) CloseSubmissionsForRound(round RoundID) error {
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
	g.CurrentRound().State = RoundStateVoting
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
func (g *Game) CompleteRound(round RoundID) error {
	if g.State != StateActive {
		return fmt.Errorf("game is not active")
	}

	r := g.CurrentRound()
	if r.ID != round {
		return fmt.Errorf("not the correct round")
	}

	r.State = RoundStateComplete
	g.scoreRound()

	if len(g.Rounds) < g.MaxRounds {
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
