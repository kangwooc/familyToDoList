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
class App extends Component {
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
            <Route path={ROUTES.add} component={AddTaskView} />
            </Switch>
      </Router>
    );
  }
}

export default App;
