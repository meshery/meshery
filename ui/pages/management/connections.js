import React from "react";
import { NoSsr, withStyles } from "@material-ui/core";
import { connect } from "react-redux";
import Head from 'next/head';

const styles = { paper : { maxWidth : '90%',
  margin : 'auto',
  overflow : 'hidden', } };

class Connections extends React.Component {
  componentDidMount () {
  }

  render () {
    return (
      <NoSsr>
        <Head>
          <title>Connections | Meshery</title>
        </Head>
      </NoSsr>
    );
  }
}


export default withStyles(styles)(connect(
  null,
)(Connections));