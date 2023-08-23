import React, { useEffect } from "react";
import { NoSsr, withStyles } from "@material-ui/core";
import MesheryPatterns from "../../../components/MesheryPatterns";
import { updatepagepath } from "../../../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import Head from "next/head";
import { getPath } from "../../../lib/path";

const styles = {
  paper : {
    maxWidth : "90%",
    margin : "auto",
    overflow : "hidden",
  },
};

function Patterns({ updatepagepath }) {
  useEffect(() => {
    console.log(`path: ${getPath()}`);
    updatepagepath({ path : getPath() });
  }, [updatepagepath]);

  return (
    <NoSsr>
      <Head>
        <title>Designs | Meshery</title>
      </Head>
      <MesheryPatterns />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch) });

export default withStyles(styles)(connect(null, mapDispatchToProps)(Patterns));
