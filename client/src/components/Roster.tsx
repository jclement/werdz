import React, { Component } from 'react';
import Table from 'react-bootstrap/Table';
import { game } from '../models/game';
import { FaCheckCircle, FaUser, FaGhost } from 'react-icons/fa';

interface RosterProps {
    playerId: string
    gameState: game
}

export class Roster extends Component<RosterProps, any> {
    render() {
        let statusMap: { [id: string]: boolean } = {}
        let inactiveMap: { [id: string]: boolean } = {}

        this.props.gameState.players.forEach((p) => {
            inactiveMap[p.id] = p.inactive
        })

        if (this.props.gameState.state === 1) {
            if (this.props.gameState.currentRound.state === 0) {
                this.props.gameState.players.forEach((p) => {
                    statusMap[p.id] = p.submitted
                })
            }
            if (this.props.gameState.currentRound.state === 1) {
                this.props.gameState.players.forEach((p) => {
                    statusMap[p.id] = p.voted
                })
            }
        }

        return (
            <Table striped bordered size="sm">
                <thead>
                    <tr><th></th><th>Player</th><th className="numeric">Score</th></tr>
                </thead>
                <tbody>
                    {this.props.gameState.players.map((p: any) => {
                        return (<tr key={p.id}>
                            <td style={{ textAlign: "center" }}>
                                {inactiveMap[p.id] ? <FaGhost style={{color: "gray"}} /> : null }
                                {this.props.playerId !== p.id && statusMap[p.id] ? <FaCheckCircle style={{ color: "green" }} /> : null}
                                {this.props.playerId === p.id && !statusMap[p.id] ? <FaUser /> : null}
                                {this.props.playerId === p.id && statusMap[p.id] ? <FaUser style={{ color: "green" }} /> : null}
                            </td>
                            <td>{p.name}</td>
                            <td className="numeric">{p.score}</td>
                        </tr>);
                    })}
                </tbody>
            </Table>
        );
    }
}