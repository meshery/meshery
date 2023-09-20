import React,{ useEffect } from 'react';
import MesherySMIResults from "../components/MesherySMIResults";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from "../lib/path";
import { Paper, withStyles } from '@material-ui/core';

const styles = { paper : { maxWidth : '90%',
  margin : 'auto',
  overflow : 'hidden', } };

const SMIResults=(props) => {

  useEffect(() => {
    console.log(`path: ${getPath()}`);
    props.updatepagepath({ path : getPath() });

  },[props.updatepagepath])



  return (
    <React.Fragment>
      <Head>
        <title>SMI Results | Meshery</title>
      </Head>
      <Paper className={props.classes.paper}>
        <MesherySMIResults />
      </Paper>
    </React.Fragment>
  );

}

const mapDispatchToProps = dispatch => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch) })

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps
)(SMIResults));