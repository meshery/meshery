import React from 'react';
import DeleteIcon from '@mui/icons-material/Delete';
import { ADDITIONAL_PROPERTY_FLAG } from '@rjsf/utils';
import { IconButton, Input, InputLabel, Grid2, FormControl } from '@sistent/sistent';
import { iconMedium } from '../../../../css/icons.styles';

const WrapIfAdditionalTemplate = ({
  children,
  classNames,
  disabled,
  id,
  label,
  onRemoveProperty,
  onKeyRename,
  onKeyRenameBlur,
  readonly,
  required,
  schema,
}) => {
  const value = label.startsWith('newKey') ? '' : label;
  const keyLabel = `Key`;
  const additional = ADDITIONAL_PROPERTY_FLAG in schema;

  if (!additional) {
    return <div className={classNames}>{children}</div>;
  }
  const handleBlur = onKeyRenameBlur;
  const handleChange = ({ target }) => onKeyRename(target.value);

  return (
    <Grid2
      container
      key={`${id}-key`}
      alignItems="center"
      spacing={2}
      className={classNames}
      size="grow"
    >
      <Grid2 size={{ xs: 12 }}>
        <FormControl fullWidth={true} required={required}>
          <InputLabel>{keyLabel}</InputLabel>
          <Input
            autoFocus={true}
            defaultValue={value}
            disabled={disabled || readonly}
            id={`${id}-key`}
            name={`${id}-key`}
            onChange={!readonly ? handleChange : undefined}
            onBlur={!readonly ? handleBlur : undefined}
            type="text"
          />
        </FormControl>
      </Grid2>
      <Grid2 style={{ alignSelf: 'flex-end' }}>
        <span style={{ fontSize: '1.25rem' }}>&nbsp;:&nbsp;</span>
      </Grid2>
      <Grid2 size={{ xs: 12 }}>{children}</Grid2>
      <Grid2>
        <IconButton component="span" disabled={disabled || readonly} onClick={onRemoveProperty}>
          <DeleteIcon style={iconMedium} />
        </IconButton>
      </Grid2>
    </Grid2>
  );
};

export default WrapIfAdditionalTemplate;
