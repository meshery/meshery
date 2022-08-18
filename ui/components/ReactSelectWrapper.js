import React from 'react';
// import Select from 'react-select';
import CreatableSelect from 'react-select/creatable';
import {Typography, TextField,Paper, Chip, MenuItem } from '@mui/material';
import { styled } from "@mui/material/styles";

const CustomCreatableSelect = styled(CreatableSelect)(({ theme }) => ({
  verticalAlign: "center",
  margin: "auto 0",
  minWidth : '250px',
  width: "100%",
  flex : '1',
  }));


function NoOptionsMessage(props) {
  return (
    <Typography
      color="textSecondary"
    >
      {props.children}
    </Typography>
  );
}

function inputComponent({ inputRef, ...props }) {
  return <div style={{display: "flex"}} ref={inputRef} {...props} />;
}

function Control(props) {
  return (
    <TextField
      fullWidth
      variant="outlined"
      
      InputProps={{ inputComponent,
       
        inputProps : {
          
          inputRef : props.innerRef,
          children : props.children,
          ...props.innerProps,
        }, }}
      {...props.selectProps.textFieldProps}
    />
  );
}

function Option(props) {
  return (
    <MenuItem
      buttonRef={props.innerRef}
      selected={props.isFocused}
      component="div"
      style={{ fontWeight : props.isSelected
        ? 500
        : 400, }}
      {...props.innerProps}
    >
      {props.children}
    </MenuItem>
  );
}

function Placeholder(props) {
  return (
    <Typography
      color="textSecondary"
      {...props.innerProps}
    >
      {props.children}
    </Typography>
  );
}

function SingleValue(props) {
  return (
    <Typography {...props.innerProps}>
      {props.children}
    </Typography>
  );
}

function ValueContainer(props) {
  return <div style={{ flex : 1,}}>{props.children}</div>;
}

function MultiValue(props) {
  return (
    <Chip
      tabIndex={-1}
      label={props.children}
      onDelete={props.removeProps.onClick}
    />
  );
}

function Menu(props) {
  return (
    <Paper square {...props.innerProps}>
      {props.children}
    </Paper>
  );
}

const components = {
  Control,
  Menu,
  MultiValue,
  NoOptionsMessage,
  Option,
  Placeholder,
  SingleValue,
  ValueContainer,
};

export default function ReactSelectWrapper ({ label, placeholder}) {

    return (
      <div >
          <CustomCreatableSelect
            textFieldProps={{ label,
              InputLabelProps : { shrink : true, } }}
            components={components}
            placeholder={placeholder}
            isClearable
          />
      </div>
    );
  }

