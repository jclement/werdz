import React, {Component} from 'react';

interface RosterProps {
    playerId: string
    players: any
}

export class Roster extends Component<RosterProps, any> {
    render() {
        return (
            <table className="table table-striped">
                <thead>
                    <tr><th>Player</th><th>Score</th></tr>
                </thead>
                <tbody>
                {this.props.players.map((p: any) => {
                return (<tr key={p.id}>
                    <td>{this.props.playerId === p.id ? <b>{p.name}</b> : p.name}</td>
                    <td>{p.score}</td>
                </tr>);
                })}
                </tbody>
            </table>
        );
    }
}