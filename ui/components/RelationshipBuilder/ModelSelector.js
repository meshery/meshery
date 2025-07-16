import React, { useState, useEffect, useRef, useCallback } from 'react';
import {
  TextField,
  Box,
  Typography,
  CircularProgress,
  Paper,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Popper,
  ClickAwayListener,
  InputAdornment,
  IconButton,
} from '@sistent/sistent';
import { styled } from '@sistent/sistent';
import { useModelSelection } from '@/utils/hooks/useModelSelection';
import { debounce } from 'lodash';
import SearchIcon from '@mui/icons-material/Search';
import ClearIcon from '@mui/icons-material/Clear';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';

const StyledPaper = styled(Paper)(({ theme }) => ({
  maxHeight: '300px',
  overflow: 'auto',
  width: '100%',
  zIndex: 1500,
  boxShadow: theme.shadows[8],
  '&::-webkit-scrollbar': {
    width: '8px',
  },
  '&::-webkit-scrollbar-track': {
    background: theme.palette.background.default,
  },
  '&::-webkit-scrollbar-thumb': {
    background: theme.palette.divider,
    borderRadius: '4px',
  },
}));

const StyledListItem = styled(ListItemButton)(({ theme }) => ({
  padding: theme.spacing(1, 2),
  '&:hover': {
    backgroundColor: theme.palette.action.hover,
  },
  '&.selected': {
    backgroundColor: theme.palette.action.selected,
    '&:hover': {
      backgroundColor: theme.palette.action.selected,
    },
  },
}));

const LoadingItem = styled(ListItem)(({ theme }) => ({
  display: 'flex',
  justifyContent: 'center',
  padding: theme.spacing(2),
  color: theme.palette.text.secondary,
}));

const ModelSelector = ({
  selectedModel,
  onModelChange,
  disabled = false,
  label = 'Select Model',
  helperText = 'Search and select a model',
  error = false,
  ...props
}) => {
  const { models, isLoading, isFetching, hasMore, handleSearch, loadNextPage } =
    useModelSelection();
  const [inputValue, setInputValue] = useState('');
  const [open, setOpen] = useState(false);
  const [highlightedIndex, setHighlightedIndex] = useState(-1);
  const textFieldRef = useRef(null);
  const listRef = useRef(null);
  const loadingRef = useRef(null);

  useEffect(() => {
    if (selectedModel) {
      setInputValue(selectedModel.displayName || selectedModel.name || '');
    } else {
      setInputValue('');
    }
  }, [selectedModel]);

  const debouncedSearch = useCallback(
    debounce((query) => {
      handleSearch(query);
      setHighlightedIndex(-1);
    }, 300),
    [handleSearch],
  );

  const handleInputChange = useCallback(
    (event) => {
      const value = event.target.value;
      setInputValue(value);
      debouncedSearch(value);
    },
    [debouncedSearch],
  );

  const handleInputFocus = useCallback(() => {
    setOpen(true);
    if (inputValue && !models.length) {
      handleSearch(inputValue);
    }
  }, [inputValue, models.length, handleSearch]);

  const handleModelSelect = useCallback(
    (model) => {
      onModelChange(model);
      setInputValue(model.displayName || model.name || '');
      setOpen(false);
      setHighlightedIndex(-1);
    },
    [onModelChange],
  );
  const handleScroll = useCallback(
    (event) => {
      const { target } = event;
      const { scrollTop, scrollHeight, clientHeight } = target;
      if (scrollHeight - scrollTop <= clientHeight + 50) {
        if (hasMore && !isLoading && !isFetching && inputValue.length === 0) {
          loadNextPage();
        }
      }
    },
    [hasMore, isLoading, isFetching, inputValue, loadNextPage],
  );
  const handleClear = useCallback(() => {
    onModelChange(null);
    setInputValue('');
    setOpen(false);
    setHighlightedIndex(-1);
  }, [onModelChange]);

  useEffect(() => {
    if (!open || !loadingRef.current || !hasMore) return;
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !isLoading && !isFetching) {
          loadNextPage();
        }
      },
      { threshold: 0.1 },
    );

    observer.observe(loadingRef.current);

    return () => {
      if (loadingRef.current) {
        observer.unobserve(loadingRef.current);
      }
    };
  }, [open, hasMore, isLoading, isFetching, loadNextPage]);

  return (
    <ClickAwayListener onClickAway={() => setOpen(false)}>
      <Box position="relative" {...props}>
        <TextField
          ref={textFieldRef}
          fullWidth
          label={label}
          helperText={helperText}
          error={error}
          variant="outlined"
          value={inputValue}
          onChange={handleInputChange}
          onFocus={handleInputFocus}
          disabled={disabled}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon color="action" />
              </InputAdornment>
            ),
            endAdornment: (
              <InputAdornment position="end">
                {inputValue && !disabled && (
                  <IconButton size="small" onClick={handleClear} edge="end">
                    <ClearIcon />
                  </IconButton>
                )}
                <IconButton
                  size="small"
                  onClick={() => setOpen(!open)}
                  edge="end"
                  disabled={disabled}
                >
                  <KeyboardArrowDownIcon
                    style={{
                      transform: open ? 'rotate(180deg)' : 'rotate(0deg)',
                      transition: 'transform 0.2s',
                    }}
                  />
                </IconButton>
              </InputAdornment>
            ),
          }}
          placeholder={selectedModel ? '' : 'Search for models...'}
        />

        <Popper
          open={open}
          anchorEl={textFieldRef.current}
          placement="bottom-start"
          style={{ width: textFieldRef.current?.offsetWidth, zIndex: 1500 }}
        >
          <StyledPaper onScroll={handleScroll} ref={listRef}>
            <List dense>
              {models.length === 0 && !isLoading ? (
                <ListItem>
                  <ListItemText
                    primary="No models found"
                    primaryTypographyProps={{
                      variant: 'body2',
                      color: 'textSecondary',
                      align: 'center',
                    }}
                  />
                </ListItem>
              ) : null}

              {models.map((model, index) => (
                <StyledListItem
                  key={model.id || `${model.name}-${index}`}
                  onClick={() => handleModelSelect(model)}
                  className={highlightedIndex === index ? 'selected' : ''}
                >
                  <ListItemText
                    primary={model.displayName || model.name}
                    secondary={
                      <Box>
                        <Typography variant="caption" color="textSecondary">
                          {model.registrant?.hostname} • {model.version}
                        </Typography>
                        {model.subCategory && (
                          <Typography variant="caption" color="textSecondary" display="block">
                            {model.subCategory}
                          </Typography>
                        )}
                      </Box>
                    }
                  />
                </StyledListItem>
              ))}

              {hasMore && (
                <LoadingItem>
                  <CircularProgress size={20} />
                  <Typography variant="caption" sx={{ ml: 1 }}>
                    {models.length === 0 ? 'Loading models...' : 'Loading more models...'}
                  </Typography>
                </LoadingItem>
              )}
              {hasMore && <Box ref={loadingRef} sx={{ height: 1 }} />}
            </List>
          </StyledPaper>
        </Popper>
      </Box>
    </ClickAwayListener>
  );
};

export default ModelSelector;
