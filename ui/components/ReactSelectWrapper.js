import React from 'react';
import PropTypes from 'prop-types';
import CreateSelect from 'react-select/creatable';
import { styled } from '@mui/material/styles';
import { Typography, TextField, Paper, Chip, MenuItem, useTheme } from '@layer5/sistent';
import NoSsr from '@mui/material/NoSsr';
import CancelIcon from '@mui/icons-material/Cancel';

const StyledNoOptionsMessage = styled(Typography)(({ theme }) => ({
  padding: '0.2rem',
  marginLeft: '0.8rem',
  color: theme.palette.text.disabled,
}));

const StyledValueContainer = styled('div')({
  display: 'flex',
  flex: 1,
  alignItems: 'center',
  overflow: 'hidden',
});

const StyledPlaceholder = styled(Typography)(({ theme }) => ({
  color: theme.palette.text.disabled,
  position: 'absolute',
  left: 16,
  fontSize: 16,
}));

const StyledPaper = styled(Paper)({
  zIndex: 9999,
  width: '100%',
  position: 'absolute',
});

const StyledChip = styled(Chip)(({ theme }) => ({
  margin: '3px 3px',
  height: 'auto',
  '&:focus': {
    backgroundColor: theme.palette.background.default,
    color: theme.palette.background.inverse,
  },
}));

function NoOptionsMessage(props) {
  return <StyledNoOptionsMessage {...props.innerProps}>{props.children}</StyledNoOptionsMessage>;
}

function inputComponent({ inputRef, ...props }) {
  return <div ref={inputRef} {...props} />;
}

function Control(props) {
  return (
    <TextField
      fullWidth
      variant="outlined"
      InputProps={{
        inputComponent,
        inputProps: {
          style: {
            display: 'flex',
          },
          inputRef: props.innerRef,
          children: props.children,
          ...props.innerProps,
        },
      }}
      {...props.selectProps.textFieldProps}
    />
  );
}

function Option(props) {
  return (
    <MenuItem
      ref={props.innerRef}
      selected={props.isFocused}
      component="div"
      style={{ fontWeight: props.isSelected ? 500 : 400, padding: '0.4rem 1rem' }}
      {...props.innerProps}
    >
      {props.children}
    </MenuItem>
  );
}

function Placeholder(props) {
  return <StyledPlaceholder {...props.innerProps}>{props.children}</StyledPlaceholder>;
}

function SingleValue(props) {
  return <Typography {...props.innerProps}>{props.children}</Typography>;
}

function ValueContainer(props) {
  return <StyledValueContainer>{props.children}</StyledValueContainer>;
}

function MultiValue(props) {
  return (
    <StyledChip
      tabIndex={-1}
      label={props.children}
      onDelete={props.removeProps.onClick}
      deleteIcon={<CancelIcon {...props.removeProps} />}
    />
  );
}

function Menu(props) {
  return (
    <StyledPaper square {...props.innerProps}>
      {props.children}
    </StyledPaper>
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

const ReactSelectWrapper = ({
  label,
  placeholder,
  onChange,
  onInputChange,
  value,
  options,
  error,
  isMulti = false,
  noOptionsMessage = 'Type to create a new Environment',
}) => {
  const theme = useTheme();
  const selectStyles = {
    input: (base) => ({
      ...base,
      color: theme.palette.text.primary,
      '& input': { font: 'inherit' },
    }),
    indicatorSeparator: () => ({
      display: 'none',
    }),
  };

  return (
    <NoSsr>
      <CreateSelect
        styles={selectStyles}
        textFieldProps={{
          label,
          InputLabelProps: { shrink: true },
          error,
        }}
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
