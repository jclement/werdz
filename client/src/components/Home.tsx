import React, { Component } from 'react';
import Button from 'react-bootstrap/Button';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import { useHistory } from "react-router-dom";
import axios from 'axios';

function HomeButton(props: { rounds: number }) {
    let history = useHistory();
    return (<Button onClick={() => {
        axios.post('/api/game/new', {
            rounds: props.rounds
        }).then((resp: any) => {
            history.push('/game/' + resp.data.id);
        }).catch((err) => {
            console.log(err)
        })
    }} variant="primary">New Game</Button>);
};

function JoinButton(props: { code: string }) {
    let history = useHistory();
    return (<Button onClick={() => {
        history.push('/game/' + props.code);
    }} variant="primary">Join Game</Button>);
};

interface HomeProps {
    playerId: string
    playerName: string
}

interface HomeState {
    rounds: number,
    gameCode: string,
}

export class Home extends Component<HomeProps, HomeState> {

    constructor(props: HomeProps) {
        super(props)
        this.state = {
            rounds: 7,
            gameCode: "",
        }
        this.setRounds = this.setRounds.bind(this)
        this.setGameCode = this.setGameCode.bind(this)
    }

    setRounds(evt: any) {
        this.setState({
            rounds: Number(evt.target.value)
        })
    }

    setGameCode(evt: any) {
        this.setState({
            gameCode: evt.target.value
        })
    }

    render() {
        if (this.props.playerId) {
            return (
                <div>
                    <Row>
                        <Col>
                            <div style={{ display: 'flex', justifyContent: 'center' }}><img className="img-fluid" src="title.png" alt="title" /></div>
                            <h1>Hello, {this.props.playerName || "Stranger!"}</h1>
                            <br />
                        </Col>
                    </Row>
                    <Row>
                        <Col>
                            <div className="card">
                                <div className="card-header">
                                    <h4>Start a New Game</h4>
                                </div>
                                <div className="card-body">
                                    <Form>
                                        <Form.Group>
                                            <Form.Label>Number of Rounds</Form.Label>
                                            <Form.Control as="select" value={this.state.rounds} onChange={this.setRounds}>
                                                <option value={3}>Short (3 rounds)</option>
                                                <option value={7}>ISO 9660 Standard Game (7 rounds)</option>
                                                <option value={11}>Die Hard (11 rounds)</option>
                                            </Form.Control>
                                        </Form.Group>
                                    </Form>
                                    <HomeButton rounds={this.state.rounds} />
                                </div>
                            </div>
                        </Col>
                        <Col>
                            <div className="card">
                                <div className="card-header">
                                    <h4>Join a Game</h4>
                                </div>
                                <div className="card-body">
                                    <Form>
                                        <Form.Group>
                                            <Form.Label>Game Code</Form.Label>
                                            <Form.Control type="text" placeholder="i.e. DR27M" value={this.state.gameCode} onChange={this.setGameCode} />
                                        </Form.Group>
                                    </Form>
                                    <JoinButton code={this.state.gameCode} />
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