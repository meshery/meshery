import MesherySettings from "../components/MesherySettings";
import { NoSsr } from "@material-ui/core";
import Paper from '@material-ui/core/Paper';
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import Head from 'next/head';
import { bindActionCreators } from 'redux';
import { getPath } from "../lib/path";
import { withStyles } from '@material-ui/core/styles';

const style = () => ({

  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',

  }
})


class Settings extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render () {
    return (
      <NoSsr>
        <Paper className={this.props.classes.paper}>
          <Head>
            <title>Settings | Meshery</title>
          </Head>
          <MesherySettings />
        </Paper>
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
)(withStyles(style)(Settings));