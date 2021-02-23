import React, { useEffect, useState } from "react";
import normalizeURI from "../utils/normalizeURI";
import dataFetch from "../lib/data-fetch";
import ExtensionPointSchemaValidator from "../utils/ExtensionPointSchemaValidator";

/**
 * getPath returns the current pathname
 * @returns {string}
 */
function getPath() {
  return window.location.pathname;
}

/**
 * getCapabilities queries the meshery server for the current providers
 * capabilities and returns the decoded capability that mathes the 
 * given type
 * @param {string} type 
 * @param {Function} cb 
 */
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
        cb(ExtensionPointSchemaValidator(type)(result?.extensions[type]));
      }
    },
    (err) => {
      console.group("extension error");
      console.error(err);
      console.groupEnd();
    }
  );
}

/**
 * getComponentURIFromPathForNavigator takes in teh navigator extensions and the current
 * path and searches recursively for the matching component
 * 
 * If there are duplicate uris then the component for first match will be returned
 * @param {import("../utils/ExtensionPointSchemaValidator").NavigatorSchema[]} extensions
 * @param {string} path
 * @returns {string}
 */
function getComponentURIFromPathForNavigator(extensions, path) {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.component || "";

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const comp = getComponentURIFromPathForNavigator(ext.children, path);
      if (comp) return comp;
    }
  }

  return "";
}

/**
 * getComponentURIFromPathForUserPrefs takes in the user_prefs extensions and returns
 * an array of all the component mappings
 * @param {import("../utils/ExtensionPointSchemaValidator").UserPrefSchema[]} extensions
 * @returns {string[]}
 */
function getComponentURIFromPathForUserPrefs(extensions) {

  if (Array.isArray(extensions)) {
    return extensions.map(ext => ext.component)
  }

  return []
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

/**
 * ExtensionSandbox takes in an extension and it's type and will handle the internal mapping
 * for the uris and components by querying the meshery server for providers capabilities
 * 
 * Only two "types" are supported by the sandbox:
 *  1. navigator - for navigator extensions
 *  2. user_prefs - for user preference extension
 * @param {{ type: "navigator" | "user_prefs", Extension: JSX.Element }} props 
 */
function ExtensionSandbox({ type, Extension }) {
  const [extensions, setExtensions] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    getCapabilities(type, (data) => {
      setExtensions(data);
      setIsLoading(false);
    });
  }, []);

  if (type === "navigator") {
    return isLoading ? null : (
      <Extension url={createPathForRemoteComponent(getComponentURIFromPathForNavigator(extensions, getPath()))} />
    );
  }

  if (type === "user_prefs") {
    return isLoading ? null : (
      getComponentURIFromPathForUserPrefs(extensions).map(uri => {
        return <Extension url={createPathForRemoteComponent(uri)} />
      })
    );
  }

  return null
}

export default ExtensionSandbox;
