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
                        Rules for this game
                    </div>
                    <div className="card-body">
                        {this.props.mode === 0 && <ul>
                            <li>This game will consist of <b>{this.props.rounds}</b> rounds.</li>
                            <li>In each round, you will be presented with an unusual <b>but real</b> word.</li>
                            <li>You will have to make up a definition that you think your friends will believe is the correct definition (hint: the real definitions are usually fairly short such as 'room for conversation').</li>
                            <li>Once everybody has submitted their best guess for what the word means, you'll be presented with a shuffled list of everybody's definitions and the correct definition.</li>
                            <li>You get one point for each of your friends who mistakenly pick your definition as the correct definition.</li>
                            <li>You get three points if you can correctly identify the real definition.</li>
                            <li>At the end of {this.props.rounds} rounds, the player with the most points is the winner.</li>
                            <li>That's it.  Have fun. </li>
                        </ul>}
                        {this.props.mode === 1 && <ul>
                            <li>This game will consist of <b>{this.props.rounds}</b> rounds.</li>
                            <li>In each round, you will be presented with a crazy random word.</li>
                            <li>You will have to make up an outlandish definition for the crazy random word.  Something that's so hilarious and awesome that your friends have to vote for it!</li>
                            <li>You get one point for each vote your crazy definition gets.</li>
                            <li>At the end of {this.props.rounds} rounds, the player with the most points is the winner.</li>
                            <li>That's it.  Have fun. </li>
                        </ul>}
                    </div>
                </div>

            </div>
        );
    }
}