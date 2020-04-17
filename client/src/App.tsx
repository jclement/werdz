import React from 'react';
import './App.css';
import { Route } from 'react-router';
import { Layout } from './components/Layout';
import { Home } from './components/Home';
import { About } from './components/About';
import { GameShell } from './components/GameShell';
import { BrowserRouter as Router } from 'react-router-dom'
import Axios from 'axios'
import { Rules } from './components/Rules';

export class App extends React.Component<{}, any> {
  constructor(props: any) {
    super(props)
    this.state = {
      playerId: localStorage.getItem("playerid"),
      playerName: localStorage.getItem("playername"),
    }

    if (!this.state.playerId) {
      Axios.get('/api/player/generate', {})
        .then((resp: any) => {
          this.setState({
            playerId: resp.data.id,
          })
          localStorage.setItem('playerid', resp.data.id)
        })
    }

    this.setName = this.setName.bind(this);
  }

  setName(name: string) {
    this.setState({
      playerName: name,
    });
    localStorage.setItem('playername', name);
  }

  render() {
    return (
      <Router>
        <Layout>
          <Route exact path='/' render={({ match }) => (
            <Home playerId={this.state.playerId} playerName={this.state.playerName} />
          )} />
          <Route path='/about' component={About} />
          <Route path='/rules' component={Rules} />
          <Route path='/game/:id' render={({ match }) => (
            <GameShell gameId={match.params.id} playerId={this.state.playerId} playerName={this.state.playerName || ""} onSetName={this.setName} /> 

          )} />
        </Layout>
      </Router>
    );
  }
}

export default App;
