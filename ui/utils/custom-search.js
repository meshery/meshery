import React, { useState, useRef } from 'react';
import {
  CustomTooltip,
  styled,
  TextField,
  IconButton,
  ClickAwayListener,
  CloseIcon,
  SearchIcon,
  useTheme,
} from '@layer5/sistent';
import debounce from './debounce';

const SearchContainer = styled('div')({
  display: 'flex',
  alignItems: 'center',
});

const SearchTextField = styled(TextField)(({ theme }) => ({
  '& .MuiOutlinedInput-root': {
    color: theme.palette.icon.secondary,
  },
  '& .MuiOutlinedInput-notchedOutline': {
    borderColor: theme.palette.icon.secondary,
    '&:hover': {
      borderColor: '#00b39f',
    },
  },
  '& .MuiInputLabel-root': {
    color: theme.palette.icon.secondary,
  },
  '& .MuiInputBase-input': {
    color: theme.palette.icon.secondary,
    caretColor: theme.palette.icon.secondary,
  },
  '& .MuiInput-underline:before': {
    borderBottomColor: theme.palette.icon.secondary,
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
}));

const SearchBar = ({ onSearch, placeholder, expanded, setExpanded, value = '' }) => {
  const [searchText, setSearchText] = useState(value);
  const searchRef = useRef(null);

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
  let searchWidth = '12.5rem';
  if (width <= 750) {
    searchWidth = '7.5rem';
  }
  const theme = useTheme();
  return (
    <SearchContainer>
      <SearchTextField
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
          <CustomTooltip title="Close">
            <IconButton onClick={handleClearIconClick}>
              <CloseIcon fill={theme.palette.icon.secondary} height={'1.5rem'} width={'1.5rem'} />
            </IconButton>
          </CustomTooltip>
        </ClickAwayListener>
      ) : (
        <CustomTooltip title="Search">
          <IconButton onClick={handleSearchIconClick}>
            <SearchIcon fill={theme.palette.icon.secondary} height={'1.5rem'} width={'1.5rem'} />
          </IconButton>
        </CustomTooltip>
      )}
    </SearchContainer>
  );
};

export default SearchBar;
