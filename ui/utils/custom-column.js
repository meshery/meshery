import React, { useState } from 'react';
import IconButton from '@mui/material/IconButton';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import Tooltip from '@mui/material/Tooltip';
import Popper from '@mui/material/Popper';
import ClickAwayListener from '@mui/material/ClickAwayListener';
import { withStyles } from '@material-ui/core/styles';
import ColumnIcon from '../assets/icons/coulmn';
// import Slide from '@mui/material/Slide';
import Box from '@mui/material/Box';
import { Card } from '@material-ui/core';
import CheckBoxOutlineBlankIcon from '@mui/icons-material/CheckBoxOutlineBlank';

const styles = (theme) => ({
  epaper: {
    '&: .MuiPaper-root': {
      background: theme.palette.secondary.headerColor,
      color: theme.palette.secondary.textMain,
    },
    padding: '1rem',
    boxShadow: `0px 4px 8px ${
      theme.palette.type === 'light'
        ? theme.palette.secondary.disabledIcon
        : '0px 4px 5px 0px rgba(0,0,0,0.14)'
    }`,
  },
  icon: {
    color: theme.palette.secondary.iconMain,
  },
});

const CustomColumnVisibilityControl = ({ columns, customToolsProps, classes }) => {
  const [anchorEl, setAnchorEl] = useState(null);

  const handleOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleColumnVisibilityChange = (columnName, isVisible) => {
    customToolsProps.setColumnVisibility((prevState) => ({
      ...prevState,
      [columnName]: isVisible,
    }));
  };

  return (
    <div>
      <Tooltip title="View Columns">
        <IconButton
          onClick={handleOpen}
          sx={{
            '&:hover': {
              borderRadius: '4px',
            },
          }}
        >
          <ColumnIcon />
        </IconButton>
      </Tooltip>

      <Box sx={{ overflow: 'hidden' }}>
        <Popper
          style={{ zIndex: 100 }}
          open={Boolean(anchorEl)}
          anchorEl={anchorEl}
          placement="bottom-start"
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'center',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'center',
          }}

          // These code is commented out now because it was having anomalies in the UI at different sections while slide in,
          // had to find better solution for this which will taken care from sistent by spring animation because the clickawaylistner
          // having issues while having inline transition style/slider in the popper.

          // transition
          // popperOptions={{
          //   modifiers: [
          //     {
          //       name: 'offset',
          //       options: {
          //         offset: [0, 8],
          //       },
          //     },
          //     {
          //       name: 'preventOverflow',
          //       options: {
          //         altAxis: true, // true by default
          //         tether: true,
          //         rootBoundary: 'document',
          //         padding: 8,
          //       },
          //     },
          //   ],
          // }}
        >
          {/* {({ TransitionProps }) => (
            <Slide
              {...TransitionProps}
              direction="down"
              in={open}
              timeout={transitionDuration}
              mountOnEnter
              unmountOnExit
              container={containerRef.current}
              easing="cubic-bezier(0.25, 0.1, 1, 1)"
            > */}
          <Box>
            <ClickAwayListener onClickAway={handleClose}>
              <Card className={classes.epaper} elevation={2}>
                <div style={{ display: 'flex', flexDirection: 'column' }}>
                  {columns.map((col) => (
                    <FormControlLabel
                      key={col.name}
                      control={
                        <Checkbox
                          className={classes.checkbox}
                          checked={customToolsProps.columnVisibility[col.name]}
                          onChange={(e) => handleColumnVisibilityChange(col.name, e.target.checked)}
                          icon={<CheckBoxOutlineBlankIcon className={classes.icon} />}
                          sx={{
                            '&.Mui-checked': {
                              color: '#00B39F',
                            },
                          }}
                        />
                      }
                      label={col.label}
                    />
                  ))}
                </div>
              </Card>
            </ClickAwayListener>
          </Box>
          {/* </Slide>
          )} */}
        </Popper>
      </Box>
    </div>
  );
};

export default withStyles(styles)(CustomColumnVisibilityControl);
