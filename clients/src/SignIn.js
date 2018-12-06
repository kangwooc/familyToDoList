
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
    
    componentDidMount() {
        let auth = window.localStorage.getItem('auth')
        if (auth !== null ) {
            this.props.history.push({pathname: '/main/' +  localStorage.getItem("roomname")})
        }
    }

    handleSubmit(evt) {
        evt.preventDefault();

        fetch("https://localhost:443/sessions", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
	            "username": this.state.userName,    
	            "password":  this.state.password
            }),
        }).then(res => {
            console.log("Sssss")
            if (!res.ok) { 
                throw Error(res.statusText + " " + res.status);
            }
            console.log(res.headers.get("Authorization"));
            localStorage.setItem("auth", res.headers.get('Authorization'));
            return res.json()
        }).then(data => {
            console.log(data)
            if (data.personrole == "Admin" || data.personrole == "Member") {
                console.log(data.roomname)
                localStorage.setItem("roomname", data.roomname);
                localStorage.setItem("userid", data.id);
                localStorage.setItem("role", data.personrole);
                this.props.history.push({pathname: '/main/' + data.roomname})
            } else {
                this.props.history.push({pathname: '/deepSign'})
            }
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
                    <div className="d-flex justify-content-center pt-4">
                        <div className="card w-50">
                            <div className="card-body">
                                <div className="container">
                                    <div id="result"></div>
                                    <form>
                                        <div className="form-group">
                                            <label htmlFor="UserName">UserName</label>
                                            <input type="text"
                                                id="UserName"
                                                className="form-control"
                                                placeholder="UserName"
                                                onInput={evt=>this.setState({userName:evt.target.value})} required/>
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
                                            <button type="submit" className="btn btn-primary" onClick={(evt) => this.handleSubmit(evt)}>Sign In</button>
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