import React, { Component } from 'react';
import Button from 'react-bootstrap/Button';
import Axios from 'axios';

interface GameVotingFormProps {
    gameId: string,
    roundId: string,
    playerId: string,
    definitions: any,
}

export class GameVotingForm extends Component<GameVotingFormProps, any> {

    constructor(props: GameVotingFormProps) {
        super(props)

        this.vote = this.vote.bind(this)
    }

    vote(definitionId: string) {
        Axios.post("/api/game/" + this.props.gameId + "/vote", {
            playerId: this.props.playerId,
            roundId: this.props.roundId,
            definitionId: definitionId,
        }).then(() => {
        })
    }


    render() {

        return (
            <div>

                <div>
                    {this.props.definitions.map((def: any) => {
                        return (
                            <div key={def.id}>
                                <Button disabled={def.ownDefinition} key={def.id} onClick={() => { this.vote(def.id); }}>
                                    {def.definition}
                                </Button>
                            </div>
                        );
                    })}
                </div>

            </div>
        );
    }
}