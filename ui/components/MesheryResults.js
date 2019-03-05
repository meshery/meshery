import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Grid } from '@material-ui/core';
import MesheryResult from './MesheryResult';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import { updateMeshResults } from '../lib/store';
import dataFetch from '../lib/data-fetch';

const styles = theme => ({
  grid: {
    padding: theme.spacing(2),
  },
});

class MesheryResults extends Component {
    state = {
        // results: {},
        // startKey: '',
        page: 1,
    }

    componentDidMount = () => {
      this.fetchResults();
    }

    fetchResults = () => {
          let self = this;
          let query = '';
          if (this.state.startKey !== ''){
            query = `?startKey=${encodeURIComponent(this.state.startKey)}`;
          }
          dataFetch(`/api/results${query}`, { 
            credentials: 'same-origin',
            method: 'GET',
            credentials: 'include',
          }, result => {
            console.log(`received results: ${JSON.stringify(result)}`);
            if (typeof result !== 'undefined'){
            //   this.setState({result, startKey: result.last_key}); // s
              this.props.updateMeshResults({startKey: result.last_key, results: result.results});
            }
          }, self.handleError);
    }

    handleError = error => {
        // this.setState({showSnackbar: true, snackbarVariant: 'error', snackbarMessage: `Load test did not run successfully with msg: ${error}`});
        console.log(`error fetching results: ${error}`);
      }

    render() {
        const { classes, results } = this.props; // data here maps to the MesheryResult model
        return (
            <NoSsr>
            <Grid container spacing={5} className={classes.grid}>
                {results.map((result) => (
                <Grid item xs={12} sm={6}>
                    <MesheryResult key={result.meshery_id} data={result} />
                </Grid>
                ))}
            </Grid>
            </NoSsr>
        );
    }
}
MesheryResults.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
    return {
        updateMeshResults: bindActionCreators(updateMeshResults, dispatch)
    }
  }
  const mapStateToProps = state => {
    const startKey = state.get("results").get('startKey');
    const results =  state.get("results").get('results').toArray();
    if (typeof results !== 'undefined'){
        return {startKey: startKey, results: results};
    }
    return {};
  }
  
export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
)(MesheryResults));
  