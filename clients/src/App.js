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
    this.state = {
      socket: null
    }
  }

  componentDidMount() {
    let auth = localStorage.getItem("auth");
    let url = "wss://localhost:443/ws?auth=" + auth;
    
    this.socket = new WebSocket(url);
    this.setState({socket : this.socket});
    this.socket.onopen = () => {
      console.log("Connection Opened");
    };

    this.socket.onclose = () => {
      console.log("Connection Closed");
    };
  }

  componentDidUpdate() {
    this.socket.onmessage = (msg) => {
      console.log("Message received " + msg.data);
    };
  }
  render() {
    
    return (
      <Router>
            <Switch>
              <Route exact path={ROUTES.signIn} component={SignInView} />
              <Route path={ROUTES.signUp} component={SignUpView} />
              <Route path={ROUTES.deepSign} component={DeepSignUpView} />
              <Route path={ROUTES.newFam} component={NewFamView} />
              <Route path={ROUTES.join} component={JoinView} />
              <Route path={ROUTES.main} component={MainView} />
              <Route path={ROUTES.member} component={MemberView} />
              <Route path={ROUTES.admin} component={AdminView} />
              <Route path={ROUTES.add} component={AddTaskView} socket={this.state.socket} />
              <Route path={ROUTES.receive} component={ReceiveView} />
              <Redirect to={ROUTES.signIn} />
            </Switch>
      </Router>
    );
  }
}

export default App;
