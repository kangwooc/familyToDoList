import React from "react";

export default class AdminView extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            href: "/admin",
            data: null,
            roomname: localStorage.getItem("roomname"),
            id: ""
        }
    }
    componentDidMount() {
        let room = localStorage.getItem("roomname")
        console.log(room)
        this.setState({ href: "/main/" + room.toLowerCase() })
    }

    handleDelete(id) {
        fetch(`https://api.kangwoo.tech/delete`, {
            method: "DELETE",
            headers: {
                "Authorization": localStorage.getItem("auth"),
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "id": id,
            }),
        }).then(res => {
            if (!res.ok) {
                throw Error(res.statusText + " " + res.status);
            }
            // localStorage.clear()
            alert("Deleted Member!")
            window.location.reload();
        }).catch(function (error) {
            localStorage.clear()
            alert(error)
        })
    }
    // this is to display all the members in this room
    componentWillMount() {
        fetch(`https://api.kangwoo.tech/memberlist/${window.localStorage.getItem("roomname")}`, {
            method: "GET",
            headers: {
                "Authorization": localStorage.getItem("auth"),
                "Content-Type": "application/json"
            }
        }).then((res) => {
            if (!res.ok) {
                throw Error(res.statusText + " " + res.status);
            }
            return res.json()
        }).then((data) => {
            if (data == null) {
                console.log("no member");

                this.setState({ data: "no member" });

            } else {
                console.log("this is data");

                console.log(data);
                let user = data.map((info) => {
                    console.log(info.firstname + " " + info.lastname);
                    return (
                        <div className="row">
                            <div className="username col-md-4">
                                <div className="container p-2">
                                    <div className="border">
                                        <p className="p-2">{info.firstname + " " + info.lastname + " Lv." + (info.score/100)}
                                            <button className="btn btn-danger my-2 my-sm-0 pull-right" onClick={() => this.handleDelete(info.id)}>
                                                Delete
                                    </button>
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    );
                });
                this.setState({ data: user });
            }
        }
        ).catch(error => {
            alert(error)
            localStorage.clear()
            this.props.history.push({ pathname: '/signin' })
        }
        )
    }

    handleSignOut() {
        fetch("https://api.kangwoo.tech/sessions/mine", {
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
                <nav className="navbar navbar-expand-lg navbar-dark bg-secondary">
                    <a className="navbar-brand" href="#">To Do App</a>
                    <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"></span>
                    </button>
                    <div className="collapse navbar-collapse" id="navbarNavAltMarkup">
                        <div className="navbar-nav">
                            <a className="nav-item nav-link" href={this.state.href}>Home</a>
                            <a className="nav-item nav-link active" href="#"><span className="sr-only">(current)</span>UserBoard</a>
                        </div>
                    </div>
                    <button type="button" className="btn btn-outline-warning mr-2 fa fa-plus" onClick={() => { this.props.history.push("/add"); console.log("clicked") }}></button>
                    <button className="btn btn btn-outline-warning mr-2 my-sm-0"
                        onClick={() => this.props.history.push({ pathname: '/receive' })}>
                        Request
                    </button>
                    <button className="btn btn-warning mr-2 my-sm-0"
                        onClick={() => this.handleSignOut()}>
                        Sign Out
                    </button>
                </nav>
                <h3 className="p-3">Current Member</h3>
                <div className="m-3">
                    {this.state.data}
                </div>
            </div>
        );
    }
}