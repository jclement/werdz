import React, { Component } from 'react';

export class Rules extends Component {
    render() {
        return (
            <div>
                <h1>Rulez for Werdz</h1>

                <p>It is recommended that you play this while video chatting with the other 3-10 other players.</p>
                <p>What works well for us is sitting in front of a laptop for the video chat, and using phones for the game.</p>

                <ul>
                    <li>Connect with your friends on video (I suggest <a href="https://meet.jit.si">meet.jit.si</a>)</li>
                    <li>One player starts a new game and shares the code with everybody via. the video chat.  That player chooses the number of rounds, and the game mode.
                    </li>
                    <li>Everybody joins the game</li>
                    <li>Some player chooses the game type, number of rounds, and presses the start button</li>
                </ul>


                <div className="card">
                    <div className="card-header">
                        Rules for <b>Real Word</b> (Normal) Mode
                    </div>
                    <div className="card-body">
                        <ul>
                        <li>This game will consist of a number rounds.</li>
                        <li>In each round, you will be presented with an unusual <b>but real</b> word.</li>
                        <li>You will have to make up a definition that you think your friends will believe is the correct definition (hint: the real definitions are usually fairly short such as 'room for conversation').</li>
                        <li>Once everybody has submitted their best guess for what the word means, you'll be presented with a shuffled list of everybody's definitions and the correct definition.</li>
                        <li>You get one point for each of your friends who mistakenly pick your definition as the correct definition.</li>
                        <li>You get three points if you can correctly identify the real definition.</li>
                        <li>At the end of the game, the player with the most points is the winner.</li>
                        <li>That's it.  Have fun. </li>
                        </ul>
                    </div>
                </div>
<br />

                <div className="card">
                    <div className="card-header">
                        Rules for <b>Fake Word</b> Mode 
                    </div>
                    <div className="card-body">
                        <ul>
                        <li>This game will consist of a number rounds.</li>
                        <li>In each round, you will be presented with a crazy random word.</li>
                        <li>You will have to make up an outlandish definition for the crazy random word.  Something that's so hilarious and awesome that your friends have to vote for it!</li>
                        <li>You get one point for each vote your crazy definition gets.</li>
                        <li>At the end of the game, the player with the most points is the winner.</li>
                        <li>That's it.  Have fun. </li>
                        </ul>
                    </div>
                </div>


            </div>
        );
    }
}