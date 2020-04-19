import React, { Component } from 'react';
import Table from 'react-bootstrap/Table'
import { game } from '../models/game';
import { FaTrophy } from 'react-icons/fa';
import { playWin, playLose } from '../Sounds';

interface GameScoreBoardProps {
    playerId: string,
    gameState: game,
}

export class GameScoreBoard extends Component<GameScoreBoardProps, any> {

    componentDidMount() {
        if (!this.props.gameState.players || this.props.gameState.players.length === 0) return;
        let firstPlaceScore = this.props.gameState.players[0].score;
        this.props.gameState.players.forEach((p) => {
            if (p.id === this.props.playerId) {
                if (p.score === firstPlaceScore) {
                    playWin()
                } else {
                    playLose()
                }
            }
        })
    }

    render() {

        let players: { [id: string]: string } = {}
        this.props.gameState.players.forEach((p) => { players[p.id] = p.name })

        let place = 0
        let placeScore = 9999999999


        return (
            <div>
                <h2>Score Board</h2>
                <Table striped bordered size="sm">
                    <thead>
                        <tr><th style={{width: "80px"}}></th><th>Player</th><th className="numeric">Score</th></tr>
                    </thead>
                    <tbody>
                        {this.props.gameState.players.map((p: any) => {
                            return (<tr key={p.id} style={p.id === this.props.playerId ? {backgroundColor: "lightblue"}: {}}>
                                <td>
                                    {
                                        (() => {
                                            if (p.score < placeScore) {
                                                placeScore = p.score
                                                place++
                                            }
                                        })()
                                    }
                                    {place ===  1 && <span><FaTrophy />&nbsp;&nbsp;1st</span>}
                                    {place ===  2 && <span><FaTrophy />&nbsp;&nbsp;2nd</span>}
                                    {place ===  3 && <span><FaTrophy />&nbsp;&nbsp;3rd</span>}
                                </td>
                                <td>{p.name}</td>
                                <td className="numeric">{p.score}</td>
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
                                                <b style={{
                                                    color: d.player ? (d.player === this.props.playerId ? "blue" : "black") : "green"
                                                }}>{d.definition}</b> {!d.player && <span> (the correct answer)</span>}
                                                {d.player && <span>&nbsp;(by {players[d.player]})</span>}
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