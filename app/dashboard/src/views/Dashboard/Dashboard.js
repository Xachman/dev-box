import React, { Component } from 'react';
import Workspace from './Workspace';
import Config from "../../Config"

class Dashboard extends Component {
  constructor(props) {
    super(props)

    this.state = {
      workspaces:[]
    }
  }
  componentDidMount() {
    this.getWorkspaces();
  }

  getWorkspaces() {
    fetch(Config.dockerHostUrl()+"workspaces").then((response) => response.json())
    .then((responseJson) => {
      this.setState({ workspaces: responseJson });
    })
  }
  render() {
    return (
      <div className="animated fadeIn">
        {this.state.workspaces.map(function(item, index) {
          console.log("item", item);
          return <Workspace key={index} name={item.Name} />
        })}
      </div>
    )
  }
}

export default Dashboard;
