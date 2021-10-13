import MesherySettings from "../components/MesherySettings";
import { NoSsr, withStyles, Paper } from "@material-ui/core";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import Head from 'next/head';
import { bindActionCreators } from 'redux';
import { getPath } from "../lib/path";

const styles = { paper : { maxWidth : '90%',
  margin : 'auto',
  overflow : 'hidden', } };

class Settings extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path : getPath() });
  }

  render () {
    return (
      <NoSsr>
        <Head>
          <title>Settings | Meshery</title>
        </Head>
        <Paper className={this.props.classes.paper}>
          <MesherySettings />
        </Paper>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({ updatepagepath : bindActionCreators(updatepagepath, dispatch) })

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps
)(Settings));