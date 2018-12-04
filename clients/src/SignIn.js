
import React from "react";
import {Link} from "react-router-dom";
import {ROUTES} from "./constants";
export default class SignInView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            userName: "",
            password: ""
        };
    }
    
    // componentDidMount() {
    //     let auth = window.localStorage.getItem('auth')
    //     if (auth !== null ) {
    //         this.props.history.push({pathname: '/users/me'})
    //     }
    // }

    handleSubmit() {

    }

    render() {
        return (     
            <div>
                <header className="container-fluid bg-secondary text-white">
                    <div className="row ">
                        <div className="col-12 col-sm-12 col-md-12 col-lg-12 col-xl-12 pt-3 my-border" >
                            <div className="text-center" >
                                <h1>To Do App</h1>
                            </div>     
                        </div>
                    </div>
                </header>
                <main>
                    <div className="d-flex justify-content-center pt-4">
                        <div className="card w-50">
                            <div className="card-body">
                                <div className="container">
                                    <div id="result"></div>
                                    <form onSubmit={evt => this.handleSubmit(evt)}>
                                        <div className="form-group">
                                            <label htmlFor="UserName">UserName</label>
                                            <input type="text"
                                                id="UserName"
                                                className="form-control"
                                                placeholder="UserName"
                                                onInput={evt=>this.setState({email:evt.target.value})} required/>
                                        </div>
                                        <div className="form-group">
                                            <label htmlFor="password">Password</label>
                                            <input type="password"
                                                id="password"
                                                className="form-control"
                                                placeholder="Password"
                                                onInput={evt=>this.setState({password:evt.target.value})} required/>
                                        </div>
                                        <div className="form-group">
                                            <button type="submit" className="btn btn-primary">Sign In</button>
                                        </div>
                                    </form>
                                    <p>Don't have an account yet? <Link to={ROUTES.signUp}> Sign Up!</Link></p>
                                </div>
                            </div>
                        </div>
                    </div>
                </main>
            </div>        
        );
    }
}