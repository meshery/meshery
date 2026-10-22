import { useState } from 'react';
import { components } from 'react-select';
import CreatableSelect from 'react-select/creatable';
import {
  ArrowDropDownIcon,
  Checkbox,
  ListItemButton,
  Paper,
  FormControlLabel,
  useTheme,
  alpha,
} from '@sistent/sistent';

/** Brand-tinted option row highlight (keyboard/hover). */
export function getMultiSelectOptionHighlight(theme) {
  return alpha(theme.palette.primary.main, theme.palette.mode === 'dark' ? 0.28 : 0.12);
}

const MultiSelectWrapper = (props) => {
  const [selectInput, setSelectInput] = useState('');
  const allOption = { value: '*' };
  const theme = useTheme();

  const filterOptions = (options, input) =>
    options?.filter(({ label }) => label?.toLowerCase().includes(input.toLowerCase()));

  const comparator = (v1, v2) => {
    if (v1.value === '*') return 1;
    else if (v2.value === '*') return -1;

    return v1.label?.localeCompare(v2.label);
  };

  let filteredOptions = filterOptions(props.options, selectInput).sort(comparator);
  let filteredSelectedOptions = filterOptions(props.value, selectInput).sort(comparator);

  const menuBackground =
    theme.palette.mode === 'dark' ? theme.palette.background.card : theme.palette.common.white;

  const controlBackground = 'transparent';

  // Opaque chip fill — action.selected is translucent and breaks getContrastText.
  const chipBackground =
    theme.palette.mode === 'dark'
      ? (theme.palette.background.elevatedComponents ??
        theme.palette.background.paper ??
        theme.palette.background.default)
      : (theme.palette.background.secondary ?? theme.palette.background.paper);
  const chipForeground =
    theme.palette.getContrastText?.(chipBackground) ?? theme.palette.text.primary;
  const chipSurface = {
    backgroundColor: chipBackground,
    color: chipForeground,
    border: 'none',
  };
  const chipType = {
    fontSize: theme.typography?.body2?.fontSize,
    lineHeight: theme.typography?.body2?.lineHeight,
  };

  const optionHighlight = getMultiSelectOptionHighlight(theme);
  const optionSelectedFill = alpha(
    theme.palette.primary.main,
    theme.palette.mode === 'dark' ? 0.14 : 0.06,
  );
  const optionTextColor = theme.palette.text.default ?? theme.palette.text.primary;

  const Option = (props) => {
    return (
      <ListItemButton
        ref={props.innerRef}
        selected={props.isSelected}
        {...props.innerProps}
        component="div"
        data-testid="multi-select-option"
        sx={{
          fontWeight: props.isSelected ? 500 : 400,
          padding: '0.4rem 1rem',
          backgroundColor: props.isFocused
            ? optionHighlight
            : props.isSelected
              ? optionSelectedFill
              : 'transparent',
          color: optionTextColor,
          '&:hover': {
            backgroundColor: optionHighlight,
          },
          '&.Mui-selected': {
            backgroundColor: props.isFocused ? optionHighlight : optionSelectedFill,
            '&:hover': {
              backgroundColor: optionHighlight,
            },
          },
          '& .MuiFormControlLabel-label': {
            color: optionTextColor,
          },
        }}
      >
        <FormControlLabel
          control={
            props.value === '*' && filteredSelectedOptions?.length > 0 ? (
              <>
                <Checkbox
                  key={props.value}
                  ref={(input) => {
                    if (input) input.indeterminate = true;
                  }}
                  style={{
                    padding: '0',
                  }}
                />
              </>
            ) : (
              <>
                <Checkbox
                  key={props.value}
                  checked={props.isSelected}
                  onChange={() => {}}
                  style={{
                    padding: '0',
                  }}
                />
              </>
            )
          }
          label={<span style={{ marginLeft: '0.5rem' }}>{props.label}</span>}
        />
      </ListItemButton>
    );
  };

  const CustomInput = (props) => (
    <components.Input autoFocus={props.selectProps.menuIsOpen} {...props}>
      {props.children}
    </components.Input>
  );

  const Menu = (props) => {
    // Custom Menu must re-apply getStyles or customStyles.menu never paints.
    const menuStyles = props.getStyles('menu', props);
    return (
      <Paper
        square
        style={{
          ...menuStyles,
          zIndex: 9999,
          width: '100%',
          position: 'absolute',
        }}
        {...props.innerProps}
      >
        {props.children}
      </Paper>
    );
  };

  const DropdownIndicator = (props) => (
    <components.DropdownIndicator {...props}>
      <ArrowDropDownIcon fontSize="small" sx={{ color: 'action.active' }} />
    </components.DropdownIndicator>
  );

  const customFilterOption = ({ value, label }, input) =>
    (value !== '*' && label?.toLowerCase().includes(input.toLowerCase())) ||
    (value === '*' && filteredOptions?.length > 0);

  const onInputChange = (inputValue, event) => {
    if (event.action === 'input-change') setSelectInput(inputValue);
    else if (event.action === 'menu-close' && selectInput !== '') setSelectInput('');
  };

  const onKeyDown = (e) => {
    if ((e.key === ' ' || e.key === 'Enter') && !selectInput) e.preventDefault();
  };

  const handleChange = (selected) => {
    if (
      selected.length > 0 &&
      (selected[selected.length - 1].value === allOption.value ||
        JSON.stringify(filteredOptions.sort(comparator)) ===
          JSON.stringify(selected.sort(comparator)))
    ) {
      setSelectInput('');
      return props.onChange(
        [
          ...(props.value ?? []),
          ...props.options.filter(
            ({ label }) =>
              label?.toLowerCase().includes(selectInput?.toLowerCase()) &&
              (props.value ?? []).filter((opt) => opt.label === label).length === 0,
          ),
        ].sort(comparator),
        [],
      );
    } else if (
      selected.length > 0 &&
      selected[selected.length - 1].value !== allOption.value &&
      JSON.stringify(selected.sort(comparator)) !== JSON.stringify(filteredOptions.sort(comparator))
    ) {
      let filteredUnselectedOptions = filteredSelectedOptions.filter(
        (opts) => !selected.some((sel) => sel.value === opts.value),
      );
      setSelectInput('');
      return props.onChange(selected, filteredUnselectedOptions);
    } else {
      setSelectInput('');
      return props.onChange(
        [
          ...props.value.filter(
            ({ label }) => !label.toLowerCase().includes(selectInput?.toLowerCase()),
          ),
        ],
        filteredSelectedOptions.length == 1 ? filteredSelectedOptions : props.options,
      );
    }
  };

  const chipRemoveHoverBg = alpha(
    theme.palette.error.main,
    theme.palette.mode === 'dark' ? 0.28 : 0.12,
  );
  const restingBorder = theme.palette.border?.strong ?? theme.palette.divider;
  const activeBorder = theme.palette.primary.main;
  const indicatorColor = {
    color: theme.palette.action.active,
    ':hover': { color: theme.palette.text.primary },
  };

  const customStyles = {
    multiValue: (base) => ({
      ...base,
      ...chipSurface,
      display: 'flex',
      alignItems: 'stretch',
      borderRadius: '4px',
      overflow: 'hidden',
      margin: '2px',
      maxWidth: '100%',
    }),
    multiValueLabel: (def) => ({
      ...def,
      ...chipSurface,
      ...chipType,
      display: 'flex',
      alignItems: 'center',
      borderRadius: 0,
      paddingTop: '4px',
      paddingBottom: '4px',
      paddingLeft: '8px',
      paddingRight: '4px',
      maxWidth: '9rem',
      overflow: 'hidden',
      textOverflow: 'ellipsis',
      whiteSpace: 'nowrap',
    }),
    multiValueRemove: (def) => ({
      ...def,
      ...chipSurface,
      alignSelf: 'stretch',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      borderRadius: 0,
      paddingTop: 0,
      paddingBottom: 0,
      paddingLeft: '4px',
      paddingRight: '6px',
      margin: 0,
      cursor: 'pointer',
      ':hover': {
        backgroundColor: chipRemoveHoverBg,
        color: theme.palette.error.main,
      },
    }),
    valueContainer: (base) => ({
      ...base,
      flexWrap: 'wrap',
      paddingLeft: '8px',
      paddingRight: '8px',
    }),
    control: (base, state) => ({
      ...base,
      backgroundColor: controlBackground,
      borderColor: state?.isFocused || state?.menuIsOpen ? activeBorder : restingBorder,
      color: optionTextColor,
      boxShadow: 'none',
      minHeight: 40,
      height: 'auto',
      cursor: 'pointer',
      '&:hover': { borderColor: activeBorder },
    }),
    dropdownIndicator: (base) => ({ ...base, ...indicatorColor, padding: '4px 8px' }),
    indicatorSeparator: (base) => ({
      ...base,
      backgroundColor: restingBorder,
      marginTop: 8,
      marginBottom: 8,
    }),
    clearIndicator: (base) => ({ ...base, ...indicatorColor, padding: '4px' }),
    menu: (base) => ({
      ...base,
      backgroundColor: menuBackground,
      paddingLeft: '4px',
      paddingRight: '4px',
    }),
    menuList: (base) => ({
      ...base,
      paddingTop: '4px',
      paddingBottom: '4px',
    }),
    placeholder: (base, state) => ({
      ...base,
      ...chipType,
      color: theme.palette.text.disabled,
      display: state?.isFocused || state?.selectProps?.inputValue ? 'none' : base.display,
    }),
    input: (base) => ({
      ...base,
      ...chipType,
      color: optionTextColor,
    }),
  };

  return (
    <CreatableSelect
      {...props}
      inputValue={selectInput}
      onInputChange={onInputChange}
      onKeyDown={onKeyDown}
      options={filterOptions(props.options, selectInput)}
      onChange={handleChange}
      components={{
        Option: Option,
        Input: CustomInput,
        Menu: Menu,
        DropdownIndicator,
        ...props.components,
      }}
      filterOption={customFilterOption}
      noOptionsMessage={props.noOptionsMessage}
      menuPlacement={props.menuPlacement ?? 'auto'}
      styles={customStyles}
      isMulti
      closeMenuOnSelect={false}
      tabSelectsValue={false}
      backspaceRemovesValue={false}
      hideSelectedOptions={false}
      isDisabled={Boolean(props.updating || props.disabled)}
      blurInputOnSelect={false}
    />
  );
};

export default MultiSelectWrapper;
