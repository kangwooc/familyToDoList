import React, { Component } from 'react';
import { BrowserRouter as Router, Switch, Redirect, Route } from "react-router-dom"
import { ROUTES } from "./constants";
import SignInView from "./SignIn";
import SignUpView from "./SignUp";
class App extends Component {
  render() {
    return (
      <Router>
            <Switch>
            <Route path={ROUTES.signIn} component={SignInView} />
            <Route path={ROUTES.signUp} component={SignUpView} />
            <Redirect to={ROUTES.signIn} />
            </Switch>
      </Router>
    );
  }
}

export default App;
