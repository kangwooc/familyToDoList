import React from "react";

export default class MainView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            // role: true,
            href: "/",
            data: null,
            progress: "",
            admin: true,
            status: "",
            role: localStorage.getItem("role"),
            msg: null
        }
    }
    componentDidMount() {
        let auth = localStorage.getItem("auth");
        let url = "wss://localhost:443/ws?auth=" + auth;
        
        this.socket = new WebSocket(url);
        this.setState({socket : this.socket});
        this.socket.onopen = () => {
          console.log("Connection Opened");
        };
    
        this.socket.onclose = () => {
          console.log("Connection Closed");
        };
      }
    
      componentDidUpdate() {
        this.socket.onmessage = (msg) => {
          console.log("Message received " + msg.data);
          this.setState({msg: msg.data})
        };
      }


    renderTask = (task) => {
        return (
            <div>

            </div>
        );
    }
    
    componentWillMount() {
        let role = localStorage.getItem("role");
        console.log(role);
        this.setState({href: "/" + role.toLowerCase()});
        fetch(`https://localhost:443/tasks/${this.props.match.params.id}`, {
            method: "GET",
            headers: {
                "Authorization": window.localStorage.getItem("auth")
            }
        }).then(res => {
            if (!res.ok) {
                throw Error(res.statusText + " " + res.status);
            }
            return res.json()
        }).then(data => {
            console.log(data);
            let users = data.map((info) => {
                console.log(info.isProgress)
                console.log(info)
                if (info.isProgress) {
                    this.setState({ progress: "Progressing" })
                } else {
                    this.setState({ progress: "Not Assigned" })
                }
                return (
                    <div className="row">
                        <div className="username col-md-4">
                            <div className="container p-2">
                                <div className="border">
                                    <p className="p-2">{info.description}
                                        <button className="btn btn-warning my-2 my-sm-0 pull-right" onClick={() => this.handleProgress(info._id)} disabled={info.progress||(localStorage.getItem("role")=="Admin")}>
                                            {this.state.progress}
                                        </button>
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                );
                
            });
            this.setState({ data: users });
        }).catch(error => {
            alert(error)
            localStorage.clear()
            this.props.history.push({ pathname: '/signin' })
        }
        )
    }

    handleProgress(id) {
        fetch(`https://localhost:443/tasks/progress/${id}`, {
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
            console.log(data)
            window.location.reload();
        }).catch(function (error) {
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
            this.props.history.push({ pathname: '/signin' })
        }).catch(function (error) {
            localStorage.clear()
        })
    }

    render() {
        return (
            <div>
                {this.state.msg}
                <nav className="navbar navbar-expand-lg navbar-dark bg-secondary">
                    <a className="navbar-brand" href="#">To Do App</a>
                    <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"></span>
                    </button>
                    <div className="collapse navbar-collapse" id="navbarNavAltMarkup">
                        <div className="navbar-nav">
                            <a className="nav-item nav-link active" href="#">Home <span className="sr-only">(current)</span></a>
                            <a className="nav-item nav-link" href={this.state.href}>UserBoard</a>
                            {console.log(this.state.href)}
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