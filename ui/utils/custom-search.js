import React, { useState, useRef } from "react";
import TextField from "@mui/material/TextField";
import { Tooltip } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import SearchIcon from "@mui/icons-material/Search";
import CloseIcon from "@mui/icons-material/Close";
import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles((theme) => ({
  icon : {
    color : theme.palette.secondary.iconMain
  },
  searchInput : {
    "& .MuiOutlinedInput-root" : {
      color : theme.palette.secondary.iconMain
    },
    "& .MuiOutlinedInput-notchedOutline" : {
      borderColor : theme.palette.secondary.iconMain
    },
    "& .MuiInputLabel-root" : {
      color : theme.palette.secondary.iconMain
    },
    "& .MuiInputBase-input" : {
      color : theme.palette.secondary.iconMain,
      caretColor : theme.palette.secondary.iconMain
    },
    "& .MuiInput-underline:before" : {
      borderBottomColor : theme.palette.secondary.iconMain
    },
    "& .MuiInput-underline:hover:before" : {
      borderBottomColor : theme.palette.secondary.iconMain
    },
    "& .MuiInput-underline:hover:after" : {
      borderBottomColor : theme.palette.secondary.iconMain
    },
    "& .MuiInput-underline.Mui-focused:before" : {
      borderBottomColor : theme.palette.type === theme.palette.secondary.iconMain
    },
    "& .MuiInput-underline.Mui-focused:after" : {
      borderBottomColor : theme.palette.type === theme.palette.secondary.iconMain
    },
  },
}));

const SearchBar = ({ onSearch, placeholder }) => {
  const [expanded, setExpanded] = useState(false);
  const [searchText, setSearchText] = useState("");
  const searchRef = useRef(null);
  const classes = useStyles();

  const handleSearchChange = (event) => {
    onSearch(event.target.value);
    setSearchText(event.target.value);
  };

  const handleClearIconClick = () => {
    setSearchText("");
    setExpanded(false);
  };

  const handleSearchIconClick = () => {
    if (expanded) {
      setSearchText("");
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
        id="standard-basic"
        variant="standard"
        value={searchText}
        onChange={handleSearchChange}
        inputRef={searchRef}
        placeholder={placeholder}
        style={{
          width : expanded ? "200px" : "0",
          opacity : expanded ? 1 : 0,
          transition : "width 0.3s ease, opacity 0.3s ease",
        }}
      />

      {expanded ? (
        <Tooltip title="Close">
          <IconButton
            onClick={handleClearIconClick}
            sx={{
              "&:hover" : {
                borderRadius : "4px",
              },
            }}
          >
            <CloseIcon className={classes.icon} />
          </IconButton>
        </Tooltip>
      ) : (
        <Tooltip title="Search">
          <IconButton
            onClick={handleSearchIconClick}
            sx={
              {
                "&:hover" : {
                  borderRadius : "4px"
                }
              }
            }
          >
            <SearchIcon className={classes.icon} />
          </IconButton>
        </Tooltip>
      )}
    </div>
  );
};

export default SearchBar;
