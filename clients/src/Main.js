import React from "react";

export default class MainView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            href: "/",
            data: [],
            progress: "",
            status: "",
            role: localStorage.getItem("role"),
            msg: null
        };
    }

    componentDidMount() {
        // Initialize WebSocket
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
            msg = JSON.parse(msg.data);

            if (msg.name === "task-new") {
                this.setState(prevState => {
                    return {
                        data: [...prevState.data, msg.task]
                    };
                });
            }
        };

        // Capture URL for later
        let role = localStorage.getItem("role");
        this.setState({ href: "/" + role.toLowerCase() });

        // Fetch data
        fetch(`https://api.kangwoo.tech/tasks/${this.props.match.params.id}`, {
            method: "GET",
            headers: {
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
            })
            .catch(error => {
                alert(error);
                localStorage.clear();
                this.props.history.push({ pathname: "/signin" });
            });
    }

    renderTask = task => {
        return (
            <div key={task._id}>
                <div className="row">
                    <div className="username col-md-4">
                        <div className="container p-2">
                            <div className="border">
                                <p className="p-2">
                                    {task.description}
                                    <button
                                        className="btn btn-warning my-2 my-sm-0 pull-right"
                                        onClick={() =>
                                            this.handleProgress(task._id)
                                        }
                                        disabled={
                                            task.progress ||
                                            localStorage.getItem("role") ==
                                                "Admin"
                                        }
                                    >
                                        {task.isProgress
                                            ? "Assigned"
                                            : "Not Assigned"}
                                    </button>
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    };

    handleProgress(id) {
        fetch(`https://api.kangwoo.tech/tasks/progress/${id}`, {
            method: "POST",
            headers: {
                Authorization: localStorage.getItem("auth")
            }
        })
            .then(res => {
                if (!res.ok) {
                    throw Error(res.statusText + " " + res.status);
                }
                return res.json();
            })
            .then(data => {
                window.location.reload();
            })
            .catch(function(error) {
                alert();
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
                {this.state.msg}
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
                            <a className="nav-item nav-link active" href="#">
                                Home <span className="sr-only">(current)</span>
                            </a>
                            <a
                                className="nav-item nav-link"
                                href={this.state.href}
                            >
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
                    <h3 className="p-3">Current Task List</h3>
                </div>

                {this.state.data.length > 0 ? (
                    this.state.data.map(task => {
                        return this.renderTask(task);
                    })
                ) : (
                    <div>No tasks</div>
                )}
            </div>
        );
    }
}
