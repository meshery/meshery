/**
 * @typedef {Object} NavigatorSchema
 * @property {string} title
 * @property {string} href
 * @property {string} component
 * @property {string} icon
 * @property {NavigatorSchema[]} children
 */

/**
 * @typedef {Object} UserPrefSchema
 * @property {string} component
 */

/**
 * ExtensionPointSchemaValidator returns the schema validator based on the
 * type passed to the function
 * @param {string} type - Type of Schema validator needed. Valid types are - navigator, userprefs
 */
export default function ExtensionPointSchemaValidator(type) {
  switch (type) {
    case "navigator":
      return NavigatorExtensionSchemaDecoder;
    case "user_prefs":
      return UserPrefsExtensionSchemaDecoder;
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
        title: item.title || "",
        href: prepareHref(item.href),
        component: item.component || "",
        icon: (item.icon && "/api/provider/extension/" + item.icon) || "",
        show: !!item.show,
        children: NavigatorExtensionSchemaDecoder(item.children),
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
      return {
        component: item.component || "",
      };
    });
  }

  return [];
}

function prepareHref(href) {
  if (href.external) return href.uri || "";

  return "/extension" + (href.uri || "");
}
