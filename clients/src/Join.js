import React from "react";
import { Link } from "react-router-dom";
import { ROUTES } from "./constants";


export default class JoinView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            // search: ""
            personrole: "",
            roomname: ""
        }
    }

    // componentWillMount() {
    //     let auth = window.localStorage.getItem('auth')
    //     if (auth === null ) {
    //         this.props.history.push({pathname: '/signin'})
    //     }
    // }

    handleSearch() {
        fetch("https://localhost:443/join", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization":localStorage.getItem("auth")
            },
            body: JSON.stringify({
	            "Role": this.state.personrole,    
                "roomname": this.state.roomname,
            }),

        }).then(res => {
            if (!res.ok) { 
                console.log(this.state.personrole)
                throw Error(res.statusText + " " + res.status);
            }
            return res.json()
        }).then(data => {
            console.log(data)
            this.setState({id: data.id})
            this.props.history.push({pathname: '/main'})    // go to main task list
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
                                    <form className="form-inline">
                                        <div className="form-group mx-sm-3 mb-2">
                                            <input type="Search" className="form-control" placeholder="Search"  onInput={evt => this.setState({roomname: evt.target.value})} />
                                        </div>
                                    </form>
                                        <button type="submit" className="btn btn-warning mb-2" onClick={() => this.handleSearch()}>Search</button>
                                    <Link to={ROUTES.signIn}>Go back to Homepage</Link>
                                </div>
                            </div>
                        </div>
                    </div>
                </main>
            </div>
        );
    }
}