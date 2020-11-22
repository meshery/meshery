import React from 'react';
import MesheryPlayComponent from '../components/MesheryPlayComponent';
import { NoSsr, Paper, withStyles } from "@material-ui/core";
import Head from 'next/head';
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import { getPath } from "../lib/path";
import { withStyles } from '@material-ui/core/styles';

const style = () => ({
  
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',

  }
})

const styles = {
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  }
}
class Manage extends React.Component {
  static getInitialProps = ({ query }) => {
    return { query };
  }
  constructor(props) {
    super(props);
    if (!props.query.adapter) {
      const urlParams = new URLSearchParams(window.location.search);
      this.props.query.adapter = urlParams.get("adapter");
    }
  }
  componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  componentDidUpdate() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render() {
    return (
      <NoSsr>
        <Head>
          <title>Management | Meshery </title>
        </Head>
        <Paper className={this.props.classes.paper}>
          <MesheryPlayComponent adapter={this.props.query.adapter} />
        </Paper>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => {
  return {
    updatepagepath: bindActionCreators(updatepagepath, dispatch)
  }
}

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps
)(withStyles(style)(Manage)));
