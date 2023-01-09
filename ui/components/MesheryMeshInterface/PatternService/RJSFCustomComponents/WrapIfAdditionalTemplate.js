import React from "react";
import FormControl from "@material-ui/core/FormControl";
import Grid from "@material-ui/core/Grid";
import DeleteIcon from "@material-ui/icons/Delete";
import InputLabel from "@material-ui/core/InputLabel";
import Input from "@material-ui/core/Input";
import {
  ADDITIONAL_PROPERTY_FLAG,
} from "@rjsf/utils";
import { IconButton } from "@material-ui/core";

const WrapIfAdditionalTemplate = ({
  children,
  classNames,
  disabled,
  id,
  label,
  onDropPropertyClick,
  onKeyChange,
  readonly,
  required,
  schema,
}) => {
  const value=label.startsWith("newKey")?"":label; // removing the default value i.e newKey.
  const keyLabel = `Key`
  const additional = ADDITIONAL_PROPERTY_FLAG in schema;

  if (!additional) {
    return <div className={classNames}>{children}</div>;
  }
  const handleChange = ({ target }) =>
    onKeyChange(target.value);

  return (
    <Grid
      container
      key={`${id}-key`}
      alignItems="center"
      spacing={2}
      className={classNames}
    >
      <Grid item xs>
        <FormControl fullWidth={true} required={required}>
          <InputLabel>{keyLabel}</InputLabel>
          <Input
            autoFocus={true}
            defaultValue={value}
            disabled={disabled || readonly}
            id={`${id}-key`}
            name={`${id}-key`}
            onChange={!readonly ? handleChange : undefined}
            type="text"
          />
        </FormControl>
      </Grid>
      <Grid item style={{ alignSelf : 'flex-end' }}>
        <span style={{ fontSize : "1.25rem", color : "#1E2117" }}>
          &nbsp;:&nbsp;
        </span>
      </Grid>
      <Grid item={true} xs>
        {children}
      </Grid>
      <Grid item={true}>
        <IconButton
          component="span"
          disabled={disabled || readonly}
          onClick={onDropPropertyClick(label)}
        >
          <DeleteIcon/>
        </IconButton>
      </Grid>
    </Grid>
  );
};

export default WrapIfAdditionalTemplate;