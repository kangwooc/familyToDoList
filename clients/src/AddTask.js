import React from "react";
import { Link } from "react-router-dom";
import { ROUTES } from "./constants";


export default class AddTaskView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            search: ""
        }
    }

    // componentWillMount() {
    //     let auth = window.localStorage.getItem('auth')
    //     if (auth !== null ) {
    //         this.props.history.push({pathname: '/users/me'})
    //     }
    // }

    handleSearch() {

    }

    render() {
        return (
            <div>
                <nav className="navbar navbar-expand-lg navbar-dark bg-secondary">
                    <a className="navbar-brand" href="/main">To Do App</a>
                    <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"></span>
                    </button>
                    <div className="collapse navbar-collapse" id="navbarNavAltMarkup">
                        <div className="navbar-nav">
                            <a className="nav-item nav-link" href='/main'>Home</a>
                            <a className="nav-item nav-link" href="/admin">UserBoard</a>
                            <a className="nav-item nav-link" href="#">LeaderBoard</a>
                        </div>
                    </div>
                    <button className="btn btn-warning my-2 my-sm-0 pull-right"
                        onClick={() => this.handleSignOut()}>
                        Sign Out
                    </button>
                </nav>
                <div className="d-flex justify-content-center pt-4 pb-5">
                    <div className="card w-50">
                        <div className="card-body">
                            <div className="container">
                                <div className="d-flex justify-content-center pt-4 pb-5">
                                    <h4>Add New Task</h4>
                                </div>
                                <form className="form-inline">
                                    <div className="form-group">
                                        <input type="New Task" className="form-control" placeholder="New Task"/>
                                    </div>
                                    <button type="submit" className="btn btn-warning mt-2 mb-2 ml-2" onClick={() => this.handleSearch()}>Submit</button>
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}