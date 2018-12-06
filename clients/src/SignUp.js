import React from "react";
import { Link } from "react-router-dom";
import { ROUTES } from "./constants";


export default class SignUpView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            password: "",
            passwordConfirm: "",
            userName: "",
            firstName: "",
            lastName: ""
        }
    }

    componentWillMount() {
        let auth = window.localStorage.getItem('auth')
        if (auth !== null) {
            console.log("auth is not null")
            console.log(auth)
        }
    }

    handleSignUp() {
        fetch("https://localhost:443/users", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
	            "Password": this.state.password,     
	            "PasswordConf": this.state.passwordConfirm,
	            "UserName": this.state.userName,    
	            "FirstName": this.state.firstName,   
	            "LastName": this.state.lastName
            }),
            mode: "cors",
            cache: "default",
        }).then(res => {
            if (!res.ok) { 
                throw Error(res.statusText + " " + res.status);
            }
            console.log(res.headers.get("Authorization"));
            localStorage.setItem("auth", res.headers.get('Authorization'));
            return res.json()
        }).then(data => {
            console.log(data)
            this.setState({id: data.id})
            this.props.history.push({pathname: '/deepSign'})
        }).catch(function(error) {
            let errorType = document.createElement("p")
            let errorMessage = document.createTextNode("Error to save your data " + error)
            errorType.appendChild(errorMessage)
            // document.getElementById("result").appendChild(errorType)
        })
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
                <div className="d-flex justify-content-center pt-4 pb-5">
                        <div className="card w-50">
                            <div className="card-body">
                                <div className="container">
                                    <div>
                                    <div id="result"></div>
                                        <form className="form-group">
                                            <p>User Name</p>
                                            {
                                                (this.state.userName === "") ?
                                                    <div className="alert alert-danger mt-2">It shouldn't be blank</div> : undefined
                                            }
                                            <input id="User Name" type="text" className="form-control"
                                                placeholder="User Name" 
                                                onInput={evt => this.setState({ userName: evt.target.value})} />
                                        </form>
                                        <form className="form-group">
                                            <p>First Name</p>
                                            {
                                                (this.state.firstName === "") ?
                                                    <div className="alert alert-danger mt-2">It shouldn't be blank</div> : undefined
                                            }
                                            <input id="First Name" type="text" className="form-control"
                                                placeholder="First Name" 
                                                onInput={evt => this.setState({ firstName: evt.target.value})} />
                                        </form>
                                        <form className="form-group">
                                            <p>Last Name</p>
                                            {
                                                (this.state.lastName === "") ?
                                                    <div className="alert alert-danger mt-2">It shouldn't be blank</div> : undefined
                                            }
                                            <input id="Last Name" type="text" className="form-control"
                                                placeholder="First Name" 
                                                onInput={evt => this.setState({ lastName: evt.target.value})} />
                                        </form>
                                        <form className="form-group">
                                            <p>Password</p>
                                            {
                                                (this.state.password.length < 6 && this.state.password.length !== 0) ?
                                                    <div className="alert alert-danger mt-2">the password should be longer than 6 digits</div> : undefined
                                            }
                                            <input id="password" type="password" className="form-control"
                                                placeholder="password" 
                                                onInput={evt => this.setState({ password: evt.target.value})}
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>Password Confirmation</p>
                                            {
                                                (this.state.password !== this.state.passwordConfirm && this.state.passwordConfirm !== "") ?
                                                    <div className="alert alert-danger mt-2">password should match</div> : undefined
                                            }
                                            <input id="passwordConfirm" type="password" className="form-control"
                                                placeholder="password" 
                                                onInput={evt => this.setState({ passwordConfirm: evt.target.value})}
                                            />

                                        </form>
                                        <button className="btn btn-primary mr-2 p-2"
                                            onClick={() => this.handleSignUp()}>
                                            Sign Up
                                        </button>

                                    </div>
                                    <p>Already have an account? <Link to={ROUTES.signIn}>Sign In!</Link></p>
                                </div>
                            </div>
                        </div>
                    </div>
                </main>
            </div>
        );
    }
}