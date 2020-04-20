import React, { Component } from 'react';
import Websocket from 'react-websocket';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { Roster } from './Roster';
import { GameStartButton } from './GameStartButton';
import { GameHeader } from './GameHeader';
import { GameSubmitForm } from './GameSubmitForm';
import { GameVotingForm } from './GameVotingForm';
import { GameRoundSummary } from './GameRoundSummary';
import { GameScoreBoard } from './GameScoreBoard';
import { game } from '../models/game';
import { Alert, Spinner } from 'react-bootstrap';
import Axios from 'axios';
import { GameRules } from './GameRules';

interface GameProps {
  gameId: string,
  playerId: string,
  playerName: string,
}

interface GameState {
  gameState: game | null,
  error: boolean,
}

export class Game extends Component<GameProps, GameState> {
  intervalId: NodeJS.Timeout | null = null;
  ws: JSX.Element | null = null;

  constructor(props: GameProps) {
    super(props)

    this.state = {
      gameState: null,
      error: false
    }

    this.timer = this.timer.bind(this)
    this.onMessage = this.onMessage.bind(this)
    this.onSocketError = this.onSocketError.bind(this)
    this.onSocketConnect = this.onSocketConnect.bind(this)
  }

  onSocketError() {
    this.setState({
      error: true
    })
  }

  onSocketConnect() {
    this.setState({
      error: false
    })
  }

  componentDidMount() {
    this.intervalId = setInterval(this.timer, 10000);
  }

  componentWillUnmount() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }

  timer() {
    Axios.post("/api/game/" + this.props.gameId + "/ping", {
      playerId: this.props.playerId
    }).catch((e) => {} )

  }

  onMessage(msg: any) {
    let gameState: game | null
    gameState = JSON.parse(msg)
    this.setState({
      gameState: gameState
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

    if (this.state.error) {
      return <Alert variant="danger">Error connecting to server.  Try reloading the page.</Alert>;
    }

    return (
      <div>
        <br />

        {!this.state.gameState && <Spinner animation="border" />}

        <Websocket url={ws_uri} onMessage={this.onMessage} onConnect={this.onSocketConnect} onError={this.onSocketError} />

        {this.state.gameState && (this.state.gameState.state === 0 || this.state.gameState.state === 1) &&
          <div>
            <Row>
              <Col md={8} style={{paddingBottom: "15px"}}>
                <GameHeader gameId={this.props.gameId} gameState={this.state.gameState} />
                {this.state.gameState.canSubmit && <GameSubmitForm gameId={this.props.gameId} playerId={this.props.playerId} roundId={this.state.gameState.currentRound.id} />}
                {this.state.gameState.canVote && <GameVotingForm gameId={this.props.gameId} playerId={this.props.playerId} roundId={this.state.gameState.currentRound.id} definitions={this.state.gameState.currentRound.definitions} />}
                {this.state.gameState.currentRound && this.state.gameState.currentRound.state === 2 && <GameRoundSummary playerId={this.props.playerId} gameState={this.state.gameState} />}
                {this.state.gameState.canStart && <div style={{paddingBottom: "15px"}}><GameStartButton gameId={this.props.gameId} /></div>}
                {this.state.gameState.state === 0 && <GameRules mode={this.state.gameState.mode} rounds={this.state.gameState.totalRounds} />}
              </Col>
              <Col md={4}>
                <Roster gameState={this.state.gameState} playerId={this.props.playerId} />
              </Col>
            </Row>
          </div>
        }

        {this.state.gameState && this.state.gameState.state === 2 &&
          <div>
            <Row>
              <Col>
                <GameHeader gameId={this.props.gameId} gameState={this.state.gameState} />
                <GameScoreBoard playerId={this.props.playerId} gameState={this.state.gameState} />
              </Col>
            </Row>
          </div>
        }

      </div>
    );
  }
}