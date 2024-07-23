import React, { Fragment, useEffect, useState } from "react";
import PropTypes from "prop-types";
import dataFetch from "../lib/data-fetch";
import ProviderLayout from "./ProviderLayout";
import {
  CustomDiv,
  CustomDialog,
  MesheryLogo,
  MenuProviderDisabled,
  CustomDialogActions,
  LearnMore,
} from "./Provider.style";
import Button from "@mui/material/Button";
import ButtonGroup from "@mui/material/ButtonGroup";
import Divider from "@mui/material/Divider";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import Popover from "@mui/material/Popover";
import MenuList from "@mui/material/MenuList";
import MenuItem from "@mui/material/MenuItem";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import IconButton from "@mui/material/IconButton";
import CircularProgress from "@mui/material/CircularProgress";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import CloseIcon from "@mui/icons-material/Close";
import ClickAwayListener from "@mui/material/ClickAwayListener";
function CustomDialogTitle(props) {
  const { children, onClose, ...other } = props;

  return (
    <DialogTitle sx={{ m : 0, p : 2 }} {...other}>
      {children}
      {onClose ? (
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{
            position : "absolute",
            right : "1rem",
            top : "1rem",
            color : (theme) => theme.palette.grey[500],
          }}
        >
          <CloseIcon />
        </IconButton>
      ) : null}
    </DialogTitle>
  );
}

CustomDialogTitle.propTypes = {
  children : PropTypes.node,
  onClose : PropTypes.func.isRequired,
};

export default function Provider() {
  const [anchorEl, setAnchorEl] = useState(null);
  const [availableProviders, setAvailableProviders] = useState({});
  const [selectedProvider, setSelectedProvider] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  /* eslint-disable no-unused-vars */
  const [openMenu, setOpenMenu] = useState(false);
  const [openModal, setModalOpen] = React.useState(false);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? "popover" : undefined;

  useEffect(() => {
    loadProvidersFromServer();
  }, []);

  function loadProvidersFromServer() {
    dataFetch(
      "/api/providers",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        if (typeof result !== "undefined") {
          Object.keys(result).forEach((key) => {
            if (result[key].ProviderType === "remote") {
              setSelectedProvider(selectedProvider);
            }
          });
          setAvailableProviders(result);
        }
      },
      (error) => {
        console.log(`there was an error fetching providers: ${error}`);
      },
    );
  }

  const handleMenuItemClick = (event, provider) => {
    event.preventDefault();
    setSelectedProvider(provider);
    setOpenMenu(false);
    setIsLoading(true);
    const existingQueryString = window.location.search;
    const existingQueryParams = new URLSearchParams(existingQueryString);
    existingQueryParams.append("provider", encodeURIComponent(provider))
    window.location.href = `/api/provider?${existingQueryParams.toString()}`;
  };

  const handleModalOpen = () => {
    setModalOpen(true);
  };

  const handleModalClose = () => {
    setModalOpen(false);
  };

  return (
    <ProviderLayout>
      <MesheryLogo src="/provider/static/img/meshery-logo/meshery-logo-light-text.png" alt="logo" />
      <CustomDiv>
        {availableProviders !== "" && (
          <Fragment>
            <ButtonGroup variant="contained" aria-label="split button" color="primary">
              <Button
                size="large"
                aria-describedby={id}
                variant="contained"
                onClick={handleClick}
                aria-label="Select Provider"
                data-cy="select_provider"
                disableElevation
              >
                {isLoading && <CircularProgress size={20} sx={{ color : "white", marginRight : 8 }} />}
                {selectedProvider !== "" ? selectedProvider : "Select your provider"}
                <ArrowDropDownIcon />
              </Button>
            </ButtonGroup>
            <Popover
              id={id}
              open={open}
              anchorEl={anchorEl}
              onClose={handleClose}
              anchorOrigin={{
                vertical : "bottom",
                horizontal : "center",
              }}
              transformOrigin={{
                vertical : "top",
                horizontal : "center",
              }}
            >
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList id="split-button-menu" autoFocusItem>
                  {Object.keys(availableProviders).map((key) => (
                    <MenuItem key={key} onClick={(e) => handleMenuItemClick(e, key)}>
                      {key}
                    </MenuItem>
                  ))}
                  <Divider sx={{ my : 0.5 }} />
                  <MenuProviderDisabled disabled={true} key="CNCF Labs">
                    CNCF Labs{"\u00A0"}
                    <span>Offline</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="Equinix US-DAL">
                    Equinix US-DAL{"\u00A0"}
                    <span>Offline</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="HPE Security">
                    HPE Security{"\u00A0"}
                    <span>Offline</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="MIT">
                    Massachusetts Institute of Technology (MIT){"\u00A0"}
                    <span>Offline</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="UT Austin">
                    The University of Texas at Austin{"\u00A0"}
                    <span> Offline</span>
                  </MenuProviderDisabled>
                </MenuList>
              </ClickAwayListener>
            </Popover>
          </Fragment>
        )}
      </CustomDiv>
      <LearnMore>
        <Typography variant="h6" sx={{ fontWeight : 500 }} gutterBottom>
          Learn more about
          <Tooltip
            title="Learn more about providers"
            placement="bottom"
            data-cy="providers-tooltip"
            sx={{
              color : "#00B39F",
              cursor : "pointer",
              fontWeight : 700,
            }}
          >
            <a
              onClick={handleModalOpen}
              style={{
                color : "#00B39F",
                cursor : "pointer",
                fontWeight : 700,
              }}
            >
              {" "}
              providers{" "}
            </a>
          </Tooltip>
        </Typography>
      </LearnMore>
      <CustomDialog
        onClose={handleModalClose}
        aria-labelledby="customized-dialog-title"
        open={openModal}
        disableScrollLock={true}
        data-cy="providers-modal"
      >
        <CustomDialogTitle id="customized-dialog-title" onClose={handleModalClose} style={{ background : "#eee" }}>
          <b>Choosing a Provider</b>
        </CustomDialogTitle>
        <DialogContent dividers>
          <DialogContentText id="customized-dialog-content">
            Login to Meshery by choosing from the available providers. Providers extend Meshery by offering various
            plugins and services, including identity services, long-term persistence, advanced performance analysis,
            multi-player user collaboration, and so on.
            <h3>Available Providers</h3>
            {Object.keys(availableProviders).map((key) => {
              return (
                <React.Fragment key={availableProviders[key].provider_name}>
                  <p style={{ fontWeight : 700 }}>{availableProviders[key].provider_name}</p>
                  <ul>
                    {availableProviders[key].provider_description?.map((desc, i) => (
                      <li key={`desc-${i}`}>{desc}</li>
                    ))}
                  </ul>
                </React.Fragment>
              );
            })}
            <p style={{ fontWeight : 700 }}>MIT</p>
            <ul>
              <li>Remote provider for performance testing</li>
              <li>Provides provenence of test results and their persistence</li>
              <li>Adaptive performance analysis - predictive optimization</li>
            </ul>
            <p style={{ fontWeight : 700 }}>The University of Texas at Austin</p>
            <ul>
              <li>Academic research and advanced studies by Ph.D. researchers</li>
              <li>Used by school of Electrical and Computer Engineering (ECE)</li>
            </ul>
            <p style={{ fontWeight : 700 }}>Cloud Native Computing Foundation Infrastructure Lab</p>
            <ul>
              <li>Performance and compatibility-centric research and validation</li>
              <li>Used by various service meshes and by the Service Mesh Performance project</li>
            </ul>
            <p style={{ fontWeight : 700 }}>HPE Security</p>
            <ul>
              <li>Istio, SPIRE, and SPIFEE integration</li>
            </ul>
            <p style={{ fontWeight : 700 }}>Equinix</p>
            <ul>
              <li>Identity services</li>
              <li>Bare-metal Kubernetes configuration</li>
            </ul>
          </DialogContentText>
        </DialogContent>
        <CustomDialogActions>
          <div className="learnmore">
            <a href="https://docs.meshery.io/extensibility/providers">
              Providers in Meshery Docs
              <img src="/provider/static/img/external-link.svg" width="16px" />
            </a>
          </div>

          <Button onClick={handleModalClose} color="primary" data-cy="providers-modal-button-ok" variant="contained">
            {" "}
            OK
          </Button>
        </CustomDialogActions>
      </CustomDialog>
    </ProviderLayout>
  );
}
