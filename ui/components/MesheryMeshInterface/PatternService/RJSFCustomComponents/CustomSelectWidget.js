import React from 'react';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import NativeSelect from '@material-ui/core/NativeSelect';
import { IconButton, InputAdornment, OutlinedInput } from '@material-ui/core';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import { ERROR_COLOR } from '../../../../constants/colors';
import { iconSmall } from '../../../../css/icons.styles';
import theme from '../../../../themes/app';
import { CustomTextTooltip } from '../CustomTextTooltip';
import { getHyperLinkDiv } from '../helper';
import {
  ariaDescribedByIds,
  enumOptionsIndexForValue,
  enumOptionsValueForIndex,
  labelValue,
} from '@rjsf/utils';

export default function CustomSelectWidget({
  schema,
  id,
  name,
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
  rawErrors = [],
  registry,
  uiSchema,
  hideError,
  formContext,
  ...nativeSelectProps
}) {
  const { enumOptions, enumDisabled, emptyValue: optEmptyVal } = options;

  multiple = typeof multiple === 'undefined' ? false : !!multiple;

  const emptyValue = multiple ? [] : '';
  const isEmpty =
    typeof value === 'undefined' ||
    (multiple && value.length < 1) ||
    (!multiple && value === emptyValue);

  const _onChange = (event) => {
    const selectedIndex = event.target.selectedIndex;
    onChange(
      enumOptionsValueForIndex(selectedIndex, enumOptions, optEmptyVal)
    );
  };
  const _onBlur = (event) => {
    const selectedIndex = event.target.selectedIndex;
    onBlur(
      id,
      enumOptionsValueForIndex(selectedIndex, enumOptions, optEmptyVal)
    );
  };
  const _onFocus = (event) => {
    const selectedIndex = event.target.selectedIndex;
    onFocus(
      id,
      enumOptionsValueForIndex(selectedIndex, enumOptions, optEmptyVal)
    );
  };
  const selectedIndexes = enumOptionsIndexForValue(
    value,
    enumOptions,
    multiple
  );

  return (
    <FormControl variant="outlined" margin='dense' style={{width : "calc(100% - 4px)"}}>
      <InputLabel>
        {labelValue(label, hideLabel || !label, false)}
      </InputLabel>
      <NativeSelect
        id={id}
        name={id}
        value={isEmpty ? emptyValue : selectedIndexes}
        required={required}
        disabled={disabled || readonly}
        input={<OutlinedInput label={label} />}
        endAdornment={
             (<InputAdornment position="start">
          {rawErrors?.length > 0 && (
            <CustomTextTooltip
              backgroundColor={ERROR_COLOR}
              flag={formContext?.overrideFlag}
              title={
                <div>
                  {rawErrors?.map((error, index) => (
                    <div key={index}>{error}</div>
                  ))}
                </div>

              }
              interactive={true}>
              <IconButton component="span" size="small">
                <ErrorOutlineIcon width="14px" height="14px" fill={ERROR_COLOR} style={{ verticalAlign : "middle", ...iconSmall }} />
              </IconButton>
            </CustomTextTooltip>
          )}
          {schema?.description && (
            <CustomTextTooltip backgroundColor="#3C494F" flag={formContext?.overrideFlag} title={getHyperLinkDiv(schema?.description)} interactive={true}>
              <IconButton component="span" size="small">
                <HelpOutlineIcon width="14px" height="14px" fill={theme.palette.type === 'dark' ? "white" : "gray"} style={{ verticalAlign : "middle", ...iconSmall }} />
              </IconButton>
            </CustomTextTooltip>
          )}
        </InputAdornment>)
        }
        autoFocus={autofocus}
        onChange={_onChange}
        onBlur={_onBlur}
        onFocus={_onFocus}
        {...nativeSelectProps}
        multiple={multiple}
        inputProps={{
          'aria-describedby': ariaDescribedByIds(id),
        }}
      >
        {Array.isArray(enumOptions) &&
          enumOptions.map(({ value, label }, i) => {
            const disabled =
              Array.isArray(enumDisabled) &&
              enumDisabled.indexOf(value) !== -1;
            return (
              <option key={i} value={String(i)} disabled={disabled}>
                {label}
              </option>
            );
          })}
      </NativeSelect>
    </FormControl>
  );
}