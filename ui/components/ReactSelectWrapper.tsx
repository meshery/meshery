import React from 'react';
import PropTypes from 'prop-types';
import CreateSelect from 'react-select/creatable';
import {
  CancelIcon,
  Typography,
  TextField,
  Paper,
  Chip,
  ListItemButton,
  useTheme,
  styled,
  NoSsr,
} from '@sistent/sistent';

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

function NoOptionsMessage(props: { innerProps: unknown; children: React.ReactNode }) {
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
    <ListItemButton
      ref={props.innerRef}
      selected={props.isFocused}
      component="div"
      style={{ fontWeight: props.isSelected ? 500 : 400, padding: '0.4rem 1rem' }}
      {...props.innerProps}
    >
      {props.children}
    </ListItemButton>
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

type ReactSelectWrapperProps = {
  label: string;
  placeholder?: string;
  onChange: (_value: any) => void;
  onInputChange?: (_value: string) => void;
  value?: any;
  options: any[];
  error?: boolean;
  helperText?: string;
  isMulti?: boolean;
  noOptionsMessage?: string;
};

const CreatableSelect = CreateSelect as unknown as React.ComponentType<any>;

const ReactSelectWrapper = ({
  label,
  placeholder = '',
  onChange,
  onInputChange,
  value,
  options,
  error,
  helperText,
  isMulti = false,
  noOptionsMessage = 'Type to create a new option',
}: ReactSelectWrapperProps) => {
  const theme = useTheme();
  const selectStyles = {
    input: (base: Record<string, unknown>) => ({
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
      <CreatableSelect
        styles={selectStyles}
        textFieldProps={{
          label,
          InputLabelProps: { shrink: true },
          error,
          helperText,
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
  onInputChange: PropTypes.func,
  value: PropTypes.oneOfType([PropTypes.object, PropTypes.array]),
  options: PropTypes.arrayOf(
    PropTypes.shape({
      value: PropTypes.string,
      label: PropTypes.string,
    }),
  ).isRequired,
  error: PropTypes.bool,
  helperText: PropTypes.string,
  isMulti: PropTypes.bool,
  noOptionsMessage: PropTypes.string,
};

export default ReactSelectWrapper;
