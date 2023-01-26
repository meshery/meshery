import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import CreateTable from 'react-select';
import { withStyles } from '@mui/styles';
import Typography from '@mui/material/Typography';
import NoSsr from '@mui/material/NoSsr';
import TextField from '@mui/material/TextField';
import Paper from '@mui/material/Paper';
import Chip from '@mui/material/Chip';
import MenuItem from '@mui/material/MenuItem';
import CancelIcon from '@mui/icons-material/Cancel';

const styles = () => ({ input : { display : 'flex', },
  valueContainer : {
    display : 'flex',
    flex : 1,
    alignItems : 'center',
    overflow : 'hidden',
  },
  placeholder : { position : 'absolute',
    left : 16,
    fontSize : 16, }, });

function NoOptionsMessage(props) {
  return (
    <Typography
      color="textSecondary"
      className={props.selectProps.classes.noOptionsMessage}
      {...props.innerProps}
    >
      {props.children}
    </Typography>
  );
}

function inputComponent({ inputRef, ...props }) {
  return <div ref={inputRef} {...props} />;
}

function Control(props) {
  return (
    <TextField
      fullWidth
      variant="outlined"
      InputProps={{ inputComponent,
        inputProps : {
          className : props.selectProps.classes.input,
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
      className={props.selectProps.classes.placeholder}
      {...props.innerProps}
    >
      {props.children}
    </Typography>
  );
}

function SingleValue(props) {
  return (
    <Typography className={props.selectProps.classes.singleValue} {...props.innerProps}>
      {props.children}
    </Typography>
  );
}

function ValueContainer(props) {
  return <div className={props.selectProps.classes.valueContainer}>{props.children}</div>;
}

function MultiValue(props) {
  return (
    <Chip
      tabIndex={-1}
      label={props.children}
      className={classNames(props.selectProps.classes.chip, { [props.selectProps.classes.chipFocused] : props.isFocused, })}
      onDelete={props.removeProps.onClick}
      deleteIcon={<CancelIcon {...props.removeProps} />}
    />
  );
}

function Menu(props) {
  return (
    <Paper square className={props.selectProps.classes.paper} {...props.innerProps}>
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

class ReactSelectWrapper extends React.Component {
  render() {
    const {
      classes, theme, label, placeholder, onChange, onInputChange, value, options, error, noOptionsMessage = "No Options"
    } = this.props;

    const selectStyles = { input : (base) => ({ ...base,
      color : theme.palette.text.primary,
      '& input' : { font : 'inherit', }, }), };

    return (
      <div className={classes.root}>
        <NoSsr>
          <CreateTable
            classes={classes}
            styles={selectStyles}
            textFieldProps={{ label,
              InputLabelProps : { shrink : true, },
              error, }}
            options={options}
            components={components}
            value={value}
            onChange={onChange}
            onInputChange={onInputChange}
            placeholder={placeholder}
            isClearable
            noOptionsMessage={() => noOptionsMessage}
          />
        </NoSsr>
      </div>
    );
  }
}

ReactSelectWrapper.propTypes = {
  classes : PropTypes.object.isRequired,
  theme : PropTypes.object.isRequired,
  label : PropTypes.string.isRequired,
  onChange : PropTypes.func.isRequired,
  onInputChange : PropTypes.func.isRequired,
  value : PropTypes.object.isRequired,
  options : PropTypes.array.isRequired,
  error : PropTypes.bool.isRequired,
  noOptionsMessage : PropTypes.string
};

export default withStyles(styles, { withTheme : true })(ReactSelectWrapper);
