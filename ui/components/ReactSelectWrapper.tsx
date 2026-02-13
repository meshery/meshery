import React from 'react';
import PropTypes from 'prop-types';
import CreateSelect from 'react-select/creatable';
import {
  Typography,
  TextField,
  Paper,
  Chip,
  MenuItem,
  useTheme,
  styled,
  NoSsr,
} from '@sistent/sistent';
import CancelIcon from '@mui/icons-material/Cancel';

// Type extension for react-select's selectProps
// selectProps is a valid runtime prop in react-select but not in TypeScript definitions
interface SelectPropsWithTextField {
  selectProps?: {
    textFieldProps?: {
      label?: string;
      InputLabelProps?: { shrink?: boolean };
      error?: boolean;
    };
  };
}

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

function NoOptionsMessage(props: any) {
  return <StyledNoOptionsMessage {...props.innerProps}>{props.children}</StyledNoOptionsMessage>;
}

function inputComponent({ inputRef, ...props }: any) {
  return <div ref={inputRef} {...props} />;
}

function Control(props: any) {
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

function Option(props: any) {
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

function Placeholder(props: any) {
  return <StyledPlaceholder {...props.innerProps}>{props.children}</StyledPlaceholder>;
}

function SingleValue(props: any) {
  return <Typography {...props.innerProps}>{props.children}</Typography>;
}

function ValueContainer(props: any) {
  return <StyledValueContainer>{props.children}</StyledValueContainer>;
}

function MultiValue(props: any) {
  return (
    <StyledChip
      tabIndex={-1}
      label={props.children}
      onDelete={props.removeProps.onClick}
      deleteIcon={<CancelIcon {...props.removeProps} />}
    />
  );
}

function Menu(props: any) {
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
}: any) => {
  const theme = useTheme();
  const selectStyles = {
    input: (base: any) => ({
      ...base,
      color: theme.palette.text.primary,
      '& input': { font: 'inherit' },
    }),
    indicatorSeparator: () => ({
      display: 'none',
    }),
  };

  const selectProps = {
    styles: selectStyles,
    selectProps: {
      textFieldProps: {
        label,
        InputLabelProps: { shrink: true },
        error,
      },
    },
    options,
    components,
    value,
    onChange,
    onInputChange,
    placeholder,
    isClearable: true,
    isMulti,
    noOptionsMessage: () => noOptionsMessage,
  };

  return (
    <NoSsr>
      {/* selectProps is a valid runtime prop in react-select used to pass custom props to component functions */}
      <CreateSelect {...(selectProps as SelectPropsWithTextField & typeof selectProps)} />
    </NoSsr>
  );
};

ReactSelectWrapper.propTypes = {
  label: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  onInputChange: PropTypes.func,
  value: PropTypes.oneOfType([PropTypes.string, PropTypes.array]),
  options: PropTypes.array.isRequired,
  error: PropTypes.bool,
  isMulti: PropTypes.bool,
  noOptionsMessage: PropTypes.string,
};

export default ReactSelectWrapper;
