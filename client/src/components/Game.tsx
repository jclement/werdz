import React, { Component } from 'react';
import Websocket from 'react-websocket';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Axios from 'axios';
import { Roster } from './Roster';
import { GameStartButton } from './GameStartButton';
import { GameHeader } from './GameHeader';
import { GameSubmitForm } from './GameSubmitForm';
import { GameVotingForm } from './GameVotingForm';
import { GameRoundSummary } from './GameRoundSummary';
import { GameScoreBoard } from './GameScoreBoard';
import { GameTimer } from './GameTimer';

interface GameProps {
  gameId: string,
  playerId: string,
  playerName: string,
}

export class Game extends Component<GameProps, any> {

  constructor(props: GameProps) {
    super(props)

    this.state = {
      gameState: null,
      definition: "",
    }

    this.onMessage = this.onMessage.bind(this)
    this.startGame = this.startGame.bind(this)
    this.submit = this.submit.bind(this)
    this.vote = this.vote.bind(this)
    this.onDefChange = this.onDefChange.bind(this)
  }

  onMessage(msg: any) {
    msg = JSON.parse(msg)
    this.setState({
      gameState: msg
    })
  }

  startGame() {
    Axios.post("/api/game/" + this.props.gameId + "/start", {
    }).then(() => {
    })
  }

  submit(evt: any) {
    Axios.post("/api/game/" + this.props.gameId + "/submit", {
      playerId: this.props.playerId,
      roundId: this.state.gameState.roundId,
      definition: this.state.definition,
    }).then(() => {
      this.setState({
        definition: ""
      })
    })
    evt.preventDefault();
  }

  vote(definitionId: string) {
    Axios.post("/api/game/" + this.props.gameId + "/vote", {
      playerId: this.props.playerId,
      roundId: this.state.gameState.roundId,
      definitionId: definitionId,
    }).then(() => {
    })
  }

  onDefChange(evt: any) {
    this.setState({
      definition: evt.target.value
    })
  }

  render() {
    // URL for web socket
    var loc = window.location, ws_uri;
    if (loc.protocol === "https:") {
      ws_uri = "wss:";
    } else {
      ws_uri = "ws:";
    }
    ws_uri += "//" + loc.host;
    ws_uri += "/api/game/" + this.props.gameId + "/ws?name=" + encodeURIComponent(this.props.playerName) + "&playerid=" + encodeURIComponent(this.props.playerId);

    if (!this.props.playerId || !this.props.playerName) {
      return null
    }

    return (
      <div>
        <Websocket url={ws_uri} onMessage={this.onMessage} />

        {this.state.gameState &&
          <div>
            <Row>
              <Col>
                <GameHeader gameId={this.props.gameId} gameState={this.state.gameState} />
                {this.state.gameState.canStart && <GameStartButton gameId={this.props.gameId} />}
                {this.state.gameState.canSubmit && <GameSubmitForm gameId={this.props.gameId} playerId={this.props.playerId} roundId={this.state.gameState.roundId} />}
                {this.state.gameState.canVote && <GameVotingForm gameId={this.props.gameId} playerId={this.props.playerId} roundId={this.state.gameState.roundId} definitions={this.state.gameState.definitions} />}
                {this.state.gameState.roundState === 2 && <GameRoundSummary playerId={this.props.playerId} gameState={this.state.gameState} />}
                {this.state.gameState.state === 2 && <GameScoreBoard playerId={this.props.playerId} gameState={this.state.gameState} />}
              </Col>
              <Col>
                <Roster players={this.state.gameState.players} playerId={this.props.playerId} />
              </Col>
            </Row>
            <GameTimer remaining={this.state.gameState.remainingTime} total={this.state.gameState.totalTime} />
          </div>
        }

        <hr />
        <pre>{JSON.stringify(this.state.gameState, null, 2)}</pre>


      </div>
    );
  }
}