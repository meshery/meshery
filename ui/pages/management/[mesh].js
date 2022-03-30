import React from "react";
import { Container, makeStyles, NoSsr, Typography } from "@material-ui/core";
import Head from "next/head";

/**
 * getPath returns the current pathname
 * @returns {string}
 */
function getPath() {
  return window.location.pathname;
}

/**
 * extractComponentName extracts the last part of the
 * given path
 * @param {string} path
 * @returns {string}
 */
function extractComponentName(path) {
  return path.substring(path.lastIndexOf("/") + 1);
}

/**
 * capitalize capitalizes the given string and returns the modified
 * string
 *
 * If the given parameter is not sting then it will return an empty
 * string
 * @param {string} string
 *
 * @returns {string}
 */
function capitalize(string) {
  if (typeof string === "string") return string.charAt(0).toUpperCase() + string.slice(1);
  return "";
}

const useStyles = makeStyles(() => ({
  container : {
    textAlign : "center",
    transform : "translateY(-20%)"
  },
  wrapper : {
    display : "flex",
    alignItems : "center",
    height : "100%"
  }
}))

function Mesh() {
  const name = extractComponentName(getPath());
  const classes =  useStyles()

  return (
    <NoSsr>
      <Head>
        <title>{capitalize(name)} Management</title>
      </Head>
      <NoSsr>
        <div className={classes.wrapper}>
          <Container className={classes.container}>
            <img src="/static/img/broken-link.png" alt="broken link freePik" height="150px" />
            <Typography variant="h4">
            Broken Adapter
            </Typography>
            <Typography>
            Your Adapter is not available
            </Typography>
          </Container>
        </div>
        {/* <MesheryMeshInterface adapter={name}/> */}
      </NoSsr>
    </NoSsr>
  );
}

export default Mesh;
