package game

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"time"
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
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 3, 600, 90)
	if err := g.StartGame(); err == nil {
		t.Error("Starting a game with no players should fail")
	}
}

func TestStartGameAlreadyStarted(t *testing.T) {
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 3, 600, 90)
	for _, state := range []State{StateActive, StateComplete} {
		g.State = state
		if err := g.StartGame(); err == nil {
			t.Errorf("Game in state '%d' should fail", state)
		}
	}
}

func TestStartGame(t *testing.T) {
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 3, 600, 90)
	if g.State != StateNew {
		t.Errorf("Game should start in New state")
	}
	if g.CanStartGame() {
		t.Errorf("Game should not be startable")
	}
	g.AddPlayer(GeneratePlayerID(), "Tester")
	g.AddPlayer(GeneratePlayerID(), "Tester 2")
	if !g.CanStartGame() {
		t.Errorf("Game should be startable")
	}
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
	if g.CanStartGame() {
		t.Errorf("Game should not be startable")
	}
}

func TestStartGameFun(t *testing.T) {
	g, _ := NewGame(testableWordGenerator(), ModeFun, 3, 600, 90)
	if g.CanStartGame() {
		t.Errorf("Game should not be startable")
	}
	g.AddPlayer(GeneratePlayerID(), "Tester")
	g.AddPlayer(GeneratePlayerID(), "Tester 2")
	if g.CanStartGame() {
		t.Errorf("Game should not be startable")
	}
	g.AddPlayer(GeneratePlayerID(), "Tester 3")
	if !g.CanStartGame() {
		t.Errorf("Game should be startable")
	}
}

func TestCreateNewRoundNormal(t *testing.T) {
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 3, 600, 90)
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

	g, _ := NewGame(testableWordGenerator(), ModeFun, 3, 600, 90)
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
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 3, 600, 90)
	p1id := GeneratePlayerID()
	p2id := GeneratePlayerID()
	g.AddPlayer(p1id, "Tester 1")
	g.AddPlayer(p2id, "Tester 2")
	g.StartGame()
	if err := g.SubmitWord(p1id, g.CurrentRound().ID, "oops"); err != nil {
		t.Error("this should work")
		return
	}
	if err := g.closeSubmissionsForCurrentRound(); err != nil {
		t.Error("Should have been able to close a new round", err)
		return
	}
	if g.Rounds[0].State != RoundStateVoting {
		t.Error("Round should be open for voting now")
		return
	}
	if err := g.closeSubmissionsForCurrentRound(); err == nil {
		t.Error("Should have been able to close a voting round")
		return
	}
}

func TestCloseSubmissionsForCurrentRoundNoVotes(t *testing.T) {
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 3, 600, 90)
	g.AddPlayer(GeneratePlayerID(), "Tester")
	g.AddPlayer(GeneratePlayerID(), "Tester 2")
	g.StartGame()
	if err := g.closeSubmissionsForCurrentRound(); err != nil {
		t.Error("Should have been able to close a new round")
		return
	}
	if g.Rounds[0].State != RoundStateComplete {
		t.Error("Round should be closed")
		return
	}
	if len(g.Rounds) != 2 {
		t.Error("Should be on round 2 now")
		return
	}
}

func TestCompleteCurrentRoundBeforeVoting(t *testing.T) {
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 3, 600, 90)
	g.AddPlayer(GeneratePlayerID(), "Tester")
	g.AddPlayer(GeneratePlayerID(), "Tester 2")
	g.StartGame()
	if err := g.completeCurrentRound(); err != nil {
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
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 3, 600, 90)
	g.AddPlayer(GeneratePlayerID(), "Tester")
	g.AddPlayer(GeneratePlayerID(), "Tester 2")
	g.StartGame()
	// add some votes in here for scoring
	if err := g.closeSubmissionsForCurrentRound(); err != nil {
		t.Error("Should have been able to close this!")
		return
	}
	if err := g.completeCurrentRound(); err != nil {
		t.Error("Should have been able to close a new round")
		return
	}
	if g.Rounds[0].State != RoundStateComplete {
		t.Error("Round should be completed now")
		return
	}
}

func TestCompleteCurrentRoundEndsGame(t *testing.T) {
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 2, 600, 90)
	g.AddPlayer(GeneratePlayerID(), "Tester")
	g.AddPlayer(GeneratePlayerID(), "Tester 2")
	g.StartGame()
	if len(g.Rounds) != 1 {
		t.Error("Expecting 1 round")
	}
	g.completeCurrentRound()
	if len(g.Rounds) != 2 {
		t.Error("Expecting 2 rounds")
	}
	g.completeCurrentRound()
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
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 2, 600, 90)
	if g.PlayerExists(p1id) {
		t.Error("Player shouldn't exist here")
	}
	if err := g.AddPlayer(p1id, "Tester 1"); err != nil {
		t.Error("Adding player should succeed here")
		return
	}
	if !g.PlayerExists(p1id) {
		t.Error("Player shouldn exist here")
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
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 2, 600, 90)
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

func TestInactivePlayer(t *testing.T) {
	p1id := GeneratePlayerID()
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 2, 600, 90)
	g.AddPlayer(p1id, "Tester 1")
	_, p, _ := g.findPlayer(p1id)
	if _, e := g.IsPlayerInactive(GeneratePlayerID()); e == nil {
		t.Error("IsPlayerInactive should return an error on non-existant players")
		return
	}
	if p.Inactive {
		t.Error("Player shouldn't be inactive yet")
		return
	}
	if b, e := g.IsPlayerInactive(p1id); e == nil && b {
		t.Error("Player shouldn't be inactive yet")
		return
	}
	if err := g.SetPlayerInactive(p1id, true); err != nil {
		t.Error("Should be able to inactivate this")
		return
	}
	if !p.Inactive {
		t.Error("Now he should be inactive")
		return
	}
	if b, e := g.IsPlayerInactive(p1id); e == nil && !b {
		t.Error("Player shouldn be inactive yet")
		return
	}
	if err := g.SetPlayerInactive(p1id, false); err != nil {
		t.Error("Should be able to activate this")
		return
	}
	if p.Inactive {
		t.Error("Now he should be active")
		return
	}
	if err := g.SetPlayerInactive("bad user", false); err == nil {
		t.Error("Shouldn't be able to non-existant user")
		return
	}
}

func TestSerializeDeserialize(t *testing.T) {
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 2, 600, 90)
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
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 2, 600, 90)
	g.AddPlayer(id, "Tester")
	g.AddPlayer(GeneratePlayerID(), "Tester 2")
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
	g.closeSubmissionsForCurrentRound()
	if err := g.SubmitWord(id, rid, "My Def"); err == nil {
		t.Error("Can't submit while voting")
	}
}

func TestVoting(t *testing.T) {
	id1 := GeneratePlayerID()
	id2 := GeneratePlayerID()
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 2, 600, 90)
	g.AddPlayer(id1, "Tester 1")
	g.AddPlayer(id2, "Tester 2")
	if err := g.Vote(id1, RoundID(""), DefinitionID("")); err == nil {
		t.Error("Can't submit words on inactive game")
	}
	g.StartGame()
	roundID := g.CurrentRound().ID
	g.SubmitWord(id1, roundID, "a")
	if err := g.Vote(id1, roundID, g.CurrentRound().Definitions[0].ID); err == nil {
		t.Error("Can't vote on round before voting starts")
	}
	g.SubmitWord(id2, roundID, "b")
	if len(g.CurrentRound().Definitions) != 3 {
		t.Error("Not enought definitions")
	}

	var defCorrect *Definition
	var def1 *Definition
	var def2 *Definition
	for _, d := range g.CurrentRound().Definitions {
		if d.Player == rightAnswerPlayerID {
			defCorrect = d
		}
		if d.Player == id1 {
			def1 = d
		}
		if d.Player == id2 {
			def2 = d
		}
	}

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
	if err := g.Vote(id1, roundID, defCorrect.ID); err == nil {
		t.Error("Can't vote twice")
	}
	if err := g.Vote(id2, roundID, defCorrect.ID); err != nil {
		t.Error("Should have worked")
	}
	g.completeCurrentRound()

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
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 2, 600, 90)
	g.AddPlayer(id1, "Tester 1")
	g.AddPlayer(GeneratePlayerID(), "Tester 2")
	g.StartGame()
	rid := g.CurrentRound().ID
	g.SubmitWord(id1, rid, "a")
	g.closeSubmissionsForCurrentRound()
	var defCorrect DefinitionID
	for _, d := range g.CurrentRound().Definitions {
		if d.Player == rightAnswerPlayerID {
			defCorrect = d.ID
		}
	}
	g.Vote(id1, rid, defCorrect)
	g.closeVotingForCurrentRound()
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

func TestNewGames(t *testing.T) {
	sets := [](struct {
		numRounds          int
		submissionDuration int
		votingDuration     int
		shouldFail         bool
	}){
		{0, 0, 0, true},
		{MinimumRounds, MinimumSubmissionDuration, MinimumVotingDuration, false},
		{MaximumRounds, MaximumSubmissionDuration, MaximumVotingDuration, false},
		{MinimumRounds - 1, MinimumSubmissionDuration, MinimumVotingDuration, true},
		{MinimumRounds, MinimumSubmissionDuration - 1, MinimumVotingDuration, true},
		{MinimumRounds, MinimumSubmissionDuration, MinimumVotingDuration - 1, true},
		{MaximumRounds + 1, MaximumSubmissionDuration, MaximumVotingDuration, true},
		{MaximumRounds, MaximumSubmissionDuration + 1, MaximumVotingDuration, true},
		{MaximumRounds, MaximumSubmissionDuration, MaximumVotingDuration + 1, true},
	}
	for _, tst := range sets {
		if _, err := NewGame(testableWordGenerator(), ModeNormal, tst.numRounds, tst.submissionDuration, tst.votingDuration); (err != nil) != tst.shouldFail {
			t.Errorf("Expected error=%t where rounds=%d, submission=%d, votes=%d", tst.shouldFail, tst.numRounds, tst.submissionDuration, tst.votingDuration)
		}

	}
}

func TestTick(t *testing.T) {
	id1 := GeneratePlayerID()
	g, _ := NewGame(testableWordGenerator(), ModeNormal, 2, 600, 90)
	g.AddPlayer(id1, "tester")
	g.AddPlayer(GeneratePlayerID(), "tester 2")
	if g.Tick() {
		t.Error("nothing")
		return
	}
	g.StartGame()
	if g.Tick() {
		t.Error("nothing")
		return
	}
	g.SubmitWord(id1, g.CurrentRound().ID, "test")
	g.CurrentRound().RoundStartTime = time.Now().Add(time.Duration(-200 * time.Second))
	if g.Tick() {
		t.Error("nothing")
		return
	}
	g.CurrentRound().RoundStartTime = time.Now().Add(time.Duration(-800 * time.Second))
	if !g.Tick() {
		t.Error("something")
		return
	}
	for _, d := range g.CurrentRound().Definitions {
		if d.Player == rightAnswerPlayerID {
			g.Vote(id1, g.CurrentRound().ID, d.ID)
		}
	}
	if _, p, err := g.findPlayer(id1); err != nil || p.Score != 0 {
		t.Error("unexpected player score")
		return
	}
	g.CurrentRound().VotingStartTime = time.Now().Add(-1 * time.Duration(2+g.CurrentRound().VotingDuration) *time.Second)
	if !g.Tick() {
		t.Error("something")
		return
	}
	if g.CurrentRound().State != RoundStateVotingComplete {
		t.Error("unexpected round state")
		return
	}
	if _, p, err := g.findPlayer(id1); err != nil || p.Score != 3 {
		t.Errorf("unexpected player score: %d", p.Score)
		return
	}
	g.CurrentRound().VotingCompleteStartTime = time.Now().Add(-1 * time.Duration(2+g.CurrentRound().VotingCompleteDuration) * time.Second)
	g.Tick()

	g.CurrentRound().RoundStartTime = time.Now().Add(time.Duration(-800 * time.Second))
	g.Tick()
	g.CurrentRound().VotingStartTime = time.Now().Add(time.Duration(-91 * time.Second))
	g.Tick()
	g.CurrentRound().VotingCompleteStartTime = time.Now().Add(-1 * time.Duration(2+g.VotingCompleteDuration) * time.Second)
	g.Tick()
	if g.State != StateComplete {
		t.Error("invalid game state")
		return
	}

}
