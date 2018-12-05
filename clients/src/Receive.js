import React from "react";

export default class ReceiveView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            roomname: ""
        }
        
    }

    componentWillMount() {
        fetch(" https://localhost:443/receive", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": window.localStorage.getItem("auth")
            }
        }).then(res => {
            if (!res.ok) { 
                throw Error(res.statusText + " " + res.status);
            }
            return res.json()
        }).then(data => {
            console.log(data)
            let users = data.map((info) => {
                this.setState({roomname: data.roomname})
                let userName = info.firstname + " " + info.lastname
                return (
                    <div className="row">
                        <div className="username col-md-4">
                            <p>{userName}</p>
                            <button className="btn btn-sucessful my-2 my-sm-0 pull-right" onClick={() => this.handleAccept(info.id, info.roomname)} disabled={info.progress}>
                                Accept
                            </button>
                            <button className="btn btn-alert my-2 my-sm-0 pull-right" onClick={() => this.handleReject(info.id, info.roomname)} disabled={info.progress}>
                                Refuse
                            </button>
                        </div>
                    </div>
                );
            });
            this.setState({data: users});
        }).catch(error => {
                alert(error)
                localStorage.clear()
                this.props.history.push({pathname: '/signin'})
            }        
        )
    }

    handleAccept(id, roomname) {
        console.log(id + " " + roomname)
        fetch(" https://localhost:443/accept", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": localStorage.getItem("auth")
            },
            body: JSON.stringify({
	            "personrole": "Member",     
	            "roomname": roomname,
	            "memberid": id 
            }),
        }).catch(function(error) {
            alert(error)
        })
    }

    handleReject(id, roomname) {
        fetch(" https://localhost:443/accept", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": localStorage.getItem("auth")
            },
            body: JSON.stringify({
	            "personrole": "default",     
	            "roomname": roomname,
	            "memberid": id 
            }),
        }).catch(function(error) {
            alert(error)
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
                            <a className="nav-item nav-link" href="/admin">UserBoard</a>
                        </div>
                    </div>
                    <button className="btn btn-warning my-2 my-sm-0 pull-right"
                        onClick={() => this.handleSignOut()}>
                        Sign Out
                    </button>
                </nav>
                <div>
                    <h3 className="p-3">Current Request</h3>
                </div>
                {this.state.data}
            </div>
        );
    }


}