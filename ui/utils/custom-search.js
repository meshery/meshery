import React, { useState, useRef } from 'react';
import TextField from '@mui/material/TextField';
import { Tooltip } from '@mui/material';
import IconButton from '@mui/material/IconButton';
// import SearchIcon from '../assets/icons/search';
// import CloseIcon from '@mui/icons-material/Close';
import { makeStyles } from '@material-ui/core/styles';
import ClickAwayListener from '@mui/material/ClickAwayListener';
import debounce from './debounce';
import { CloseIcon, SearchIcon } from '@layer5/sistent';

const useStyles = makeStyles((theme) => ({
  icon: {
    fill: theme.palette.secondary.iconMain,
    width: '1.5rem',
    height: '1.5rem',
  },
  searchInput: {
    '& .MuiOutlinedInput-root': {
      color: theme.palette.secondary.iconMain,
    },
    '& .MuiOutlinedInput-notchedOutline': {
      borderColor: theme.palette.secondary.iconMain,
      '&:hover': {
        borderColor: '#00b39f',
      },
    },
    '& .MuiInputLabel-root': {
      color: theme.palette.secondary.iconMain,
    },
    '& .MuiInputBase-input': {
      color: theme.palette.secondary.iconMain,
      caretColor: theme.palette.secondary.iconMain,
    },
    '& .MuiInput-underline:before': {
      borderBottomColor: theme.palette.secondary.iconMain,
    },
    '& .MuiInput-underline:hover:before': {
      borderBottomColor: '#00b39f',
    },
    '& .MuiInput-underline:hover:after': {
      borderBottomColor: '#00b39f',
    },
    '& .MuiInput-underline.Mui-focused:before': {
      borderBottomColor: '#00b39f',
    },
    '& .MuiInput-underline.Mui-focused:after': {
      borderBottomColor: '#00b39f',
    },
  },
}));

const SearchBar = ({ onSearch, placeholder, expanded, setExpanded, value = '' }) => {
  const [searchText, setSearchText] = useState(value);
  const searchRef = useRef(null);
  const classes = useStyles();

  const debouncedOnSearch = debounce(onSearch, 500);

  const handleSearchChange = (event) => {
    debouncedOnSearch(event.target.value);
    setSearchText(event.target.value);
  };

  const handleClearIconClick = () => {
    setSearchText('');
    debouncedOnSearch('');
    setExpanded(false);
  };

  const handleSearchIconClick = () => {
    if (expanded) {
      setSearchText('');
      setExpanded(false);
    } else {
      setExpanded(true);
      setTimeout(() => {
        searchRef.current.focus();
      }, 300);
    }
  };

  const width = window.innerWidth;
  let searchWidth = '200px';
  if (width <= 750) {
    searchWidth = '120px';
  }

  return (
    <div>
      <TextField
        className={classes.searchInput}
        id="searchClick"
        variant="standard"
        value={searchText}
        onChange={handleSearchChange}
        inputRef={searchRef}
        placeholder={placeholder}
        style={{
          width: expanded ? searchWidth : '0',
          opacity: expanded ? 1 : 0,
          transition: 'width 0.3s ease, opacity 0.3s ease',
        }}
      />

      {expanded ? (
        <ClickAwayListener
          onClickAway={(event) => {
            const isTable = event.target.closest('#ref');

            if (searchText !== '') {
              return;
            }
            if (isTable) {
              handleClearIconClick(); // Close the search bar as needed
            }
          }}
        >
          <Tooltip title="Close">
            <IconButton
              onClick={handleClearIconClick}
              sx={{
                '&:hover': {
                  borderRadius: '4px',
                },
              }}
            >
              <CloseIcon className={classes.icon} />
            </IconButton>
          </Tooltip>
        </ClickAwayListener>
      ) : (
        <Tooltip title="Search">
          <IconButton
            onClick={handleSearchIconClick}
            sx={{
              '&:hover': {
                borderRadius: '4px',
              },
            }}
          >
            <SearchIcon className={classes.icon} />
          </IconButton>
        </Tooltip>
      )}
    </div>
  );
};

export default SearchBar;
