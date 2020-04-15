import React, {Component} from 'react';

interface GameProps {
  id:string,
}

export class Game extends Component<GameProps, any> {

    constructor(props : GameProps) {
      super(props)
    }

    render() {
        return (
            <div>
                <p>Game : {this.props.id}</p>
            </div>
        );
    }
}