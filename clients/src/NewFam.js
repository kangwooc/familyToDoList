import React from "react";
import { Link } from "react-router-dom";
import { ROUTES } from "./constants";


export default class NewFamView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            familyRoomName: ""
        }
    }

    // componentWillMount() {
    //     let auth = window.localStorage.getItem('auth')
    //     if (auth !== null ) {
    //         this.props.history.push({pathname: '/users/me'})
    //     }
    // }

    handleMakeRoom() {
        fetch("https://localhost:443/create", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization":localStorage.getItem("auth")
            },
            body: JSON.stringify({
                "RoomName": this.state.familyRoomName,     
            }),
        }).then(res => {
            console.log(res)

            if (!res.ok) { 
                console.log(res.headers.get('Authorization'))

                console.log("111")

                throw Error(res.statusText + " " + res.status);
            }

            localStorage.getItem("auth");
            return res.json()
        }).then(data => {
            console.log(data)
            this.setState({id: data.id})
            this.props.history.push({pathname: '/main'})    // go to main task list
        }).catch(function(error) {
            let errorType = document.createElement("p")
            let errorMessage = document.createTextNode("Error to save your data " + error)
            errorType.appendChild(errorMessage)
            document.getElementById("result").appendChild(errorType)
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
                                            <p>Family UserName</p>
                                            {
                                                (this.state.familyRoomName === "") ?
                                                    <div className="alert alert-danger mt-2">It shouldn't be blank</div> : undefined
                                            }
                                            <input id="User Name" type="text" className="form-control"
                                                placeholder="Family UserName" 
                                                onInput={evt => this.setState({ familyRoomName: evt.target.value})} />
                                        </form>
                                        <button className="btn btn-success mr-2 p-2"
                                            onClick={() => this.handleMakeRoom()}>
                                            Sign Up
                                        </button>

                                    </div>
                                    <Link to={ROUTES.signIn}>Back to HomePage</Link>
                                </div>
                            </div>
                        </div>
                    </div>
                </main>
            </div>
        );
    }
}