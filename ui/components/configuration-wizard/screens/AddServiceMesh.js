import React from 'react'
import { makeStyles, Container, Fade, Button } from '@material-ui/core/'

import ConfigCard from './ConfigCard'
import ConsulIcon from '../icons/ConsulIcon.js'
import LinkerdIcon from '../icons/LinkerdIcon.js'
import OpenServiceMeshIcon from '../icons/OpenServiceMeshIcon.js'
import InfoContainer from '../InfoContainer'

const useStyles = makeStyles({
  contentContainer: {
    position: 'relative',
    display: 'flex',
    flexDirection: 'row',
    height: '32.5rem',
    justifyContent: 'space-evenly',
    padding: '2rem 6rem 0rem 6rem',
  },
  dataContainer: {
    display: 'flex',
    flexDirection: 'column',
    alignContent: 'center',
    marginRight: '1rem',
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

const ServiceMesh = () => {
  const [state, setState] = React.useState({})
  const classes = useStyles()

  const handleSwitch = (name, value) => {
    console.log(state)
    setState({
      ...state,
      [name]: value,
    })
  }
  return (
    <Fade timeout={{ enter: '500ms' }} in='true'>
      <Container className={classes.contentContainer}>
        <div className={classes.dataContainer}>
          <ConfigCard
            handleSwitch={handleSwitch}
            name='openServiceMesh'
            Icon={OpenServiceMeshIcon}
            topInputPlaceholder='URL'
            bottomInputPlaceholder=''
          />
        </div>
        <div className={classes.dataContainer}>
          <ConfigCard
            handleSwitch={handleSwitch}
            name='consul'
            Icon={ConsulIcon}
            topInputPlaceholder='URL'
            bottomInputPlaceholder='Context'
          />
          {!state.consul ? null : (
            <InfoContainer controlPlane='2' dataPlane='9' />
          )}
        </div>
        <div className={classes.dataContainer}>
          <ConfigCard
            handleSwitch={handleSwitch}
            name='linkerd'
            Icon={LinkerdIcon}
            topInputPlaceholder='URL'
            bottomInputPlaceholder='Context'
          />
          {!state.linkerd ? null : (
            <InfoContainer controlPlane='4' dataPlane='6' />
          )}
        </div>
        <Button className={classes.settingsButton}>Advanced Settings</Button>
      </Container>
    </Fade>
  )
}

export default ServiceMesh
