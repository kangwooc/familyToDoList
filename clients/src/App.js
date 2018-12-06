import React, { Component } from 'react';
import { BrowserRouter as Router, Switch, Redirect, Route } from "react-router-dom"
import { ROUTES } from "./constants";
import SignInView from "./SignIn";
import SignUpView from "./SignUp";
import DeepSignUpView from './DeepSignUp';
import NewFamView from "./NewFam"
import JoinView from "./Join"
import MainView from './Main';
import MemberView from './Member';
import AdminView from './Admin';
import AddTaskView from './AddTask';
import ReceiveView from './Receive';
class App extends Component {
  constructor(props) {
    super(props);
    let auth = localStorage.getItem("auth");
    let url = "wss://localhost:443/ws?auth=" + auth;
    this.socket = new WebSocket(url);
  }

  componentDidMount() {
    this.socket.onopen = () => {
      console.log("Connection Opened");
    };
    this.socket.onclose = () => {
      console.log("Connection Closed");
    };
    this.socket.onmessage = (msg) => {
      console.log("Message received " + msg.data);
    };
  }

  render() {
    return (
      <Router>
            <Switch>
              <Route exact path={ROUTES.signIn} component={SignInView} props={this.socket}/>
              <Route path={ROUTES.signUp} component={SignUpView} props={this.socket}/>
              <Route path={ROUTES.deepSign} component={DeepSignUpView} props={this.socket}/>
              <Route path={ROUTES.newFam} component={NewFamView} props={this.socket}/>
              <Route path={ROUTES.join} component={JoinView} props={this.socket}/>
              <Route path={ROUTES.main} component={MainView} props={this.socket} />
              <Route path={ROUTES.member} component={MemberView} props={this.socket} />
              <Route path={ROUTES.admin} component={AdminView} props={this.socket} />
              <Route path={ROUTES.add} component={AddTaskView} props={this.socket} />
              <Route path={ROUTES.receive} component={ReceiveView} props={this.socket} />
              <Redirect to={ROUTES.signIn} />
            </Switch>
      </Router>
    );
  }
}

export default App;
