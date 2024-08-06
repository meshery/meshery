import React from 'react';
import {
  IconButton,
  InputAdornment,
  ListItemText,
  MenuItem,
  TextField,
  InputLabel,
} from '@material-ui/core';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import { ERROR_COLOR } from '../../../../constants/colors';
import { iconSmall } from '../../../../css/icons.styles';
import theme from '../../../../themes/app';
import { CustomTextTooltip } from '../CustomTextTooltip';
import {
  ariaDescribedByIds,
  enumOptionsIndexForValue,
  enumOptionsValueForIndex,
  labelValue,
} from '@rjsf/utils';
import { Checkbox } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

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

  multiple = typeof multiple === 'undefined' ? false : !!multiple;
  const emptyValue = multiple ? [] : '';
  const isEmpty =
    typeof value === 'undefined' ||
    (multiple && value.length < 1) ||
    (!multiple && value === emptyValue);

  const _onChange = ({ target: { value } }) =>
    onChange(enumOptionsValueForIndex(value, enumOptions, optEmptyVal));
  const _onBlur = ({ target: { value } }) =>
    onBlur(id, enumOptionsValueForIndex(value, enumOptions, optEmptyVal));
  const _onFocus = ({ target: { value } }) =>
    onFocus(id, enumOptionsValueForIndex(value, enumOptions, optEmptyVal));
  const selectedIndexes = enumOptionsIndexForValue(value, enumOptions, multiple);

  return (
    <>
      {xRjsfGridArea && (
        <InputLabel required={required} htmlFor={id}>
          {labelValue(label, hideLabel || !label, false)}
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
        label={xRjsfGridArea ? '' : labelValue(label, hideLabel || !label, false)}
        error={rawErrors?.length > 0}
        onChange={_onChange}
        onBlur={_onBlur}
        onFocus={_onFocus}
        InputProps={{
          style: { paddingRight: '0px' },
          endAdornment: (
            <InputAdornment position="start" style={{ position: 'absolute', right: '1rem' }}>
              {rawErrors?.length > 0 && (
                <CustomTextTooltip
                  bgColor={ERROR_COLOR}
                  flag={formContext?.overrideFlag}
                  title={rawErrors?.join('  ')}
                  interactive={true}
                >
                  <IconButton component="span" size="small">
                    <ErrorOutlineIcon
                      width="14px"
                      height="14px"
                      fill={ERROR_COLOR}
                      style={{ verticalAlign: 'middle', ...iconSmall }}
                    />
                  </IconButton>
                </CustomTextTooltip>
              )}
              {schema?.description && (
                <CustomTextTooltip
                  flag={formContext?.overrideFlag}
                  title={schema?.description}
                  interactive={true}
                >
                  <IconButton component="span" size="small" style={{ marginRight: '4px' }}>
                    <HelpOutlineIcon
                      width="14px"
                      height="14px"
                      fill={theme.palette.type === 'dark' ? 'white' : 'gray'}
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
            if (multiple) {
              return selected.map((index) => enumOptions[index].label).join(', ');
            }
            return enumOptions[selected].label;
          },
          multiple,
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
                {multiple && (
                  <UsesSistent>
                    <Checkbox checked={selectedIndexes?.indexOf(String(i)) !== -1} />
                  </UsesSistent>
                )}
                <ListItemText primary={label} />
              </MenuItem>
            );
          })}
      </TextField>
    </>
  );
}
