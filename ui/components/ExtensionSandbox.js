import React, { useEffect, useState } from 'react';
import normalizeURI from '../utils/normalizeURI';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import LoadingScreen from './LoadingComponents/LoadingComponent';
import {
  useLazyGetExtensionsByTypeQuery,
  useLazyGetFullPageExtensionsQuery,
} from '@/rtk-query/user';
import { useDispatch, useSelector } from 'react-redux';
import { toggleDrawer } from '@/store/slices/mesheryUi';

// Custom Hook for getCapabilities
function useCapabilities(type, cb) {
  const [getExtensionsByType] = useLazyGetExtensionsByTypeQuery();

  useEffect(() => {
    if (type && cb) {
      getExtensionsByType(type)
        .unwrap()
        .then((data) => {
          if (typeof data !== 'undefined') {
            cb(data);
          }
        })
        .catch((err) => {
          console.group('extension error');
          console.error(err);
          console.groupEnd();
        });
    }
  }, [type, cb, getExtensionsByType]);
}

// Custom Hook for getFullPageExtensions
function useFullPageExtensions(cb) {
  const [getFullPageExtensions] = useLazyGetFullPageExtensionsQuery();

  useEffect(() => {
    if (cb) {
      getFullPageExtensions()
        .unwrap()
        .then((data) => {
          cb(data);
        })
        .catch((err) => {
          console.group('extension error');
          console.error(err);
          console.groupEnd();
        });
    }
  }, [cb, getFullPageExtensions]);
}

/**
 * getPath returns the current pathname
 * @returns {string}
 */
function getPath() {
  return window.location.pathname;
}

/**
 * getCapabilities queries the meshery server for the current providers
 * capabilities and returns the decoded capability that matches the
 * given type
 * @param {string} type
 * @param {Function} cb
 */
export function getCapabilities(type, cb) {
  // This function now relies on the useCapabilities Hook
  // Call this from a component that uses the Hook
  console.log('getCapabilities called with:', type, cb);
}

/**
 * getFullPageExtensions queries the meshery server for the current providers
 * capabilities and returns all the extensions names and URIs having full_page type as true
 * @param {Function} cb
 */
export function getFullPageExtensions(cb) {
  console.log('getFullPageExtensions called with:', cb);
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
    if (fext) return fext.component || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const comp = getComponentURIFromPathForNavigator(ext.children, path);
      if (comp) return comp;
    }
  }

  return '';
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
    if (fext) return fext.title || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const title = getComponentURIFromPathForNavigator(ext.children, path);
      if (title) return title;
    }
  }

  return '';
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
    if (fext) return fext.component || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const comp = getComponentURIFromPathForAccount(ext.children, path);
      if (comp) return comp;
    }
  }

  return '';
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
    if (fext) return fext.title || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const title = getComponentURIFromPathForAccount(ext.children, path);
      if (title) return title;
    }
  }

  return '';
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
    if (fext) return fext.component || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const comp = getComponentURIFromPath(ext.children, path);
      if (comp) return comp;
    }
  }

  return '';
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
    if (fext) return fext.title || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const title = getComponentURIFromPath(ext.children, path);
      if (title) return title;
    }
  }

  return '';
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
    if (extension) return extension.isBeta ?? false;
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
    return extensions.map((ext) => ext.component);
  }

  return [];
}

/**
 * getComponentURIFromPathForCollaborator takes in the user_prefs extensions and returns
 * an array of all the component mappings
 * @param {import("../utils/ExtensionPointSchemaValidator").CollaboratorSchema[]} extensions
 * @returns {string[]}
 */
function getComponentURIFromPathForCollaborator(extensions) {
  if (Array.isArray(extensions)) {
    return extensions.map((ext) => ext.component);
  }

  return [];
}

/**
 * createPathForRemoteComponent takes in the name of the component and
 * returns a path for making the http request for that path
 *
 * this path must be registered on the backend
 * @param {string} componentName
 * @returns {string} url
 */
export function createPathForRemoteComponent(componentName) {
  let prefix = '/api/provider/extension';
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
const ExtensionSandbox = React.memo(
  function MemoizedExtensionSandbox({ type, Extension }) {
    const [extension, setExtension] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const { capabilitiesRegistry } = useSelector((state) => state.ui);
    const { isDrawerCollapsed } = useSelector((state) => state.ui);
    const dispatch = useDispatch();

    // Use custom Hooks to fetch data
    useCapabilities(type, () => {
      if (capabilitiesRegistry && capabilitiesRegistry.extensions) {
        try {
          const extensionData = capabilitiesRegistry.extensions[type];
          const processedData = ExtensionPointSchemaValidator(type)(extensionData);
          setExtension(processedData);
          setIsLoading(false);
        } catch {
          setExtension([]);
          setIsLoading(false);
        }
      }
    });

    useFullPageExtensions((extensionData) => {
      if (extensionData) {
        try {
          const processedData = ExtensionPointSchemaValidator('full_page')(extensionData);
          setExtension(processedData);
        } catch (err) {
          console.error('Error processing full page extension data:', err);
          setExtension([]);
        } finally {
          setIsLoading(false);
        }
      } else {
        setExtension([]);
        setIsLoading(false);
      }
    });

    useEffect(() => {
      if (type === 'navigator' && !isDrawerCollapsed) {
        dispatch(toggleDrawer({ isDrawerCollapsed: !isDrawerCollapsed }));
      }

      // Cleanup states on unmount to prevent memory leaks
      return () => {
        setExtension([]);
        setIsLoading(true);
      };
    }, [type, isDrawerCollapsed, dispatch]);

    const renderContent = () => {
      if (isLoading) {
        return type === 'collaborator' ? null : (
          <LoadingScreen animatedIcon="AnimatedMeshery" message="Establishing Remote Connection" />
        );
      }

      switch (type) {
        case 'collaborator': {
          const collaboratorUri = getComponentURIFromPathForCollaborator(extension, getPath());
          return collaboratorUri.map((uri) => (
            <Extension url={createPathForRemoteComponent(uri)} key={uri} />
          ));
        }
        case 'navigator': {
          const navigatorUri = getComponentURIFromPathForNavigator(extension, getPath());
          return navigatorUri ? (
            <Extension url={createPathForRemoteComponent(navigatorUri)} />
          ) : null;
        }
        case 'user_prefs': {
          const userPrefUris = getComponentURIFromPathForUserPrefs(extension);
          return userPrefUris.map((uri) => (
            <Extension url={createPathForRemoteComponent(uri)} key={uri} />
          ));
        }
        case 'account': {
          const accountUri = getComponentURIFromPathForAccount(extension, getPath());
          return accountUri ? <Extension url={createPathForRemoteComponent(accountUri)} /> : null;
        }

        default:
          return null;
      }
    };

    return <>{renderContent()}</>;
  },
  (prevProps, nextProps) =>
    prevProps.type === nextProps.type &&
    prevProps.capabilitiesRegistry === nextProps.capabilitiesRegistry,
);

export default ExtensionSandbox;
