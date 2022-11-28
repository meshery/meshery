import React from "react";
import FormControl from "@material-ui/core/FormControl";
import Grid from "@material-ui/core/Grid";
import InputLabel from "@material-ui/core/InputLabel";
import Input from "@material-ui/core/Input";
import {
  ADDITIONAL_PROPERTY_FLAG,
} from "@rjsf/utils";

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
  uiSchema,
  registry,
}) => {
  const { RemoveButton } = registry.templates.ButtonTemplates;
  const value=label.startsWith("newKey")?"":label; // removing the default value i.e newKey.
  const keyLabel = `Key`
  const additional = ADDITIONAL_PROPERTY_FLAG in schema;
  const btnStyle= {
    flex : 1,
    paddingLeft : 6,
    paddingRight : 6,
    fontWeight : "bold",
  };

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
        <RemoveButton
          iconType="default"
          style={btnStyle}
          disabled={disabled || readonly}
          onClick={onDropPropertyClick(label)}
          uiSchema={uiSchema}
        />
      </Grid>
    </Grid>
  );
};

export default WrapIfAdditionalTemplate;