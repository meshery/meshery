import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import normalizeURI from '../utils/normalizeURI';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import LoadingScreen from './LoadingComponents/LoadingComponent';
import { toggleDrawer } from '@/store/slices/mesheryUi';
import type {
  NavigatorSchema,
  AccountSchema,
  UserPrefSchema,
  CollaboratorSchema,
  FullPageExtensionSchema,
} from '../utils/ExtensionPointSchemaValidator';
import type { RootState } from '../store';

/**
 * getPath returns the current pathname
 */
function getPath(): string {
  return window.location.pathname;
}

/**
 * getComponentURIFromPathForNavigator takes in the navigator extensions and the current
 * path and searches recursively for the matching component
 *
 * If there are duplicate uris then the component for first match will be returned
 */
function getComponentURIFromPathForNavigator(extensions: NavigatorSchema[], path: string): string {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.component || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const comp = getComponentURIFromPathForNavigator(ext.children || [], path);
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
 */
export function getComponentTitleFromPathForNavigator(
  extensions: NavigatorSchema[],
  path: string,
): string {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.title || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const title = getComponentTitleFromPathForNavigator(ext.children || [], path);
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
 */
function getComponentURIFromPathForAccount(extensions: AccountSchema[], path: string): string {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.component || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const comp = getComponentURIFromPathForAccount(ext.children || [], path);
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
 */
export function getComponentTitleFromPathForAccount(
  extensions: AccountSchema[],
  path: string,
): string {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.title || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const title = getComponentTitleFromPathForAccount(ext.children || [], path);
      if (title) return title;
    }
  }

  return '';
}

/**
 * getComponentTitleFromPath takes in the extensions and the current
 * path and searches recursively for the matching component and returns title
 *
 * If there are duplicate uris then the component for first match will be returned
 */
export function getComponentTitleFromPath(
  extensions: FullPageExtensionSchema[],
  path: string,
): string {
  path = normalizeURI(path);

  if (Array.isArray(extensions)) {
    const fext = extensions.find((item) => item?.href === path);
    if (fext) return fext.title || '';

    // If not found then start searching in the child of each extension
    for (const ext of extensions) {
      const title = getComponentTitleFromPath(ext.children || [], path);
      if (title) return title;
    }
  }

  return '';
}

/**
 * getComponentIsBetaFromPath takes in the extensions and the current
 * path and searches for the matching component and returns isBeta
 */
export function getComponentIsBetaFromPath(
  extensions: FullPageExtensionSchema[],
  path: string,
): boolean {
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
 */
function getComponentURIFromPathForUserPrefs(extensions: UserPrefSchema[]): string[] {
  if (Array.isArray(extensions)) {
    return extensions.map((ext) => ext.component);
  }

  return [];
}

/**
 * getComponentURIFromPathForCollaborator takes in the collaborator extensions and returns
 * an array of all the component mappings
 */
function getComponentURIFromPathForCollaborator(extensions: CollaboratorSchema[]): string[] {
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
 */
export function createPathForRemoteComponent(componentName: string): string {
  const prefix = '/api/provider/extension';
  return prefix + normalizeURI(componentName);
}

type ExtensionType = 'navigator' | 'user_prefs' | 'account' | 'collaborator';

interface ExtensionComponentProps {
  url: string;
}

interface ExtensionSandboxProps {
  type: ExtensionType;
  Extension: React.ComponentType<ExtensionComponentProps>;
}

/**
 * ExtensionSandbox takes in an extension and its type and will handle the internal mapping
 * for the uris and components by querying the meshery server for providers capabilities
 *
 * Only four "types" are supported by the sandbox:
 *  1. navigator - for navigator extensions
 *  2. user_prefs - for user preference extension
 *  3. account - for user account extension
 *  4. collaborator - for collaborator extension
 */
const ExtensionSandbox = React.memo<ExtensionSandboxProps>(
  function MemoizedExtensionSandbox({ type, Extension }) {
    const [extension, setExtension] = useState<
      NavigatorSchema[] | UserPrefSchema[] | AccountSchema[] | CollaboratorSchema[]
    >([]);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const { capabilitiesRegistry, isDrawerCollapsed } = useSelector((state: RootState) => state.ui);
    const dispatch = useDispatch();

    useEffect(() => {
      if (type === 'navigator' && !isDrawerCollapsed) {
        dispatch(toggleDrawer({ isDrawerCollapsed: !isDrawerCollapsed }));
      }

      if (capabilitiesRegistry && capabilitiesRegistry.extensions) {
        try {
          const extensionData = (capabilitiesRegistry.extensions as any)[type];
          const processedData = ExtensionPointSchemaValidator(type)(extensionData);
          setExtension(processedData);
          setIsLoading(false);
        } catch {
          setExtension([]);
          setIsLoading(false);
        }
      } else {
        setIsLoading(true);
      }

      // necessary to cleanup states on each unmount to prevent memory leaks and unwanted clashes between extension points
      return () => {
        setExtension([]);
        setIsLoading(true);
      };
    }, [type, capabilitiesRegistry, isDrawerCollapsed, dispatch]);

    const renderContent = (): React.ReactNode => {
      if (isLoading) {
        return type === 'collaborator' ? null : (
          <LoadingScreen animatedIcon="AnimatedMeshery" message="Establishing Remote Connection" />
        );
      }

      switch (type) {
        case 'collaborator': {
          const collaboratorUri = getComponentURIFromPathForCollaborator(
            extension as CollaboratorSchema[],
          );
          return collaboratorUri.map((uri) => (
            <Extension url={createPathForRemoteComponent(uri)} key={uri} />
          ));
        }
        case 'navigator': {
          const navigatorUri = getComponentURIFromPathForNavigator(
            extension as NavigatorSchema[],
            getPath(),
          );
          return navigatorUri ? (
            <Extension url={createPathForRemoteComponent(navigatorUri)} />
          ) : null;
        }
        case 'user_prefs': {
          const userPrefUris = getComponentURIFromPathForUserPrefs(extension as UserPrefSchema[]);
          return userPrefUris.map((uri) => (
            <Extension url={createPathForRemoteComponent(uri)} key={uri} />
          ));
        }
        case 'account': {
          const accountUri = getComponentURIFromPathForAccount(
            extension as AccountSchema[],
            getPath(),
          );
          return accountUri ? <Extension url={createPathForRemoteComponent(accountUri)} /> : null;
        }

        default:
          return null;
      }
    };

    return <>{renderContent()}</>;
  },
  (prevProps, nextProps) => prevProps.type === nextProps.type,
);

export default ExtensionSandbox;
