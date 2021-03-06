import React, { Component } from 'react';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import Axios from 'axios';
import { playClick } from '../Sounds';

interface GameSubmitFormProps {
    gameId: string,
    roundId: string,
    playerId: string,
}

interface GameSubmitFormState {
    definition: string,
}

export class GameSubmitForm extends Component<GameSubmitFormProps, GameSubmitFormState> {

    constructor(props: GameSubmitFormProps) {
        super(props)

        this.state = {
            definition: "",
        }

        this.onSubmit = this.onSubmit.bind(this)
        this.onDefChange = this.onDefChange.bind(this)
    }

    onSubmit(evt: any) {
        playClick()
        Axios.post("/api/game/" + this.props.gameId + "/submit", {
            playerId: this.props.playerId,
            roundId: this.props.roundId,
            definition: this.state.definition,
        }).then(() => {
        })
        evt.preventDefault()
    }

    onDefChange(evt: any) {
        this.setState({
            definition: evt.target.value
        })
    }

    render() {

        return (
            <div>
                <Form onSubmit={this.onSubmit}>
                    <Form.Group controlId="def">
                        <Form.Label>Definition:</Form.Label>
                        <Form.Control value={this.state.definition} onChange={this.onDefChange} type="text" placeholder="Enter definition" />
                    </Form.Group>
                    <Button variant="primary" disabled={!this.state.definition} type="submit">
                        Submit
                    </Button>
                </Form>

            </div>
        );
    }
}