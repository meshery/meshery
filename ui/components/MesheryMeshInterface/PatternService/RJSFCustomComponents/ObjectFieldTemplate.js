import React from 'react';

import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/styles';

import { utils } from '@rjsf/core';

import AddButton from "@material-ui/icons/Add"
import { Box, IconButton, Typography } from '@material-ui/core';

const { canExpand } = utils;

const useStyles = makeStyles({
  root: {
    marginTop: 10,
  },
});

const ObjectFieldTemplate = ({
  DescriptionField,
  description,
  TitleField,
  title,
  properties,
  required,
  disabled,
  readonly,
  uiSchema,
  idSchema,
  schema,
  formData,
  onAddClick,
}) => {
  const classes = useStyles();
  console.log(schema)

const CustomTitleField = ({ title }) => (
  <>
    <Box mb={1} mt={1}>
      <Typography variant="h6">{title.charAt(0).toUpperCase() + title.slice(1)}</Typography>
    </Box>
  </>
);

  return (
    <>
      {(uiSchema['ui:title'] || title) && (
        <CustomTitleField
          id={`${idSchema.$id}-title`}
          title={title}
          required={required}
        />
      )}
      {description && (
        <DescriptionField
          id={`${idSchema.$id}-description`}
          description={description}
        />
      )}
      <Grid container={true} spacing={2} className={classes.root}>
        {properties.map((element, index) =>
          // Remove the <Grid> if the inner element is hidden as the <Grid>
          // itself would otherwise still take up space.
          element.hidden ? (
            element.content
          ) : (
            <Grid
              item={true}
              xs={12}
              key={index}
              style={{ marginBottom: "10px" }}>
              {element.content}
            </Grid>
          )
        )}
        {canExpand(schema, uiSchema, formData) && (
          <Grid container justify='flex-end'>
            <Grid item={true}>
              <IconButton
                className='object-property-expand'
                onClick={onAddClick(schema)}
                disabled={disabled || readonly}
              >
              <AddButton/>
              </IconButton>
            </Grid>
          </Grid>
        )}
      </Grid>
    </>
  );
};

export default ObjectFieldTemplate;
