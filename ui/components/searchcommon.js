import React from 'react';
import TextField from '@mui/material/TextField';
import Box from '@mui/material/Box';
import { InputAdornment } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import { withStyles } from "@material-ui/core/styles";

const styles = ((theme) => ({
  searchInput : {
    '& .MuiOutlinedInput-root' : {
      color : theme.palette.type === 'dark'? theme.palette.common.white: theme.palette.common.grey,
      backgroundColor : theme.palette.type === 'dark' ? theme.palette.grey[800] : 'rgba(102, 102, 102, 0.12)',
    },
    '& .MuiOutlinedInput-notchedOutline' : {
      borderColor : theme.palette.type === 'dark'? 'rgba(255, 255, 255, 0.7)': theme.palette.grey[800],
    },
    '& .MuiInputLabel-root' : {
      color : theme.palette.type === 'dark'? theme.palette.common.white: theme.palette.common.grey,
    },
    '& .MuiInputBase-input' : {
      caretColor : theme.palette.type === 'dark'? theme.palette.common.white: theme.palette.common.grey,
    },
    '& .MuiInputAdornment-root .MuiSvgIcon-root' : {
      color : theme.palette.type === 'dark'? theme.palette.common.white: theme.palette.common.grey,
    },
  }
}));

const SearchBar = ({ onChange, value, classes, width }) => {
  // const [searchValue, setSearchValue] = useState(value);
  // const router = useRouter();

  // const handleChange = (event) => {
  //   const newValue = event.target.value;
  //   setSearchValue(newValue);
  //   onChange(newValue);
  // };

  // const handleSubmit = (event) => {
  //   event.preventDefault();
  //   search(searchValue);
  //   onRequestSearch();
  //   router.push('/search');
  // };

  // const handleCancel = () => {
  //   setSearchValue('');
  //   onCancelSearch();
  // };

  return (
    <>
      <Box
        component="form"
        sx={{
          '& > :not(style)' : {  width : width },
        }}
        noValidate
        autoComplete="on"
      >
        <TextField
          id="outlined-basic"
          label="Search Applications"
          variant="outlined"
          fullWidth
          type="search"
          value={value}
          onChange={onChange}
          sx={{
            margin : 'auto',
          }}
          placeholder="Search"
          className={classes.searchInput}
          InputProps={{
            endAdornment : (
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