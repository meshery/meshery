import React from 'react'
import { makeStyles, Container, Typography } from '@material-ui/core/'
import BackupIcon from '@material-ui/icons/Backup'

import KubernetesIcon from '../icons/KubernetesIcon'
import ConfigCard from './ConfigCard'

const useStyles = makeStyles({
  cardContainer: {
    position: 'relative',
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '2rem 6rem',
  },
  infoContainer: {
    position: 'relative',
    display: 'flex',
    flexDirection: 'column',
    width: '20rem',
    padding: '5rem 2rem',
    marginTop: '2rem',
    boxShadow: '0px 1px 6px 1px rgba(0,0,0,0.75)',
  },
  infoStatus: {
    position: 'absolute',
    bottom: '10.50rem',
    right: '10rem',
    color: '#647881',
    background: '#F1F3F4',
    padding: '.5rem 5rem .75rem 1.5rem',
    borderRadius: '0.25rem',
    fontSize: '.8rem',
  },
  infoContext: {
    fontSize: '.9rem',
  },
  infoKind: {
    fontSize: '.75rem',
    color: '#CACACA',
  },
})

const Kubernetes = ({ handleConnectToKubernetes }) => {
  const [state, setState] = React.useState(false)
  const classes = useStyles()

  const handleSwitch = (name, checked) => {
    setState(checked)
    if (handleConnectToKubernetes) {
      handleConnectToKubernetes(checked)
    }
  }

  return (
    <Container className={classes.cardContainer}>
      {' '}
      <ConfigCard
        handleSwitch={handleSwitch}
        name='Kubernetes'
        Icon={KubernetesIcon}
        topInputPlaceholder='Upload Kubeconfig'
        TopInputIcon={BackupIcon}
        bottomInputPlaceholder='Current-Context'
      />
      {!state ? null : (
        <div className={classes.infoContainer}>
          <Typography className={classes.infoStatus}>Status</Typography>
          <Typography className={classes.infoContext}>
            Current-Context: bob-us-east
          </Typography>
          <Typography className={classes.infoContext}>Clust</Typography>
          <Typography className={classes.infoContext}>
            Current Context:{' '}
          </Typography>
        </div>
      )}
    </Container>
  )
}

export default Kubernetes
