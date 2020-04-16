import React, {Component} from 'react';
import Button from 'react-bootstrap/Button';
import { useHistory } from "react-router-dom";
import axios from 'axios';

function HomeButton(props: {rounds: number}) {
    let history = useHistory();
    return (<Button onClick={()=>{

                axios.post('/api/game/new', {
                    rounds:props.rounds 
                }).then((resp: any) => {
                    history.push('/game/' + resp.data.id);
                }).catch((err) => {
                    console.log(err)
                })
            }} variant="primary">New Game</Button>);
  };

export class Home extends Component {
    render() {
        return (
            <div>
                <HomeButton rounds={3} />
            </div>
        );
    }
}