import React, { Component } from 'react';

interface GameScoreBoardProps {
    playerId: string,
    gameState: any,
}

export class GameScoreBoard extends Component<GameScoreBoardProps, any> {
    render() {

        return (
            <div>
                Game Summary
            </div>
        );
    }
}