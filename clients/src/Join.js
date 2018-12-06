import React from "react";
import { Link } from "react-router-dom";
import { ROUTES } from "./constants";


export default class JoinView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            roomname: "",
            role: "",
        }
    }

    componentWillMount() {
        let auth = window.localStorage.getItem('auth')
        if (auth === null ) {
            this.props.history.push({pathname: '/signin'})
        }
    }

    handleSignOut() {
        fetch("https://localhost:443/sessions/mine", {
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

    handleSearch() {
        fetch("https://localhost:443/join", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": localStorage.getItem("auth")
            },
            body: JSON.stringify({
                "roomname": this.state.roomname,
            }),

        }).then(res => {
            if (!res.ok) {
                console.log(localStorage.getItem("auth"))
                console.log(this.state.roomname)
                throw Error(res.statusText + " " + res.status);
            }

        }).then(() => {
            // console.log("proceed")
            alert("please wait for admin's decision, log out and re-log in to check your status")
            // this.setState({ id: data.id })
            // this.props.history.push({ pathname: '/main/' + data.id })    // go to main task list
        }).catch(function (error) {
            alert("RoomNotFound Please double check the room name")
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
                                    <form className="form-inline">
                                        <div className="form-group mx-sm-3 mb-2">
                                            <input type="Search" className="form-control" placeholder="Search"
                                             onInput={evt => this.setState({ roomname: evt.target.value})} />
                                             {console.log(this.state.roomname)}
                                        </div>
                                    </form>
                                    <button className="btn btn-primary"
                                            onClick={() => this.handleSearch()}>
                                            Search
                                    </button>
                                    <button className="btn btn-warning"
                                            onClick={() => this.handleSignOut()}>
                                            SignOut
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </main>
            </div>
        );
    }
}