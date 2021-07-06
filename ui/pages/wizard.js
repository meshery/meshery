import React from 'react';
<<<<<<< HEAD
import { withStyles, NoSsr } from '@material-ui/core';
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import Head from 'next/head';

import { updatepagepath, updatepagetitle } from "../lib/store";
import { getPath } from "../lib/path";
import ConfigurationWizard from '../components/ConfigurationWizard'

const styles = {
=======
import { updatepagepath, updatepagetitle } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import Head from 'next/head';
import { getPath } from "../lib/path";
import { Paper, withStyles, NoSsr } from '@material-ui/core';

const styles = {
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  }
>>>>>>> dec1c397 (Connect PopUp to config wizard path)
}

class Wizard extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
    this.props.updatepagetitle({ title: 'Configuration Wizard' });
  }

  render () {
    return (
      <NoSsr>
        <Head>
          <title>Configuration Wizard | Meshery</title>
        </Head>
<<<<<<< HEAD
        <ConfigurationWizard/>
=======
        <Paper className={this.props.classes.paper}>
          <h1>Wizard</h1>
        </Paper>
>>>>>>> dec1c397 (Connect PopUp to config wizard path)
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
  updatepagetitle: bindActionCreators(updatepagetitle, dispatch),
})

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps
)(Wizard));