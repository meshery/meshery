import UserPreference from "../components/UserPreference";
import { NoSsr } from "@material-ui/core";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import { getPath } from "../lib/path";

class UserPref extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render () {
    return (
      <NoSsr>
        <UserPreference />
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
  )(UserPref);