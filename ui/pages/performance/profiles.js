import { NoSsr, Paper, withStyles } from "@material-ui/core";
import PerformanceProfiles from "../../components/PerformanceProfiles";
import { updatepagepath } from "../../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux'
import Head from 'next/head';
import { getPath } from "../../lib/path";

const styles = {
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  }
}

class Results extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render () {
    return (
      <NoSsr>
        <Head>
          <title>Results | Meshery</title>
        </Head>
        <Paper className={this.props.classes.paper}>
          <PerformanceProfiles />
        </Paper>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch)
})

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps
)(Results));