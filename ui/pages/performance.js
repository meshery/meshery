import { NoSsr, Paper } from "@material-ui/core";
import MesheryPerformanceComponent from "../components/MesheryPerformanceComponent";
import { updatepagepath } from "../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from "../lib/path";
import { withStyles } from '@material-ui/core/styles';

const style = () => ({

  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',

  }
})

class Performance extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render () {
    return (
      <NoSsr>
        <Paper className={this.props.classes.paper}>
          <Head>
            <title>Performance | Meshery</title>
          </Head>
          <MesheryPerformanceComponent />
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
)(withStyles(style)(Performance));