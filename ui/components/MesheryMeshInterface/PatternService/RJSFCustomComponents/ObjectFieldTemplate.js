import React from 'react';

import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/styles';

import { utils } from '@rjsf/core';

import AddButton from "@material-ui/icons/Add";
import { Box, IconButton, Typography } from '@material-ui/core';
import EnlargedTextTooltip from '../EnlargedTextTooltip';
import HelpOutlineIcon from '../HelpOutlineIcon';

const { canExpand } = utils;

const useStyles = makeStyles({
  objectFieldGrid : {
    marginTop : 10,
    // paddingLeft: "0.6rem",
    padding : "0.6rem",
    border : '1px solid rgba(0, 0, 0, .125)',
  },
});

const ObjectFieldTemplate = ({
  description,
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

  const CustomTitleField = ({ title }) => (
    <Box mb={1} mt={1}>
      <Grid container justify="space-between" alignItems="center">
        <Grid item mb={1} mt={1}>
          <Typography variant="body1" style={{ fontWeight : "bold", display : "inline" }}>{title.charAt(0).toUpperCase() + title.slice(1)}{" "}</Typography>
          {description &&
            <EnlargedTextTooltip title={description}>
              <HelpOutlineIcon />
            </EnlargedTextTooltip>}
        </Grid>

        {canExpand(schema, uiSchema, formData) && (
          <Grid item={true}>
            <IconButton
              className="object-property-expand"
              onClick={onAddClick(schema)}
              disabled={disabled || readonly}
            >
              <AddButton />
            </IconButton>
          </Grid>
        )}
      </Grid>

    </Box>
  );

  return (
    <>
      {(uiSchema['ui:title'] || title) && (
        <CustomTitleField
          id={`${idSchema.$id}-title`}
          title={title}
          description={description}
          required={required}
        />
      )}

      {/* {description && (
        <CustomDescriptionField
          id={`${idSchema.$id}-description`}
          description={description}
        />
      )} */}

      <Grid container={true} spacing={2} className={classes.objectFieldGrid} style={Object.keys(properties).length === 0 || schema["$schema"] ? { border : "none" } : null}>
        {properties.map((element, index) => {
          // console.log("eke", element)
          // Remove the <Grid> if the inner element is hidden as the <Grid>
          // itself would otherwise still take up space.
          return (
            element.hidden ? (
              element.content
            ) : (
              <Grid
                item={true}
                xs={element.name === "name" || element.name === "namespace" ? 6 : 12}
                key={index}
              >
                {element.content}
              </Grid>
            )
          );
        })}
        {/** canExpand(schema, uiSchema, formData) && (
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
        ) */}
      </Grid>
    </>
  );
};

export default ObjectFieldTemplate;
