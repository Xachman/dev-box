import React, { Component } from 'react';
import Workspace from './Workspace';

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
    fetch("http://192.168.1.150:9080/workspaces").then((response) => response.json())
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
