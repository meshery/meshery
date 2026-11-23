import normalizeURI from './normalizeURI';

/**
 * @typedef {Object} NavigatorSchema
 * @property {string} title
 * @property {number} onClickCallback
 * @property {string} href
 * @property {string} component
 * @property {string} icon
 * @property {NavigatorSchema[]} children
 * @property {string} type
 */

/**
 * @typedef {Object} UserPrefSchema
 * @property {string} component
 * @property {string} type
 */

/**
 * @typedef {Object} CollaboratorSchema
 * @property {string} component
 * @property {string} type
 */

/**
 * @typedef {Object} AccountSchema
 * @property {string} title
 * @property {number} onClickCallback
 * @property {string} href
 * @property {string} component
 * @property {AccountSchema[]} children
 * @property {string} type
 */

/**
 * @typedef {Object} FullPageExtensionSchema
 * @property {string} title
 * @property {number} onClickCallback
 * @property {string} href
 * @property {string} component
 * @property {AccountSchema[]} children
 * @property {string} type
 */

/**
 * ExtensionPointSchemaValidator returns the schema validator based on the
 * type passed to the function
 * @param {string} type - Type of Schema validator needed. Valid types are - navigator, userprefs, account
 */
export default function ExtensionPointSchemaValidator(type) {
  switch (type) {
    case 'navigator':
      return NavigatorExtensionSchemaDecoder;
    case 'userPrefs':
      return UserPrefsExtensionSchemaDecoder;
    case 'collaborator':
      return CollaboratorExtensionSchemaDecoder;
    case 'account':
      return AccountExtensionSchemaDecoder;
    default:
      return () => {};
  }
}

/**
 * NavigatorExtensionSchemaDecoder
 * @param {*} content
 * @returns {NavigatorSchema[]}
 */
function NavigatorExtensionSchemaDecoder(content) {
  if (Array.isArray(content)) {
    return content
      .filter((item) => item && typeof item === 'object')
      .map((item) => {
        return {
          title: item.title || '',
          href: prepareHref(item.href),
          component: item.component || '',
          onClickCallback: item?.onClickCallback || 0,
          // Mirror createPathForRemoteComponent's prefix-then-normalizeURI shape
          // so the icon URL never picks up a double slash. Without this, an
          // item.icon starting with "/" concatenated onto a prefix ending in "/"
          // produced "/api/provider/extension//provider/..." which the browser
          // treated as a distinct URL from the single-slash form, causing two
          // fetches (one of which stalled) for the same asset.
          icon:
            (typeof item.icon === 'string' &&
              item.icon &&
              '/api/provider/extension' + normalizeURI(item.icon)) ||
            '',
          show: !!item.show,
          children: NavigatorExtensionSchemaDecoder(item.children),
          full_page: item.full_page,
          isBeta: item.isBeta ?? false,
        };
      });
  }

  return [];
}

/**
 *
 * @param {*} content
 * @returns {UserPrefSchema[]}
 */
function UserPrefsExtensionSchemaDecoder(content) {
  if (Array.isArray(content)) {
    return content
      .filter((item) => item && typeof item === 'object')
      .map((item) => {
        return { component: item.component || '' };
      });
  }

  return [];
}

/**
 *
 * @param {*} content
 * @returns {CollaboratorSchema[]}
 */
function CollaboratorExtensionSchemaDecoder(content) {
  if (Array.isArray(content)) {
    return content
      .filter((item) => item && typeof item === 'object')
      .map((item) => {
        return { component: item.component || '' };
      });
  }

  return [];
}

/**
 * AccountExtensionSchemaDecoder
 * @param {*} content
 * @returns {AccountSchema[]}
 */
function AccountExtensionSchemaDecoder(content) {
  if (Array.isArray(content)) {
    return content
      .filter((item) => item && typeof item === 'object')
      .map((item) => {
        return {
          title: item.title || '',
          href: prepareHref(item.href),
          component: item.component || '',
          onClickCallback: item?.onClickCallback || 0,
          show: !!item.show,
          children: AccountExtensionSchemaDecoder(item.children),
          full_page: item.full_page,
          isHtmlLink: item.link,
        };
      });
  }

  return [];
}

/**
 * Resolves an extension href object or string into a fully-qualified route.
 * Returns an empty string for nullish or invalid href inputs.
 *
 * @param {Object|string|null|undefined} href
 * @returns {string}
 */
function prepareHref(href) {
  if (!href) return '';

  if (typeof href === 'string') {
    return href;
  }

  if (typeof href === 'object') {
    if (typeof href.uri !== 'string' || !href.uri) {
      return '';
    }

    if (href.external) return href.uri;

    return '/extension' + href.uri;
  }

  return '';
}
