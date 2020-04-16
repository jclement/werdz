import React, { Component } from 'react';
import Button from 'react-bootstrap/Button';
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

export class Home extends Component<HomeProps, any> {
    render() {
        if (this.props.playerId) {
            return (
                <div>
                    <h1>Hello, {this.props.playerName}</h1>
                    <HomeButton rounds={3} />
                </div>
            );
        }
        return <p>Loading</p>
    }
}