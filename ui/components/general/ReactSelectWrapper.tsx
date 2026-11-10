import React, { useId, useMemo } from 'react';
import PropTypes from 'prop-types';
import CreateSelect from 'react-select/creatable';
import {
  CancelIcon,
  Typography,
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

const StyledLabel = styled('label')(({ theme }) => ({
  display: 'block',
  marginBottom: theme.spacing(1),
  color: theme.palette.text.primary,
  fontSize: '1rem',
}));

function NoOptionsMessage(props) {
  return <StyledNoOptionsMessage {...props.innerProps}>{props.children}</StyledNoOptionsMessage>;
}

function Option(props) {
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

function Placeholder(props) {
  return <StyledPlaceholder {...props.innerProps}>{props.children}</StyledPlaceholder>;
}

function SingleValue(props) {
  return (
    <Typography {...props.innerProps} sx={{ pl: 2 }}>
      {props.children}
    </Typography>
  );
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

  const inputId = useId();

  const selectStyles = useMemo(
    () => ({
      control: (base, state) => ({
        ...base,
        backgroundColor: theme.palette.background.paper,
        borderColor: error
          ? theme.palette.error.main
          : state.isFocused
            ? theme.palette.primary.main
            : theme.palette.divider,
        minHeight: 56,
        borderRadius: 4,
        boxShadow: 'none',
        '&:hover': {
          borderColor: theme.palette.primary.main,
        },
      }),

      input: (base) => ({
        ...base,
        color: theme.palette.text.primary,
        paddingLeft: '16px',
        '& input': { font: 'inherit' },
      }),

      singleValue: (base) => ({
        ...base,
        color: theme.palette.text.primary,
      }),

      placeholder: (base) => ({
        ...base,
        color: theme.palette.text.disabled,
      }),

      menu: (base) => ({
        ...base,
        zIndex: 9999,
        backgroundColor: theme.palette.background.paper,
      }),

      option: (base, state) => ({
        ...base,
        backgroundColor: state.isFocused
          ? theme.palette.action.hover
          : theme.palette.background.paper,
        color: theme.palette.text.primary,
        cursor: 'pointer',
      }),

      indicatorSeparator: () => ({
        display: 'none',
      }),
    }),
    [theme, error],
  );

  return (
    <NoSsr>
      <div>
        {label && <StyledLabel htmlFor={inputId}>{label}</StyledLabel>}

        <CreateSelect
          inputId={inputId}
          aria-label={label}
          aria-invalid={!!error}
          styles={selectStyles}
          options={options}
          components={components}
          value={value}
          onChange={onChange}
          onInputChange={onInputChange}
          placeholder={placeholder}
          isClearable
          isMulti={isMulti}
          backspaceRemovesValue={!isMulti}
          noOptionsMessage={() => noOptionsMessage}
        />
      </div>
    </NoSsr>
  );
};

ReactSelectWrapper.propTypes = {
  label: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  onInputChange: PropTypes.func,
  value: PropTypes.oneOfType([PropTypes.string, PropTypes.array, PropTypes.object]),
  options: PropTypes.array.isRequired,
  error: PropTypes.bool,
  isMulti: PropTypes.bool,
  noOptionsMessage: PropTypes.string,
};

export default ReactSelectWrapper;
