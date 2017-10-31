import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Config from "../../Config"
class TerminalComponent extends Component {
  constructor(props) {
    super(props)
    console.log("props", props)
    this.state = {name: props.match.params.name}
  }
  componentDidMount() {
            var websocket = new WebSocket("ws://"+Config.host()+":"+Config.port()+"/workspaces/exec/"+this.state.name);
            websocket.onopen = function(evt) {
                console.log(Terminal)
                var term = new Terminal({
                    cols: 100,
                    rows: 30,
                    screenKeys: true,
                    useStyle: true,
                    cursorBlink: true,
                });
                term.on('data', function(data) {
                    console.log(data);
                    websocket.send(data);
                });
                term.on('title', function(title) {
                    document.title = title;
                });
                console.log(this.refs.term)
                term.open(this.refs.term);
                websocket.onmessage = function(evt) {
                    term.write(evt.data);
                }
                websocket.onclose = function(evt) {
                    term.write("Session terminated");
                    term.destroy();
                }
                websocket.onerror = function(evt) {
                    if (typeof console.log == "function") {
                        console.log(evt)
                    }
                }
            }.bind(this)
  }
  
  render() {
    var divStyle = {
        background: "#eee",
        padding: "20px",
        margin: "20px"
      };
      
    return (
      <div className="animated fadeIn">
        <div id="container-terminal" ref="term"></div>
      </div>
    )
  }
}

Terminal.PropTypes = {
    name: PropTypes.string.isRequired
}

export default TerminalComponent;
        // </style>
        // <script type="text/javascript">
        // console.log('test');
        // $(function() {
        // });