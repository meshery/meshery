import { makeStyles } from '@material-ui/core';
import React from 'react';

const useStyles = makeStyles(() => ({
  root: {
    width: '100%',
    color: '#fff',
    position: 'fixed',
    top: 80,
    // left: 0,
    backgroundColor: '#477E96',
    padding: '4px 50px',
    transform: 'translateX(-40px)',
  },
  span: {
    // marginLeft: 300,
    color: '#fff',
    fontStyle: 'italic',
    '&:hover': {
      cursor: 'pointer',
      textDecoration: 'underline'
    }
  }
}));

export default function CustomBreadCrumb({ title, onBack }) {
  const classes = useStyles();

  return (
    <div
      className={classes.root}
    >
      {'> '}
      <span
        className={classes.span}
        onClick={onBack}
      >
        Patterns
      </span>
      {' > '}
      <span
        className={classes.span}
      >{title}</span>
    </div>
  );
}