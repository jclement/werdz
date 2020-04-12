import React from 'react';
import logo from './logo.svg';
import './App.css';
import Button from 'react-bootstrap/Button';

function App() {
  return (
    <div>
      Test<br />
      <Button onClick={() => {console.log("click")}}>Test Button</Button>
    </div>
  );
}

export default App;
