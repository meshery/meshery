import React from 'react'
import {
  makeStyles,
  Container,
  Fade,
  Button
} from '@material-ui/core/'

import GrafanaIcon from '../icons/GrafanaIcon'
import PrometheusIcon from '../icons/PrometheusIcon'
import ConfigCard from './ConfigCard'

const useStyles = makeStyles({
  container: {
    position: 'relative',
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-start',
    height: '35rem',
    padding: '2rem 6rem',
  },
  cardContainer: {
    display: 'flex',
    flexDirection: 'column',
    alignContent: 'center',
    marginRight: '2rem',
  },
  infoContainer: {
    boxSizing: 'border-box',
    position: 'relative',
    width: '12rem',
    height: '8rem',
    padding: '1rem',
    margin: '3.5rem auto',
    boxShadow: '0px 1px 6px 1px rgba(0,0,0,0.75)',
  },
  infoContext: {
    display: 'inline',
    padding: '0.25rem 0.5rem',
    fontSize: '.75rem',
    fontWeight: '300',
    background: 'lightgray',
  },
  settingsButton: {
    position: 'absolute',
    bottom: '-2.5rem',
    left: '1rem',
    color: '#647881',
    border: '1px solid #647881',
    borderRadius: '.8rem',
    fontSize: '.8rem',
    fontWeight: '400',
  },
})

const External = () => {
  const classes = useStyles()
  const [state, setState] = React.useState({})

  const handleSwitch = (name, checked) => {
    setState({
      ...state,
      [name]: checked,
    })
  }

  return (
    <Fade timeout={{ enter: '500ms' }} in='true'>
      <Container className={classes.container}>
        <div className={classes.cardContainer}>
          <ConfigCard
            name='grafana'
            Icon={GrafanaIcon}
            topInputPlaceholder='URL'
            bottomInputPlaceholder='API Key'
            handleSwitch={handleSwitch}
          />
        </div>
        <div className={classes.cardContainer}>
          <ConfigCard
            name='prometheus'
            Icon={PrometheusIcon}
            topInputPlaceholder=''
            bottomInputPlaceholder=''
            handleSwitch={handleSwitch}
          />
        </div>
        <Button className={classes.settingsButton}>Advanced Settings</Button>
      </Container>
    </Fade>
  )
}

export default External
