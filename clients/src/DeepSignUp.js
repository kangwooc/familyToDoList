import React from "react";
import { Link } from "react-router-dom";
import { ROUTES } from "./constants";


export default class DeepSignUpView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
          token: window.localStorage.getItem("auth")
        }
    }
    
    componentWillMount() {
        let auth = window.localStorage.getItem('auth')
        if (auth === null ) {
            this.props.history.push({pathname: '/signin'})
        }
    }

    handleNewFam() {
        this.props.history.push({pathname: '/newFam'})

    }
    handleJoin() {
        this.props.history.push({pathname: '/join'})

    }

    handleSignOut() {
        fetch("https://api.kangwoo.tech/sessions/mine", {
            method: "DELETE",
            headers: {
                "Authorization": localStorage.getItem("auth")
            }
        }).then(res => {
            if (!res.ok) { 
                throw Error(res.statusText + " " + res.status);
            }
            localStorage.clear()
            this.props.history.push({pathname: '/signin'})
        }).catch(function(error) {
            localStorage.clear()
        })   
    }

    render() {
        return (
            <div>
                <header className="container-fluid bg-secondary text-white">
                    <div className="row">
                        <div className="col-12 col-sm-12 col-md-12 col-lg-12 col-xl-12 pt-3 my-border"xw>
                            <div className="text-center">
                                <h1>To Do App</h1>
                            </div>     
                        </div>
                    </div>
                </header>
                <main>
                <div className="d-flex justify-content-center pt-4 pb-5">
                        <div className="card w-50">
                            <div className="card-body">
                                <div className="container d-flex justify-content-center">
                                    <div className="p-3">
                                        <button type="button" className="btn btn-outline-success btn-lg" onClick={() => this.handleNewFam()}>New Family</button>
                                    </div>
                                    <div className="p-3">
                                        <button type="button" className="btn btn-outline-warning btn-lg" onClick={() => this.handleJoin()}>Join the Family</button>
                                    </div>
                                </div>
                                <div className="p-3">
                                        <button type="button" className="btn btn-outline-warning btn-lg" onClick={() => this.handleSignOut()}>SignOut</button>
                                </div>
                            </div>
                        </div>
                    </div>
                </main>
            </div>
        );
    }
}