import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import Config from "../../Config"

class Workspace extends Component {
    constructor(props) {
        super(props)

        this.state = {
            name: props.name,
            status: ""
        }
        console.log(new Config())
        this.getStatus(this.state.name)
    }
    getStatus(id) {
        fetch(Config.dockerHostUrl()+"workspaces/status/"+id).then((response) => response.json())
        .then((responseJson) => {
        this.setState({ status: responseJson.Status });
        })

    }
    startContainer() {
        fetch(Config.dockerHostUrl()+"workspaces/start/"+this.state.name,{method: "POST"}).then(() => {
            this.getStatus(this.state.name)
        })
    }
    stopContainer() {
        fetch(Config.dockerHostUrl()+"workspaces/stop/"+this.state.name, {method:"POST"}).then(() => this.getStatus(this.state.name))
    }
    removeContainer() {
        fetch(Config.dockerHostUrl()+"workspaces/remove/"+this.state.name, {method:"POST"}).then(() => this.getStatus(this.state.name))
    }
    render() {
        console.log("state", this.state)
        return (
        <div className="workspace">
            <span className="Name">{this.state.name}</span> &nbsp;
            <span className="status">{this.state.status}</span> &nbsp;
            <button className="start btn btn-primary" onClick={this.startContainer.bind(this)}>Start</button> &nbsp;
            <button className="stop btn btn-primary" onClick={this.stopContainer.bind(this)}>Stop</button> &nbsp;
            <button className="stop btn btn-primary" onClick={this.removeContainer.bind(this)}>Remove</button> &nbsp;
            <Link to={"/terminal/"+this.state.name} className="stop btn btn-primary" >Terminal</Link>
        </div>
        )
    }
    }
Workspace.propTypes = {
    name: PropTypes.string.isRequired
}

export default Workspace;
