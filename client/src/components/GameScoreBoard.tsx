import React, { Component } from 'react';
import Table from 'react-bootstrap/Table'
import { game } from '../models/game';

interface GameScoreBoardProps {
    playerId: string,
    gameState: game,
}

export class GameScoreBoard extends Component<GameScoreBoardProps, any> {
    render() {

        let players : { [id: string] : string} = {}
        this.props.gameState.players.forEach((p) => {players[p.id] = p.name})

        return (
            <div>
                <h2>Score Board</h2>
                <Table striped bordered size="sm">
                    <thead>
                        <tr><th>Player</th><th>Score</th></tr>
                    </thead>
                    <tbody>
                        {this.props.gameState.players.map((p: any) => {
                            return (<tr key={p.id}>
                                <td>{this.props.playerId === p.id ? <b>{p.name}</b> : p.name}</td>
                                <td>{p.score}</td>
                            </tr>);
                        })}
                    </tbody>
                </Table>

                <h2>Rounds</h2>

                {this.props.gameState.rounds.map((r) => {

                    return (
                        <div key={r.id}>
                            <div className="card">
                                <div className="card-header">
                                    <h4>Round {r.num} - <b>{r.word}</b></h4>
                                </div>
                                <div className="card-body">
                                    <ul>
                                    {r.definitions.map((d) => (
                                        <li key={d.id}>
                                            {d.definition}  
                                            <ul>
                                                {d.player && <li>Written by {players[d.player]}</li>}
                                                {!d.player && <li>THE CORRECT ANSWER</li>}
                                                {(d.votes || []).map((v) => (
                                                    <li key={v}>Voted for by <i>{players[v]}</i></li>
                                                ))}
                                            </ul>
                                        </li>
                                    ))}
                                    </ul>
                                </div>
                            </div>
                            <br />
                        </div>
                    );

                })}

            </div>
        );
    }
}