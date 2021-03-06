import React from "react";

export default class ReceiveView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            roomname: "",
            data: []
        };
    }
    componentDidMount() {
        let auth = localStorage.getItem("auth");
        let url = "wss://api.kangwoo.tech/ws?auth=" + auth;

        this.socket = new WebSocket(url);
        this.setState({ socket: this.socket });
        this.socket.onopen = () => {
            console.log("Connection Opened");
        };

        this.socket.onclose = () => {
            console.log("Connection Closed");
        };

        this.socket.onmessage = msg => {
            console.log("Message received " + msg.data);
            msg = JSON.parse(msg.data);
            console.log(msg.username !== null);
            if (msg.username !== null) {
                this.setState(prevState => {
                    return {
                        data: [...prevState.data, msg]
                    };
                });
            }
        };
    }
    
    componentWillMount() {
        fetch(" https://api.kangwoo.tech/receive", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                Authorization: window.localStorage.getItem("auth")
            }
        })
            .then(res => {
                if (!res.ok) {
                    throw Error(res.statusText + " " + res.status);
                }
                return res.json();
            })
            .then(data => {
                this.setState({
                    data
                });
                this.setState({ roomname: data.roomname });
            })
            .catch(error => {
                alert(error);
                localStorage.clear();
                this.props.history.push({ pathname: "/signin" });
            });
    }

    renderUser = info => {
        let userName = info.firstname + " " + info.lastname;
        return (
            <div className="row">
                <div class="col-md-4">
                    <div className="container p-2">
                        <div className="border">
                            <div className="p-2">
                                <p>
                                    {userName}
                                    <div className="my-2 my-sm-0 pull-right">
                                        <button
                                            className="btn btn-success mr-2 my-sm-0"
                                            onClick={() =>
                                                this.handleAccept(
                                                    info.id,
                                                    info.roomname
                                                )
                                            }
                                            disabled={info.progress}
                                        >
                                            Accept
                                        </button>
                                        <button
                                            className="btn btn-danger mr-2 my-sm-0"
                                            onClick={() =>
                                                this.handleReject(
                                                    info.id,
                                                    info.roomname
                                                )
                                            }
                                            disabled={info.progress}
                                        >
                                            Refuse
                                        </button>
                                    </div>
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    };

    handleAccept(id, roomname) {
        console.log(id + " " + roomname)
        fetch(" https://api.kangwoo.tech/accept", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: localStorage.getItem("auth")
            },
            body: JSON.stringify({
                personrole: "Member",
                roomname: roomname,
                memberid: id
            })
        })
            .then(() => {
                // this.props.history.push({pathname: '/receive'})
                window.location.reload();
            })
            .catch(function(error) {
                alert(error);
            });
    }

    handleReject(id, roomname) {
        fetch(" https://api.kangwoo.tech/accept", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: localStorage.getItem("auth")
            },
            body: JSON.stringify({
                personrole: "default",
                roomname: roomname,
                memberid: id
            })
        })
            .then(() => {
                // this.props.history.push({pathname: '/receive'})
                window.location.reload();
            })
            .catch(function(error) {
                alert(error);
            });
    }

    handleSignOut() {
        fetch("https://api.kangwoo.tech/sessions/mine", {
            method: "DELETE",
            headers: {
                Authorization: localStorage.getItem("auth")
            }
        })
            .then(res => {
                if (!res.ok) {
                    throw Error(res.statusText + " " + res.status);
                }
                localStorage.clear();
                this.props.history.push({ pathname: "/signin" });
            })
            .catch(function(error) {
                localStorage.clear();
            });
    }

    render() {
        return (
            <div>
                <nav className="navbar navbar-expand-lg navbar-dark bg-secondary">
                    <a className="navbar-brand" href="#">
                        To Do App
                    </a>
                    <button
                        className="navbar-toggler"
                        type="button"
                        data-toggle="collapse"
                        data-target="#navbarNavAltMarkup"
                        aria-controls="navbarNavAltMarkup"
                        aria-expanded="false"
                        aria-label="Toggle navigation"
                    >
                        <span className="navbar-toggler-icon" />
                    </button>
                    <div
                        className="collapse navbar-collapse"
                        id="navbarNavAltMarkup"
                    >
                        <div className="navbar-nav">
                            <a className="nav-item nav-link" href="/admin">
                                UserBoard
                            </a>
                        </div>
                    </div>
                    <button
                        className="btn btn-warning my-2 my-sm-0 pull-right"
                        onClick={() => this.handleSignOut()}
                    >
                        Sign Out
                    </button>
                </nav>
                <div>
                    <h3 className="p-3">Current Request</h3>
                </div>
                {this.state.data.length > 0 ? (
                    this.state.data.map(user => {
                        return this.renderUser(user);
                    })
                ) : (
                    <div>No Member</div>
                )}
            </div>
        );
    }
}
