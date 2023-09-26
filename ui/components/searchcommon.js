import React from 'react';
import { TextField, Box, InputAdornment } from '@material-ui/core';
import SearchIcon from '@mui/icons-material/Search';
import { withStyles } from '@material-ui/core/styles';

const styles = (theme) => ({
  searchInput: {
    '& .MuiOutlinedInput-root': {
      color: theme.palette.type === 'dark' ? theme.palette.common.white : theme.palette.common.grey,
      backgroundColor:
        theme.palette.type === 'dark' ? theme.palette.grey[800] : 'rgba(102, 102, 102, 0.12)',
      height: '4.5ch',
    },
    '& .MuiOutlinedInput-notchedOutline': {
      borderColor:
        theme.palette.type === 'dark' ? 'rgba(255, 255, 255, 0.7)' : theme.palette.grey[800],
      height: '5ch',
    },
    '& .MuiInputLabel-root': {
      color: theme.palette.type === 'dark' ? theme.palette.common.white : theme.palette.common.grey,
      lineHeight: '0.975rem',
    },
    '& .MuiInputBase-input': {
      caretColor:
        theme.palette.type === 'dark' ? theme.palette.common.white : theme.palette.common.grey,
    },
    '& .MuiInputAdornment-root .MuiSvgIcon-root': {
      color: theme.palette.type === 'dark' ? theme.palette.common.white : theme.palette.common.grey,
    },
    '& .MuiInputBase-root:hover .MuiOutlinedInput-notchedOutline': {
      borderColor: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.grey[800],
    },
    '& .MuiInputBase-root.Mui-focused .MuiOutlinedInput-notchedOutline': {
      borderColor: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.grey[800],
    },
  },
});

const SearchBar = ({ onChange, value, classes, width, label }) => {
  return (
    <>
      <Box
        component="form"
        sx={{
          '& > :not(style)': { width: width },
        }}
        noValidate
        autoComplete="on"
      >
        <TextField
          id="outlined-basic"
          label={label}
          variant="outlined"
          fullWidth
          type="search"
          value={value}
          onChange={onChange}
          sx={{
            margin: 'auto',
            height: '5ch',
          }}
          placeholder="Search"
          className={classes.searchInput}
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <SearchIcon />
              </InputAdornment>
            ),
          }}
        />
      </Box>
    </>
  );
};

export default withStyles(styles)(SearchBar);
