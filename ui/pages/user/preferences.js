import UserPreferences from "../../components/UserPreferences";
import { NoSsr, Paper, withStyles } from "@material-ui/core";
import { updatepagepath, updatepagetitle } from "../../lib/store";
import { connect } from "react-redux";
import { bindActionCreators } from 'redux';
import { getPath } from "../../lib/path";
import Head from 'next/head';
import dataFetch from '../../lib/data-fetch';
import { ctxUrl } from "../../utils/multi-ctx";
import React from "react";

const styles = { paper : { maxWidth : '90%',
  margin : 'auto',
  overflow : 'hidden', } };

class UserPref extends React.Component {
  constructor(props){
    super(props);
    this.state={};
  }

  async componentDidMount() {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path : getPath() });
    this.props.updatepagetitle({ title : "User Preferences" });

    await new Promise(resolve => {
      dataFetch(
        ctxUrl('/api/user/prefs', this.props.selectedK8sContexts),
        {
          method : 'GET',
          credentials : 'include',
        }, (result) => {
          resolve();
          console.log(result);
          if (typeof result !== 'undefined') {
            this.setState({
              anonymousStats : result.anonymousUsageStats||false,
              perfResultStats : result.anonymousPerfResults||false,
            });
          }
        },
        // Ignore error because we will fallback to default state
        // and to avoid try catch due to async await functions
        resolve);
    });
  }

  render () {
    const { anonymousStats, perfResultStats }=this.state;
    console.log(this.state)
    if (anonymousStats===undefined){
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
          <UserPreferences anonymousStats={anonymousStats} perfResultStats={perfResultStats}  theme={this.props.theme} themeSetter={this.props.themeSetter}/>
        </Paper>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  updatepagepath : bindActionCreators(updatepagepath, dispatch),
  updatepagetitle : bindActionCreators(updatepagetitle, dispatch),
})
const mapStateToProps = (state) => {
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    selectedK8sContexts,
  };
};

export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps
)(UserPref));