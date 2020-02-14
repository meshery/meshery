import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { NoSsr, Chip, IconButton, Card, CardContent, Typography, CardHeader, Tooltip, Avatar, Button, ButtonGroup } from '@material-ui/core';
import dataFetch from '../lib/data-fetch';
import blue from '@material-ui/core/colors/blue';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import Grow from '@material-ui/core/Grow';
import Paper from '@material-ui/core/Paper';
import Popper from '@material-ui/core/Popper';
import MenuItem from '@material-ui/core/MenuItem';
import MenuList from '@material-ui/core/MenuList';


const styles = theme => ({
  root: {
    padding: theme.spacing(10),
    textAlign: 'center',
  },
  container: {
    width: '60%',
    marginLeft: 'auto',
    marginRight: 'auto',
    marginTop: theme.spacing(5),
  },
  logo: {
    width: theme.spacing(50),
  },
});

class ProviderComponent extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            availableProviders: {},
            selectedRemote: '',
            selectedLocal: '',
            open: false,
          };
        this.anchorRef = null;
      }
    
      loadProvidersFromServer() {
        const self = this;
          dataFetch('/api/providers', { 
            credentials: 'same-origin',
            method: 'GET',
            credentials: 'include',
          }, result => {
          if (typeof result !== 'undefined'){
            let selectedRemote = '';
            let selectedLocal = '';
            Object.keys(result).forEach(key => {
              if(result[key] === 'remote'){
                selectedRemote = key;
              } else {
                selectedLocal = key;
              }
            })
            self.setState({availableProviders: result, selectedRemote, selectedLocal});
            }
          }, error => {
            console.log(`there was an error fetching providers: ${error}`);
          });
      }

      componentDidMount = () => {
        this.loadProvidersFromServer();
      }

  handleError = (msg) => error => {
    const self = this;
  }

  handleMenuItemClick = (index) => {
    this.setState({selectedRemote: index});
  };

  handleToggle(){
    const self = this;
    return () => {
      self.setState({open: !self.state.open});
    };
  }

  handleClose(){
    const self = this;
    return event => {
      if (self.anchorRef && self.anchorRef.contains(event.target)) {
        return;
      }
      self.setState({open: false});
    };
  }
  
  render(){
    const { classes } = this.props;
    const { availableProviders, selectedRemote, selectedLocal, open } = this.state;
    const self = this;
      return (
    <NoSsr>
    <div className={classes.root}>
    <img className={classes.logo} src={'/provider/static/img/meshery-logo/meshery-logo-light-text.png'} />
        <Typography variant="h6" gutterBottom className={classes.chartTitle}>
          Please choose a provider to continue
        </Typography>  
        <div className={classes.container}>
        <Grid container spacing={10} >
        <Grid item xs={12} sm={6} justify="flex-end" alignItems="center">
          {selectedLocal !== '' &&
        <Button
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            // onClick={self.handleLocalSubmit}
            href={`/api/provider?provider=${encodeURIComponent(selectedLocal)}`}
            // className={classes.button}
          >
           {selectedLocal}
          </Button>}
        </Grid>
        <Grid item xs={12} sm={6}>
        {selectedRemote !== '' && 
        <React.Fragment>
        <ButtonGroup variant="contained" color="primary" ref={ref => self.anchorRef = ref} aria-label="split button">
          <Button size="large" 
            // onClick={self.handleRemoteSubmit(selectedRemote)}
            href={`/api/provider?provider=${encodeURIComponent(selectedRemote)}`}
            >{selectedRemote}</Button>
          <Button
            color="primary"
            size="small"
            aria-controls={open ? 'split-button-menu' : undefined}
            aria-expanded={open ? 'true' : undefined}
            aria-label="select remote provider"
            aria-haspopup="menu"
            onClick={self.handleToggle()}
          >
            <ArrowDropDownIcon />
          </Button>
        </ButtonGroup>
        <Popper open={open} anchorEl={self.anchorRef} role={undefined} transition disablePortal>
          {({ TransitionProps, placement }) => (
            <Grow
              {...TransitionProps}
              style={{
                transformOrigin: placement === 'bottom' ? 'center top' : 'center bottom',
              }}
            >
              <Paper>
                <ClickAwayListener onClickAway={self.handleClose()}>
                  <MenuList id="split-button-menu">
                    {Object.keys(availableProviders).filter(key => availableProviders[key] === 'remote').map(key => (
                      <MenuItem
                        key={key}
                        // disabled={index === 2}
                        selected={key === selectedRemote}
                        onClick={ev => self.handleMenuItemClick(key)}
                      >
                        {key}
                      </MenuItem>
                    ))}
                  </MenuList>
                </ClickAwayListener>
              </Paper>
            </Grow>
          )}
        </Popper>
        </React.Fragment>}
        </Grid>
        </Grid>
        </div>      
    </div>
    </NoSsr>
  );
    
}
}

ProviderComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ProviderComponent);
