import React from "react";
import { Link } from "react-router-dom";
import { ROUTES } from "./constants";


export default class JoinView extends React.Component {
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
                                            <label for="Search" className="sr-only">Search</label>
                                            <input type="Search" className="form-control" placeholder="Search"/>
                                        </div>
                                        <button type="submit" className="btn btn-warning mb-2" onClick={() => this.handleSearch()}>Search</button>
                                    </form>
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