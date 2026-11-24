import React, { useState, useRef, useMemo } from 'react';
import {
  CustomTooltip,
  styled,
  TextField,
  IconButton,
  ClickAwayListener,
  CloseIcon,
  SearchIcon,
  useTheme,
} from '@sistent/sistent';
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

const SearchBar = ({
  onSearch,
  placeholder,
  expanded,
  setExpanded,
  value = '',
  setModelsFilters,
}) => {
  const [searchText, setSearchText] = useState(value);
  const searchRef = useRef(null);

  const debouncedOnSearch = useMemo(() => {
    return typeof onSearch === 'function' ? debounce(onSearch, 500) : () => {};
  }, [onSearch]);

  const handleSearchChange = (event) => {
    const val = event.target.value;
    debouncedOnSearch(val);
    setSearchText(val);
  };

  const handleClearIconClick = () => {
    if (typeof setModelsFilters === 'function') {
      setModelsFilters({ page: 0 });
    }
    setSearchText('');
    debouncedOnSearch('');
    if (typeof setExpanded === 'function') {
      setExpanded(false);
    }
  };

  const handleSearchIconClick = () => {
    if (expanded) {
      setSearchText('');
      if (typeof setExpanded === 'function') {
        setExpanded(false);
      }
    } else {
      if (typeof setExpanded === 'function') {
        setExpanded(true);
      }
      setTimeout(() => {
        searchRef.current?.focus?.();
      }, 300);
    }
  };

  const width = typeof window !== 'undefined' ? window.innerWidth : 1200;
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
            <IconButton onClick={handleClearIconClick} data-testid="close-icon" aria-label="close">
              <CloseIcon fill={theme.palette.icon.secondary} height={'1.5rem'} width={'1.5rem'} />
            </IconButton>
          </CustomTooltip>
        </ClickAwayListener>
      ) : (
        <CustomTooltip title="Search">
          <IconButton onClick={handleSearchIconClick} data-testid="search-icon">
            <SearchIcon fill={theme.palette.icon.secondary} height={'1.5rem'} width={'1.5rem'} />
          </IconButton>
        </CustomTooltip>
      )}
    </SearchContainer>
  );
};

export default SearchBar;
