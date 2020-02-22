import MesherySettings from "../components/MesherySettings";
import { NoSsr } from "@material-ui/core";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import { getPath } from "../lib/path";

class Settings extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render () {
    return (
      <NoSsr>
        <MesherySettings />
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({
    updatepagepath: bindActionCreators(updatepagepath, dispatch)
  })

export default connect(
    null,
    mapDispatchToProps
  )(Settings);