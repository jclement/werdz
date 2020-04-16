import React from 'react';
import './App.css';
import { Route } from 'react-router';
import { Layout } from './components/Layout';
import { Home } from './components/Home';
import { About } from './components/About';
import { Game } from './components/Game';
import { BrowserRouter as Router } from 'react-router-dom'
import Axios from 'axios'

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
            playerName: resp.data.id,
          })
          localStorage.setItem('playerid', resp.data.id)
          localStorage.setItem('playername', resp.data.id)
        })
    }
  }
  render() {
    return (
      <Router>
        <Layout>
          <Route exact path='/' render={({ match }) => (
            <Home playerId={this.state.playerId} playerName={this.state.playerName} />
          )} />
          <Route path='/about' component={About} />
          <Route path='/game/:id' render={({ match }) => (
            <Game gameId={match.params.id} playerId={this.state.playerId} playerName={this.state.playerName} />
          )} />
        </Layout>
      </Router>
    );
  }
}

export default App;
