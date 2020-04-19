import React, { Component } from 'react';

export class About extends Component {
    render() {
        return (
            <div>
                <h1>About Werdz</h1>

                <p>Werdz is an online word guessing game where you try and make up plausible definitions for escoteric words to fool your friends.</p>

                <p>Werdz was created during the COVID-19 isolation because I wanted more games that were fun to play with friends remotely.</p>

                <p>If you have questions, suggestions, bug reports, etc. please submit them to <a href="mailto:werdz@werd.ca">werdz@werdz.ca</a> or <a href="https://twitter.com/werdzgame">@werdzgame</a>.</p>

                <div style={{textAlign: "center"}}><img className="img-fluid" src="title.png" alt="title" /></div>

                <h4>Credits</h4>
                <ul>
                    <li>Sound Effects:
                        <ul>
                            <li><a href="https://freesound.org/people/LittleRobotSoundFactory/sounds/270402/">game winning sound</a> by LittleRobotSoundFactory</li>
                            <li><a href="https://freesound.org/people/Timbre/sounds/73750/">game losing sound</a> by Timbre</li>
                        </ul>
                    </li>
                </ul>
            </div>
        );
    } }