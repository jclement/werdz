import React, { Component } from 'react';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import { Game } from './Game';
import Axios from 'axios';

interface GameShellProps {
    gameId: string,
    playerId: string,
    playerName: string,
    onSetName: (name: string) => void,
}

interface GameShellState {
    name: string,
    hasJoined: boolean,
    nameIsAvailable: boolean,
    nameIsAvailableRunning: boolean,
    loading: boolean,
}

export class GameShell extends Component<GameShellProps, GameShellState> {
    constructor(props: GameShellProps) {
        super(props)

        this.state = {
            name: props.playerName,
            loading: true,
            hasJoined: false,
            nameIsAvailable: false,
            nameIsAvailableRunning: props.playerName !== "",
        }

        this.onNameUpdate = this.onNameUpdate.bind(this)
        this.onJoin = this.onJoin.bind(this)
    }

    componentDidMount() {
        Axios.post('/api/game/' + this.props.gameId + '/has_player', {
            playerId: this.props.playerId,
        }).then((r) => {
            this.setState({
                loading: false,
                hasJoined: r.data,
            })
        })

        if (this.props.playerName.length > 0) {
            Axios.post('/api/game/' + this.props.gameId + '/name_available', {
                name:  this.props.playerName,
            }).then((r) => {
                this.setState({
                    nameIsAvailable: r.data,
                    nameIsAvailableRunning: false,
                })
            })
        }
    }

    onNameUpdate(evt: any) {
        var name = evt.target.value.substring(0,20);
        this.setState({
            name: name,
            nameIsAvailable: false,
            nameIsAvailableRunning: true,
        })
        if (name.length > 0) {
            Axios.post('/api/game/' + this.props.gameId + '/name_available', {
                name:  name,
            }).then((r) => {
                this.setState({
                    nameIsAvailable: r.data,
                    nameIsAvailableRunning: false,
                })
            })
        }
    }

    onJoin(evt: any) {
        if (this.state.name && this.state.name.length > 0) {
            (new Audio("/sounds/click.mp3")).play()
            this.props.onSetName(this.state.name)
            this.setState({
                hasJoined: true
            })
        }
        evt.preventDefault()
    }

    render() {

        if (this.state.loading) return null;

        if (!this.state.hasJoined) {
            return (
                <div>
                    <h1>Joining game <b>{this.props.gameId}</b></h1>
                    <Form onSubmit={this.onJoin}>
                        <Form.Group controlId="name">
                            <Form.Label>Name</Form.Label>
                            <Form.Control value={this.state.name} onChange={this.onNameUpdate} type="text" placeholder="Who are you?" isInvalid={(!this.state.nameIsAvailable && !this.state.nameIsAvailableRunning) || this.state.name.length === 0} />
                            {!this.state.name && <Form.Control.Feedback type="invalid">A name is required to join a game</Form.Control.Feedback>}
                            {this.state.name && !this.state.nameIsAvailable && !this.state.nameIsAvailableRunning && <Form.Control.Feedback type="invalid">This name is already in use in this game</Form.Control.Feedback>}
                        </Form.Group>
                        <Button disabled={!this.state.nameIsAvailable} variant="primary" type="submit">
                            Join Game 
                        </Button>
                    </Form>
                </div>
            );
        } else {
            return (<div><Game playerName={this.props.playerName} playerId={this.props.playerId} gameId={this.props.gameId} /></div>);
        }
    }
}