import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import { NoSsr, Typography, Button, ButtonGroup, Tooltip } from "@material-ui/core";
import Dialog from "@material-ui/core/Dialog";
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import MuiDialogContent from "@material-ui/core/DialogContent";
import MuiDialogActions from "@material-ui/core/DialogActions";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import ArrowDropDownIcon from "@material-ui/icons/ArrowDropDown";
import ClickAwayListener from "@material-ui/core/ClickAwayListener";
import Grow from "@material-ui/core/Grow";
import Paper from "@material-ui/core/Paper";
import Popper from "@material-ui/core/Popper";
import MenuItem from "@material-ui/core/MenuItem";
import MenuList from "@material-ui/core/MenuList";
import dataFetch from "../lib/data-fetch";
import Divider from "@material-ui/core/Divider";

const styles = (theme) => ({
  root: {
    padding: "80px 0px",
    textAlign: "center",
  },
  container: {
    width: "60%",
    marginLeft: "auto",
    marginRight: "auto",
    marginTop: theme.spacing(5),
  },
  logo: {
    width: theme.spacing(50),
    maxWidth: "100%",
    height: "auto",
  },
  modalHeading: {
    margin: 0,
    padding: theme.spacing(2),
  },
  closeButton: {
    position: "absolute",
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
  providerLink: {
    color: "darkcyan",
    cursor: "pointer",
    fontWeight: 700,
  },
  providerDesc: {
    whiteSpace: "pre",
  },
  chartTitle:{
    fontWeight: 700,
  }
});

const DialogTitle = withStyles(styles)((props) => {
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

const DialogContent = withStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
  },
}))(MuiDialogContent);

const DialogActions = withStyles((theme) => ({
  root: {
    padding: theme.spacing(1),
  },
}))(MuiDialogActions);

class ProviderComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      availableProviders: {},
      selectedProvider: "",
      open: false,
      modalOpen: false,
    };
    this.anchorRef = null;
  }

  loadProvidersFromServer() {
    const self = this;
    dataFetch(
      "/api/providers",
      {
        credentials: "same-origin",
        method: "GET",
        credentials: "include",
      },
      (result) => {
        if (typeof result !== "undefined") {
          let selectedProvider = "";
          Object.keys(result).forEach((key) => {
            if (result[key]["ProviderType"] === "remote") {
              selectedProvider = key;
            }
          });
          self.setState({ availableProviders: result, selectedProvider });
        }
      },
      (error) => {
        console.log(`there was an error fetching providers: ${error}`);
      }
    );
  }
  componentDidMount = () => {
    this.loadProvidersFromServer();
  };

  handleMenuItemClick = (index) => {
    this.setState({ selectedProvider: index });
  };

  handleToggle() {
    const self = this;
    return () => {
      self.setState({ open: !self.state.open });
    };
  }

  handleClose() {
    const self = this;
    return (event) => {
      if (self.anchorRef && self.anchorRef.contains(event.target)) {
        return;
      }
      self.setState({ open: false });
    };
  }

  handleModalOpen() {
    const self = this;
    return () => {
      self.setState({ modalOpen: true });
    };
  }

  handleModalClose() {
    const self = this;
    return () => {
      self.setState({ modalOpen: false });
    };
  }

  render() {
    const { classes } = this.props;
    const { availableProviders, selectedProvider, open, modalOpen } = this.state;
    const self = this;
    return (
      <NoSsr>
        <div data-cy="root" className={classes.root}>
          <img
            className={classes.logo}
            src="/provider/static/img/meshery-logo/meshery-logo-light-text.png"
            alt="logo"
          />
          <Typography variant="h6" gutterBottom className={classes.chartTitle}>
            Please choose a
            <Tooltip title="Learn more about providers" placement="bottom" data-cy="providers-tooltip">
              <a className={classes.providerLink} onClick={self.handleModalOpen()}>
                {" "}
                provider{" "}
              </a>
            </Tooltip>
            to continue
          </Typography>
          <Dialog
            onClose={self.handleModalClose()}
            aria-labelledby="customized-dialog-title"
            open={modalOpen}
            disableScrollLock={true}
            data-cy="providers-modal"
          >
            <DialogTitle id="customized-dialogs-title" onClose={self.handleModalClose()}>
              <b>Choosing a provider</b>
            </DialogTitle>
            <DialogContent dividers>
              <Typography gutterBottom>
                <p>
                  Login to Meshery by choosing from the available providers. Providers offer authentication, session
                  management and long-term persistence of user preferences, performance tests, service mesh adapter
                  configurations and so on.
                </p>
                {Object.keys(availableProviders).map((key) => {
                  return (
                    <React.Fragment key={availableProviders[key]["provider_name"]}>
                      <p>{availableProviders[key]["provider_name"]}</p>
                      <ul>
                        {availableProviders[key]["provider_description"]?.map((desc, i) => <li key={`desc-${i}`}>{desc}</li>)}
                      </ul>
                    </React.Fragment>
                  );
                })}
                <p>Provider: SMI Conformance</p>
                <ul>
                  <li>Remote provider for SMI Conformance Testing</li>
                  <li>Provides provenence of test results and their persistence</li>
                </ul>
                <p>Provider: The University of Texas at Austin</p>
                <ul>
                  <li>Academic research and advanced studies by Ph.D. researchers</li>
                  <li>Used by school of Electrical and Computer Engineering (ECE)</li>
                </ul>
              </Typography>
            </DialogContent>
            <DialogActions>
              <Button autoFocus onClick={self.handleModalClose()} color="primary" data-cy="providers-modal-button-ok">
                OK
              </Button>
            </DialogActions>
          </Dialog>
          <div className={classes.container}>
            <Grid item xs={12} justify="center">
              {availableProviders !== "" && (
                <>
                  <ButtonGroup
                    variant="contained"
                    color="primary"
                    ref={(ref) => (self.anchorRef = ref)}
                    aria-label="split button"
                  >
                    <Button
                      size="large"
                      href={
                        selectedProvider == "" ? "" : `/api/provider?provider=${encodeURIComponent(selectedProvider)}`
                      }
                    >
                      {selectedProvider !== "" ? selectedProvider : "Select Your Provider"}
                    </Button>
                    <Button
                      color="primary"
                      size="small"
                      aria-controls={open ? "split-button-menu" : undefined}
                      aria-expanded={open ? "true" : undefined}
                      aria-label="Select Provider"
                      data-cy="select_provider"
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
                          transformOrigin: placement === "bottom" ? "center top" : "center bottom",
                        }}
                      >
                        <Paper>
                          <ClickAwayListener onClickAway={self.handleClose()}>
                            <MenuList id="split-button-menu">
                              {Object.keys(availableProviders).map((key) => (
                                <MenuItem key={key} onClick={() => self.handleMenuItemClick(key)}>
                                  {key}
                                </MenuItem>
                              ))}
                              <Divider />
                              <MenuItem disabled={true} key="SMI">
                                SMI Conformance
                              </MenuItem>
                              <MenuItem disabled={true} key="UT Austin">
                                The University of Texas at Austin
                              </MenuItem>
                            </MenuList>
                          </ClickAwayListener>
                        </Paper>
                      </Grow>
                    )}
                  </Popper>
                </>
              )}
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
