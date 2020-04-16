import React, { Component } from 'react';
import Alert from 'react-bootstrap/Alert';

interface GameHeaderProps {
    gameState: any,
    gameId: string,
}

export class GameHeader extends Component<GameHeaderProps, any> {
    render() {
        return (
            <div>
                <h1>Game <b>{this.props.gameId}</b></h1>
                {this.props.gameState.state === 1 &&
                    <div>
                        <p><b>Current Word: </b> {this.props.gameState.word}</p>
                        <p><b>Round: </b> {this.props.gameState.round}</p>
                    </div>

                }
                {this.props.gameState.state === 0 && !this.props.gameState.canStart &&
                    <Alert variant="secondary">Not Started.  Waiting for Players.</Alert>
                }
                {this.props.gameState.state === 0 && this.props.gameState.canStart &&
                    <Alert variant="secondary">Ready to start!</Alert>
                }
                {this.props.gameState.canSubmit &&
                    <Alert variant="primary">Submit your definition for <b>{this.props.gameState.word}</b></Alert>
                }
                {this.props.gameState.state === 1 && !this.props.gameState.canSubmit && this.props.gameState.roundState === 0 && 
                    <Alert variant="secondary">Waiting for other players to submit their definitions...</Alert>
                }
                {this.props.gameState.canVote &&
                    <Alert variant="danger">Select the best definition for <b>{this.props.gameState.word}</b></Alert>
                }
                {this.props.gameState.state === 1 && !this.props.gameState.canVote && this.props.gameState.roundState === 1 && 
                    <Alert variant="secondary">Waiting for other players to vote...</Alert>
                }
                {this.props.gameState.state === 2 &&
                    <Alert variant="warning">Game over man!</Alert>
                }

            </div>
        );
    }
}