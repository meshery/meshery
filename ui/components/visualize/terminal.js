import React from 'react';
import { Terminal as Term } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { SearchAddon } from 'xterm-addon-search';
import { withStyles } from '@material-ui/core/styles';
import 'xterm/css/xterm.css';
import { AppBar, IconButton, Toolbar } from '@material-ui/core';
import SearchIcon from '@material-ui/icons/Search';
import TextField from '@material-ui/core/TextField';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpandLessIcon from '@material-ui/icons/ExpandLess';

const useStyles = () => ({
  root: {
  },
  appBar: {
    flexGrow: 1,
    maxHeight: '40px'
  },
  terminal: {
    maxHeight: '120px', 
  },
  
});

class Terminal extends React.Component {
  constructor(props) {
    super(props);
    this.line = ""
    this.entries = []
    this.fitAddon = new FitAddon();
    this.SearchAddon = new SearchAddon();
    this.state = {
      searchTerm: ""
    }
  }

  getContent = async function (term, endpoint) {
    var url = `http://127.0.0.1:8001/${endpoint}`;

    await fetch(url, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
    }).then(function (response) {
      return response.text();
    })
      .then(function (data) {
        term.write(`${data}\r\n$ `);
      })
      .catch(err => {
        term.write(`${err}\r\n$ `);
      });
  };

  componentDidMount() {
    this.mountTerminal();
  }

  mountTerminal() {
    this.term = new Term({
      cursorBlink: true,
      fontFamily: '"Roboto Mono", "Courier New", "Courier", monospace',
      fontSize: 12,
      scrollback: 10000,
    });
    this.term.loadAddon(this.fitAddon);
    this.term.loadAddon(this.SearchAddon);
    this.term.open(document.getElementById('terminal'));
    this.fitAddon.fit();
    this.term.write('$ ');

    this.term.prompt = () => {
      if (this.line) {
        this.getContent(this.term, this.line);
      }
    };
    // this.term.prompt();

    this.term.onKey(key => {
      const char = key;
      if (char.domEvent.key === "Enter") {
        if (this.line) {
          this.entries.push(this.line);
          this.term.write("\r\n");
          this.term.prompt();
          this.line = "";
        }
      } else if (char.domEvent.key === "Backspace") {
        if (this.line) {
          this.line = this.line.slice(0, this.line.length - 1);
          this.term.write("\b \b");
        }
      } else {
        this.line += char.key;
        this.term.write(char.key);
      }
    });
  }

  changeSearch = event => {
    this.setState({searchTerm: event.target.value});
  }
  searchLogs(){
    const { searchTerm } = this.state
    this.SearchAddon.findNext(searchTerm);
  }

  findNext(){
    const { searchTerm } = this.state
    this.SearchAddon.findNext(searchTerm);
  }

  findPrevious(){
    const { searchTerm } = this.state
    this.SearchAddon.findPrevious(searchTerm);
  }

  render() {
    const { classes } = this.props;
    const { searchTerm } = this.state;
    return (
      <div className={classes.root}>
        <AppBar position="static" className={classes.appBar}>
          <Toolbar>
            <div className={classes.search}>
              <TextField
                className={classes.searchBox}
                placeholder="search..."
                value={searchTerm}
                onChange={this.changeSearch}
                size="small"
              />
              <IconButton onClick={this.searchLogs.bind(this)} size="small">
                <SearchIcon />
              </IconButton>
              <IconButton onClick={this.findNext.bind(this)} size="small">
                <ExpandMoreIcon />
              </IconButton>
              <IconButton onClick={this.findPrevious.bind(this)} size="small">
                <ExpandLessIcon />
              </IconButton>
            </div>
          </Toolbar>
        </AppBar>
        <div id="terminal" className={classes.terminal} />
      </div>
    )
  }
}

export default withStyles(useStyles, { withTheme: true })(Terminal);