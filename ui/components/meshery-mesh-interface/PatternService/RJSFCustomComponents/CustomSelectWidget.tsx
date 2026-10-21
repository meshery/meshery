import React from 'react';
import {
  IconButton,
  InputAdornment,
  ListItemText,
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
  // hideError,
  formContext,
  ...textFieldProps
}) {
  const { enumOptions, enumDisabled, emptyValue: optEmptyVal } = options;
  const xRjsfGridArea = schema?.['x-rjsf-grid-area']; // check if the field is used in different modal (e.g. publish)

  // RJSF may sometimes pass multiple=false or undefined even for array enums.
  // We force it to true if the schema type is array or it has items.
  const isMultiple =
    multiple || schema?.type === 'array' || schema?.items !== undefined || options?.multiple;

  const emptyValue = isMultiple ? [] : '';
  const isEmpty =
    typeof value === 'undefined' ||
    (isMultiple && value.length < 1) ||
    (!isMultiple && value === emptyValue);

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
        error={rawErrors?.length > 0}
        onChange={_onChange}
        onBlur={_onBlur}
        onFocus={_onFocus}
        size="small"
        InputProps={{
          style: { paddingRight: '0px' },
          endAdornment: (
            <InputAdornment position="start" style={{ position: 'absolute', right: '1rem' }}>
              {rawErrors?.length > 0 && (
                <CustomTextTooltip
                  bgColor={theme.palette.error.main}
                  flag={formContext?.overrideFlag}
                  title={rawErrors?.join('  ')}
                  interactive={true}
                >
                  <IconButton component="span" size="small">
                    <ErrorOutlineIcon
                      width="14px"
                      height="14px"
                      fill={theme.palette.error.main}
                      style={{ verticalAlign: 'middle', ...iconSmall }}
                    />
                  </IconButton>
                </CustomTextTooltip>
              )}
              {typeof schema?.description === 'string' && schema.description && (
                <CustomTextTooltip
                  flag={formContext?.overrideFlag}
                  title={schema.description}
                  interactive={true}
                >
                  <IconButton component="span" size="small" style={{ marginRight: '4px' }}>
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
          ),
        }}
        {...textFieldProps}
        select
        InputLabelProps={{
          ...textFieldProps.InputLabelProps,
          shrink: !isEmpty,
        }}
        SelectProps={{
          ...textFieldProps.SelectProps,
          renderValue: (selected) => {
            if (isMultiple && Array.isArray(selected)) {
              return selected.map((i) => safeDisplayValue(enumOptions?.[i]?.label)).join(', ');
            }
            const idx = selected as number;
            return safeDisplayValue(enumOptions?.[idx]?.label);
          },
          multiple: isMultiple,
          MenuProps: {
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
            transformOrigin: {
              vertical: 'top',
              horizontal: 'left',
            },
            getContentAnchorEl: null,
            PaperProps: {
              style: {
                maxHeight: '400px',
              },
            },
          },
        }}
        aria-describedby={ariaDescribedByIds(id)}
      >
        {Array.isArray(enumOptions) &&
          enumOptions.map(({ value, label }, i) => {
            const disabled = Array.isArray(enumDisabled) && enumDisabled?.indexOf(value) !== -1;
            return (
              <MenuItem key={i} value={String(i)} disabled={disabled}>
                {isMultiple && <Checkbox checked={selectedIndexes?.indexOf(String(i)) !== -1} />}
                <ListItemText primary={safeDisplayValue(label)} />
              </MenuItem>
            );
          })}
      </TextField>
    </>
  );
}
