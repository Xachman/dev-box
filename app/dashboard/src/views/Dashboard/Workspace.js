import React, { Component } from 'react';
import PropTypes from 'prop-types';

class Workspace extends Component {
    constructor(props) {
        super(props)

        this.state = {
            name: props.name,
            status: ""
        }
        this.getStatus(this.state.name)
    }
    getStatus(id) {
        fetch("http://192.168.1.150:9080/workspaces/status/"+id).then((response) => response.json())
        .then((responseJson) => {
        this.setState({ status: responseJson.Status });
        })

    }
    startContainer() {
        fetch("http://192.168.1.150:9080/workspaces/start/"+this.state.name,{method: "POST"}).then(() => {
            this.getStatus(this.state.name)
        })
    }
    stopContainer() {
        fetch("http://192.168.1.150:9080/workspaces/stop/"+this.state.name, {method:"POST"}).then(() => this.getStatus(this.state.name))
    }
    removeContainer() {
        fetch("http://192.168.1.150:9080/workspaces/remove/"+this.state.name, {method:"POST"}).then(() => this.getStatus(this.state.name))
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
        </div>
        )
    }
    }
Workspace.propTypes = {
    name: PropTypes.string.isRequired
}

export default Workspace;
