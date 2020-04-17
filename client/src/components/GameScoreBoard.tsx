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
                                            <b style={!d.player ? {color: "green"}: {}}>{d.definition}</b>
                                            {d.player &&  <span>&nbsp;(by {players[d.player]})</span>}
                                            <ul>
                                                {(d.votes || []).map((v) => (
                                                    <li key={v}>Voted for by <i>{players[v]}</i></li>
                                                ))}
                                            </ul>
                                            <br />
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