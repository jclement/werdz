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
            remaining: props.remaining,
        }
        this.timer = this.timer.bind(this);
    }

    componentDidMount() {
        this.intervalId = setInterval(this.timer, 100);
    }

    componentWillUnmount() {
        if (this.intervalId) {
            clearInterval(this.intervalId);
        }
    }

    timer() {
        let rem = this.state.remaining - 0.1;
        if (rem < 0) {
            rem = 0;
        }
        this.setState({
            remaining: rem
        })
    }

    componentDidUpdate(prevProps: GameTimerProps, prevState: any, snapshot: any) {
        if (prevProps.remaining !== this.props.remaining) {
            this.setState({
                remaining: this.props.remaining
            })
        }
    }

    render() {
        return (
            <div>
                <br />
                {this.state.remaining > 0 &&
                    <ProgressBar variant={this.state.remaining < 5 ? "danger" : (this.state.remaining < 10 ? "warning" : undefined)} now={(this.state.remaining / this.props.total) * 100} label={Math.round(this.state.remaining) + "s"} />
                }
            </div>
        );
    }
}