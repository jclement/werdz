import React, { Component } from 'react';
import { game } from '../models/game';

interface GameRoundSummaryProps {
    playerId: string,
    gameState: game,
}

export class GameRoundSummary extends Component<GameRoundSummaryProps, any> {
    render() {

        let players: { [id: string]: string } = {}
        this.props.gameState.players.forEach((p) => { players[p.id] = p.name })

        return (
            <div>
                <div className="card">
                    <div className="card-header">
                        Round {this.props.gameState.currentRound.num} Summary - <b>{this.props.gameState.currentRound.word}</b>
                    </div>
                    <div className="card-body">
                        <ul>
                            {this.props.gameState.currentRound.definitions.map((d) => (
                                <li key={d.id}>
                                    <b style={{ 
                                        color: d.player ? (d.player === this.props.playerId ? "blue": "black") : "green" 
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

            </div>
        );
    }
}