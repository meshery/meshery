import React, { useEffect, useState } from 'react'
import dataFetch from '../lib/data-fetch'
import ProviderLayout from './ProviderLayout'
import { Div, MesheryLogo, MenuProviderDisabled } from './Provider.style'
import Typography from '@mui/material/Typography'
import Tooltip from '@mui/material/Tooltip'
import Dialog from '@mui/material/Dialog'
import MuiDialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'
import Button from '@mui/material/Button'
import ButtonGroup from '@mui/material/ButtonGroup'
import Divider from '@mui/material/Divider'
import Grid from '@mui/material/Grid'
import IconButton from '@mui/material/IconButton'
import MenuList from '@mui/material/MenuList'
import MenuItem from '@mui/material/MenuItem'
import Menu from '@mui/material/Menu'
import Paper from '@mui/material/Paper'
import CircularProgress from '@mui/material/CircularProgress'
import ClickAwayListener from '@mui/material/ClickAwayListener'
import ArrowDropIcon from '@mui/icons-material/ArrowDropDown'
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
    <MuiDialogTitle disableTypography {...other}>
      <Typography variant="h6">{children}</Typography>
      {onClose
        ? (
          <IconButton aria-label="close" onClick={onClose}>
            <CloseIcon />
          </IconButton>
          )
        : null}
    </MuiDialogTitle>
  )
}

export default function ProviderComponent () {
  const [availableProviders, setAvailableProviders] = useState({})
  const [selectedProvider, setSelectedProvider] = useState('')
  const [open, setOpen] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [anchorEl, setAnchorEl] = React.useState(null);
  const anchorRef = React.useRef(null)

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
  const handleToggle = () => setOpen(true)
  const handleClickAway = () => setOpen(false)

  function handleMenuItemClick (e, provider) {
    e.preventDefault()
    setSelectedProvider(provider)
    setOpen(false)
    setIsLoading(true)
    window.location.href = `/api/provider?provider=${encodeURIComponent(provider)}`
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
          <Button autoFocus onClick={handleModalClose} color="primary" data-cy="providers-modal-button-ok" variant="contained">
            OK
          </Button>
        </DialogActions>
      </Dialog>
      <Div>
        <Grid item xs={12} justify="center">
          {availableProviders !== '' && (
            <>
              <ButtonGroup
                variant="contained"
                color="primary"
                aria-label="split button"
              >
                <Button
                  size="large"
                  id="basic-button"
                  aria-controls={open ? 'split-button-menu' : undefined}
                  aria-expanded={open ? 'true' : undefined}
                  aria-label="Select your provider"
                  data-cy="select_provider"
                  aria-haspopup="menu"
                  onClick={handleToggle}
                >
                  {isLoading && <CircularProgress
                    size={20} sx={{ color: 'white', marginRight: 8 }} />}
                  {selectedProvider !== '' ? selectedProvider : 'Select your provider'}
                  <ArrowDropIcon />
                </Button>
              </ButtonGroup>
              <Menu
                id="basic-menu"
                anchorEl={anchorEl}
                open={open}
                onClose={handleToggle}
                MenuListProps={{
                  'aria-labelledby': 'basic-button'
                }}
              >
                <Paper sx={{ justify: "center" }}>
                  <ClickAwayListener onClickAway={handleClickAway}>
                    <MenuList id="split-button-menu">
                      {Object.keys(availableProviders).map((key) => (
                        <MenuItem key={key} onClick={(e) => handleMenuItemClick(e, key)}>
                          {key}
                        </MenuItem>
                      ))}
                      <Divider sx={{
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
              </Menu>
            </>
          )}
        </Grid>
      </Div>
    </ProviderLayout>
  )
}
