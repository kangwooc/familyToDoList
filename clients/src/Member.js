import React from "react";

export default class MemberView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            done: "",
            data: null,
        }
    }
    
    componentWillMount() {
        fetch(`https://localhost:443/tasks/${window.localStorage.getItem("roomid")}`, {
            method: "GET",
            headers: {
                "Authorization": window.localStorage.getItem("auth")
            }
        }).then((res) => {
            if (!res.ok) {
                throw Error(res.statusText + " " + res.status);
            }
            return res.json()
        })
        .then((data) => {
            console.log(data);
            let user = data.map((info) => {
                console.log(info.isProgress);
                console.log(info.userid);
                // check condition
                if (info.isProgress && localStorage.getItem('userid') === info.userid) {
                //render 
                return ( 
                    <div className="row">
                        <div className="username col-md-4">
                            <p>{info.description}</p>
                            <button className="btn btn-warning my-2 my-sm-0 pull-right" onClick={() => this.handleDone(info._id)} disabled={!info.progress}>
                                {this.state.done}
                            </button>
                        </div>
                    </div>
                  );
                }
            });
            this.setState({data: user});
        }).catch(

        );
    }
    handleDone(id) {
        fetch(`https://localhost:443/tasks/done/${id}`, {
            method: "POST",
            headers: {
                "Authorization": localStorage.getItem("auth")
            }
        }).then(res => {
            if (!res.ok) { 
                throw Error(res.statusText + " " + res.status);
            }
            return res.json()
        }).then(data => {
            console.log(data);
        }).catch(function(error) {
            alert()
        })
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
                </div>
                {this.state.data}
            </div>
        );
    }
}