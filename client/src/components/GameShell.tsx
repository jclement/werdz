import React, { Component } from 'react';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import { Game } from './Game';

interface GameShellProps {
    gameId: string,
    playerId: string,
    playerName: string,
    onSetName: (name: string) => void,
}

export class GameShell extends Component<GameShellProps, any> {
    constructor(props: GameShellProps) {
        super(props)

        this.state = {
            name: "",
        }

        this.onNameUpdate = this.onNameUpdate.bind(this)
        this.onSetName = this.onSetName.bind(this)
    }

    onNameUpdate(evt: any) {
        this.setState({
            name: Number(evt.target.value)
        })
    }

    onSetName(evt: any) {
        if (this.state.name && this.state.name.length > 1) {
            this.props.onSetName(this.state.name)
        }
        evt.preventDefault()
    }

    render() {
        if (!this.props.playerName) {
            return (
                <div>
                    <h1>Joining game <b>{this.props.gameId}</b></h1>
                    <Form onSubmit={this.onSetName}>
                        <Form.Group controlId="name">
                            <Form.Label>Name</Form.Label>
                            <Form.Control value={this.state.name} onChange={this.onNameUpdate} type="text" placeholder="Who are you?" />
                        </Form.Group>
                        <Button variant="primary" type="submit">
                            Ok
                        </Button>
                    </Form>
                </div>
            );
        } else {
            return (<div><Game playerName={this.props.playerName} playerId={this.props.playerId} gameId={this.props.gameId} /></div>);
        }
    }
}