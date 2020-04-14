package game

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestGameIDUnique(t *testing.T) {
	seen := make(map[GID]bool)
	for i := 0; i < 100; i++ {
		id := generateGameID()
		if _, exists := seen[id]; exists {
			t.Error("Generated Game IDs must be unique")
			return
		}
		seen[id] = true
	}
}

func TestGameIDSize(t *testing.T) {
	if len(generateGameID()) != 5 {
		t.Error("Generated Game IDs must be 5 characters")
	}
}

func TestGameIDCase(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := string(generateGameID())
		if id != strings.ToUpper(id) {
			t.Error("Generated Game ID must be upper case")
			return
		}
	}
}

func TestPlayerIDUnique(t *testing.T) {
	seen := make(map[PlayerID]bool)
	for i := 0; i < 100; i++ {
		id := GeneratePlayerID()
		if _, exists := seen[id]; exists {
			t.Error("Generated Player IDs must be unique")
			return
		}
		seen[id] = true
	}
}

func TestPlayerIDSize(t *testing.T) {
	if len(GeneratePlayerID()) != 20 {
		t.Error("Generated Player IDs must be 20 characters")
	}
}

func testableWordGenerator() func() (word, definition string) {
	i := 0
	words := []string{"hello", "world"}
	return func() (word string, definition string) {
		if i == len(words) {
			i = 0
		}
		w := words[i]
		i++
		return w, "definition of " + w
	}
}

func TestStartGameNoPlayers(t *testing.T) {
	g := NewGame(testableWordGenerator(), ModeNormal, 3, 600)
	if err := g.StartGame(); err == nil {
		t.Error("Starting a game with no players should fail")
	}
}

func TestStartGameAlreadyStarted(t *testing.T) {
	g := NewGame(testableWordGenerator(), ModeNormal, 3, 600)
	for _, state := range []State{StateActive, StateComplete} {
		g.State = state
		if err := g.StartGame(); err == nil {
			t.Errorf("Game in state '%d' should fail", state)
		}
	}
}

func TestStartGame(t *testing.T) {
	g := NewGame(testableWordGenerator(), ModeNormal, 3, 600)
	if g.State != StateNew {
		t.Errorf("Game should start in New state")
	}
	g.AddPlayer(GeneratePlayerID(), "Tester")
	if err := g.StartGame(); err != nil {
		t.Errorf("Game should have started")
	}
	if g.State != StateActive {
		t.Error("Game state should be active")
	}
	if g.StartTime.IsZero() {
		t.Error("Start time should be set")
	}
	if len(g.Rounds) != 1 {
		t.Error("There should be a round!")
	}
}

func TestCreateNewRoundNormal(t *testing.T) {
	g := NewGame(testableWordGenerator(), ModeNormal, 3, 600)
	r := g.createNewRound()
	if r.State != RoundStateOpen {
		t.Error("Rounds should start as open")
	}
	if r.Word != "hello" {
		t.Error("Expecting 'hello' from the generator")
	}
	if len(r.Definitions) != 1 {
		t.Error("Expecting a single definition")
	} else {
		if r.Definitions[0].Definition != "definition of hello" {
			t.Error("Expecting 'definition of hello' from the generator")
		}
	}
}

func TestCreateNewRoundFun(t *testing.T) {

	g := NewGame(testableWordGenerator(), ModeFun, 3, 600)
	r := g.createNewRound()
	if r.State != RoundStateOpen {
		t.Error("Rounds should start as open")
	}
	if r.Word != "hello" {
		t.Error("Expecting 'hello' from the generator")
	}
	if len(r.Definitions) != 0 {
		t.Error("Expecting a no definitions")
	}
}

func TestCloseSubmissionsForCurrentRound(t *testing.T) {
	g := NewGame(testableWordGenerator(), ModeNormal, 3, 600)
	g.AddPlayer(GeneratePlayerID(), "Tester")
	if err := g.CloseSubmissionsForRound(RoundID("")); err == nil {
		t.Error("Can't close submissions on an inactive game")
		return
	}
	g.StartGame()
	if err := g.CloseSubmissionsForRound(RoundID("Bad ID")); err == nil {
		t.Error("Bad Round ID")
		return
	}
	if err := g.CloseSubmissionsForRound(g.CurrentRound().ID); err != nil {
		t.Error("Should have been able to close a new round")
		return
	}
	if g.Rounds[0].State != RoundStateVoting {
		t.Error("Round should be open for voting now")
		return
	}
	if err := g.CloseSubmissionsForRound(g.CurrentRound().ID); err == nil {
		t.Error("Should have been able to close a voting round")
		return
	}
}

func TestCompleteCurrentRoundBeforeVoting(t *testing.T) {
	g := NewGame(testableWordGenerator(), ModeNormal, 3, 600)
	g.AddPlayer(GeneratePlayerID(), "Tester")
	if err := g.CompleteRound(RoundID("")); err == nil {
		t.Error("Can't complete a round on an inactive game")
		return
	}
	g.StartGame()
	if err := g.CompleteRound(RoundID("Bad Round ID")); err == nil {
		t.Error("Bad Round ID")
		return
	}
	if err := g.CompleteRound(g.CurrentRound().ID); err != nil {
		t.Error("Should have been able to close a new round")
		return
	}
	if g.Rounds[0].State != RoundStateComplete {
		t.Error("Round should be completed now")
		return
	}
	if len(g.Rounds) != 2 {
		t.Error("Should have started a new round")
		return
	}
}

func TestCompleteCurrentRoundAfterVoting(t *testing.T) {
	g := NewGame(testableWordGenerator(), ModeNormal, 3, 600)
	g.AddPlayer(GeneratePlayerID(), "Tester")
	if err := g.CompleteRound(RoundID("")); err == nil {
		t.Error("Can't complete a round on an inactive game")
		return
	}
	g.StartGame()
	rid := g.CurrentRound().ID
	// add some votes in here for scoring
	if err := g.CloseSubmissionsForRound(rid); err != nil {
		t.Error("Should have been able to close this!")
		return
	}
	if err := g.CompleteRound(rid); err != nil {
		t.Error("Should have been able to close a new round")
		return
	}
	if g.Rounds[0].State != RoundStateComplete {
		t.Error("Round should be completed now")
		return
	}
}

func TestCompleteCurrentRoundEndsGame(t *testing.T) {
	g := NewGame(testableWordGenerator(), ModeNormal, 2, 600)
	g.AddPlayer(GeneratePlayerID(), "Tester")
	g.StartGame()
	if len(g.Rounds) != 1 {
		t.Error("Expecting 1 round")
	}
	rid := g.CurrentRound().ID
	g.CompleteRound(rid)
	if len(g.Rounds) != 2 {
		t.Error("Expecting 2 rounds")
	}
	rid = g.CurrentRound().ID
	g.CompleteRound(rid)
	if len(g.Rounds) != 2 {
		t.Error("Expecting 2 rounds")
	}
	if g.State != StateComplete {
		t.Error("Expecting game to be complete")
	}

}

func TestAddingPlayers(t *testing.T) {
	p1id := GeneratePlayerID()
	p2id := GeneratePlayerID()
	g := NewGame(testableWordGenerator(), ModeNormal, 2, 600)
	if err := g.AddPlayer(p1id, "Tester 1"); err != nil {
		t.Error("Adding player should succeed here")
		return
	}
	if err := g.AddPlayer(p1id, "Tester 2"); err == nil {
		t.Error("Adding player should fail here")
		return
	}
	if err := g.AddPlayer(p2id, "Tester 1"); err == nil {
		t.Error("Adding player should fail here")
		return
	}
	if err := g.AddPlayer(p2id, "Tester 2"); err != nil {
		t.Error("Adding player should succeed here")
		return
	}
	if _, p, err := g.findPlayer(p1id); err != nil || p.Name != "Tester 1" || p.ID != p1id {
		t.Error("Can't find Tester 1")
	}
	if _, p, err := g.findPlayer(p2id); err != nil || p.Name != "Tester 2" || p.ID != p2id {
		t.Error("Can't find Tester 2")
	}
	if _, p, err := g.findPlayer("bad user"); err == nil || p != nil {
		t.Error("Shouldn't find non-existant user")
	}

}

func TestRemovePlayer(t *testing.T) {
	p1id := GeneratePlayerID()
	g := NewGame(testableWordGenerator(), ModeNormal, 2, 600)
	if err := g.AddPlayer(p1id, "Tester 1"); err != nil {
		t.Error("Adding player should succeed here")
		return
	}
	_, p, _ := g.findPlayer(p1id)
	if p.Deleted {
		t.Error("Player shouldn't be deleted yet")
		return
	}
	if err := g.RemovePlayer(p1id); err != nil {
		t.Error("Should be able to remove this")
		return
	}
	if !p.Deleted {
		t.Error("Now he should")
		return
	}
	if err := g.RemovePlayer(p1id); err == nil {
		t.Error("Shouldn't be able to remove this twice")
		return
	}
	if err := g.RemovePlayer("bad user"); err == nil {
		t.Error("Shouldn't be able to non-existant user")
		return
	}
}

func TestSerializeDeserialize(t *testing.T) {
	g := NewGame(testableWordGenerator(), ModeNormal, 2, 600)
	g.AddPlayer(PlayerID("123"), "Bob")
	g.StartGame()
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	g.Serialize(w)
	w.Flush()
	t.Log(b.String())

	var g2 Game
	if err := Unserialize(bufio.NewReader(&b), &g2); err != nil {
		t.Error("Failed to deserialize")
	}

	if g.ID != g2.ID {
		t.Error("Not the same game back")
	}
}

func TestSubmitWord(t *testing.T) {
	id := GeneratePlayerID()
	g := NewGame(testableWordGenerator(), ModeNormal, 2, 600)
	g.AddPlayer(id, "Tester")
	if err := g.SubmitWord(id, RoundID(""), "My Def"); err == nil {
		t.Error("Can't submit words on inactive game")
	}
	g.StartGame()
	rid := g.CurrentRound().ID
	if err := g.SubmitWord(GeneratePlayerID(), rid, "My Def"); err == nil {
		t.Error("Non-existent users can't submit words")
	}
	if err := g.SubmitWord(id, RoundID("Bad ID"), "My Def"); err == nil {
		t.Error("Bad Round ID")
	}
	if err := g.SubmitWord(id, rid, "My Def"); err != nil {
		t.Error("Submit should have worked here")
	}
	if err := g.SubmitWord(id, rid, "My Def"); err == nil {
		t.Error("No double submissions")
	}
	if len(g.CurrentRound().Definitions) != 2 {
		t.Error("Should be two definitions here")
	}
	if g.CurrentRound().Definitions[1].Player != id || g.CurrentRound().Definitions[1].Definition != "My Def" {
		t.Error("This should be my definition")
	}
	g.CloseSubmissionsForRound(rid)
	if err := g.SubmitWord(id, rid, "My Def"); err == nil {
		t.Error("Can't submit while voting")
	}
}

func TestVoting(t *testing.T) {
	id1 := GeneratePlayerID()
	id2 := GeneratePlayerID()
	g := NewGame(testableWordGenerator(), ModeNormal, 2, 600)
	g.AddPlayer(id1, "Tester 1")
	g.AddPlayer(id2, "Tester 2")
	if err := g.Vote(id1, RoundID(""), DefinitionID("")); err == nil {
		t.Error("Can't submit words on inactive game")
	}
	g.StartGame()
	roundID := g.CurrentRound().ID
	g.SubmitWord(id1, roundID, "a")
	g.SubmitWord(id2, roundID, "b")
	if len(g.CurrentRound().Definitions) != 3 {
		t.Error("Not enought definitions")
	}
	if err := g.Vote(id1, roundID, DefinitionID("")); err == nil {
		t.Error("Can't submit words on inactive rounds")
	}
	g.CloseSubmissionsForRound(roundID)

	defCorrect := g.CurrentRound().Definitions[0].ID
	def1 := g.CurrentRound().Definitions[1]
	def2 := g.CurrentRound().Definitions[2]

	t.Log("def1", id1, def1)
	t.Log("def2", id2, def2)

	if err := g.Vote(PlayerID("bogus"), roundID, def1.ID); err == nil {
		t.Error("Non-existend users can't vote")
	}
	if err := g.Vote(id1, roundID, def1.ID); err == nil {
		t.Error("Can't vote for self")
	}
	if err := g.Vote(id1, RoundID("not real"), def2.ID); err == nil {
		t.Error("this shouldn't work.  bad round ID")
	}
	if err := g.Vote(id1, roundID, def2.ID); err != nil {
		t.Error("Vote should have worked")
	}
	if err := g.Vote(id1, roundID, def2.ID); err == nil {
		t.Error("Can't vote twice")
	}
	if err := g.Vote(id1, roundID, defCorrect); err == nil {
		t.Error("Can't vote twice")
	}
	if err := g.Vote(id2, roundID, defCorrect); err != nil {
		t.Error("Should have worked")
	}
	g.CompleteRound(roundID)

	if _, p, err := g.findPlayer(id1); err == nil {
		if p.Score != 0 {
			t.Error("Player 1 score should be zero")
		}
	}

	if _, p, err := g.findPlayer(id2); err == nil {
		if p.Score != 4 {
			t.Error("Player 2 score should be four")
		}
	}

}

func TestEndGame(t *testing.T) {
	id1 := GeneratePlayerID()
	g := NewGame(testableWordGenerator(), ModeNormal, 2, 600)
	g.AddPlayer(id1, "Tester 1")
	g.StartGame()
	rid := g.CurrentRound().ID
	g.SubmitWord(id1, rid, "a")
	g.CloseSubmissionsForRound(rid)
	defCorrect := g.CurrentRound().Definitions[0].ID
	g.Vote(id1, rid, defCorrect)
	if err := g.EndGame(); err != nil {
		t.Error("This should work")
	}
	if _, p, err := g.findPlayer(id1); err == nil {
		if p.Score != 3 {
			t.Error("Player 1 score should be three")
		}
	}
	if err := g.EndGame(); err == nil {
		t.Error("This should not work")
	}
}
