import React, { Component } from 'react';

interface GameRoundSummaryProps {
    playerId: string,
    gameState: any,
}

export class GameRoundSummary extends Component<GameRoundSummaryProps, any> {
    render() {

        return (
            <div>
                Round Summary
            </div>
        );
    }
}