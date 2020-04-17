import React, {Component} from 'react';
import Table from 'react-bootstrap/Table';

interface RosterProps {
    playerId: string
    players: any
}

export class Roster extends Component<RosterProps, any> {
    render() {
        return (
            <Table striped bordered size="sm">
                <thead>
                    <tr><th>Player</th><th className="numeric">Score</th></tr>
                </thead>
                <tbody>
                {this.props.players.map((p: any) => {
                return (<tr key={p.id}>
                    <td>{this.props.playerId === p.id ? <b>{p.name}</b> : p.name}</td>
                    <td className="numeric">{p.score}</td>
                </tr>);
                })}
                </tbody>
            </Table>
        );
    }
}