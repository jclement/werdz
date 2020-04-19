import React, { Component } from 'react';
import Button from 'react-bootstrap/Button';
import Axios from 'axios';
import { playClick } from '../Sounds';

interface GameStartButtonProps {
    gameId: string,
}

export class GameStartButton extends Component<GameStartButtonProps, any> {

    constructor(props: GameStartButtonProps) {
        super(props)

        this.startGame = this.startGame.bind(this)
    }

    startGame() {
        playClick()
        Axios.post("/api/game/" + this.props.gameId + "/start", {
        }).then(() => {
        })
    }
    render() {
        return (
            <div>
                <h3>Once all players are ready...</h3>
                <Button onClick={this.startGame}>Start Game</Button>
            </div>
        );
    }
}