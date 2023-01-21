import React, { useEffect, useState } from 'react'
import dataFetch from '../lib/data-fetch'
import ProviderLayout from './ProviderLayout'
import { styled } from '@mui/material'
import { Div, MesheryLogo, CustomDialogTitle, MenuProviderDisabled, CustomIconButton } from './Provider.style'
import Typography from '@mui/material/Typography'
import Tooltip from '@mui/material/Tooltip'
import Dialog from '@mui/material/Dialog'
// import MuiDialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'
import Button from '@mui/material/Button'
import ButtonGroup from '@mui/material/ButtonGroup'
import Divider from '@mui/material/Divider'
import Grid from '@mui/material/Grid'
import IconButton from '@mui/material/IconButton'
import MenuList from '@mui/material/MenuList'
import MenuItem from '@mui/material/MenuItem'
import Paper from '@mui/material/Paper'
// import Popper from '@mui/material/Popper'
import PopperUnstyled from '@mui/base/PopperUnstyled';
import Grow from '@mui/material/Grow'
import CircularProgress from '@mui/material/CircularProgress'
import ClickAwayListener from '@mui/material/ClickAwayListener'
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import CloseIcon from '@mui/icons-material/Close'

// Call providers via getServerSideProps
// This does not do anything at the moment
export async function getServerSideProps () {
  const res = await fetch('http://localhost:9081/api/providers')
  const data = res.json()
  return { props: { providers: data } }
}

function DialogTitle ({ children, onClose, ...other }) {
  return (
    <CustomDialogTitle disableTypography {...other}>
      <Typography variant="h6">{children}</Typography>
      {onClose
        ? (
          <CustomIconButton aria-label="close" onClick={onClose}>
            <CloseIcon />
          </CustomIconButton>
          )
        : null}
    </CustomDialogTitle>
  )
}

const Popper = styled(PopperUnstyled)`
  z-index: 1;
`;

export default function ProviderComponent () {
  const [availableProviders, setAvailableProviders] = useState({})
  const [selectedProvider, setSelectedProvider] = useState('')
  // const [open, setOpen] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    loadProvidersFromServer()
  }, [])

  function loadProvidersFromServer () {
    dataFetch(
      '/api/providers',
      {
        method: 'GET',
        credentials: 'include'
      },
      (result) => {
        if (typeof result !== 'undefined') {
          let selectedProvider = ''
          Object.keys(result).forEach((key) => {
            if (result[key].ProviderType === 'remote') {
              selectedProvider = key
            }
          })
          setAvailableProviders(result)
          setSelectedProvider(selectedProvider)
        }
      },
      (error) => {
        console.log(`there was an error fetching providers: ${error}`)
      }
    )
  }

  const handleModalOpen = () => setModalOpen(true)
  const handleModalClose = () => setModalOpen(false)
  // const handleToggle = () => setOpen(true)
  // const handleClickAway = () => setOpen(false)

  const [open, setOpen] = React.useState(false)
  const anchorRef = React.useRef(null)

  const handleMenuItemClick = (event, provider) => {
    event.preventDefault()
    setSelectedProvider(provider)
    setOpen(false)
    setIsLoading(true)
    window.location.href = `/api/provider?provider=${encodeURIComponent(provider)}`
  }

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen)
  }

  const handleClose = (event) => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return
    }

    setOpen(false)
  }

  return (
    <ProviderLayout>
      <MesheryLogo src="/provider/static/img/meshery-logo/meshery-logo-light-text.png"
        alt="logo"
      />
      <Typography variant="h6" sx={{ fontWeight: 700 }} gutterBottom>
        Please choose a
        <Tooltip title="Learn more about providers" placement="bottom" data-cy="providers-tooltip">
          <a onClick={handleModalOpen} style={{
            color: 'darkcyan',
            cursor: 'pointer',
            fontWeight: 700
          }}>
            {' '}
            provider{' '}
          </a>
        </Tooltip>
        to continue
      </Typography>
      <Dialog
        onClose={handleModalClose}
        aria-labelledby="customized-dialog-title"
        open={modalOpen}
        disableScrollLock={true}
        data-cy="providers-modal"
      >
        <DialogTitle id="customized-dialog-title" onClose={handleModalClose}>
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
                <React.Fragment key={availableProviders[key].provider_name}>
                  <p style={{ fontWeight: 700 }}>{availableProviders[key].provider_name}</p>
                  <ul>
                    {availableProviders[key].provider_description?.map((desc, i) => <li key={`desc-${i}`}>{desc}</li>)}
                  </ul>
                </React.Fragment>
              )
            })}
            <p style={{ fontWeight: 700 }}>SMI Conformance</p>
            <ul>
              <li>Remote provider for SMI Conformance Testing</li>
              <li>Provides provenence of test results and their persistence</li>
            </ul>
            <p style={{ fontWeight: 700 }}>The University of Texas at Austin</p>
            <ul>
              <li>Academic research and advanced studies by Ph.D. researchers</li>
              <li>Used by school of Electrical and Computer Engineering (ECE)</li>
            </ul>
            <p style={{ fontWeight: 700 }}>Cloud Native Computing Foundation Infrastructure Lab</p>
            <ul>
              <li>Performance and compatibility-centric research and validation</li>
              <li>Used by various service meshes and by the Service Mesh Performance project</li>
            </ul>
          </Typography>
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
      </Dialog>
      <Div>
        <Grid item xs={12} justify="center">
          {availableProviders !== '' && (
            <>
              <ButtonGroup variant="contained" ref={anchorRef} aria-label="split button">
                <Button
                  size="large"
                  aria-controls={open ? 'split-button-menu' : undefined}
                  aria-expanded={open ? 'true' : undefined}
                  aria-label="select your provider"
                  data-cy="select_provider"
                  aria-haspopup="menu"
                  onClick={handleToggle}
                >
                  {isLoading && <CircularProgress
                    size={20} sx={{ color: 'white', marginRight: 8 }} />}
                  {selectedProvider !== '' ? selectedProvider : 'Select your provider'}
                  <ArrowDropDownIcon />
                </Button>
              </ButtonGroup>
              <Popper
                open={open}
                anchorEl={anchorRef.current}
                role={undefined}
                transition
                disablePortal
              >
                {({ TransitionProps, placement }) => (
                  <Grow
                    {...TransitionProps}
                    style={{
                      transformOrigin: placement === 'bottom' ? 'center top' : 'center bottom'
                    }}
                  >
                    <Paper>
                      <ClickAwayListener onClickAway={handleClose}>
                        <MenuList id="split-button-menu" autoFocusItem>
                        {Object.keys(availableProviders).map((key) => (
                             <MenuItem key={key} onClick={(e) => handleMenuItemClick(e, key)}>
                              {key}
                            </MenuItem>
                        ))}
                          <Divider
                            sx={{
                              backgroundColor: '#c1c8d2',
                              marginLeft: '10px',
                              marginRight: '10px'
                            }} />
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
                    </Paper>
                  </Grow>
                )}
              </Popper>
            </>)}
        </Grid>
      </Div>
    </ProviderLayout>
  )
}
