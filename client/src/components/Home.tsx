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

interface HomeProps {
    playerId: string
    playerName: string
}

interface HomeState {
    rounds: number
}

export class Home extends Component<HomeProps, HomeState> {

    constructor(props: HomeProps) {
        super(props)
        this.state = {
            rounds:7 
        }
        this.setRounds = this.setRounds.bind(this)
    }

    setRounds(evt: any) {
        this.setState({
            rounds: Number(evt.target.value)
        })
    }

    render() {
        if (this.props.playerId) {
            return (
                <Row>
                    <Col>
                        <div style={{display: 'flex', justifyContent: 'center'}}><img className="img-fluid" src="title.png" alt="title"/></div>
                        <h1>Hello, {this.props.playerName || "Stranger!"}</h1>
                        <br />
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
                    </Col>
                </Row>
            );
        }
        return <p>Loading</p>
    }
}