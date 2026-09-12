import React from 'react';
import { FormControlLabel, FormLabel, Radio, RadioGroup } from '@sistent/sistent';
import {
  ariaDescribedByIds,
  enumOptionsIndexForValue,
  enumOptionsValueForIndex,
  optionId,
} from '@rjsf/utils';
import { safeDisplayValue } from '../helper';

export default function CustomRadioWidget(props) {
  const {
    id,
    htmlName,
    options = {},
    value,
    required,
    disabled,
    readonly,
    label,
    hideLabel,
    onChange,
    onBlur,
    onFocus,
  } = props;

  const { enumOptions = [], enumDisabled, emptyValue } = options;

  const _onChange = (_e, val) => onChange(enumOptionsValueForIndex(val, enumOptions, emptyValue));
  const _onBlur = ({ target }) =>
    onBlur(id, enumOptionsValueForIndex(target?.value, enumOptions, emptyValue));
  const _onFocus = ({ target }) =>
    onFocus(id, enumOptionsValueForIndex(target?.value, enumOptions, emptyValue));

  const row = options?.inline === true;
  const selectedIndex = enumOptionsIndexForValue(value, enumOptions) ?? null;

  const safeLabel =
    label != null && typeof label === 'object' && !React.isValidElement(label)
      ? String(label)
      : label;

  return (
    <>
      {!hideLabel && (
        <FormLabel required={required} htmlFor={id}>
          {safeLabel ?? undefined}
        </FormLabel>
      )}
      <RadioGroup
        id={id}
        name={htmlName || id}
        value={selectedIndex}
        row={row}
        onChange={_onChange}
        onBlur={_onBlur}
        onFocus={_onFocus}
        aria-describedby={ariaDescribedByIds(id)}
      >
        {Array.isArray(enumOptions) &&
          enumOptions.map((option, index) => {
            const itemDisabled =
              Array.isArray(enumDisabled) && enumDisabled.indexOf(option.value) !== -1;
            return (
              <FormControlLabel
                key={index}
                control={<Radio name={htmlName || id} id={optionId(id, index)} color="primary" />}
                label={safeDisplayValue(option.label)}
                value={String(index)}
                disabled={disabled || itemDisabled || readonly}
              />
            );
          })}
      </RadioGroup>
    </>
  );
}
