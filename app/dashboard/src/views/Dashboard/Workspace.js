import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import Config from "../../Config"

class Workspace extends Component {
    constructor(props) {
        super(props)

        this.state = {
            name: props.name,
            status: "",
            projecturl: ""
        }
        console.log(new Config())
        this.getStatus(this.state.name)
        this.getUrl(this.state.name)
    }
    getStatus(id) {
        fetch(Config.dockerHostUrl()+"workspaces/status/"+id).then((response) => response.json())
        .then((responseJson) => {
        this.setState({ status: responseJson.Status });
        })

    }
    getUrl(id) {
        fetch(Config.dockerHostUrl()+"workspaces/ports/"+id).then((response) => response.json())
        .then((responseJson) => {
        this.setState({ projecturl: "http://"+responseJson.PortDomain });
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

    launchIde() {
        fetch(Config.dockerHostUrl()+"workspaces/ide/cloud9/"+this.state.name, {method:"POST"}).then(() => this.getStatus(this.state.name))
    }
    render() {
        console.log("state", this.state)
        return (
        <div className="workspace">
            <div className="item info">
                <div className="name">
                    <span className="Name">{this.state.name}</span>
                </div>
                <div className="status">
                    <span className="status">{this.state.status}</span>
                </div>
            </div>
            <div className="item actions">
                <button className="start btn btn-primary" onClick={this.startContainer.bind(this)}>Start</button>
                <button className="stop btn btn-primary" onClick={this.stopContainer.bind(this)}>Stop</button>
                <button className="stop btn btn-primary" onClick={this.removeContainer.bind(this)}>Remove</button>
                <Link to={"/terminal/"+this.state.name} className="stop btn btn-primary" >Terminal</Link>
                <a href={this.state.projecturl} className="stop btn btn-primary" >See Project</a>
                <button className="stop btn btn-primary" onClick={this.launchIde.bind(this)}>IDE</button>
            </div>
        </div>
        )
    }
    }
Workspace.propTypes = {
    name: PropTypes.string.isRequired
}

export default Workspace;
