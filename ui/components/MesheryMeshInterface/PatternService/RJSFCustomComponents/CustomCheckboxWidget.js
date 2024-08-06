import { FormControlLabel, IconButton, useTheme } from '@material-ui/core';
import { labelValue, schemaRequiresTrueValue } from '@rjsf/utils';
import React from 'react';
import { CustomTextTooltip } from '../CustomTextTooltip';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import { iconSmall } from '../../../../css/icons.styles';
import { UsesSistent } from '@/components/SistentWrapper';
import { Checkbox } from '@layer5/sistent';

export const CustomCheckboxWidget = (props) => {
  const {
    schema,
    id,
    value,
    disabled,
    readonly,
    label = '',
    hideLabel,
    autofocus,
    onChange,
  } = props;

  const required = schemaRequiresTrueValue(schema);
  const theme = useTheme();

  const _onChange = ({ target: { checked } }) => onChange(checked);

  return (
    <>
      <FormControlLabel
        control={
          <UsesSistent>
            <Checkbox
              id={id}
              name={id}
              checked={typeof value === 'undefined' ? false : Boolean(value)}
              required={required}
              disabled={disabled || readonly}
              autoFocus={autofocus}
              onChange={_onChange}
            />
          </UsesSistent>
        }
        label={
          <>
            {labelValue(label, hideLabel, required)}
            {schema.description && (
              <CustomTextTooltip
                flag={props?.formContext?.overrideFlag}
                title={schema?.description}
                interactive={true}
              >
                <IconButton component="span" size="small">
                  <HelpOutlineIcon
                    width="14px"
                    height="14px"
                    fill={theme.palette.type === 'dark' ? 'white' : 'gray'}
                    style={{ verticalAlign: 'middle', ...iconSmall }}
                  />
                </IconButton>
              </CustomTextTooltip>
            )}
          </>
        }
      />
    </>
  );
};
