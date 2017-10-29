import React, { Component } from 'react';

class Terminal extends Component {
  constructor(props) {
    super(props)
    console.log(props)
    this.state = {name: props.match.params.name}
  }
  render() {
    return (
      <div className="animated fadeIn">
        terminal {this.state.name}
        {`<style>
            body {
                background-color: #000;
            }
            .terminal {
                border: #000 solid 5px;
                font-family: "DejaVu Sans Mono", "Liberation Mono", monospace;
                font-size: 11px;
                color: #f0f0f0;
                background: #000;
            }
            .terminal-cursor {
                color: #000;
                background: #f0f0f0;
            }
        </style>`}
    
        {`<script type="text/javascript">
            $(function() {
                var websocket = new WebSocket("ws://" + window.location.hostname + ":" + window.location.port + "/workspaces/exec/" + prompt("cid"));
                websocket.onopen = function(evt) {
                    var term = new Terminal({
                        cols: 100,
                        rows: 30,
                        screenKeys: true,
                        useStyle: true,
                        cursorBlink: true,
                    });
                    term.on('data', function(data) {
                        websocket.send(data);
                    });
                    term.on('title', function(title) {
                        document.title = title;
                    });
                    term.open(document.getElementById('container-terminal'));
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
                }
            });
        </script>
    </head>
    <div id="container-terminal"></div>`}
      </div>
    )
  }
}

export default Terminal;
