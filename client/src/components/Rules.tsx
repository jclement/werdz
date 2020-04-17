import React, { Component } from 'react';

export class Rules extends Component {
    render() {
        return (
            <div>
                <h1>Rulez for Werdz</h1>

                <h3>Game Play:</h3>
                <p>It is recommended that you play this while video chatting with the other 3-10 other players.</p>
                <p>What works well for us is sitting in front of a laptop for the video chat, and phones for the game.</p>

                <ul>
                    <li>Connect with your friends on video (I suggest <a href="https://meet.jit.si">meet.jit.si</a></li>
                    <li>One player starts a new game and shares the code with everybody via. the video chat</li>
                    <li>Everybody joins the game</li>
                    <li>Some player presses the start button</li>
                </ul>

                Once the game has started, you'll play through a series of rounds.  Each round is the same.  

                <ul>
                    <li>The system will select a weird random word and present it to the players</li>
                    <li>The players each have some period of time (the blue bar indicates how much time is left) to make an awesome definition and submit it</li>
                    <li>Once all players have submitted their definition, or the timer has elapsed, the round enters voting mode</li>
                    <li>During voting mode each player will see the real definition and all of the player submitted definitions and votes on which one they thing is correct</li>
                    <li>Once all players have voted, the round is scored.
                        <li>You get one point for each friend who votes for your definition, instead of the correct one. </li>
                        <li>You get three points for guessing the correct definition.</li>
                    </li>
                </ul>

            </div>
        );
    }
}