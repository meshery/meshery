/* eslint-disable no-unused-vars */

import Extension from "../../components/NavigatorExtension";
import ExtensionSandbox, { getCapabilities, getFullPageExtensions, getComponentTitleFromPath } from "../../components/ExtensionSandbox";
import { NoSsr } from "@material-ui/core";
import { updatepagepath, updatepagetitle, updateExtensionType } from "../../lib/store";
import { connect } from "react-redux";
import Head from "next/head";
import { bindActionCreators } from "redux";
import React, { useEffect, useState } from "react";
import RemoteComponent from "../../components/RemoteComponent";
import { useRouter } from "next/router";


/**
 * getPath returns the current pathname
 * @returns {string}
 */
function getPath() {
  return window.location.pathname;
}

/**
 * extractComponentURI extracts the last part of the
 * given path
 * @param {string} path
 * @returns {string}
 */
function extractComponentURI(path) {
  return path.substring(path.lastIndexOf("/"));
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

class RemoteExtension extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      componentTitle : ''
    }
  }

  componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path : getPath() });
    console.log("pid(useEffect): ", this.props.pid);
    this.renderExtension();
  }

  componentDidUpdate(prevProps) {
    if (this.props.extensionType !== prevProps.extensionType) {
      this.renderExtension();
    }
  }

  renderExtension = () => {
    getFullPageExtensions(extNames => {
      extNames.forEach((ext) => {
        const currentURI = extractComponentURI(getPath());
        console.log("currentURI: ", currentURI);
        console.log("ext: ", ext);
        if (currentURI === ext.uri) {
          console.log("inside")
          console.log("ext")
          this.props.updateExtensionType({ extensionType : ext.name })
          getCapabilities(ext.name, extensions => {
            this.setState({ componentTitle : getComponentTitleFromPath(extensions, getPath()) });
            this.props.updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) });
          });
        }
      })
    });

  }

  render() {
    console.log("props: ", this.props)
    const extensionType = this.props.extensionType;

    return (
      <NoSsr>
        <Head>
          <title>{`${this.state.componentTitle} | Meshery` || ""}</title>
        </Head>
        {
          extensionType &&
          <NoSsr>
            {console.log("ext(UI)", extensionType)}
            {
              (extensionType === 'navigator') ?
                <ExtensionSandbox type={extensionType} Extension={Extension} />
                :
                <ExtensionSandbox type={extensionType} Extension={(url) => RemoteComponent({ url })} />
            }
          </NoSsr>
        }
      </NoSsr>
    )
  }

}

// const RemoteExtension = (props) => {
//   const { updatepagetitle, updatepagepath, pid } = props;
//   // const extensionType = "navigator"
//   console.log("props: ", props)

//   const [componentTitle, setComponentTitle] = useState("");
//   // const [currentURI, setCurrentURI] = useState("");
//   const extensionType = pid;
//   // const updateExtensionType = props.setExtensionType;
//   const router = useRouter();
//   console.log("router: ", router);
//   // const pid = router?.query?.component;
//   console.log("pid (init):", pid);

//   const renderExtension = () => {
//     getFullPageExtensions(extNames => {
//       extNames.forEach((ext) => {
//         const currentURI = extractComponentURI(getPath());
//         console.log("currentURI: ", currentURI);
//         console.log("ext: ", ext);
//         if (currentURI === ext.uri) {
//           console.log("inside")
//           console.log("ext")
//           // updateExtensionType({ extensionType : pid })
//           getCapabilities(ext.name, extensions => {
//             setComponentTitle(getComponentTitleFromPath(extensions, getPath()))
//             updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) })
//           });
//         }
//       })
//     });
//   }


//   // useEffect(() => {
//   //   console.log(`path: ${getPath()}`);
//   //   updatepagepath({ path : getPath() });
//   //   console.log("pid(useEffect): ", pid);
//   //   renderExtension();
//   // }, [])

//   // componentDidUpdate(prevProps, prevState) {
//   //   console.log("inside update");
//   //   console.log("prevState", prevState);
//   //   console.log("this.state.currentURI", this.state.currentURI);

//   //   if (prevState.currentURI !== getPath()) {
//   //     this.setState({ currentURI : getPath() });
//   //     getFullPageExtensions(extNames => {
//   //       extNames.forEach((ext) => {
//   //         const currentURI = extractComponentURI(getPath());
//   //         if (currentURI === ext.uri) {
//   //           this.setState({ extensionType : ext.name })
//   //           getCapabilities(ext.name, extensions => {
//   //             this.setState({ componentTitle : getComponentTitleFromPath(extensions, getPath()) });
//   //             this.props.updatepagetitle({ title : getComponentTitleFromPath(extensions, getPath()) });
//   //           });
//   //           console.log(`path: ${getPath()}`);
//   //           this.props.updatepagepath({ path : getPath() });
//   //         }
//   //       })
//   //     });
//   //   }
//   // }

//   // componentDidMount() {
//   //   this.setState({ currentURI : getPath() });
//   // }

//   return (
//     <NoSsr>
//       <Head>
//         <title>{`${componentTitle} | Meshery` || ""}</title>
//       </Head>
//       {
//         extensionType &&
//           <NoSsr>
//             {console.log("ext(UI)", extensionType)}
//             {
//               (extensionType === 'navigator') ?
//                 <ExtensionSandbox type={extensionType} Extension={Extension} />
//                 :
//                 <ExtensionSandbox type={extensionType} Extension={(url) => RemoteComponent({ url })} />
//             }
//           </NoSsr>
//       }
//     </NoSsr>
//   );
// }

const mapStateToProps = (state) => ({
  extensionType : state.get('extensionType'),
});

const mapDispatchToProps = (dispatch) => ({
  updatepagepath : bindActionCreators(updatepagepath, dispatch),
  updatepagetitle : bindActionCreators(updatepagetitle, dispatch),
  updateExtensionType : bindActionCreators(updateExtensionType, dispatch)
});

export default connect(mapStateToProps, mapDispatchToProps)(RemoteExtension);

// export async function getStaticProps({ params }) {
//   console.log("param: ", params);
//   let ext, exts = [];
//   // getFullPageExtensions(extNames => extNames.forEach((ext) => exts.push(ext)));
//   console.log("exts: ", exts);
//   ext = params.component;
//   console.log("ext: ", ext);
//   return {
//     props : {
//       pid : ext
//     }
//   }
// }
// export const getStaticPaths = async () => {
//   const pathsWithParams = [
//     { params : { component : 'meshmap' } },
//     { params : { component : 'account' } },
//   ]

//   return {
//     paths : pathsWithParams,
//     fallback : true
//   }
// }

