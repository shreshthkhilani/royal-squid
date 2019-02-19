import React, { Component } from 'react';
import axios from 'axios';
import date from 'date-and-time';
import './App.css';

let formatDate = function (d) {
  return date.format(date.parse(d, 'YYYY-MM-DDTHH:mm:ssZ'), 'ddd, MMM DD');
}

const initalState = {
  signupClicked: false,
  dinners: [],
  dinnersChecked: false,
  dinnerSelection: undefined,
  slots: 0,
  name: '',
  email: '',
  dietary: '',
  awaitingConfirmation: false,
  otp: '',
  confirmed: false,
  dgae: false,
  error: false
};

class App extends Component {
  constructor(props) {
    super(props);
    this.state = initalState;
    this.escFunction = this.escFunction.bind(this);
    this.signupClickedHandler = this.signupClickedHandler.bind(this);
    this.nameChangedHandler = this.nameChangedHandler.bind(this);
    this.emailChangedHandler = this.emailChangedHandler.bind(this);
    this.dietaryChangedHandler = this.dietaryChangedHandler.bind(this);
    this.resetClickedHandler = this.resetClickedHandler.bind(this);
    this.selectDinner = this.selectDinner.bind(this);
    this.selectSlots = this.selectSlots.bind(this);
    this.buttonClickedHandler = this.buttonClickedHandler.bind(this);
    this.otpChangedHandler = this.otpChangedHandler.bind(this);
    this.confirmedHandler = this.confirmedHandler.bind(this);
    this.dgaeHandler = this.dgaeHandler.bind(this);
    this.attemptConfirmation = this.attemptConfirmation.bind(this);
  }
  componentDidMount(){
    document.addEventListener("keydown", this.escFunction, false);
  }
  componentWillUnmount(){
    document.removeEventListener("keydown", this.escFunction, false);
  }
  escFunction(e) {
    if(e.keyCode === 27) {
      this.resetClickedHandler();
    }
  }
  signupClickedHandler() {
    this.setState({signupClicked: true});
    document.body.style.backgroundColor = '#fff';
    document.body.style.color = '#aaa';
    axios.get('/api/dinners/')
      .then((data) => {
        this.setState({dinners: data.data.dinners, dinnersChecked: true});
      })
      .catch((error) => {
        this.setState({dinners: [], dinnersChecked: true});
        // this.setState({dinners, dinnersChecked: true});
      });
  }
  nameChangedHandler(e) {
    this.setState({name: e.target.value});
  }
  emailChangedHandler(e) {
    this.setState({email: e.target.value});
  }
  dietaryChangedHandler(e) {
    this.setState({dietary: e.target.value});
  }
  resetClickedHandler() {
    this.setState(initalState);
    document.body.style.backgroundColor = '#fbf268';
    document.body.style.color = null;
  }
  selectDinner(e) {
    const dinnerSelection = this.state.dinners.find((dinner) => dinner.id === parseInt(e.target.value));
    this.setState({dinnerSelection});
  }
  selectSlots(e) {
    const slots = parseInt(e.target.value);
    this.setState({slots});
  }
  buttonClickedHandler() {
    const postData = {
      dinnerId: this.state.dinnerSelection.id,
      slots: this.state.slots,
      name: this.state.name,
      email: this.state.email,
      dietary: this.state.dietary
    };
    axios.post('/api/reserve/', postData)
      .then((data) => {
        this.setState({
          awaitingConfirmation: true,
          confirmed: data.data.reservation.confirmed,
          dgae: data.data.reservation.dgae
        });
      })
      .catch((error) => {
        this.setState({error: true});
      });
    // this.setState({awaitingConfirmation: true});
  }
  otpChangedHandler(e) {
    this.setState({otp: e.target.value});
  }
  attemptConfirmation() {
    const postData = {
      dinnerId: this.state.dinnerSelection.id,
      otp: this.state.otp
    };
    axios.post('/api/confirm/', postData)
      .then((data) => {
        this.setState({
          confirmed: data.data.reservation.confirmed,
          dgae: data.data.reservation.dgae
        });
      })
      .catch((error) => {
        this.setState({error: true});
      });
  }
  confirmedHandler() {
    this.attemptConfirmation();
    // this.setState({confirmed: true});
  }
  dgaeHandler() {
    this.attemptConfirmation();
    // this.setState({dgae: true});
  }
  render() {
    let slotted;
    if (this.state.dinners.length !== 0) {
      const dinners = this.state.dinners.map((dinner) => {
        return (<option key={dinner.id} value={dinner.id} onClick={this.selectDinner}>{formatDate(dinner.dinnerTime)} &mdash; {dinner.available} available.</option>);
      });
      let slotOptions = [];
      if (this.state.dinnerSelection) {
        for (let i = 1; i <= this.state.dinnerSelection.available; i++) {
          slotOptions.push(i);
        }
      }
      const slots = slotOptions.map((so) => (<option key={so} value={so} onClick={this.selectSlots}>{so}</option>));
      const inputs = (
        <div>
          <label>name *</label><input type="text" autoComplete="name" id="name" value={this.state.name} placeholder="john doe" onChange={this.nameChangedHandler} /> 
          <br/>
          <label>email address *</label><input type="text" autoComplete="email" id="email" value={this.state.email} placeholder="john@gmail.com" onChange={this.emailChangedHandler} /> 
          <br/>
          <label>any dietary restrictions?</label><input type="text" id="dietary" value={this.state.dietary} onChange={this.dietaryChangedHandler} />
          <br/>
          <div align="center">
          {this.state.name !== '' && this.state.email !== '' ? <button onClick={this.buttonClickedHandler}>signup now!</button> : undefined}
          </div>
        </div>
      );
      const otpsubmit = (
        <div>
          <label>One-time pin (sent to you via email) *</label><input type="text" id="otp" value={this.state.otp} onChange={this.otpChangedHandler} /> 
          <br/>
          <div align="center">
          {this.state.otp !== '' ? <button onClick={this.confirmedHandler}>confirm</button> : undefined}
          <button id="dgae" onClick={this.dgaeHandler}>i didn't get an email...</button>
          </div>
        </div>
      );
      const selections = (
        <div>
          {this.state.dinnerSelection ? <p className="selections">{formatDate(this.state.dinnerSelection.dinnerTime)}</p> : undefined}
          {this.state.slots !== 0 ? <p className="selections">{this.state.slots}</p> : undefined}
        </div>
      );
      if (this.state.error) {
        slotted = (
          <div id="slotted">
            <p className="question">sorry, there was an error processing your request...</p>
          </div>
        );
      } else if (this.state.dgae) {
        slotted = (
          <div id="slotted">
            <p className="question">you're confirmed&mdash;expect an email! see you soon...</p>
          </div>
        );
      } else if (this.state.confirmed) {
        slotted = (
          <div id="slotted">
            <p className="question">thanks for confirming&mdash;see you soon...</p>
          </div>
        );
      } else if (this.state.awaitingConfirmation) {
        slotted = (
          <div id="slotted">
            {otpsubmit}
            <p className="question">check your email! you need to enter the pin sent to you to confirm.</p>
            {selections}
            <p className="selections">{this.state.name}</p>
            <p className="selections">{this.state.email}</p>
            {this.state.dietary ? <p className="selections">{this.state.dietary}</p> : undefined}
          </div>
        );
      } else if (this.state.dinnerSelection && this.state.slots !== 0) {
        slotted = (
          <div id="slotted">
            {inputs}
            <p className="question">who are you?</p>
            {selections}
          </div>
        );
      } else if (this.state.dinnerSelection) {
        slotted = (
          <div id="slotted">
            {slots}
            <p className="question">how many people?</p>
            {selections}
          </div>
        );
      } else {
        slotted = (
          <div id="slotted">
            {dinners}
            <p className="question">pick a dinner.</p>
            {selections}
          </div>
        );
      }
    } else if (this.state.dinnersChecked) {
      slotted = (
        <div id="slotted">
          <p className="question">sorry, there are currently no dinners available.</p>
        </div>
      );
    }
    return (
      <div className="App">
        {slotted ? <div id="reset" onClick={this.resetClickedHandler}>[x]</div> : undefined}
        <div id="txt">
          <p><span className="imptxt">silent&amp;counter</span>&nbsp;is a $5-per-person 2-person-per-night dinner in silence.</p>
          <p>I am a performance artist, but this isn’t a performance. Come to my home, eat the food I’ve cooked, and leave.</p>
          <p>Your stated dietary restrictions will be accommodated, but you will not know what the meal is until it is in front of you.</p>
          <p>Dinners take place Tuesdays and Thursdays 10pm in West Philly.</p>
          <span className="imptxt" id="signup" onClick={this.signupClickedHandler}>signup</span>
          {slotted}
        </div>
        <div id="img" className={slotted ? 'grayscale' : ''}></div>
      </div>
    );
  }
}

export default App;
