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
import { GameTimer } from './GameTimer';
import { game } from '../models/game';
import { Alert } from 'react-bootstrap';

interface GameProps {
  gameId: string,
  playerId: string,
  playerName: string,
}

interface GameState{
  gameState : game | null
  
}

export class Game extends Component<GameProps, GameState> {

  constructor(props: GameProps) {
    super(props)

    this.state = {
      gameState: null,
    }

    this.onMessage = this.onMessage.bind(this)
  }

  onMessage(msg: any) {
    let gameState : game | null 
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

    return (
      <div>
        <br />
        <Websocket url={ws_uri} onMessage={this.onMessage}  />

        {!this.state.gameState && <Alert variant="secondary">Loading or a bad room code.  Who knows!</Alert>}

        {this.state.gameState && (this.state.gameState.state === 0 || this.state.gameState.state === 1) &&
          <div>
            <Row>
              <Col sm={8}>
                <GameHeader gameId={this.props.gameId} gameState={this.state.gameState} />
                {this.state.gameState.canStart && <GameStartButton gameId={this.props.gameId} />}
                {this.state.gameState.canSubmit && <GameSubmitForm gameId={this.props.gameId} playerId={this.props.playerId} roundId={this.state.gameState.currentRound.id} />}
                {this.state.gameState.canVote && <GameVotingForm gameId={this.props.gameId} playerId={this.props.playerId} roundId={this.state.gameState.currentRound.id} definitions={this.state.gameState.currentRound.definitions} />}
                {this.state.gameState.currentRound && this.state.gameState.currentRound.state === 2 && <GameRoundSummary playerId={this.props.playerId} gameState={this.state.gameState} />}
                <GameTimer remaining={this.state.gameState.remainingTime} total={this.state.gameState.totalTime} />
                <br />
              </Col>
              <Col sm={4}>
                <Roster players={this.state.gameState.players} playerId={this.props.playerId} />
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