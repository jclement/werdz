import React, { Component } from 'react';
import Axios from 'axios';
import { definition } from '../models/definition';
import { ListGroup } from 'react-bootstrap';

interface GameVotingFormProps {
    gameId: string,
    roundId: string,
    playerId: string,
    definitions: definition[],
}

export class GameVotingForm extends Component<GameVotingFormProps, any> {

    constructor(props: GameVotingFormProps) {
        super(props)

        this.vote = this.vote.bind(this)
    }

    vote(definitionId: string) {
        (new Audio("/sounds/click.mp3")).play()
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
                    <ListGroup>
                        {this.props.definitions.map((def: any) => {
                            return (
                                <ListGroup.Item key={def.id} disabled={def.ownDefinition} onClick={() => { this.vote(def.id); }}>
                                    {def.definition}
                                </ListGroup.Item>
                            );
                        })}
                    </ListGroup>
                </div>

            </div>
        );
    }
}