import React, { Component } from 'react';
import Websocket from 'react-websocket';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import Axios from 'axios';

interface GameProps {
  gameId: string,
  playerId: string,
  playerName: string,
}

export class Game extends Component<GameProps, any> {
  ws: WebSocket | null = null;

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

  vote(definitionId : string) {
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

    return (
      <div>
        <Websocket url={ws_uri} onMessage={this.onMessage} />
        <p>Game : {this.props.gameId}</p>
        <hr />
        <pre>{JSON.stringify(this.state.gameState, null, 2)}</pre>

        {this.state.gameState && this.state.gameState.state === 0 && <Button onClick={this.startGame}>Start Game</Button>}

        {this.state.gameState && this.state.gameState.canSubmit && (
          <Form onSubmit={this.submit}>
            <Form.Group controlId="def">
              <Form.Label>Definition</Form.Label>
              <Form.Control value={this.state.definition} onChange={this.onDefChange} type="text" placeholder="Enter definition" />
            </Form.Group>
            <Button variant="primary" type="submit">
              Submit
          </Button>
          </Form>
        )}

        {this.state.gameState && this.state.gameState.canVote && (
          <div>
            {this.state.gameState.definitions.map((def: any) => {
              return (<Button key={def.id} onClick={() => {this.vote(def.id);}}>
                {def.definition}
              </Button>);
            })}
          </div>
        )}


      </div>
    );
  }
}