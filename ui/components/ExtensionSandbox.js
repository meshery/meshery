import React, { useEffect, useState } from "react";
import { toggleDrawer } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
// import { CircularProgress, Typography } from "@material-ui/core";
import normalizeURI from "../utils/normalizeURI";
import dataFetch from "../lib/data-fetch";
import ExtensionPointSchemaValidator from "../utils/ExtensionPointSchemaValidator";
import LoadingScreen from "./LoadingComponents/LoadingComponent";

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
export function getCapabilities(type, cb) {
  dataFetch(
    "/api/provider/capabilities",
    {
      method : "GET",
      credentials : "include",
    },
    (result) => {
      if (typeof result !== "undefined") {
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
 * getFullPageExtensions queries the meshery server for the current providers
 * capabilities and returns all the extensions names and URIs having full_page type as true
 * @param {Function} cb
 */
export function getFullPageExtensions(cb) {
  let extNames = [];
  dataFetch(
    "/api/provider/capabilities",
    {
      method : "GET",
      credentials : "include",
    },
    (result) => {
      for (var key of Object.keys(result?.extensions)) {
        if (Array.isArray(result?.extensions[key])) {
          result?.extensions[key].forEach((comp) => {
            if (comp?.type === "full_page") {
              let ext = {
                name : key,
                uri : comp?.href?.uri
              }
              extNames.push(ext)
            }
          })
        }
      }
      cb(extNames)
    },
    (err) => {
      console.group("extension error");
      console.error(err);
      console.groupEnd();
    }
  );
}

/**
 * getComponentURIFromPathForNavigator takes in the navigator extensions and the current
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
 * getComponentTitleFromPathForNavigator takes in the navigator extensions and the current
 * path and searches recursively for the matching component and returns title
 *
 * If there are duplicate uris then the component for first match will be returned
 * @param {import("../utils/ExtensionPointSchemaValidator").NavigatorSchema[]} extensions
 * @param {string} path
 * @returns {string}
 */
export function getComponentTitleFromPathForNavigator(extensions, path) {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.title || "";

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const title = getComponentURIFromPathForNavigator(ext.children, path);
      if (title) return title;
    }
  }

  return "";
}

/**
 * getComponentURIFromPathForAccount takes in the account extensions and the current
 * path and searches recursively for the matching component
 *
 * If there are duplicate uris then the component for first match will be returned
 * @param {import("../utils/ExtensionPointSchemaValidator").AcccountSchema[]} extensions
 * @param {string} path
 * @returns {string}
 */
function getComponentURIFromPathForAccount(extensions, path) {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.component || "";

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const comp = getComponentURIFromPathForAccount(ext.children, path);
      if (comp) return comp;
    }
  }

  return "";
}

/**
 * getComponentTitleFromPathForAccount takes in the account extensions and the current
 * path and searches recursively for the matching component and returns title
 *
 * If there are duplicate uris then the component for first match will be returned
 * @param {import("../utils/ExtensionPointSchemaValidator").AccountSchema[]} extensions
 * @param {string} path
 * @returns {string}
 */
export function getComponentTitleFromPathForAccount(extensions, path) {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.title || "";

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const title = getComponentURIFromPathForAccount(ext.children, path);
      if (title) return title;
    }
  }

  return "";
}

/**
 * getComponentURIFromPath takes in the extensions and the current
 * path and searches recursively for the matching component
 *
 * If there are duplicate uris then the component for first match will be returned
 * @param {import("../utils/ExtensionPointSchemaValidator").FullPageExtensionSchema[]} extensions
 * @param {string} path
 * @returns {string}
 */
function getComponentURIFromPath(extensions, path) {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.component || "";

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const comp = getComponentURIFromPath(ext.children, path);
      if (comp) return comp;
    }
  }

  return "";
}

/**
 * getComponentTitleFromPath takes in the extensions and the current
 * path and searches recursively for the matching component and returns title
 *
 * If there are duplicate uris then the component for first match will be returned
 * @param {import("../utils/ExtensionPointSchemaValidator").FullPageExtensionSchema[]} extensions
 * @param {string} path
 * @returns {string}
 */
export function getComponentTitleFromPath(extensions, path) {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.title || "";

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const title = getComponentURIFromPath(ext.children, path);
      if (title) return title;
    }
  }

  return "";
}


/**
 * getComponentIsBetaFromPath takes in the extensions and the current
 * path and searches for the matching component and returns isBeta
 *
 * @param {import("../utils/ExtensionPointSchemaValidator").FullPageExtensionSchema[]} extensions
 * @param {string} path
 * @returns {boolean}
 */
export function getComponentIsBetaFromPath(extensions, path) {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const extension = extensions.find((item) => item?.href === path);
    if (extension) return extension.isBeta ?? false
  }

  return false;
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
 * getComponentURIFromPathForCollaborator takes in the user_prefs extensions and returns
 * an array of all the component mappings
 * @param {import("../utils/ExtensionPointSchemaValidator").CollaboratorSchema[]} extensions
 * @returns {string[]}
 */
function getComponentURIFromPathForCollaborator(extensions) {

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
 *  3. account - for user account extension
 *  4. collaborator - for collaborator extension
 * @param {{ type: "navigator" | "user_prefs" | "account" | "collaborator", Extension: JSX.Element }} props
 */
function ExtensionSandbox({ type, Extension, isDrawerCollapsed, toggleDrawer, capabilitiesRegistry }) {
  const [extension, setExtension] = useState([]);
  const [isLoading, setIsLoading] = useState(true);


  useEffect(() => {
    if (type === "navigator" && !isDrawerCollapsed) {
      toggleDrawer({ isDrawerCollapsed : !isDrawerCollapsed });
    }
    if (capabilitiesRegistry) {
      const data = ExtensionPointSchemaValidator(type)(capabilitiesRegistry?.extensions[type]);
      if (data !== undefined) {
        setExtension(data);
        setIsLoading(false);
      }
    }
    // necessary to cleanup states on each unmount to prevent memory leaks and unwanted clashes between extension points
    return () => {
      setExtension([]);
      setIsLoading(true);
    }
  }, []);


  useEffect(() => {
    if (type === "navigator" && !isDrawerCollapsed) {
      toggleDrawer({ isDrawerCollapsed : !isDrawerCollapsed });
    }
    if (capabilitiesRegistry) {
      const data = ExtensionPointSchemaValidator(type)(capabilitiesRegistry?.extensions[type]);
      if (data !== undefined) {
        setExtension(data);
        setIsLoading(false);
      }
    }
    // necessary to cleanup states on each unmount to prevent memory leaks and unwanted clashes between extension points
    return () => {
      setExtension([]);
      setIsLoading(true);
    }
  }, [type]);

  return (
    <>
      {
        (
          isLoading ? (
            <LoadingScreen animatedIcon="AnimatedMeshery" message="Establishing Remote Connection" />
          ) :
            (
              (type === "navigator")?
                (
                  <Extension url={createPathForRemoteComponent(getComponentURIFromPathForNavigator(extension, getPath()))} />
                )
                : (type === "user_prefs")?
                  (
                    getComponentURIFromPathForUserPrefs(extension).map(uri => {
                      return <Extension url={createPathForRemoteComponent(uri)} key={uri} />
                    })
                  )
                  : (type === "collaborator")?
                    (
                      getComponentURIFromPathForCollaborator(extension).map(uri => {
                        return <Extension url={createPathForRemoteComponent(uri)} key={uri} />
                      })
                    )
                    : (type === "account")?
                      (
                        <Extension url={createPathForRemoteComponent(getComponentURIFromPathForAccount(extension, getPath()))} />
                      )
                      : null
            )
        )
      }
    </>
  )
}

const mapDispatchToProps = (dispatch) => ({
  toggleDrawer : bindActionCreators(toggleDrawer, dispatch),
});

const mapStateToProps = (state) => ({
  isDrawerCollapsed : state.get("isDrawerCollapsed"),
  capabilitiesRegistry : state.get("capabilitiesRegistry"),
});

export default connect(mapStateToProps, mapDispatchToProps)(ExtensionSandbox);
