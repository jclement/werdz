import React, { Component } from 'react';
import Alert from 'react-bootstrap/Alert';
import { game } from '../models/game';
import { GameTimer } from './GameTimer';

interface GameHeaderProps {
    gameState: game,
    gameId: string,
}

export class GameHeader extends Component<GameHeaderProps, any> {
    render() {

        return (
            <div>
                <h2>Game <b>{this.props.gameId}</b></h2>
                {this.props.gameState.state === 1 &&
                    <table style={{width: "100%", marginBottom: "20px"}}>
                        <tbody>
                        <tr><th>Current&nbsp;Word:&nbsp;</th><td style={{width: "100%"}}>{this.props.gameState.currentRound.word}</td></tr>
                        <tr><th>Round: </th><td>{this.props.gameState.currentRound.num} of {this.props.gameState.totalRounds}</td></tr>
                        <tr><td colSpan={2}><GameTimer remaining={this.props.gameState.remainingTime} total={this.props.gameState.totalTime} /></td></tr>
                        </tbody>
                    </table>

                }
                {this.props.gameState.state === 0 && !this.props.gameState.canStart &&
                    <Alert variant="secondary">Not Started.  Waiting for Players.</Alert>
                }
                {this.props.gameState.state === 0 && this.props.gameState.canStart &&
                    <Alert variant="secondary">Ready to start!</Alert>
                }
                {this.props.gameState.canSubmit &&
                    <Alert variant="primary">Submit your definition for <b>{this.props.gameState.currentRound.word}</b></Alert>
                }
                {this.props.gameState.state === 1 && !this.props.gameState.canSubmit && this.props.gameState.currentRound.state === 0 && 
                    <Alert variant="secondary">Waiting for other players to submit their definitions...</Alert>
                }
                {this.props.gameState.canVote && this.props.gameState.mode === 0 &&
                    <Alert variant="danger">Select the correct definition for <b>{this.props.gameState.currentRound.word}</b></Alert>
                }
                {this.props.gameState.canVote && this.props.gameState.mode === 1 &&
                    <Alert variant="danger">Select the best definition for <b>{this.props.gameState.currentRound.word}</b></Alert>
                }
                {this.props.gameState.state === 1 && !this.props.gameState.canVote && this.props.gameState.currentRound.state === 1 && 
                    <Alert variant="secondary">Waiting for other players to vote...</Alert>
                }
                {this.props.gameState.state === 2 &&
                    <Alert variant="warning">Game over man!</Alert>
                }

            </div>
        );
    }
}