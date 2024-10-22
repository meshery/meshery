import React, { useState, useEffect } from 'react';
import { NoSsr, Typography, makeStyles, Link } from '@material-ui/core';

const styles = makeStyles((theme) => ({
  rootClass: {
    backgroundColor: theme.palette.type === 'dark' ? theme.palette.secondary.headerColor : '#fff',
    padding: '2rem',
    textAlign: 'center',
    borderRadius: 4,
    height: '100%',
  },
  errorSection: {},
  message: {
    fontSize: '3rem',
    lineHeight: '2rem',
    marginBottom: '2rem',
  },
  errMessage: {
    fontWeight: '400',
    fontSize: '1.5rem',
    color: 'gray',
    fontStyle: 'italic',
    marginTop: '2.5rem',
  },

  discussionforumlink: {
    color: theme.palette.type === 'dark' ? '#00B39F' : '#',
  },

  mesh: {
    display: 'block',
    margin: 'auto',
    marginTop: '3.125rem',
    maxWidth: '50%',
    height: '45%',
  },
  helpMessage: {
    marginTop: '5rem',
    color: theme.palette.type === 'dark' ? '#fff' : '#222',
  },
}));

const customMessages = [
  'Oh, no. Please pardon our meshy app.',
  'Oops. Please excuse the mesh.',
  'Things tend to get a bit meshy around here.',
  'Please pardon our mesh.',
  'How did this mesh happen?',
  "Well, isn't this a mesh?",
  'Yikes. Things are a mesh here.',
];

/**
 * CustomErrorMessage component is used to display a custom error message when a page is not found.
 * @returns {JSX.Element} JSX.Element
 */
function CustomErrorMessage({ message }) {
  const classes = styles();
  const [customMessage, setCustomMessage] = useState(customMessages[0]);

  useEffect(() => {
    setCustomMessage(customMessages[Math.floor(Math.random() * customMessages.length)]);
  }, []);

  return (
    <NoSsr>
      <div className={classes.rootClass}>
        <div className={classes.errorSection}>
          <Typography variant="h1">
            <div className={classes.message}>{customMessage}</div>
          </Typography>
          <Typography variant="h5">
            <div className={classes.errMessage}> {message || 'Page does not exist.'}</div>
          </Typography>
        </div>
        <img src="/static/img/service-mesh.svg" alt="service meshed" className={classes.mesh} />
        <Typography variant="body1">
          <p className={classes.helpMessage}>
            Start a conversation at Layer5 community{' '}
            <Link
              className={classes.discussionforumlink}
              underline="none"
              href="http://discuss.meshery.io/c/meshery/5"
              target="_blank"
            >
              discussion forum
            </Link>
            .
          </p>
        </Typography>
      </div>
    </NoSsr>
  );
}

export default CustomErrorMessage;
