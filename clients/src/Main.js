import React from "react";

export default class MainView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            role: true,
            href: "/admin"
        }
        
    }

    handleSignOut() {

    }

    render() {
        return (
            <div>
                <nav className="navbar navbar-expand-lg navbar-dark bg-secondary">
                    <a className="navbar-brand" href="#">To Do App</a>
                    <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"></span>
                    </button>
                    <div className="collapse navbar-collapse" id="navbarNavAltMarkup">
                        <div className="navbar-nav">
                            <a className="nav-item nav-link active" href="#">Home <span className="sr-only">(current)</span></a>
                            <a className="nav-item nav-link" href={this.state.href}>UserBoard</a>
                            <a className="nav-item nav-link" href="#">LeaderBoard</a>
                        </div>
                    </div>
                    <button className="btn btn-warning my-2 my-sm-0 pull-right"
                        onClick={() => this.handleSignOut()}>
                        Sign Out
                    </button>
                </nav>
                <div>
                    <h3 className="p-3">Current Task List</h3>
                    <button className="btn btn-warning my-2 my-sm-0 pull-right"
                        onClick={() => this.setState({href: '/member'})}>
                        Sign Out
                    </button>
                </div>
            </div>
        );
    }


}