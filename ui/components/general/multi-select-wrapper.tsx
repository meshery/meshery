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

/** Brand-tinted option row highlight; exported for contrast / regression tests. */
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

  // Closed control should sit flush with the table row (transparent over the
  // cell/row surface). The open menu keeps the solid card fill separately.
  const controlBackground = 'transparent';

  // Opaque chip fill (action.selected is translucent and breaks getContrastText /
  // paint contrast). Prefer elevated/secondary; fall back to paper so dark mode
  // never paints a near-black chip with an invisible remove control.
  const chipBackground =
    theme.palette.mode === 'dark'
      ? (theme.palette.background.elevatedComponents ??
        theme.palette.background.paper ??
        theme.palette.background.default)
      : (theme.palette.background.secondary ?? theme.palette.background.paper);
  const chipForeground =
    theme.palette.getContrastText?.(chipBackground) ?? theme.palette.text.primary;

  const optionHighlight = getMultiSelectOptionHighlight(theme);
  // Quieter fill for selected (checked) rows so keyboard/hover focus (optionHighlight)
  // stays visually distinct from "already selected".
  const optionSelectedFill = alpha(
    theme.palette.primary.main,
    theme.palette.mode === 'dark' ? 0.14 : 0.06,
  );
  // Sistent body token; fall back to MUI text.primary if a theme omits default.
  const optionTextColor = theme.palette.text.default ?? theme.palette.text.primary;

  const Option = (props) => {
    return (
      <ListItemButton
        ref={props.innerRef}
        // Checkbox + fontWeight indicate selection; Mui-selected is not used for focus
        // so keyboard focus stays distinct among many already-selected options.
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
          // react-select sets isFocused for both keyboard and mouse-hovered rows.
          // Keep `&:hover` as a defensive fallback so MUI default hover styles
          // cannot wash the row back to a faint action.hover.
          '&:hover': {
            backgroundColor: optionHighlight,
          },
          // Selected-but-unfocused: quiet wash (checkbox already marks selection).
          // When also focused, keep the brand wash so keyboard position stays visible.
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

  // Keep Input as a direct child of the value container so react-select's
  // grid-area stacking with the placeholder stays intact (no caret overlap).
  const CustomInput = (props) => (
    <components.Input autoFocus={props.selectProps.menuIsOpen} {...props}>
      {props.children}
    </components.Input>
  );

  const Menu = (props) => {
    // Apply react-select menu styles (background, horizontal padding, etc.).
    // Custom components.Menu replaces the default node, so getStyles must be
    // re-applied or customStyles.menu is never painted. Spread first, then force
    // layout so styles.menu cannot clobber position/zIndex.
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

  // Same caret as MUI Select / Connection status dropdown (Sistent ArrowDropDown).
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

  // Soft error wash for the remove control — never action.hover (washes white
  // in dark mode) and never the default react-select coral flash.
  const chipRemoveHoverBg = alpha(
    theme.palette.error.main,
    theme.palette.mode === 'dark' ? 0.28 : 0.12,
  );
  const chipRemoveHoverFg = theme.palette.error.main;

  const customStyles = {
    // Container must share the same solid fill as label/remove. react-select's
    // default multiValue is light grey/white and shows as bright border seams
    // when only the children are recolored.
    multiValue: (base) => ({
      ...base,
      display: 'flex',
      // Stretch label + remove so remove hover wash fills full chip height.
      alignItems: 'stretch',
      backgroundColor: chipBackground,
      border: 'none',
      borderRadius: '4px',
      overflow: 'hidden',
      margin: '2px',
      maxWidth: '100%',
    }),
    multiValueLabel: (def) => ({
      ...def,
      display: 'flex',
      alignItems: 'center',
      backgroundColor: chipBackground,
      color: chipForeground,
      border: 'none',
      borderRadius: 0,
      // body2 matches Connections chips / Sistent table density; stock multiValue is smaller.
      fontSize: theme.typography.body2.fontSize,
      lineHeight: theme.typography.body2.lineHeight,
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
      // Stretch to chip height so :hover paints the full remove column, not a
      // short icon-sized pill (react-select defaults leave this centered).
      alignSelf: 'stretch',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      backgroundColor: chipBackground,
      color: chipForeground,
      border: 'none',
      borderRadius: 0,
      // Match label vertical padding so the strip is full-height without growing
      // taller than the text side.
      paddingTop: 0,
      paddingBottom: 0,
      paddingLeft: '4px',
      paddingRight: '6px',
      margin: 0,
      cursor: 'pointer',
      // Kill react-select's default :hover (light coral / near-white) entirely.
      ':hover': {
        backgroundColor: chipRemoveHoverBg,
        color: chipRemoveHoverFg,
      },
    }),
    // Grow with chips instead of a fixed max-height + scrollbar.
    valueContainer: (base) => ({
      ...base,
      flexWrap: 'wrap',
      paddingLeft: '8px',
      paddingRight: '8px',
    }),
    // Resting outline matches Sistent OutlinedInput (`palette.border.strong`),
    // not faint `divider` — brand/primary only on hover, focus, or open menu.
    // See sistent outlinedinput.modifier.ts + palette.border tokens.
    control: (base, state) => {
      const restingBorder = theme.palette.border?.strong ?? theme.palette.divider;
      const activeBorder = theme.palette.primary.main;
      return {
        ...base,
        backgroundColor: controlBackground,
        borderColor: state.isFocused || state.menuIsOpen ? activeBorder : restingBorder,
        color: optionTextColor,
        boxShadow: 'none',
        minHeight: 40,
        height: 'auto',
        cursor: 'pointer',
        '&:hover': {
          borderColor: activeBorder,
        },
      };
    },
    dropdownIndicator: (base) => ({
      ...base,
      color: theme.palette.action.active,
      padding: '4px 8px',
      ':hover': {
        color: theme.palette.text.primary,
      },
    }),
    // Default react-select mid-bar between clear and caret (table multi-select affordance).
    indicatorSeparator: (base) => ({
      ...base,
      backgroundColor: theme.palette.border?.strong ?? theme.palette.divider,
      marginTop: 8,
      marginBottom: 8,
    }),
    clearIndicator: (base) => ({
      ...base,
      color: theme.palette.action.active,
      padding: '4px',
      ':hover': {
        color: theme.palette.text.primary,
      },
    }),
    menu: (base) => ({
      ...base,
      backgroundColor: menuBackground,
      paddingLeft: '4px',
      paddingRight: '4px',
    }),
    menuList: (base) => ({
      ...base,
      // Small padding keeps the create-option from sitting flush against the border
      // while removing the large empty bands above first / below last option.
      paddingTop: '4px',
      paddingBottom: '4px',
    }),
    placeholder: (base, state) => ({
      ...base,
      color: theme.palette.text.disabled,
      fontSize: theme.typography.body2.fontSize,
      // Hide once focused or while typing so the caret never overlaps ghost text.
      display: state.isFocused || state.selectProps.inputValue ? 'none' : base.display,
    }),
    input: (base) => ({
      ...base,
      color: optionTextColor,
      fontSize: theme.typography.body2.fontSize,
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
