package game

import (
	"encoding/json"
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

// PlayerID is the unique ID for a player
type PlayerID string

// RoundState represents the status of a round
type RoundState int

const (
	// RoundStateOpen is for a round that is open for submissions
	RoundStateOpen State = iota
	// RoundStateVoting for a round that is open for voting
	RoundStateVoting
	// RoundStateComplete for a completed round
	RoundStateComplete
)

// PlayerState represents the state of a player in the game
type PlayerState struct {
	ID     PlayerID
	Name   string
	Score  int
	Active bool
}

// Definition represents a possible definition
type Definition struct {
	Definition string
	Votes      []PlayerID
}

// PlayerDefinition represents a definition provided by a player
type PlayerDefinition struct {
	Player PlayerID
	Definition
}

// RoundData represents data collected during a round of the game
type RoundData struct {
	Round       int
	State       RoundState
	Word        string
	RoundTime   time.Time
	Definitions []Definition
}

// Game represents the state of an instance of a game
type Game struct {
	ID            string
	State         State
	Mode          Mode
	RoundDuration int
	Players       []PlayerState
	MaxRounds     int
	CurrentRound  int
	RoundData     []RoundData
	StartTime     time.Time
}

func generateRoomID() string {
	letters := []rune("23456789ABCDEFGHJKMNPQRSTUVWXYZ")
	id := make([]rune, 5)

	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}

	return string(id)
}

// GeneratePlayerID generates a unique ID for a player
func GeneratePlayerID() PlayerID {
	letters := []rune("123456789ABCDEFGHJKMNPQRSTUVWXYZ")
	id := make([]rune, 20)

	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}

	return PlayerID(id)
}

// NewGame returns a new game instance
func NewGame(mode Mode, maxRounds int, roundDuration int) Game {
	// bail on unreasonable maxRounds or roundDuration
	return Game{
		ID:            generateRoomID(),
		Mode:          mode,
		MaxRounds:     maxRounds,
		RoundDuration: roundDuration,
		CurrentRound:  0,
		State:         StateNew,
	}
}

// AddPlayer adds a player to the game
func (g *Game) AddPlayer(id PlayerID, name string) {
	// make sure player doesn't exist already
	g.Players = append(g.Players, PlayerState{ID: id, Name: name, Active: true})

}

// StartGame begins the game
func (g *Game) StartGame() {
	// bail if game already started
	// bail if no players
	g.State = StateComplete
	g.StartTime = time.Now()
}

// CloseSubmissionsForCurrentRound closes the round for new definions from players and starts voting
func (g *Game) CloseSubmissionsForCurrentRound() {

}

// CloseCurrentRound closes voting on the current round and tallies up the scores
func (g *Game) CloseCurrentRound() {

}

// SubmitWord submits a definition
func (g *Game) SubmitWord(player PlayerID, round int, definition string) {

}

// EndGame sets a game status to complete
func (g *Game) EndGame() {
	g.State = StateComplete
}

// Serialize the game to a writer
func (g *Game) Serialize(writer io.Writer) {
	enc := json.NewEncoder(writer)
	enc.Encode(g)
}

// Unserialize a game from a reader
func Unserialize(reader io.Reader, g *Game) error {
	dec := json.NewDecoder(reader)
	return dec.Decode(&g)
}
