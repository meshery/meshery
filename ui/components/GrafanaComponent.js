import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Snackbar, TextField, Grid, Button, Chip, FormControl, InputLabel, Select, MenuItem, OutlinedInput } from '@material-ui/core';
import MesherySnackbarWrapper from './MesherySnackbarWrapper';
import dataFetch from '../lib/data-fetch';
import ReactDOM from 'react-dom';
import GrafanaConfigComponent from './GrafanaConfigComponent';
import GrafanaSelectionComponent from './GrafanaSelectionComponent';
import GrafanaDisplaySelection from './GrafanaDisplaySelection';

const grafanaStyles = theme => ({
    root: {
      padding: theme.spacing(10),
    },
    buttons: {
      display: 'flex',
    //   justifyContent: 'flex-end',
    },
    button: {
      marginTop: theme.spacing(3),
    //   marginLeft: theme.spacing(1),
    },
    margin: {
      margin: theme.spacing(1),
    },
    chartTitle: {
      textAlign: 'center',
    },
    icon: {
        width: theme.spacing(2.5),
    },
    alignRight: {
        textAlign: 'right',
    },
    formControl: {
        margin: theme.spacing(1),
        minWidth: 180,
    },
    panelChips: {
        display: 'flex',
        flexWrap: 'wrap',
    },
    panelChip: {
        margin: theme.spacing(0.25),
    }
  });

class GrafanaComponent extends Component {
    constructor(props) {
        super(props);
        
        this.state = {
            grafanaURL: '',
            grafanaAPIKey: '',
      
            urlError: false,

            showSnackbar: false,
            snackbarVariant: '',
            snackbarMessage: '',
            
            grafanaConfigSuccess: false,

            grafanaBoardSearch: '',
            grafanaBoards: [],

            selectedBoardsConfigs: [],
          };
    }

    handleSnackbarClose = (event, reason) => {
        if (reason === 'clickaway') {
          return;
        }
    
        this.setState({ showSnackbar: false });
      };
    
      handleChange = name => event => {
        if (name === 'grafanaURL' && event.target.value !== ''){
          this.setState({urlError: false});
        }
        if (name === 'grafanaBoardSearch') {
            this.getGrafanaBoards();
        }

        this.setState({ [name]: event.target.value });
      };
    
      handleGrafanaConfigure = () => {
    
        const { grafanaURL } = this.state;
        if (grafanaURL === ''){
          this.setState({urlError: true})
          return;
        }
        this.submitGrafanaConfigure();
      }
    
      submitGrafanaConfigure = () => {
        const {grafanaURL, grafanaAPIKey} = this.state;
        const data = {
            grafanaURL,
            grafanaAPIKey
        }
        const params = Object.keys(data).map((key) => {
          return encodeURIComponent(key) + '=' + encodeURIComponent(data[key]);
        }).join('&');
        // console.log(`data to be submitted for load test: ${params}`);
        let self = this;
        dataFetch('/api/grafana/config', { 
          credentials: 'same-origin',
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
          },
          body: params
        }, result => {
          if (typeof result !== 'undefined'){
            // TODO: push the grafana info to store
            this.setState({result, grafanaConfigSuccess: true, showSnackbar: true, snackbarVariant: 'success', snackbarMessage: 'Grafana configured successfully!'});
            this.getGrafanaBoards();
            // this.props.updateLoadTestData({loadTest: {
            //   url,
            //   qps,
            //   c,
            //   t, 
            //   result,
            // }});
          }
        }, self.handleError);
      }

      getGrafanaBoards = () => {
        const {grafanaBoardSearch} = this.state;
        
        let self = this;
        dataFetch(`/api/grafana/boards?dashboardSearch=${grafanaBoardSearch}`, { 
          credentials: 'same-origin',
          method: 'GET',
          credentials: 'include',
        }, result => {
          if (typeof result !== 'undefined'){
            this.setState({grafanaBoards: result});
            // TODO: push the grafana info to store
            // this.setState({result, grafanaConfigSuccess: true, showSnackbar: true, snackbarVariant: 'success', snackbarMessage: 'Grafana boards retrieved successfully!'});
            // this.props.updateLoadTestData({loadTest: {
            //   url,
            //   qps,
            //   c,
            //   t, 
            //   result,
            // }});
          }
        }, self.handleError);
      }
    
      handleError = error => {
        // this.setState({timerDialogOpen: false });
        this.setState({showSnackbar: true, snackbarVariant: 'error', snackbarMessage: `Error communicating with Grafana: ${error}`});
      }

      handleGrafanaChipDelete = () => {
        this.setState({grafanaConfigSuccess: false});
      }

    snackbarTmpl = (showSnackbar, snackbarVariant, snackbarMessage) => {
        return (
            <Snackbar
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                open={showSnackbar}
                autoHideDuration={6000}
                onClose={this.handleSnackbarClose}
                >
                <MesherySnackbarWrapper 
                variant={snackbarVariant}
                message={snackbarMessage}
                onClose={this.handleSnackbarClose}
                />
            </Snackbar>
            );
    }

    addSelectedBoardPanelConfig = (boardsSelection) => {
      // %s/d-solo/%s/%s?orgId=%d&panelId=%d&refresh=10s\n\n", grafanaURL, board.UID, slug.Make(board.Title), board.orgID, panel.ID)
      // http://10.199.75.64:30487/d/LJ_uJAvmk/istio-service-dashboard?refresh=10s&orgId=1&var-service=istio-policy.istio-system.svc.cluster.local&var-srcns=All&var-srcwl=All&var-dstns=All&var-dstwl=All
      let {selectedBoardsConfigs} = this.state;
      selectedBoardsConfigs.push(boardsSelection);
      this.setState({selectedBoardsConfigs});
    }

    render() {
        const { grafanaURL, grafanaAPIKey, urlError, grafanaBoards, grafanaBoardSearch, showSnackbar, 
            snackbarVariant, snackbarMessage, grafanaConfigSuccess, selectedBoardsConfigs } = this.state;
        if (grafanaConfigSuccess) {
            let displaySelec = '';
            if (selectedBoardsConfigs.length > 0) {
              displaySelec = (<GrafanaDisplaySelection boardPanelConfigs={selectedBoardsConfigs} />);
            }

            return (
              <NoSsr>
              <React.Fragment>
                <GrafanaSelectionComponent
                  grafanaURL={grafanaURL}
                  grafanaBoards={grafanaBoards}
                  grafanaBoardSearch={grafanaBoardSearch}
                  handleGrafanaBoardSearchChange={this.handleChange}
                  handleGrafanaChipDelete={this.handleGrafanaChipDelete}
                  addSelectedBoardPanelConfig={this.addSelectedBoardPanelConfig}
                  handleError={this.handleError}
                />
                {displaySelec}
                {this.snackbarTmpl(showSnackbar, snackbarVariant, snackbarMessage)}
              </React.Fragment>
              </NoSsr>
            );
        }
        return (
          <NoSsr>
          <React.Fragment>
            <GrafanaConfigComponent
              grafanaURL={grafanaURL}
              grafanaAPIKey={grafanaAPIKey}
              urlError={urlError}
              handleChange={this.handleChange}
              handleGrafanaConfigure={this.handleGrafanaConfigure}
            />
            {this.snackbarTmpl(showSnackbar, snackbarVariant, snackbarMessage)}
          </React.Fragment>
          </NoSsr>
        );
    }
}

GrafanaComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(grafanaStyles)(GrafanaComponent);