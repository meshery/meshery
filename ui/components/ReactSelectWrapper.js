import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import CreateSelect from 'react-select/creatable';
import { useTheme } from '@mui/styles';
import { UsesSistent } from './SistentWrapper';
import Typography from '@mui/material/Typography';
import NoSsr from '@mui/material/NoSsr';
import TextField from '@mui/material/TextField';
import Paper from '@mui/material/Paper';
import Chip from '@mui/material/Chip';
import MenuItem from '@mui/material/MenuItem';
import CancelIcon from '@mui/icons-material/Cancel';

function NoOptionsMessage(props) {
  return (
    <Typography
      sx={{
        padding: '0.2rem',
        marginLeft: '0.8rem',
      }}
      {...props.innerProps}
    >
      {props.children}
    </Typography>
  );
}

const InputComponent = React.forwardRef(function InputComponent(props, ref) {
  return <div ref={ref} {...props} />;
});

const Control = React.forwardRef(function Control(props, ref) {
  return (
    <TextField
      fullWidth
      variant="outlined"
      InputProps={{
        inputComponent: InputComponent,
        inputProps: {
          children: props.children,
          ...props.innerProps,
          sx: {
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          },
        },
        inputRef: ref,
      }}
      {...props.selectProps.textFieldProps}
    />
  );
});

const Option = React.forwardRef(function Option(props, ref) {
  const theme = useTheme();
  return (
    <MenuItem
      ref={ref}
      selected={props.isFocused}
      component="div"
      sx={{
        fontWeight: props.isSelected ? 500 : 400,
        padding: '0.4rem 1rem',
        backgroundColor: theme.palette.background.secondary,
      }}
      {...props.innerProps}
    >
      {props.children}
    </MenuItem>
  );
});

function Placeholder(props) {
  return (
    <Typography
      sx={{
        position: 'absolute',
        left: 16,
        fontSize: 16,
      }}
      {...props.innerProps}
    >
      {props.children}
    </Typography>
  );
}

function SingleValue(props) {
  return <Typography {...props.innerProps}>{props.children}</Typography>;
}

function ValueContainer(props) {
  return (
    <div
      style={{
        display: 'flex',
        flex: 1,
        alignItems: 'center',
        overflow: 'hidden',
      }}
    >
      {props.children}
    </div>
  );
}

function MultiValue(props) {
  return (
    <Chip
      tabIndex={-1}
      label={props.children}
      className={classNames(props.selectProps.classes.chip, {
        [props.selectProps.classes.chipFocused]: props.isFocused,
      })}
      onDelete={props.removeProps.onClick}
      deleteIcon={<CancelIcon {...props.removeProps} />}
    />
  );
}

function Menu(props) {
  const theme = useTheme();
  return (
    <Paper
      square
      {...props.innerProps}
      sx={{
        zIndex: 9999,
        width: '100%',
        position: 'absolute',
        backgroundColor: theme.palette.background.secondary,
      }}
    >
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

// NOTE: This is a wrapper for react-select
// It is used to customize the look and feel of the react-select component
// NOTE: Migrate to functional component and move to sistent
const ReactSelectWrapper = (props) => {
  const {
    label,
    placeholder,
    onChange,
    onInputChange,
    value,
    options,
    error,
    isMulti = false,
    noOptionsMessage = 'Type to create a new Environment',
  } = props;
  const theme = useTheme();
  const selectStyles = {
    input: (base) => ({
      ...base,
      color: theme?.palette.text,
      '& input': { font: 'inherit' },
    }),
  };

  return (
    <UsesSistent>
      <div>
        <NoSsr>
          <CreateSelect
            styles={selectStyles}
            textFieldProps={{ label, InputLabelProps: { shrink: true }, error }}
            options={options}
            components={components}
            value={value}
            onChange={onChange}
            onInputChange={onInputChange}
            placeholder={placeholder}
            isClearable
            isMulti={isMulti}
            noOptionsMessage={() => noOptionsMessage}
          />
        </NoSsr>
      </div>
    </UsesSistent>
  );
};

ReactSelectWrapper.propTypes = {
  label: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  onInputChange: PropTypes.func.isRequired,
  value: PropTypes.object.isRequired,
  options: PropTypes.array.isRequired,
  error: PropTypes.bool.isRequired,
  isMulti: PropTypes.bool,
  noOptionsMessage: PropTypes.string,
};

export default ReactSelectWrapper;
