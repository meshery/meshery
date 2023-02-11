import React, { Fragment, useEffect, useState } from "react";
import PropTypes from "prop-types";
import dataFetch from "../lib/data-fetch";
import ProviderLayout from "./ProviderLayout"
import { CustomDiv, CustomDialog, MesheryLogo, MenuProviderDisabled } from "./Provider.style";
import Button from "@mui/material/Button"
import ButtonGroup from "@mui/material/ButtonGroup"
import Divider from "@mui/material/Divider"
import DialogActions from "@mui/material/DialogActions"
import DialogContent from "@mui/material/DialogContent"
import DialogContentText from "@mui/material/DialogContentText"
import DialogTitle from "@mui/material/DialogTitle"
import Popover from "@mui/material/Popover"
import MenuList from "@mui/material/MenuList"
import MenuItem from "@mui/material/MenuItem"
import Tooltip from "@mui/material/Tooltip"
import Typography from "@mui/material/Typography"
import IconButton from "@mui/material/IconButton"
import CircularProgress from '@mui/material/CircularProgress'
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import CloseIcon from "@mui/icons-material/Close"
import ClickAwayListener from '@mui/material/ClickAwayListener'

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
            position : 'absolute',
            right : 8,
            top : 8,
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
  const [availableProviders, setAvailableProviders] = useState({})
  const [selectedProvider, setSelectedProvider] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  /* eslint-disable no-unused-vars */
  const [openMenu, setOpenMenu] = useState(false)
  const [openModal, setModalOpen] = React.useState(false)

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  const open = Boolean(anchorEl)
  const id = open ? 'popover' : undefined

  useEffect(() => {
    loadProvidersFromServer()
  }, [])

  function loadProvidersFromServer() {
    dataFetch(
      '/api/providers',
      {
        method : 'GET',
        credentials : 'include'
      },
      (result) => {
        if (typeof result !== 'undefined') {
          Object.keys(result).forEach((key) => {
            if (result[key].ProviderType === 'remote') {
              setSelectedProvider(selectedProvider)
            }
          })
          setAvailableProviders(result)
        }
      },
      (error) => {
        console.log(`there was an error fetching providers: ${error}`)
      }
    )
  }

  const handleMenuItemClick = (event, provider) => {
    event.preventDefault()
    setSelectedProvider(provider)
    setOpenMenu(false)
    setIsLoading(true)
    window.location.href = `/api/provider?provider=${encodeURIComponent(provider)}`
  }

  const handleModalOpen = () => {
    setModalOpen(true)
  }

  const handleModalClose = () => {
    setModalOpen(false);
  };

  return (
    <ProviderLayout>
      <MesheryLogo src="/provider/static/img/meshery-logo/meshery-logo-light-text.png"
        alt="logo"
      />
      <Typography variant="h6" sx={{ fontWeight : 700 }} gutterBottom>
                Please choose a
        <Tooltip
          title="Learn more about providers"
          placement="bottom"
          data-cy="providers-tooltip"
          sx={{
            color : '#00B39F',
            cursor : 'pointer',
            fontWeight : 700
          }}
        >
          <a onClick={handleModalOpen} style={{
            color : 'darkcyan',
            cursor : 'pointer',
            fontWeight : 700
          }}>
            {' '}
                        provider{' '}
          </a>
        </Tooltip>
      </Typography>
      <CustomDialog
        onClose={handleModalClose}
        aria-labelledby="customized-dialog-title"
        open={openModal}
        disableScrollLock={true}
        data-cy="providers-modal"
      >
        <CustomDialogTitle id="customized-dialog-title" onClose={handleModalClose}>
          <b>Choosing a provider</b>
        </CustomDialogTitle>
        <DialogContent dividers>
          <DialogContentText id="customized-dialog-content">
                        Login to Meshery by choosing from the available providers. Providers offer authentication, session
                        management and long-term persistence of user preferences, performance tests, service mesh adapter
                        configurations and so on.

            {Object.keys(availableProviders).map((key) => {
              return (
                <React.Fragment key={availableProviders[key].provider_name}>
                  <p style={{ fontWeight : 700 }}>{availableProviders[key].provider_name}</p>
                  <ul>
                    {availableProviders[key].provider_description?.map((desc, i) => <li key={`desc-${i}`}>{desc}</li>)}
                  </ul>
                </React.Fragment>
              )
            })}
            <p style={{ fontWeight : 700 }}>SMI Conformance</p>
            <ul>
              <li>Remote provider for SMI Conformance Testing</li>
              <li>Provides provenence of test results and their persistence</li>
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
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={handleModalClose}
            color="primary"
            data-cy="providers-modal-button-ok"
            variant="contained"
          >
                        OK
          </Button>
        </DialogActions>
      </CustomDialog>
      <CustomDiv>
        {availableProviders !== '' && (
          <Fragment>
            <ButtonGroup
              variant="contained"
              aria-label="split button"
              color="primary"
            >
              <Button
                size="large"
                aria-describedby={id}
                variant="contained"
                onClick={handleClick}
                aria-label="Select Provider"
                data-cy="select_provider"
                disableElevation
              >
                {isLoading && <CircularProgress
                  size={20} sx={{ color : 'white', marginRight : 8 }} />}
                {selectedProvider !== '' ? selectedProvider : 'Select your provider'}
                <ArrowDropDownIcon />
              </Button>
            </ButtonGroup>
            <Popover
              id={id}
              open={open}
              anchorEl={anchorEl}
              onClose={handleClose}
              anchorOrigin={{
                vertical : 'bottom',
                horizontal : 'center',
              }}
              transformOrigin={{
                vertical : 'top',
                horizontal : 'center',
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
                  <MenuProviderDisabled disabled={true} key="SMI">
                                        SMI Conformance <span>Disabled</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="UT Austin">
                                        The University of Texas at Austin{'\u00A0'}<span>Disabled</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="CNCF Cluster">
                                        CNCF Cluster{'\u00A0'}<span>Disabled</span>
                  </MenuProviderDisabled>
                </MenuList>
              </ClickAwayListener>
            </Popover>
          </Fragment>
        )}
      </CustomDiv>
    </ProviderLayout>
  )
}