import React, { Component } from 'react';
import ProgressBar from 'react-bootstrap/ProgressBar';

interface GameTimerProps {
    remaining: number,
    total: number,
}

export class GameTimer extends Component<GameTimerProps, any> {
    intervalId: NodeJS.Timeout | null = null;

    constructor(props: GameTimerProps) {
        super(props)
        this.state = {
            remaining: 0,
        }
        this.timer = this.timer.bind(this);
    }

    componentDidMount() {
        this.intervalId = setInterval(this.timer, 1000);
     }
     
     componentWillUnmount() {
        if (this.intervalId) {
            clearInterval(this.intervalId);
        }
     }

     componentDidUpdate(prevProps: GameTimerProps, prevState: any, snapshot: any) {
         if (prevProps.remaining !== this.props.remaining) {
             this.setState({
                 remaining: this.props.remaining
             })
         }
     }

     timer() {
         let rem = this.state.remaining- 1;
         if (rem < 0) {
             rem = 0;
         }
         this.setState({
            remaining: rem
         })
     }
     
    render() {
        return (
            <div>
                <br />
                {this.state.remaining> 0 &&
                <ProgressBar now={(this.state.remaining/this.props.total) * 100} />
                }
            </div>
        );
    }
}