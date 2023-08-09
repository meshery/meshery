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
    case "navigator":
      return NavigatorExtensionSchemaDecoder;
    case "user_prefs":
      return UserPrefsExtensionSchemaDecoder;
    case "collaborator":
      return CollaboratorExtensionSchemaDecoder;
    case "account":
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
    return content.map((item) => {
      return {
        title : item.title || "",
        href : prepareHref(item.href),
        component : item.component || "",
        onClickCallback : item?.on_click_callback || 0,
        icon : (item.icon && "/api/provider/extension/" + item.icon) || "",
        show : !!item.show,
        children : NavigatorExtensionSchemaDecoder(item.children),
        full_page : item.full_page,
        isBeta : item.isBeta ?? false
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
    return content.map((item) => {
      return { component : item.component || "", };
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
  console.log("content", content)
  if (Array.isArray(content)) {
    return content.map((item) => {
      return { component : item.component || "", };
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
    return content.map((item) => {
      return {
        title : item.title || "",
        href : prepareHref(item.href),
        component : item.component || "",
        onClickCallback : item?.on_click_callback || 0,
        show : !!item.show,
        children : AccountExtensionSchemaDecoder(item.children),
        full_page : item.full_page
      };
    });
  }

  return [];
}

function prepareHref(href) {
  if (href.external) return href.uri || "";

  return "/extension" + (href.uri || "");
}
