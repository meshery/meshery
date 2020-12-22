import React, { useEffect, useState } from "react";
import normalizeURI from "../utils/normalizeURI";
import dataFetch from "../lib/data-fetch";
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator'

/**
 * getPath returns the current pathname
 * @returns {string}
 */
function getPath() {
  return window.location.pathname;
}

function getCapabilities(type, cb) {
  dataFetch(
    "/api/provider/capabilities",
    {
      credentials: "same-origin",
      method: "GET",
      credentials: "include",
    },
    (result) => {
      if (result) {
        console.log(result)
        cb(ExtensionPointSchemaValidator(type)(result?.extensions[type]));
      }
    },
    (err) => {
      console.group('extension error')
      console.error(err)
      console.groupEnd()
    }
  );
}

/**
 *
 * @param {import("../utils/ExtensionPointSchemaValidator").UserPrefSchema[] |
 * import("../utils/ExtensionPointSchemaValidator").NavigatorSchema[]
 * } extension
 * @param {string} path
 * @returns {string}
 */
function getComponentURIFromPath(extension, path) {
  path = normalizeURI(path);

  if (Array.isArray(extension)) {
    const ext = extension.find((item) => item?.href === path);
    return ext ? ext.component : "";
  }

  return "";
}

/**
 * createPathForRemoteComponent takes in the name of the component and
 * returns a path for making the http request for that path
 *
 * this path must be registered on the backend
 * @param {string} componentName
 * @returns {string} url
 */
function createPathForRemoteComponent(componentName) {
  let prefix = "/api/provider/extension";
  return prefix + normalizeURI(componentName);
}

function ExtensionSandbox({ type, Extension }) {
  const [extension, setExtension] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    getCapabilities(type, (data) => {
      setExtension(data);
      setIsLoading(false);
    });
  }, []);

  return isLoading ? null : (
    <Extension url={createPathForRemoteComponent(getComponentURIFromPath(extension, getPath()))} />
  );
}

export default ExtensionSandbox;
