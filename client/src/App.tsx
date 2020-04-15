import React from 'react';
import './App.css';
import { Route } from 'react-router';
import { Layout } from './components/Layout';
import { Home } from './components/Home';
import { About } from './components/About';
import { Game } from './components/Game';
import { BrowserRouter as Router } from 'react-router-dom'

export class App extends React.Component<{}, any> {
  render() {
    return (
      <Router>
        <Layout>
          <Route exact path='/' component={Home} />
          <Route path='/about' component={About} />
          <Route path='/game/:id' render={({match}) => (
            <Game id={match.params.id} />
          )}/>
        </Layout>
      </Router>
    );
  }
}

export default App;
