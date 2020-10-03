/* eslint-disable */
import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { Autocomplete } from '@material-ui/lab'
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import {
  NoSsr, Tooltip, MenuItem, CircularProgress, FormControl, FormLabel, RadioGroup, FormControlLabel, Radio, Divider, ExpansionPanel, ExpansionPanelSummary, ExpansionPanelDetails,
} from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import { withSnackbar } from 'notistack';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import AddCircleIcon from '@material-ui/icons/AddCircle';
import { updateLoadTestData, updateStaticPrometheusBoardConfig, updateLoadTestPref } from '../lib/store';
import LoadTestTimerDialog from './load-test-timer-dialog';
import { durationOptions } from '../lib/prePopulatedOptions';

const meshes = [
  'Istio',
  'Linkerd',
  'App Mesh',
  'Aspen Mesh',
  'Citrix Service Mesh',
  'Consul Connect',
  'Grey Matter',
  'Kong',
  'Mesher',
  'Network Service Mesh',
  'Octarine',
  'Rotor',
  'SOFAMesh',
  'Open Service Mesh',
  'Zuul',
];

const loadGenerators = [
  'fortio',
  'wrk2',
];

const styles = (theme) => ({
  root: {
    padding: theme.spacing(10),
    paddingBottom: 0,
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  expansionPanel: {
    boxShadow:'none',
    border: '1px solid rgb(196,196,196)',
  },
  margin: {
    margin: theme.spacing(1),
  },
  icon:{
    fontSize: 25,
    marginRight: '8px',
  },
});

class MesheryPerformanceFormComponent extends React.Component {
  constructor(props) {
    super(props);
    const {
      testName, meshName, url, qps, c, t, result, staticPrometheusBoardConfig, k8sConfig, loadTestPrefs,
    } = props;

    this.state = {
      testName,
      meshName,
      url,
      qps,
      c,
      t,
      tValue: t,
      loadGenerator: 'fortio',
      result,
      headers: "",
      cookies: "",
      reqBody: "",
      contentType: "",

      timerDialogOpen: false,
      blockRunTest: false,
      urlError: false,
      tError: '',
      selectedMesh: '',
      availableAdapters: [],
    };
  }

  handleChange = (name) => (event) => {
    if (name === 'url' && event.target.value !== '') {
      let urlPattern = event.target.value;
      let val = urlPattern.match(/(http(s)?:\/\/.)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)/g);
      if ( !val ){
        this.setState({ disableTest: true });
        this.setState({ urlError: true });
      }
      else{
        this.setState({ disableTest: false });
        this.setState({ urlError: false });
      }
    }
    else
      this.setState({ urlError: false });
    this.setState({ [name]: event.target.value });
  };

  handleDurationChange = (event, newValue) => {
    this.setState({tValue: newValue})
    if (newValue !== null) {
      this.setState({ tError: '' })
    }
  };

  handleInputDurationChange = (event, newValue) => {
    this.setState({t: newValue})
  };

  render() {
    const { classes } = this.props;
    const {
      timerDialogOpen, blockRunTest, url, qps, c, t, loadGenerator, testName, meshName, urlError,
      tError, selectedMesh, availableAdapters, headers, cookies, reqBody, contentType, tValue, disableTest
    } = this.state;

    availableAdapters.forEach((item) => {
      const index = meshes.indexOf(item);
      if (index !== -1) meshes.splice(index, 1);
    });

    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.root}>
            <Grid container spacing={1}>
              <Grid item xs={12} md={6}>
                <Tooltip title="If a test name is not provided, a random one will be generated for you.">
                  <TextField
                    id="testName"
                    name="testName"
                    label="Test Name"
                    fullWidth
                    value={testName}
                    margin="normal"
                    variant="outlined"
                    onChange={this.handleChange('testName')}
                    inputProps={{ maxLength: 300 }}
                  />
                </Tooltip>
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  select
                  id="meshName"
                  name="meshName"
                  label="Service Mesh"
                  fullWidth
                  value={meshName === '' && selectedMesh !== '' ? selectedMesh : meshName}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange('meshName')}
                >

                  {availableAdapters && availableAdapters.map((mesh) => (
                    <MenuItem key={`mh_-_${mesh}`} value={mesh.toLowerCase()}>{mesh}</MenuItem>
                  ))}
                  {availableAdapters && (availableAdapters.length > 0) && <Divider />}
                  <MenuItem key="mh_-_none" value="None">None</MenuItem>
                  {meshes && meshes.map((mesh) => (
                    <MenuItem key={`mh_-_${mesh}`} value={mesh.toLowerCase()}>{mesh}</MenuItem>
                  ))}
                </TextField>
              </Grid>
              <Grid item xs={12}>
                <TextField
                  required
                  id="url"
                  name="url"
                  label="URL to test"
                  type="url"
                  fullWidth
                  value={url}
                  error={urlError}
                  helperText={urlError ? "Please enter a valid URL along with protocol" : ""}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange('url')}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <TextField
                  required
                  id="c"
                  name="c"
                  label="Concurrent requests"
                  type="number"
                  fullWidth
                  value={c}
                  inputProps={{ min: '0', step: '1' }}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange('c')}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <TextField
                  required
                  id="qps"
                  name="qps"
                  label="Queries per second"
                  type="number"
                  fullWidth
                  value={qps}
                  inputProps={{ min: '0', step: '1' }}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleChange('qps')}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <Tooltip title={"Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."}>
                  <Autocomplete
                    required
                    id="t"
                    name="t"
                    freeSolo
                    label="Duration*"
                    fullWidth
                    variant="outlined"
                    className={classes.errorValue}
                    classes={{ root: tError }}
                    value={tValue}
                    inputValue={t}
                    onChange={this.handleDurationChange}
                    onInputChange={this.handleInputDurationChange}
                    options={durationOptions}
                    style={{ marginTop: '16px', marginBottom: '8px' }}
                    renderInput={(params) => <TextField {...params} label="Duration*" variant="outlined" />}
                  />
                </Tooltip>
              </Grid>
              <Grid item xs={12} md={12} gutterBottom>
                <ExpansionPanel className={classes.expansionPanel}>
                  <ExpansionPanelSummary expanded={true} expandIcon={<ExpandMoreIcon/>}>
                    <Typography align="center" color="textSecondary" varient="h6">Advanced Options</Typography>
                  </ExpansionPanelSummary>
                  <ExpansionPanelDetails>
                    <Grid container spacing={1}>
                      <Grid item xs={12}>
                        <TextField
                          id="headers"
                          name="headers"
                          label="Request Headers"
                          fullWidth
                          value={headers}
                          multiline
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange('headers')}
                        >
                        </TextField>
                      </Grid>
                      <Grid item xs={12}>
                        <TextField
                          id="cookies"
                          name="cookies"
                          label="Request Cookies"
                          fullWidth
                          value={cookies}
                          multiline
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange('cookies')}
                        >
                        </TextField>
                      </Grid>
                      <Grid item xs={12}>
                        <TextField
                          id="contentType"
                          name="contentType"
                          label="Content Type"
                          fullWidth
                          value={contentType}
                          multiline
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange('contentType')}
                        >
                        </TextField>
                      </Grid>
                      <Grid item xs={12} md={12}>
                        <TextField
                          id="cookies"
                          name="cookies"
                          label="Request Body"
                          fullWidth
                          value={reqBody}
                          multiline
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange('reqBody')}
                        >
                        </TextField>
                      </Grid>
                    </Grid>
                  </ExpansionPanelDetails>
                </ExpansionPanel>
              </Grid>
              <Grid item xs={12} md={4}>
                <FormControl component="loadGenerator" className={classes.margin}>
                  <FormLabel component="loadGenerator">Load generator</FormLabel>
                  <RadioGroup aria-label="loadGenerator" name="loadGenerator" value={loadGenerator} onChange={this.handleChange('loadGenerator')} row>
                    {loadGenerators.map((lg) => (
                      <FormControlLabel value={lg} control={<Radio color="primary" />} label={lg} />
                    ))}
                  </RadioGroup>
                </FormControl>
              </Grid>
            </Grid>
          </div>
        </React.Fragment>
      </NoSsr>
    );
  }
}

MesheryPerformanceFormComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateLoadTestData: bindActionCreators(updateLoadTestData, dispatch),
  updateStaticPrometheusBoardConfig: bindActionCreators(updateStaticPrometheusBoardConfig, dispatch),
  updateLoadTestPref: bindActionCreators(updateLoadTestPref, dispatch),
});

export default withStyles(styles)(connect(
  null,
  mapDispatchToProps,
)(withSnackbar(MesheryPerformanceFormComponent)));
