import React, {Component} from 'react';
import Button from 'react-bootstrap/Button';
import { useHistory } from "react-router-dom";

function HomeButton() {
    let history = useHistory();
    return (<Button onClick={()=>{
                history.push('/game/test')
            }} variant="primary">New Game</Button>);
  };

export class Home extends Component {

    render() {
        return (
            <div>
                <HomeButton/>
            </div>
        );
    }
}