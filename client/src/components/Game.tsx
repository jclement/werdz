import React, { Component } from 'react';
import Websocket from 'react-websocket';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import Axios from 'axios';

interface GameProps {
  gameId: string,
  userId: string,
  name: string,
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
    this.onDefChange = this.onDefChange.bind(this)
  }

  onMessage(msg: any) {
    msg = JSON.parse(msg)
    console.log(msg)
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
      playerid: this.props.userId,
      roundid: this.state.gameState.roundId,
      definition: this.state.definition,
    }).then(() => {
    })
    evt.preventDefault();
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
    ws_uri += "/api/game/" + this.props.gameId + "/ws?name=" + encodeURIComponent(this.props.name) + "&playerid=" + encodeURIComponent(this.props.userId);

    return (
      <div>
        <Websocket url={ws_uri} onMessage={this.onMessage} />
        <p>Game : {this.props.gameId}</p>
        <hr />
        <pre>{JSON.stringify(this.state.gameState, null, 2)}</pre>

        { this.state.gameState && this.state.gameState.state === 0 && <Button onClick={this.startGame}>Start Game</Button>}

        { this.state.gameState && this.state.gameState.state === 1 && (
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


      </div>
    );
  }
}