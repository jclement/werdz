import React from 'react';
import './App.css';
import Button from 'react-bootstrap/Button';
import axios from 'axios'

export class App extends React.Component<{}, any> {

  constructor(props: any) {
    super(props)
    this.state = {
      msg : "Nothing Here"
    }

    this.handleButton = this.handleButton.bind(this)
  }

  handleButton() {
    axios.post("/api/test/", {
      Name: "World"
    }).then((response) => {
      this.setState({
        msg: response.data.Msg
      })
    })
  }

  render() {
    return (
      <div>
        Message: <b>{this.state.msg}</b><br/>
        <Button onClick={this.handleButton}>Test Button</Button>
      </div>
    );
  }
}

export default App;
