import React, { useState, useRef } from 'react';
import { ClickAwayListener, TextField, IconButton, Tooltip } from '@layer5/sistent-components';
import SearchIcon from '../assets/icons/search';
import CloseIcon from '@mui/icons-material/Close';
import { makeStyles } from '@material-ui/core/styles';
import debounce from './debounce';

const useStyles = makeStyles((theme) => ({
  icon: {
    color: theme.palette.secondary.iconMain,
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

const SearchBar = ({ onSearch, placeholder }) => {
  const [expanded, setExpanded] = useState(false);
  const [searchText, setSearchText] = useState('');
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
          width: expanded ? '200px' : '0',
          opacity: expanded ? 1 : 0,
          transition: 'width 0.3s ease, opacity 0.3s ease',
        }}
      />

      {expanded ? (
        <ClickAwayListener
          onClickAway={(event) => {
            //when user clicks on actions menu, search bar should not close
            const isSearchBar = event.target.closest('#your-search-bar-id');
            const isTable = event.target.closest('#searchClick');

            if (!isSearchBar && !isTable) {
              // The click is outside the search bar and table, so you can close the search bar
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
            <SearchIcon />
          </IconButton>
        </Tooltip>
      )}
    </div>
  );
};

export default SearchBar;
