/* eslint-disable */
import React from 'react';
import PropTypes from 'prop-types';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import { withRouter } from 'next/router';
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
import CreateIcon from '@material-ui/icons/Create';
import { bindActionCreators } from 'redux';
import MesheryPerformanceComponent from './MesheryPerformanceComponent';

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
  createProfile: {
    textAlign: 'center',
    padding: theme.spacing(20),
  },
  icon:{
    fontSize: 20,
    marginRight: '8px',
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

  handleCreateProfile = () => {
    this.props.router.push('/performance');
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

    if(profiles.length <= 0){
      return (
        <NoSsr>
          <React.Fragment>
          <div className={classes.createProfile}>
              <Button variant="contained" color="primary" size="large" onClick={this.handleCreateProfile}>
                <CreateIcon className={classes.icon} />
                  Create a Test Profile
              </Button>
            </div>
          </React.Fragment>
        </NoSsr>
      )
    }
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
                    <MesheryPerformanceComponent editProfile={"editProfile"}/>
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
)(withRouter(MesheryPerfProfileComponent)));
