import React from 'react';
import {
  IconButton,
  InputAdornment,
  Box,
  MenuItem,
  TextField,
  InputLabel,
  useTheme,
} from '@sistent/sistent';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import ErrorOutlineIcon from '../../../../assets/icons/ErrorOutlineIcon';
import { iconSmall } from '../../../../css/icons.styles';
import { CustomTextTooltip } from '../CustomTextTooltip';
import {
  ariaDescribedByIds,
  enumOptionsIndexForValue,
  enumOptionsValueForIndex,
  labelValue,
} from '@rjsf/utils';
import { safeDisplayValue } from '../helper';
import { Checkbox } from '@sistent/sistent';
// Sistent does not export MUI's slot composition utility.
// eslint-disable-next-line no-restricted-imports
import { mergeSlotProps } from '@mui/material/utils';

const resolveSlot = (props, ownerState) =>
  (typeof props === 'function' ? props(ownerState) : props) || {};

// Resolve each callback before merging: MUI's callback merger also copies
// ownerState.className, which can apply input sizing to an unrelated slot.
const mergeSlots = (ownerState, ...sources) =>
  sources.reduce((merged, source) => mergeSlotProps(resolveSlot(source, ownerState), merged), {});

export default function CustomSelectWidget({
  schema,
  id,
  options,
  label,
  hideLabel,
  required,
  disabled,
  readonly,
  placeholder,
  value,
  multiple,
  autofocus,
  onChange,
  onBlur,
  onFocus,
  rawErrors,
  // registry,
  // uiSchema,
  hideError,
  formContext,
  ...textFieldProps
}) {
  const {
    InputProps: legacyInputProps,
    InputLabelProps: legacyInputLabelProps,
    SelectProps: legacySelectProps,
    slotProps: incomingSlotProps = {},
    ...cleanTextFieldProps
  } = textFieldProps;

  const { enumOptions, enumDisabled, emptyValue: optEmptyVal } = options;
  const xRjsfGridArea = schema?.['x-rjsf-grid-area']; // check if the field is used in different modal (e.g. publish)

  // RJSF may sometimes pass multiple=false or undefined even for array enums.
  // We force it to true if the schema type is array or it has items.
  const isMultiple =
    multiple || schema?.type === 'array' || schema?.items !== undefined || options?.multiple;

  const emptyValue = isMultiple ? [] : '';
  const isEmpty =
    value == null || (isMultiple && value.length < 1) || (!isMultiple && value === emptyValue);

  const _onChange = ({ target: { value } }) => {
    let safeValue = value;
    if (isMultiple && typeof value === 'string') {
      safeValue = value.split(',').filter(Boolean);
    }
    const nextValue = enumOptionsValueForIndex(safeValue, enumOptions, optEmptyVal);
    onChange(isMultiple && !Array.isArray(nextValue) ? (nextValue ? [nextValue] : []) : nextValue);
  };
  const _onBlur = ({ target: { value } }) =>
    onBlur(id, enumOptionsValueForIndex(value, enumOptions, optEmptyVal));
  const _onFocus = ({ target: { value } }) =>
    onFocus(id, enumOptionsValueForIndex(value, enumOptions, optEmptyVal));
  const selectedIndexes = enumOptionsIndexForValue(value, enumOptions, isMultiple);
  const theme = useTheme();
  const showError = !hideError && rawErrors?.length > 0;
  const hasDescription = typeof schema?.description === 'string' && !!schema.description;

  const labelContent = labelValue(label, hideLabel || !label, false);
  const safeLabel = safeDisplayValue(labelContent);

  return (
    <>
      {xRjsfGridArea && (
        <InputLabel required={required} htmlFor={id}>
          {safeLabel}
        </InputLabel>
      )}
      <TextField
        id={id}
        name={id}
        value={isEmpty ? emptyValue : selectedIndexes}
        required={required}
        disabled={disabled || readonly}
        autoFocus={autofocus}
        placeholder={placeholder}
        label={xRjsfGridArea ? '' : safeLabel}
        error={showError}
        onChange={_onChange}
        onBlur={_onBlur}
        onFocus={_onFocus}
        size="small"
        slotProps={{
          ...incomingSlotProps,
          input: (ownerState: any) => {
            const inputProps = mergeSlots(ownerState, legacyInputProps, incomingSlotProps.input);
            const callerEndAdornment = inputProps.endAdornment;
            const hasWidgetAdornment = showError || hasDescription;

            const widgetAdornment = hasWidgetAdornment ? (
              <InputAdornment position="end" sx={{ flexShrink: 0, ml: 0.5 }}>
                {showError && (
                  <CustomTextTooltip
                    bgColor={theme.palette.error.main}
                    flag={formContext?.overrideFlag}
                    title={rawErrors?.join('  ')}
                    disableInteractive={false}
                  >
                    <IconButton
                      component="span"
                      size="small"
                      aria-label="Validation errors"
                      tabIndex={0}
                    >
                      <ErrorOutlineIcon
                        width="14px"
                        height="14px"
                        fill={theme.palette.error.main}
                        style={{ verticalAlign: 'middle', ...iconSmall }}
                      />
                    </IconButton>
                  </CustomTextTooltip>
                )}
                {hasDescription && (
                  <CustomTextTooltip
                    flag={formContext?.overrideFlag}
                    title={schema.description}
                    disableInteractive={false}
                  >
                    <IconButton component="span" size="small" aria-label="Field help" tabIndex={0}>
                      <HelpOutlineIcon
                        width="14px"
                        height="14px"
                        fill={theme.palette.mode === 'dark' ? 'white' : 'gray'}
                        style={{ verticalAlign: 'middle', ...iconSmall }}
                      />
                    </IconButton>
                  </CustomTextTooltip>
                )}
              </InputAdornment>
            ) : null;

            const finalEndAdornment =
              callerEndAdornment && widgetAdornment ? (
                <>
                  {callerEndAdornment}
                  {widgetAdornment}
                </>
              ) : (
                callerEndAdornment || widgetAdornment
              );

            return {
              ...inputProps,
              // Keep the chevron in the select's own padding box; adornments follow it.
              sx: [
                {
                  minWidth: 0,
                  '& .MuiSelect-select': { minWidth: 0 },
                  '& .MuiSelect-icon': { position: 'static', flexShrink: 0, ml: -3, mr: 1 },
                },
                ...(Array.isArray(inputProps.sx) ? inputProps.sx : [inputProps.sx]),
              ],
              endAdornment: finalEndAdornment,
            };
          },
          inputLabel: (ownerState) => ({
            ...mergeSlots(ownerState, legacyInputLabelProps, incomingSlotProps.inputLabel),
            shrink: !isEmpty,
          }),
          select: (ownerState: any) => {
            const legacy = resolveSlot(legacySelectProps, ownerState);
            const incoming = resolveSlot(incomingSlotProps.select, ownerState);
            const { PaperProps: _paperProps, ...menuProps } = mergeSlotProps(
              incoming.MenuProps,
              legacy.MenuProps || {},
            );

            return {
              ...mergeSlotProps(incoming, legacy),
              multiple: isMultiple,
              renderValue: (selected: any) => {
                const indexes = isMultiple
                  ? Array.isArray(selected)
                    ? selected
                    : String(selected).split(',').filter(Boolean)
                  : [selected];
                const labels = indexes.map((i, index) => (
                  <React.Fragment key={i}>
                    {index > 0 && ', '}
                    {safeDisplayValue(enumOptions?.[i]?.label)}
                  </React.Fragment>
                ));
                return (
                  <CustomTextTooltip title={<>{labels}</>} flag={formContext?.overrideFlag}>
                    <Box
                      component="span"
                      sx={{
                        display: 'block',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                      }}
                    >
                      {labels}
                    </Box>
                  </CustomTextTooltip>
                );
              },
              MenuProps: {
                anchorOrigin: {
                  vertical: 'bottom',
                  horizontal: 'left',
                },
                transformOrigin: {
                  vertical: 'top',
                  horizontal: 'left',
                },
                ...menuProps,
                slotProps: {
                  ...legacy.MenuProps?.slotProps,
                  ...incoming.MenuProps?.slotProps,
                  paper: (paperOwnerState) =>
                    mergeSlots(
                      paperOwnerState,
                      { style: { maxHeight: '400px' } },
                      legacy.MenuProps?.PaperProps,
                      legacy.MenuProps?.slotProps?.paper,
                      incoming.MenuProps?.PaperProps,
                      incoming.MenuProps?.slotProps?.paper,
                    ),
                },
              },
            };
          },
        }}
        {...cleanTextFieldProps}
        select
        aria-describedby={ariaDescribedByIds(id)}
      >
        {Array.isArray(enumOptions) &&
          enumOptions.map(({ value, label }, i) => {
            const disabled = Array.isArray(enumDisabled) && enumDisabled?.indexOf(value) !== -1;
            const optionLabel = React.isValidElement(label) ? label : safeDisplayValue(label);
            return (
              <MenuItem
                key={i}
                value={String(i)}
                disabled={disabled}
                sx={{
                  display: 'flex',
                  flexDirection: 'row',
                  alignItems: 'center',
                  flexWrap: 'nowrap',
                  gap: 1,
                  minWidth: 0,
                }}
              >
                {isMultiple && (
                  <Checkbox
                    tabIndex={-1}
                    disableRipple
                    slotProps={{ input: { 'aria-hidden': true } }}
                    checked={
                      Array.isArray(selectedIndexes) ? selectedIndexes.includes(String(i)) : false
                    }
                    sx={{
                      padding: 0,
                      flexShrink: 0,
                    }}
                  />
                )}
                <CustomTextTooltip
                  flag={formContext?.overrideFlag}
                  title={optionLabel}
                  disableInteractive={false}
                >
                  <Box
                    component="span"
                    sx={{
                      minWidth: 0,
                      flex: 1,
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                    }}
                  >
                    {optionLabel}
                  </Box>
                </CustomTextTooltip>
              </MenuItem>
            );
          })}
      </TextField>
    </>
  );
}
