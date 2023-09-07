import React, {useEffect} from "react";
import { NoSsr, withStyles } from "@material-ui/core";
import { connect } from "react-redux";
import Head from 'next/head';
import MesheryConnections from "../../components/connections";

const styles = { 
  paper : { 
  maxWidth : '90%',
  margin : 'auto',
  overflow : 'hidden', } };


function Connections() {

  useEffect(() => {
  },[]);

    return (
      <NoSsr>
        <Head>
          <title>Connections | Meshery</title>
        </Head>
        <MesheryConnections />
      </NoSsr>
    );
  }


export default withStyles(styles)(connect(
  null,
)(Connections));