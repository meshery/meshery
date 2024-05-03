import React, { useState } from 'react';
import { Popper } from '@mui/material';
import {
  Button,
  Card,
  ClickAwayListener,
  InputLabel,
  MenuItem,
  Select,
  Tooltip,
  IconButton,
} from '@material-ui/core';
import { withStyles } from '@material-ui/core/styles';
import { FilterIcon } from '@layer5/sistent';

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
    fill: theme.palette.secondary.iconMain,
  },
});

const UniversalFilter = ({
  filters,
  selectedFilters,
  setSelectedFilters,
  handleApplyFilter,
  showAllOption = true,
  style,
  classes,
  conditionForMaxHeight = false,
}) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const [open, setOpen] = useState(false);
  const handleFilterChange = (event, columnName) => {
    const value = event.target.value;

    setSelectedFilters((filters) => ({
      ...filters,
      [columnName]: value,
    }));
  };

  const handleApplyOnClick = () => {
    handleClose();
    handleApplyFilter();
  };

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
    setOpen((previousOpen) => !previousOpen);
    // setOpenDialog(true);
  };

  const canBeOpen = open && Boolean(anchorEl);
  const id = canBeOpen ? 'transition-popper' : undefined;

  const handleClose = () => {
    // setAnchorEl(null);
    setOpen(false);
    // setOpenDialog(false);
  };

  return (
    <div>
      <Tooltip title="Filter" arrow>
        <IconButton
          onClick={handleClick}
          sx={{
            '&:hover': {
              '& svg': {
                borderRadius: '4px',
              },
            },
            ...style,
          }}
          disableRipple
        >
          <FilterIcon className={classes.icon} />
        </IconButton>
      </Tooltip>
      <Popper
        //
        id={id}
        open={open}
        anchorEl={anchorEl}
        style={{ zIndex: 100 }}
        placement="bottom-start"
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
      >
        <div>
          <ClickAwayListener
            onClickAway={handleClose}
            mouseEvent="onMouseDown"
            touchEvent="onTouchStart"
          >
            <div>
              <Card elevation={2} className={classes.epaper}>
                {Object.keys(filters).map((filter_column) => {
                  const options = filters[filter_column].options;
                  return (
                    <div key={filter_column} role="presentation">
                      <InputLabel id={filters[filter_column].name} style={{ marginTop: '5px' }}>
                        {filters[filter_column].name}
                      </InputLabel>
                      <Select
                        defaultValue="All"
                        key={filter_column}
                        value={selectedFilters[filter_column]}
                        onChange={(e) => handleFilterChange(e, filter_column)}
                        style={{
                          width: '15rem',
                          marginBottom: '1rem',
                        }}
                        variant="outlined"
                        inputProps={{ 'aria-label': 'Without label' }}
                        displayEmpty
                        MenuProps={
                          conditionForMaxHeight
                            ? {
                                PaperProps: {
                                  style: {
                                    maxHeight: '350px',
                                  },
                                },
                              }
                            : {}
                        }
                      >
                        {showAllOption && <MenuItem value="All">All</MenuItem>}
                        {options.map((option) => {
                          return (
                            <MenuItem key={option.value} value={option.value}>
                              {option.label}
                            </MenuItem>
                          );
                        })}
                      </Select>
                    </div>
                  );
                })}
                <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                  <Button color="primary" variant="contained" onClick={handleApplyOnClick}>
                    Apply
                  </Button>
                </div>
              </Card>
            </div>
          </ClickAwayListener>
        </div>
      </Popper>
    </div>
  );
};

export default withStyles(styles)(UniversalFilter);
