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
        slotProps={{
          ...incomingSlotProps,
          input: {
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
            ...legacyInputProps,
            ...incomingSlotProps.input,
          },
          inputLabel: {
            shrink: !isEmpty,
            ...legacyInputLabelProps,
            ...incomingSlotProps.inputLabel,
          },
          select: {
            multiple: isMultiple,
            renderValue: (selected) => {
              if (isMultiple && Array.isArray(selected)) {
                return selected.map((i, index) => {
                  const rawLabel = enumOptions?.[i]?.label;
                  const labelNode = React.isValidElement(rawLabel)
                    ? rawLabel
                    : safeDisplayValue(rawLabel);
                  return (
                    <React.Fragment key={i}>
                      {labelNode}
                      {index < selected.length - 1 ? ', ' : ''}
                    </React.Fragment>
                  );
                });
              }
              const idx = selected as number;
              const rawLabel = enumOptions?.[idx]?.label;
              return React.isValidElement(rawLabel) ? rawLabel : safeDisplayValue(rawLabel);
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
              PaperProps: {
                style: {
                  maxHeight: '400px',
                },
              },
            },
            ...legacySelectProps,
            ...incomingSlotProps.select,
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
                style={{
                  display: 'flex',
                  flexDirection: 'row',
                  alignItems: 'center',
                  flexWrap: 'nowrap',
                  gap: '0.5rem',
                  paddingRight: '2rem',
                }}
                sx={{
                  display: 'flex !important',
                  flexDirection: 'row !important',
                  alignItems: 'center !important',
                  flexWrap: 'nowrap !important',
                  gap: '0.5rem !important',
                  paddingRight: '2rem !important',
                }}
              >
                {isMultiple && (
                  <Checkbox
                    checked={
                      Array.isArray(selectedIndexes) ? selectedIndexes.includes(String(i)) : false
                    }
                    style={{ padding: 0, flexShrink: 0, marginRight: '0.25rem' }}
                    sx={{
                      padding: '0 !important',
                      flexShrink: '0 !important',
                      marginRight: '0.25rem !important',
                    }}
                  />
                )}
                <CustomTextTooltip
                  flag={formContext?.overrideFlag}
                  title={optionLabel}
                  interactive={true}
                >
                  <ListItemText
                    primary={optionLabel}
                    primaryTypographyProps={{
                      noWrap: true,
                      style: {
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                      },
                    }}
                    style={{
                      margin: 0,
                      minWidth: 0,
                      flex: '1 1 auto',
                      overflow: 'hidden',
                    }}
                  />
                </CustomTextTooltip>
              </MenuItem>
            );
          })}
      </TextField>
    </>
  );
}
