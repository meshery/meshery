import React from 'react';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux'
import { withRouter } from 'next/router';
import { withSnackbar } from 'notistack';
import { updateUser, updateProgress } from '../lib/store';
import { withStyles } from '@material-ui/core/styles';
import dataFetch from '../lib/data-fetch';
import CloseIcon from '@material-ui/icons/Close';
import { IconButton, FormControlLabel, Switch } from '@material-ui/core';
import NoSsr from '@material-ui/core/NoSsr';


const styles = theme => ({
  formContainer: {
    margin:50,
  },
});

class UserPreference extends React.Component {

  constructor(props){
    super(props);
    this.state = {
      anonymousStats: false,
      perfResultStats: false,
    }
  }

  handleToggle = name => () => {
    const self = this;
    if(name == "anonymousUsageStats"){
      //val=anonymousStats;
      self.setState(state => ({ anonymousStats : !state.anonymousStats }));
    } else{
      //val=perfResultStats;
      self.setState(state => ({ perfResultStats : !state.perfResultStats }));
    }
    
    this.handleChange(name);
  };

  handleError = msg => error => {
    const self = this;
    this.props.updateProgress({showProgress: false});
    this.props.enqueueSnackbar(msg, {
      variant: 'error',
      action: (key) => (
        <IconButton
          key="close"
          aria-label="Close"
          color="inherit"
          onClick={() => self.props.closeSnackbar(key) }
        >
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration: 8000,
    });
  }

  handleChange = (name) => {
    const self = this;
    const { anonymousStats , perfResultStats } = this.state;
    var val;
    if(name == "anonymousUsageStats"){
      val=anonymousStats;
    } else{
      val=perfResultStats;
    }
    const params = encodeURIComponent(name) + '=' + encodeURIComponent(val);
    // console.log(`data to be submitted for load test: ${params}`);
    this.props.updateProgress({showProgress: true});
    //let self = this;
    dataFetch('/api/user/stats', { 
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
      },
      body: params
    }, result => {
      this.props.updateProgress({showProgress: false});
      if (typeof result !== 'undefined'){
        this.props.enqueueSnackbar('Preference was successfully updated!', {
          variant: 'success',
          autoHideDuration: 2000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key) }
            >
              <CloseIcon />
            </IconButton>
          ),
        });
      }
    }, self.handleError(`There was an error sending your preference`));
  }

  render() {
    const { anonymousStats , perfResultStats } = this.state;
    const { classes } = this.props;
    
    return (
      <NoSsr>
        <div className={classes.formContainer}>
          <FormControlLabel
            key="UsageStatsPreference"
            control={
              <Switch
                checked={anonymousStats}
                onChange={this.handleToggle('anonymousUsageStats')}
                color="default"
              />
            }
            labelPlacement="end"
            label="Send Anonymous Usage Statistics"
          />
        </div>
        <div className={classes.formContainer}>
          <FormControlLabel
            key="PerfResultPreference"
            control={
              <Switch
                checked={perfResultStats}
                onChange={this.handleToggle('anonymousPerfResults')}
                color="default"
              />
            }
            labelPlacement="end"
            label="Send Anonymous Performance Results"
          />
        </div>
      </NoSsr>
    )
  }
}

const mapDispatchToProps = dispatch => {
  return {
    updateUser: bindActionCreators(updateUser, dispatch),
    updateProgress: bindActionCreators(updateProgress, dispatch),
  }
}

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps,
)(withRouter(withSnackbar(UserPreference))));