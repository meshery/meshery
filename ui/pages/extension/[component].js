import Extension from "../../components/NavigatorExtension";
import ExtensionSandbox from "../../components/ExtensionSandbox"
import { NoSsr } from "@material-ui/core";
import { updatepagepath, updatepagetitle } from "../../lib/store";
import { connect } from "react-redux";
import Head from "next/head";
import { bindActionCreators, compose } from "redux";

 
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

/**
  * converts "Meshmap" to "MeshMap"
  *
  * @param {string} string
  *
  * @returns {string}
*/
const capitalizeMeshmap = (string) => string == "Meshmap" ? "MeshMap" : string;

const getComponentName = compose(capitalizeMeshmap,capitalize,extractComponentName,getPath);

class Settings extends React.Component {
  componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
    this.props.updatepagetitle({ title: getComponentName() });
  }

  render() {
    return (
      <NoSsr>
        <Head>
          <title>{getComponentName()}</title>
        </Head>
        <NoSsr>
          <ExtensionSandbox type="navigator" Extension={Extension} />
        </NoSsr>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
  updatepagetitle: bindActionCreators(updatepagetitle, dispatch),
});

export default connect(null, mapDispatchToProps)(Settings);
