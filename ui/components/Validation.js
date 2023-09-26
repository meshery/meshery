import React from 'react';
import {
  Divider,
  List,
  ListItem,
  ListItemText,
  ListSubheader,
  Typography,
  Collapse,
  withStyles,
} from '@material-ui/core';
import { ExpandLess, ExpandMore } from '@material-ui/icons';

const styles = (theme) => ({
  nested: {
    padding: theme.spacing(0.5),
    backgroundColor: '#e0e0e0',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  root: {
    width: '100%',
    maxHeight: '18rem',
  },
  subHeader: {
    marginTop: '1rem',
    display: 'flex',
    justifyContent: 'space-around',
  },
  error: {
    position: 'relative',
    left: '-50px',
    borderRadius: '0.4rem',
    padding: '0.5rem',
    top: '-0.45rem',
  },
});

const ValidationComponent = ({ errors, compCount, classes, handleClose }) => {
  const [open, setOpen] = React.useState([false]);

  const handleClick = (index) => {
    let updatedState = [...open];
    updatedState[index] = !updatedState[index];
    setOpen(updatedState);
  };

  let errorCount =
    errors?.reduce((count, ele) => {
      return ele.errors.length + count;
    }, 0) || 0;

  return (
    <List
      aria-labelledby="nested-list-subheader"
      subheader={
        <ListSubheader
          disableSticky="true"
          component="div"
          id="nested-list-subheader"
          className={classes.subHeader}
        >
          <Typography varaint="h6" disablePadding style={{ position: 'relative', left: '35px' }}>
            {compCount} component{compCount > 1 ? 's' : ''}
            {
              <Divider
                style={{
                  transform: 'rotate(90deg)',
                  width: '33%',
                  top: '-10px',
                  position: 'relative',
                  left: '120px',
                }}
              />
            }
          </Typography>

          <Typography
            varaint="h6"
            disablePadding
            className={classes.error}
            style={{
              border: `2px solid ${errorCount > 0 ? '#F0A303' : '#3fc6b6'}`,
              color: `${errorCount > 0 ? '#F0A303' : '#3fc6b6'}`,
            }}
          >
            error{errorCount > 1 ? 's' : ''}: {errorCount}
          </Typography>
        </ListSubheader>
      }
      className={classes.root}
    >
      {errors?.length > 0 ? (
        errors?.map((err, index) => (
          <div style={{ margin: '0.6rem 0rem' }} key={index}>
            <ListItem
              button
              onClick={() => handleClick(index)}
              style={{ backgroundColor: 'transparent' }}
            >
              <ListItemText primary={err?.service} />({err?.errors.length})
              {open[index] ? <ExpandLess /> : <ExpandMore />}
            </ListItem>
            <Collapse
              in={open[index]}
              timeout="auto"
              unmountOnExit
              onClick={() => {
                handleClose();
                // err.openRJSF();
              }}
            >
              {err?.errors.map((description, index) => (
                <Typography
                  variant="subtitle2"
                  disablePadding
                  className={classes.nested}
                  key={index}
                >
                  {description?.charAt(0).toUpperCase() + description?.slice(1)}
                  {index !== err?.errors.length - 1 ? ', ' : ''}
                </Typography>
              ))}
            </Collapse>
          </div>
        ))
      ) : (
        <Typography varaint="h6" align="center" disablePadding>
          No Validation errors.
        </Typography>
      )}
    </List>
  );
};

export default withStyles(styles)(ValidationComponent);
