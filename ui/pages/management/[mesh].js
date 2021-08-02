import React from "react";
import { NoSsr } from "@material-ui/core";
import Head from "next/head";
import MesheryMeshInterface from "../../components/MesheryMeshInterface";

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

function Mesh() {
  const name = extractComponentName(getPath())

  return (
    <NoSsr>
      <Head>
        <title>{capitalize(name)} Management</title>
      </Head>
      <NoSsr>
        <MesheryMeshInterface adapter={name}/>
      </NoSsr>
    </NoSsr>
  );
}

export default Mesh;
