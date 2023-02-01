import React from 'react'
import Link from 'next/link'
import { Button, Typography, Container, withStyles } from '@material-ui/core/'
import CheckCircleIcon from '@material-ui/icons/CheckCircle'

const styles = () => ({
  container : {
    position : 'relative',
    height : '30rem',
    textAlign : 'center',
    marginTop : "3rem"
  },
  checkCircleIcon : {
    color : '#00B39F',
    padding : '1rem',
    height : 'auto',
    width : '4rem',
  },
  subtitle : {
    fontWeight : '300',
  },
  link : {
    textDecoration : 'none',
  },
  returnButton : {
    background : '#647881',
    color : 'white',
    marginTop : '2rem',
    padding : '1rem 2rem',
    '&:hover' : {
      background : '#647881',
      color : 'white',
    }
  },
  startOverButton : {
    position : 'absolute',
    right : '1rem',
    bottom : '2rem',
    padding : '0.5rem 2rem',
    textDecoration : 'none',
    background : 'white',
    color : '#647881',
    border : '1.5px solid #647881',
    '&:hover' : {
      backgroundColor : '#647881',
      color : 'white',
    },
  },
});

const ConfigurationDoneScreen = ({ handleUserClick, classes }) => {
  return (
    <Container className={classes.container}>
      <CheckCircleIcon className={classes.checkCircleIcon} />
      <Typography variant="h4" gutterBottom="true">
        Configuration done
      </Typography>
      <Typography
        variant="subtitle1"
        paragraph="true"
        gutterBottom="true"
        className={classes.subtitle}
      >
        You are ready to manage your cloud native infrastructure
      </Typography>
      <Link href="/" className={classes.link}>
        <Button className={classes.returnButton}>Dashboard</Button>
      </Link>
      <Button
        onClick={() => handleUserClick(0)}
        className={classes.startOverButton}
      >
        Start Over
      </Button>
    </Container>
  )
}

export default withStyles(styles)(ConfigurationDoneScreen);
