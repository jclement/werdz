import React, { Component } from 'react';

interface GameRulesProps {
    mode: number,
    rounds: number,
}

export class GameRules extends Component<GameRulesProps, any> {
    render() {
        return (
            <div>
                <div className="card">
                    <div className="card-header">
                        <h4>Rules for this game</h4>
                    </div>
                    <div className="card-body">
                    </div>
                </div>

            </div>
        );
    }
}