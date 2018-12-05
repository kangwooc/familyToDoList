import React from "react";

export default class MemberView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {

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
                            <a className="nav-item nav-link" href='/main'>Home</a>
                            <a className="nav-item nav-link active" href="#"><span className="sr-only">(current)</span>UserBoard</a>
                            <a className="nav-item nav-link" href="#">LeaderBoard</a>
                        </div>
                    </div>
                    <button className="btn btn-warning my-2 my-sm-0 pull-right"
                        onClick={() => this.handleSignOut()}>
                        Sign Out
                    </button>
                </nav>
                <div>
                    <h3 className="p-3">Current Member</h3>
                </div>
            </div>
        );
    }
}