import { useEffect } from "react";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import ConfigWizComp from "../components/ConfigurationWizard/index.js"
import { updatepagepath } from "../lib/store";
import { getPath } from "../lib/path";
import { NoSsr } from "@material-ui/core";

const ConfigurationWizard = (props) => {

  useEffect(() => {
    console.log(`path: ${getPath()}`);
    props.updatepagepath({ path: getPath() });
  },[]) 
    
  return (
    <NoSsr>
      <ConfigWizComp />
    </NoSsr>
  )
}

const mapDispatchToProps = dispatch => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch)
})

export default connect(
  null,
  mapDispatchToProps
)(ConfigurationWizard);
