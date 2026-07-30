import { useState } from 'react';
import { components } from 'react-select';
import CreatableSelect from 'react-select/creatable';
import {
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
  // paint contrast). elevatedComponents / background.secondary are mode-aware
  // solids — not grey[300].
  const chipBackground =
    theme.palette.mode === 'dark'
      ? theme.palette.background.elevatedComponents
      : theme.palette.background.secondary;
  const chipForeground = theme.palette.getContrastText(chipBackground);

  const optionHighlight = getMultiSelectOptionHighlight(theme);

  const Option = (props) => {
    return (
      <ListItemButton
        ref={props.innerRef}
        selected={props.isFocused || props.isSelected}
        {...props.innerProps}
        component="div"
        data-testid="multi-select-option"
        sx={{
          fontWeight: props.isSelected ? 500 : 400,
          padding: '0.4rem 1rem',
          backgroundColor: props.isFocused ? optionHighlight : 'transparent',
          color: theme.palette.text.primary,
          // Intentional redundancy: react-select sets isFocused for both keyboard and
          // mouse-hovered rows, so the ternary already covers hover. Keep `&:hover`
          // (and `&.Mui-selected:hover`) as a defensive fallback so MUI default hover
          // styles cannot wash the row back to a faint action.hover.
          '&:hover': {
            backgroundColor: optionHighlight,
          },
          '&.Mui-selected': {
            backgroundColor: optionHighlight,
            '&:hover': {
              backgroundColor: optionHighlight,
            },
          },
          '& .MuiFormControlLabel-label': {
            color: theme.palette.text.primary,
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
    return (
      <Paper
        square
        style={{
          zIndex: 9999,
          width: '100%',
          position: 'absolute',
          backgroundColor: menuBackground,
        }}
        {...props.innerProps}
      >
        {props.children}
      </Paper>
    );
  };

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

  const customStyles = {
    multiValueLabel: (def) => ({
      ...def,
      backgroundColor: chipBackground,
      color: chipForeground,
    }),
    multiValueRemove: (def) => ({
      ...def,
      backgroundColor: chipBackground,
      color: chipForeground,
      // Keep the x readable on hover; only the wash changes.
      ':hover': {
        backgroundColor: theme.palette.action.hover,
        color: chipForeground,
      },
    }),
    valueContainer: (base) => ({
      ...base,
      maxHeight: '65px',
      overflow: 'auto',
      paddingLeft: '8px',
      paddingRight: '8px',
    }),
    control: (base) => ({
      ...base,
      backgroundColor: controlBackground,
      borderColor: theme.palette.primary.main,
      color: theme.palette.primary.main,
      boxShadow: 'none',
      '&:hover': {
        borderColor: theme.palette.primary.main,
      },
      '&$focused': {
        borderColor: theme.palette.primary.main,
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
      // Hide once focused or while typing so the caret never overlaps ghost text.
      display: state.isFocused || state.selectProps.inputValue ? 'none' : base.display,
    }),
    input: (base) => ({
      ...base,
      color: theme.palette.text.primary,
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
      isDisabled={props.updating}
      blurInputOnSelect={false}
    />
  );
};

export default MultiSelectWrapper;
