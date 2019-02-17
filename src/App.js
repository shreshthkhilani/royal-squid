import React, { Component } from 'react';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <div id="txt">
          <p><span className="imptxt">silent&amp;counter</span>&nbsp;is a $5-per-person 2-person-per-night dinner in silence.</p>
          <p>I am a performance artist, but this isn’t a performance. Come to my home, eat the food I’ve cooked, and leave.</p>
          <p>Your stated dietary restrictions will be accommodated, but you will not know what the meal is until it is in front of you.</p>
          <p>Dinners take place Tuesdays and Thursdays 10pm in West Philly.</p>
          <span className="imptxt" id="signup">signup</span>
          <div id="slotted"></div>
        </div>
        <div id="img"></div>
      </div>
    );
  }
}

export default App;
