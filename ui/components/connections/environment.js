import React, { useState } from 'react';
// import PropTypes from 'prop-types';
import Popper from '@material-ui/core/Popper';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import CloseIcon from '@material-ui/icons/Close';
import DoneIcon from '@material-ui/icons/Done';
import Autocomplete, { createFilterOptions } from '@material-ui/lab/Autocomplete';
// import { autocompleteClasses } from '@material-ui/lab/Autocomplete';
import ButtonBase from '@material-ui/core/ButtonBase';
import InputBase from '@material-ui/core/InputBase';
import Box from '@material-ui/core/Box';
import Tooltip from '@material-ui/core/Tooltip';
// import AddIcon from '../../assets/icons/AddIcon';
import { withStyles } from '@material-ui/core/styles';
import { styled } from '@material-ui/core/styles';

const filter = createFilterOptions();

const styles = () => ({
  iconSelected: {
    width: 17,
    height: 17,
    marginRight: 5,
    marginLeft: -2,
  },
  color: {
    width: 14,
    height: 14,
    flexShrink: 0,
    borderRadius: 3,
    marginRight: 8,
    marginTop: 2,
  },
  text: {
    flexGrow: 1,
  },
  close: {
    opacity: 0.6,
    width: 18,
    height: 18,
    // marginLeft: 15,
  },
});

const StyledPopper = styled(Popper)(({ theme }) => ({
  // border: `1px solid ${theme.palette.mode === 'light' ? '#e1e4e8' : '#30363d'}`,
  boxShadow: `0 4px 10px ${
    theme.palette.mode === 'light' ? 'rgba(149, 157, 165, 0.2)' : 'rgb(1, 4, 9)'
  }`,
  borderRadius: 6,
  width: 270,
  // zIndex: theme.zIndex.modal,
  // fontSize: 13,
  color: theme.palette.secondary.text,
  backgroundColor:
    theme.palette.type === 'dark'
      ? theme.palette.secondary.toolbarBg2
      : theme.palette.secondary.toolbarBg1,
}));

const StyledInput = styled(InputBase)(({ theme }) => ({
  padding: 10,
  width: '100%',
  // borderBottom: `1px solid ${theme.palette.mode === 'light' ? '#eaecef' : '#30363d'}`,
  '& input': {
    borderRadius: 4,
    color: theme.palette.secondary.text,
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
    padding: 8,
    // transition: theme.transitions.create(['border-color', 'box-shadow']),
    border: `1px solid ${theme.palette.mode === 'light' ? '#eaecef' : '#30363d'}`,
    // fontSize: 14,
    '&:focus': {
      boxShadow: `0px 0px 0px 1px #00b39f`,
      borderColor: '#00b39f',
    },
  },
}));

const Button = styled(ButtonBase)(({ theme }) => ({
  // fontSize: 13,
  width: '100%',
  textAlign: 'left',
  paddingBottom: 8,
  color: theme.palette.secondary.text,
  fontWeight: 600,
  '&:hover,&:focus': {
    // color: theme.palette.mode === 'light' ? '#0366d6' : '#58a6ff',
    '& .arrow-icon': {
      display: 'inline', // Show the icon on hover
    },
  },
  '& span': {
    width: '100%',
  },
  '& svg': {
    width: 16,
    height: 16,
  },
}));

const StyledArrowDropDownIcon = styled(ArrowDropDownIcon)({
  display: 'none', // Initially hide the icon
  verticalAlign: 'middle',
  fontSize: 24,
});

function GitHubLabel({ classes }) {
  const [anchorEl, setAnchorEl] = useState(null);
  const [value, setValue] = useState([]);
  const [inputValue, setInputValue] = useState('');

  const handleClick = (event) => {
    // setPendingValue(value);
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    // setValue(pendingValue);
    if (anchorEl) {
      anchorEl.focus();
    }
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'github-label' : undefined;

  return (
    <React.Fragment>
      <Box sx={{ width: 221 }}>
        <Tooltip
          title={
            value.length > 1
              ? value
                  .slice(1)
                  .map((option) => option.title)
                  .join(', ')
              : ''
          }
        >
          <Button disableRipple aria-describedby={id} onClick={handleClick}>
            <span>
              {value.length > 0 ? (
                <>
                  {value.map((option, index) => {
                    // Display the first selected option and count of remaining options
                    if (index === 0) {
                      return option.inputValue ? option.inputValue : option.title;
                    } else if (index === 1) {
                      return ` +${value.length - 1}`;
                    }
                    return null; // Skip the rest of the options
                  })}
                </>
              ) : (
                'Choose an environment'
              )}
            </span>
            <StyledArrowDropDownIcon className="arrow-icon" />
          </Button>
        </Tooltip>
      </Box>
      <StyledPopper id={id} open={open} anchorEl={anchorEl} placement="bottom-start">
        <ClickAwayListener onClickAway={handleClose}>
          <div>
            <Box
              sx={{
                borderBottom: `1px solid`,
                padding: '8px 10px',
                // fontWeight: 600,
              }}
            >
              Search or create environments
            </Box>
            <Autocomplete
              // className={classes.autocompleteOptions}
              value={value}
              onChange={(_, newValue) => {
                setValue(newValue);
                setInputValue('');
              }}
              filterOptions={(options, params) => {
                const filtered = filter(options, params);

                if (params.inputValue !== '') {
                  filtered.push({
                    inputValue: params.inputValue,
                    title: `Create:  ${params.inputValue}`,
                  });
                }

                return filtered;
              }}
              options={environmentOptions}
              getOptionLabel={(option) => {
                if (typeof option === 'string') {
                  return option;
                }
                if (option.inputValue) {
                  return option.inputValue;
                }
                return option.title;
              }}
              selectOnFocus
              clearOnBlur
              handleHomeEndKeys
              isOptionEqualToValue={(option, value) => option.title === value.title}
              renderOption={(option, { selected }) => (
                <React.Fragment>
                  <DoneIcon
                    className={classes.iconSelected}
                    style={{ visibility: selected ? 'visible' : 'hidden' }}
                  />
                  <span className={classes.color} />
                  <div className={classes.text}>{option.title}</div>
                  <CloseIcon
                    className={classes.close}
                    style={{ visibility: selected ? 'visible' : 'hidden' }}
                  />
                </React.Fragment>
              )}
              multiple
              sx={{ '& fieldset': { border: 'none' } }}
              inputValue={inputValue}
              onInputChange={(_, newInputValue) => setInputValue(newInputValue)}
              renderInput={(params) => (
                <StyledInput
                  ref={params.InputProps.ref}
                  inputProps={params.inputProps}
                  autoFocus
                  placeholder="Search environments"
                />
              )}
            />
          </div>
        </ClickAwayListener>
      </StyledPopper>
    </React.Fragment>
  );
}

const environmentOptions = [
  { title: 'Environment 1' },
  { title: 'Environment 2' },
  { title: 'Environment 3' },
  { title: 'Environment 4' },
  { title: 'Environment  ' },
];

export default withStyles(styles)(GitHubLabel);
