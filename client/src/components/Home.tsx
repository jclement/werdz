import React, { Component } from 'react';
import Button from 'react-bootstrap/Button';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import { useHistory } from "react-router-dom";
import Axios from 'axios';

function HomeButton(props: { mode: number, rounds: number }) {
    let history = useHistory();
    return (<Button onClick={() => {
        Axios.post('/api/game/new', {
            rounds: props.rounds,
            mode: props.mode
        }).then((resp: any) => {
            history.push('/game/' + resp.data.id);
        }).catch((err) => {
            console.log(err)
        })
    }} variant="primary">New Game</Button>);
};

function JoinButton(props: { disabled: boolean, code: string }) {
    let history = useHistory();
    return (<Button disabled={props.disabled} onClick={() => {
        history.push('/game/' + props.code);
    }} variant="primary">Join Game</Button>);
};

interface HomeProps {
    playerId: string
    playerName: string
}

interface HomeState {
    rounds: number,
    mode : number,
    gameCode: string,
    gameCodeExists: boolean,
}

export class Home extends Component<HomeProps, HomeState> {

    constructor(props: HomeProps) {
        super(props)
        this.state = {
            rounds: 7,
            mode: 0,
            gameCode: "",
            gameCodeExists: false,
        }
        this.setRounds = this.setRounds.bind(this)
        this.setMode = this.setMode.bind(this)
        this.setGameCode = this.setGameCode.bind(this)
    }

    setRounds(evt: any) {
        this.setState({
            rounds: Number(evt.target.value)
        })
    }

    setMode(evt: any) {
        this.setState({
            mode: Number(evt.target.value)
        })
    }

    setGameCode(evt: any) {
        var code = evt.target.value.toUpperCase().replace(/[^a-zA-Z0-9]/g, "");
        this.setState({
            gameCode: code
        })
        if (code.length === 5) {
            Axios.get("/api/game/" + code + "/exists", {})
                .then((r: any) => {
                    this.setState({
                        gameCodeExists: r.data,
                    })
                })
        } else {
            this.setState({
                gameCodeExists: false
            })
        }
    }

    render() {
        if (this.props.playerId) {
            return (
                <div>
                    <Row>
                        <Col>
                            <div style={{textAlign: "center"}}><img className="img-fluid" src="title.png" alt="title" /></div>
                            <h1>Hello, {this.props.playerName || "Stranger!"}</h1>
                            <br />
                        </Col>
                    </Row>
                    <Row>
                        <Col sm={6}>
                            <div className="card">
                                <div className="card-header">
                                    <b>Start a New Game</b>
                                </div>
                                <div className="card-body">
                                    <Form>
                                        <Form.Group>
                                            <Form.Label>Number of Rounds:</Form.Label>
                                            <Form.Control as="select" value={this.state.rounds} onChange={this.setRounds}>
                                                <option value={3}>Short (3 rounds)</option>
                                                <option value={7}>ISO 9660 Standard Game (7 rounds)</option>
                                                <option value={11}>Die Hard (11 rounds)</option>
                                            </Form.Control>
                                        </Form.Group>
                                        <Form.Group>
                                            <Form.Label>Game Type:</Form.Label>
                                            <Form.Control as="select" value={this.state.mode} onChange={this.setMode}>
                                                <option value={0}>Real Words (normal)</option>
                                                <option value={1}>Fake Words</option>
                                            </Form.Control>
                                        </Form.Group>
                                    </Form>
                                    <HomeButton mode={this.state.mode} rounds={this.state.rounds} />
                                </div>
                            </div>
                            <br />
                        </Col>
                        <Col sm={6}>
                            <div className="card">
                                <div className="card-header">
                                    <b>Join a Game</b>
                                </div>
                                <div className="card-body">
                                    <Form>
                                        <Form.Group>
                                            <Form.Label>Game Code:</Form.Label>
                                            <Form.Control type="text" isInvalid={!this.state.gameCodeExists && !!this.state.gameCode} maxLength={5} minLength={5} placeholder="i.e. DR27M" value={this.state.gameCode} onChange={this.setGameCode} />
                                            {this.state.gameCode && !this.state.gameCodeExists && this.state.gameCode.length !== 5 && <Form.Control.Feedback type="invalid">Invalid game code</Form.Control.Feedback>}
                                            {this.state.gameCode && !this.state.gameCodeExists && this.state.gameCode.length === 5 && <Form.Control.Feedback type="invalid">Game does not exist</Form.Control.Feedback>}
                                        </Form.Group>
                                    </Form>
                                    <JoinButton disabled={!this.state.gameCodeExists} code={this.state.gameCode} />
                                </div>
                            </div>
                        </Col>

                    </Row>
                </div>
            );
        }
        return <p>Loading</p>
    }
}