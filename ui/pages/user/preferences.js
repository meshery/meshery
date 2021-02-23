import UserPreferences from "../../components/UserPreferences";
import { NoSsr, Paper, withStyles } from "@material-ui/core";
import { updatepagepath } from "../../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import { getPath } from "../../lib/path";
import Head from 'next/head';
import dataFetch from '../../lib/data-fetch';

const styles = {
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  }
}

class UserPref extends React.Component {
  constructor(props){
    super(props);
    this.state={};
  }

  async componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });

    await new Promise(resolve => {
      dataFetch('/api/user/stats', {
        credentials: 'same-origin',
        method: 'GET',
        credentials: 'include',
      }, (result) => {
        resolve();
        if (typeof result !== 'undefined') {
          this.setState({
            anonymousStats: result.anonymousUsageStats||false,
            perfResultStats: result.anonymousPerfResults||false,
            startOnZoom: result.startOnZoom||false, //meshmap specific user preferences are not stored in any db as of now
          })
        }
      },
      // Ignore error because we will fallback to default state
      // and to avoid try catch due to async await functions
      resolve);
    })
  }

  render () {
    const {anonymousStats, perfResultStats, startOnZoom}=this.state;
    if(anonymousStats==undefined){
      // Skip rendering till data is not loaded
      return <div></div>
    }
    return (
      <NoSsr>
        <Head>
          <title>Preferences | Meshery</title>
        </Head>
        <Paper className={this.props.classes.paper}>
          {/* {should meshmap specific user preferences be placed along with general preferences or from the remote provider} */}
          <UserPreferences anonymousStats={anonymousStats} perfResultStats={perfResultStats} startOnZoom={startOnZoom}/>
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
)(UserPref));