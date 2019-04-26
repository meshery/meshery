import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { NoSsr,  Chip } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import Snackbar from '@material-ui/core/Snackbar';
import MesherySnackbarWrapper from '../components/MesherySnackbarWrapper';
import dataFetch from '../lib/data-fetch';
import blue from '@material-ui/core/colors/blue';
import { updateAdaptersInfo } from '../lib/store';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import CreatableSelect from 'react-select/lib/Creatable';
import ReactSelectWrapper from './ReactSelectWrapper';


const styles = theme => ({
  root: {
    padding: theme.spacing(5),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  margin: {
    margin: theme.spacing(1),
  },
  alreadyConfigured: {
    textAlign: 'center',
    padding: theme.spacing(20),
  },
  colorSwitchBase: {
    color: blue[300],
    '&$colorChecked': {
      color: blue[500],
      '& + $colorBar': {
        backgroundColor: blue[500],
      },
    },
  },
  colorBar: {},
  colorChecked: {},
  fileLabel: {
    width: '100%',
  },
  fileLabelText: {
  },
  inClusterLabel: {
    paddingRight: theme.spacing(2),
  },
  alignCenter: {
    textAlign: 'center',
  },
  alignRight: {
    textAlign: 'right',
    marginBottom: theme.spacing(2),
  },
  fileInputStyle: {
    opacity: '0.01',
  },
  icon: {
    width: theme.spacing(2.5),
  },
});

class MeshAdapterConfigComponent extends React.Component {

  constructor(props) {
    super(props);
    const {meshAdapters} = props;
    this.state = {
        showSnackbar: false,
        snackbarVariant: '',
        snackbarMessage: '',
    
        meshAdapters,
        availableAdapters: [],
    
        meshLocationURLError: false,
      };
  }

  componentDidMount = () => {
    this.fetchAvailableAdapters();
  }

  fetchAvailableAdapters = () => {
    let self = this;
    dataFetch('/api/mesh/adapters', { 
      credentials: 'same-origin',
      method: 'GET',
      credentials: 'include',
    }, result => {
      if (typeof result !== 'undefined'){
        const options = result.map(res => {
          return {
            value: res,
            label: res,
          };
        });
        this.setState({availableAdapters: options});
      }
    }, self.handleError("Unable to fetch available adapters"));
  }

  handleSnackbarClose = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }

    this.setState({ showSnackbar: false });
  };

  handleChange = name => event => {
    if (name === 'meshLocationURL' && event.target.value !== '') {
        this.setState({meshLocationURLError: false})
    }
    this.setState({ [name]: event.target.value });
  };

  handleMeshLocURLChange = (newValue, actionMeta) => {
    // console.log(newValue);
    // console.log(`action: ${actionMeta.action}`);
    // console.groupEnd();
    if (typeof newValue !== 'undefined'){
      this.setState({meshLocationURL: newValue, meshLocationURLError: false});
    }
  };
  handleInputChange = (inputValue, actionMeta) => {
    console.log(inputValue);
    // console.log(`action: ${actionMeta.action}`);
    // console.groupEnd();

    // TODO: try to submit it and get 
    // if (typeof inputValue !== 'undefined'){
    //   this.setState({meshLocationURL: inputValue});
    // }
  }

  handleSubmit = () => {
    const { meshLocationURL } = this.state;
    
    if (!meshLocationURL || !meshLocationURL.value || meshLocationURL.value === ''){
        this.setState({meshLocationURLError: true})
        return;
      }

    this.submitConfig();
  }

  submitConfig = () => {
    const { meshLocationURL } = this.state;
    
    const data = {meshLocationURL: meshLocationURL.value};

    const params = Object.keys(data).map((key) => {
      return encodeURIComponent(key) + '=' + encodeURIComponent(data[key]);
    }).join('&');

    let self = this;
    dataFetch('/api/mesh/manage', { 
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
      },
      body: params
    }, result => {
      if (typeof result !== 'undefined'){
        this.setState({meshAdapters: result, meshLocationURL: '', showSnackbar: true, snackbarVariant: 'success', snackbarMessage: 'Adapter was successfully configured!'});
        this.props.updateAdaptersInfo({meshAdapters: result});
        this.fetchAvailableAdapters();
      }
    }, self.handleError("Adapter was not configured due to an error"));
  }

  handleDelete = (adapterID) => () => {
    // const { meshAdapters } = this.state;

    let self = this;
    dataFetch(`/api/mesh/manage?adapterID=${adapterID}`, { 
      credentials: 'same-origin',
      method: 'DELETE',
      credentials: 'include',
    }, result => {
      if (typeof result !== 'undefined'){
        this.setState({meshAdapters: result, showSnackbar: true, snackbarVariant: 'success', snackbarMessage: 'Adapter was successfully removed!'});
        this.props.updateAdaptersInfo({meshAdapters: result});
      }
    }, self.handleError("Adapter was not removed due to an error"));
  }

  handleError = (msg) => (error) => {
    this.setState({showSnackbar: true, snackbarVariant: 'error', snackbarMessage: `${msg}: ${error}`});
  }

  configureTemplate = () => {
    const { classes } = this.props;
    const { availableAdapters, meshAdapters, meshLocationURL, meshLocationURLError, showSnackbar, 
        snackbarVariant, snackbarMessage, clusterConfigured } = this.state;
    
    let showAdapters = '';
    const self = this;
    if (meshAdapters.length > 0) {
      showAdapters = (
        <div className={classes.alignRight}>
          {meshAdapters.map((adapter, ind) => {
            let image = "/static/img/meshery-logo.png";
            switch (adapter.name.toLowerCase()){
              case 'istio':
                image = "/static/img/istio.svg";
                break;
              case 'linkerd':
                image = "/static/img/linkerd.svg";
                break;
              // default:
            } 
            return (
            <Chip 
            label={adapter.adapter_location}
            onDelete={self.handleDelete(ind)} 
            icon={<img src={image} className={classes.icon} />} 
            variant="outlined" />
          );
          })}
          
        </div>
      )
    }


      return (
    <NoSsr>
    <div className={classes.root}>
    
    {showAdapters}
    
    <Grid container spacing={1} alignItems="flex-end">
      <Grid item xs={12}>

        {/* <CreatableSelect
          isClearable
          onChange={this.handleMeshLocURLChange}
          onInputChange={this.handleInputChange}
          options={availableAdapters}
        /> */}

        <ReactSelectWrapper
          onChange={this.handleMeshLocURLChange}
          onInputChange={this.handleInputChange}
          options={availableAdapters}
          value={meshLocationURL}
          // placeholder={'Mesh Adapter URL'}
          label={'Mesh Adapter URL'}
          error={meshLocationURLError}
        />

        {/* <TextField
          required
          id="meshLocationURL"
          name="meshLocationURL"
          label="Mesh Adapter Location"
          type="url"
          fullWidth
          value={meshLocationURL}
          error={meshLocationURLError}
          margin="normal"
          variant="outlined"
          onChange={this.handleChange('meshLocationURL')}
        /> */}
      </Grid>
    </Grid>
    <React.Fragment>
      <div className={classes.buttons}>
        <Button
          type="submit"
          variant="contained"
          color="primary"
          size="large"
          onClick={this.handleSubmit}
          className={classes.button}
        >
         Submit
        </Button>
      </div>
    </React.Fragment>
    </div>
  
  {/* <LoadTestTimerDialog open={timerDialogOpen} 
    t={t}
    onClose={this.handleTimerDialogClose} 
    countDownComplete={this.handleTimerDialogClose} />

  <Typography variant="h6" gutterBottom className={classes.chartTitle}>
      Results
    </Typography>
  <MesheryChart data={result} />     */}
  
  <Snackbar
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'right',
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
    </NoSsr>
  );
    }

  render() {
    const { reconfigureCluster } = this.state;
    // if (reconfigureCluster) {
    return this.configureTemplate();
    // }
    // return this.alreadyConfiguredTemplate();
  }
}

MeshAdapterConfigComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
    return {
        updateAdaptersInfo: bindActionCreators(updateAdaptersInfo, dispatch),
    }
}
const mapStateToProps = state => {
    const meshAdapters = state.get("meshAdapters").toJS();
    return {meshAdapters};
}

export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
  )(withRouter(MeshAdapterConfigComponent)));
