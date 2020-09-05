/* eslint-disable */
import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import {
  NoSsr, Card, CardHeader, CardContent, Grid, IconButton, Tooltip,
} from '@material-ui/core';

import Dialog from '@material-ui/core/Dialog';
import MuiDialogTitle from '@material-ui/core/DialogTitle';
import MuiDialogContent from '@material-ui/core/DialogContent';
import MuiDialogActions from '@material-ui/core/DialogActions';
import TextField from '@material-ui/core/TextField';
import { connect } from 'react-redux';
import dataFetch from '../lib/data-fetch';
import PlayCircleFilledIcon from '@material-ui/icons/PlayCircleFilled';
import TimerIcon from '@material-ui/icons/Timer';
import CloseIcon from '@material-ui/icons/Close';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import { bindActionCreators } from 'redux';

const styles = (theme) => ({
  root: {
    padding: theme.spacing(10),
  },
  expansionPanel: {
    boxShadow:'none',
    border: '1px solid rgb(196,196,196)',
  },
  margin: {
    margin: theme.spacing(1),
  },
  button: {
    marginRight: '10px',
    [theme.breakpoints.down(1350)]: {
      marginTop: '10px',
      width: '160px',
    },
    [theme.breakpoints.down(1070)]: {
      marginTop: '10px',
      width:'100px',
    },
    [theme.breakpoints.down('sm')]: {
      marginTop: 0,
      width: 'auto',
    },
    [theme.breakpoints.down(900)]: {
      marginTop: '10px',
      width: '160px',
    },
    [theme.breakpoints.down(780)]: {
      marginTop: '10px',
      width: '100px',
    },
    [theme.breakpoints.down(600)]: {
      marginTop: '10px',
      width: 'auto',
    },
  },
  profileButtons: {
    padding: '15px'
  },
  closeButton: {
    position: 'absolute',
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
});

const DialogTitle = withStyles(styles)(props => {
  const { children, classes, onClose, ...other } = props;
  return (
    <MuiDialogTitle disableTypography className={classes.modalHeading} {...other}>
      <Typography variant="h6">{children}</Typography>
      {onClose ? (
        <IconButton aria-label="close" className={classes.closeButton} onClick={onClose}>
          <CloseIcon />
        </IconButton>
      ) : null}
    </MuiDialogTitle>
  );
});

const DialogContent = withStyles(theme => ({
  root: {
    padding: theme.spacing(2),
  },
}))(MuiDialogContent);

const DialogActions = withStyles(theme => ({
  root: {
    padding: theme.spacing(1),
  },
}))(MuiDialogActions);

class MesheryPerfProfileComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state={
      modalOpen : false,
      updatedProfile : true,
    }
  }

  handleSubmit = () => {
    const data = {
      rps:'10',
      t:'10s',
      c:'10',
      gen:'fortio',
      protocol:'TCP',
      headers:{"h1":"v1"},
      cookies: {"h1":"v1"},
      reqBody: 'reqBody',
      contentType: 'type',
      endpoint:"https://google.com",
      labels:{"h1":"v1"}
    };

    const params = Object.keys(data).map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`).join('&');
    dataFetch('/api/user/test-prefs', {
      credentials: 'same-origin',
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        'Cookie': 'meshery-provider=None',
      },
      body: params,
    }, (result) => {
      console.log(result);
      alert(result);
    }, console.error("Fetch Fail"));
  }
  async startEventStream(url) {
    this.closeEventStream();
    this.eventStream = new EventSource(url);
    this.eventStream.onmessage = this.handleEvents();
    this.eventStream.onerror = this.handleError('Connection to the server got disconnected. Load test might be running in the background. Please check the results page in a few.');
    this.props.enqueueSnackbar('Load test has been successfully submitted', {
      variant: 'info',
      autoHideDuration: 1000,
      action: (key) => (
        <IconButton
          key="close"
          aria-label="Close"
          color="inherit"
          onClick={() => this.props.closeSnackbar(key)}
        >
          <CloseIcon />
        </IconButton>
      ),
    });
  }

  handleChange = (val) => {
    const self = this;
    return () => {
      self.setState({
        updatedProfile : false,
      });
    };
  }

  handleModalOpen() {
    const self = this;
    return () => {
      self.setState({ modalOpen : true });
    };
  }

  handleModalClose() {
    const self=this;
    return () => {
      self.setState({ modalOpen : false });
    };
  }

  render() {
    const { classes } = this.props;
    const self = this;
    const { modalOpen, updatedProfile } = this.state;
    const profiles = [
      {
        name: 'Test1',
        id: '1001'
      },
      {
        name: 'Test2',
        id: '1002'
      },
      {
        name: 'Test3',
        id: '1002'
      },
    ]
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.root}>
          <Grid container spacing={3}>
            {profiles.map( ({ name,id }, index) => {
              return (
                <>
                  <Grid item xs={12} md={6} key={index} gutterBottom>
                    <Card>
                      <CardHeader 
                        title={name}
                        subheader={id}
                        action={
                          <Tooltip title="View or Edit Profile">
                          <IconButton onClick={self.handleModalOpen()}>
                            <EditIcon />
                          </IconButton>
                          </Tooltip>
                        }
                      />
                      <CardContent className={classes.profileButtons}>
                      <Button
                            type="submit"
                            variant="contained"
                            color="primary"
                            size="medium"
                            className={classes.button}
                            onClick={this.handleSubmit}
                            startIcon={<PlayCircleFilledIcon />}
                        >
                          Run Test
                        </Button>
                        <Button
                            type="submit"
                            variant="contained"
                            color="primary"
                            size="medium"
                            className={classes.button}
                            startIcon={<TimerIcon />}
                        >
                          Schedule Test
                        </Button>
                        <Tooltip title="Delete Profile">
                          <IconButton style={{ float: 'right' }}>
                            <DeleteIcon />
                          </IconButton>
                        </Tooltip>
                      </CardContent>
                    </Card>
                  </Grid>
                  <Dialog onClose={self.handleModalClose()} aria-labelledby="customized-dialog-title" open={modalOpen} disableScrollLock={true}>
                  <DialogTitle id="customized-dialogs-title" onClose={self.handleModalClose()}>
                    <b>Profile Details</b>
                  </DialogTitle>
                  <DialogContent dividers>
                    <Grid container spacing={1}>
                    <Grid item xs={12} md={6}>
                      <Tooltip title="If a test name is not provided, a random one will be generated for you.">
                        <TextField
                          id="testName"
                          name="testName"
                          label="Test Name"
                          fullWidth
                          value={name}
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
                        value="Service Mesh"
                        margin="normal"
                        variant="outlined"
                        onChange={this.handleChange('meshName')}
                      >
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
                        value="URL"
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
                        value="Concurrent Req"
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
                        value="QPS"
                        inputProps={{ min: '0', step: '1' }}
                        margin="normal"
                        variant="outlined"
                        onChange={this.handleChange('qps')}
                      />
                    </Grid>
                    <Grid item xs={12} md={4}>
                      <Tooltip title={"Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."}>
                        <TextField
                          required
                          id="t"
                          name="t"
                          label="Duration"
                          fullWidth
                          value="Duration"
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange('t')}
                        />
                      </Tooltip>
                    </Grid>
                    <Grid item xs={12} md={12} gutterBottom>
                          <Grid container spacing={1}>
                            <Grid item xs={12}>
                              <TextField
                                id="headers"
                                name="headers"
                                label="Request Headers"
                                fullWidth
                                value="Req Headers"
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
                                value=""
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
                                value="Content Type"
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
                                value="Req Body"
                                multiline
                                margin="normal"
                                variant="outlined"
                                onChange={this.handleChange('reqBody')}
                              >
                              </TextField>
                            </Grid>
                          </Grid>
                    </Grid>
                  </Grid>
                </DialogContent>
                <DialogActions>
                  <Button color="primary" disabled={updatedProfile}>
                    Update Profile
                  </Button>
                  <Button autoFocus onClick={self.handleModalClose()} color="primary">
                    OK
                  </Button>
                </DialogActions>
              </Dialog>
              </>
            )
            })}
            </Grid>
          </div>
        </React.Fragment>
      </NoSsr>
    );
  }
}

MesheryPerfProfileComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(connect(
  //mapStateToProps,
  //mapDispatchToProps,
  null,
  null
)(MesheryPerfProfileComponent));
