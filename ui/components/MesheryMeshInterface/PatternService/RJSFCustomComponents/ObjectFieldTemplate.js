import React from 'react';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/styles';
import { canExpand } from '@rjsf/utils';
import AddButton from "@material-ui/icons/Add";
import { Box, IconButton, Typography } from '@material-ui/core';
import { EnlargedTextTooltip } from '../EnlargedTextTooltip';
import HelpOutlineIcon from '../HelpOutlineIcon';
import ArrowDown from '@material-ui/icons/KeyboardArrowDown';
import ArrowUp from '@material-ui/icons/KeyboardArrowUp';

const useStyles = makeStyles({
  objectFieldGrid : {
    // paddingLeft: "0.6rem",
    padding : ".5rem",
    // margin : ".5rem",
    backgroundColor : "#f4f4f4",
    border : '1px solid rgba(0, 0, 0, .125)',
  },
});

const ObjectFieldTemplate = ({
  description,
  title,
  properties,
  // required,
  disabled,
  readonly,
  uiSchema,
  idSchema,
  schema,
  formData,
  onAddClick,
}) => {
  const additional = schema?.__additional_property; // check if the object is additional
  const classes = useStyles();
  // If the parent type is an `array`, then expand the current object.
  const [show, setShow] = React.useState(schema?.p_type ? true : false);
  properties.forEach((property, index) => {
    if (schema.properties[property.name].type) {
      properties[index].type = schema.properties[property.name].type;
      properties[index].__additional_property =
        schema.properties[property.name]?.__additional_property || false;
    }
  });
  const CustomTitleField = ({ title, id, description, properties }) => {
    return <Box mb={1} mt={1} id={id} >
      <Grid container justify="flex-start" alignItems="center">
        {canExpand(schema, uiSchema, formData) ? (
          <Grid item={true} onClick={() => {
            if (!show) setShow(true);
          }}>
            <IconButton
              className="object-property-expand"
              onClick={onAddClick(schema)}
              disabled={disabled || readonly}
            >
              <AddButton style={{ backgroundColor : "#647881", width : "1.25rem", height : "1.25rem", color : "#ffffff", borderRadius : ".2rem" }} />
            </IconButton>
          </Grid>
        ) : (
          Object.keys(properties).length > 0 && (
            <Grid item={true}>
              <IconButton
                className="object-property-expand"
                onClick={() => setShow(!show)}
              >
                {show ? <ArrowUp /> : <ArrowDown />}
              </IconButton>
            </Grid>
          )
        )}

        <Grid item mb={1} mt={1}>
          <Typography variant="body1" style={{ fontWeight : "bold", display : "inline" }}>{title.charAt(0).toUpperCase() + title.slice(1)}{" "}
          </Typography>
          {description &&
            <EnlargedTextTooltip title={description}>
              <HelpOutlineIcon />
            </EnlargedTextTooltip>}
        </Grid>


      </Grid>
    </Box>
  };

  const Properties = (<Grid container={true} spacing={2} className={classes.objectFieldGrid} style={Object.keys(properties).length === 0 || schema["$schema"] ? { border : "none" } : null}>
    {properties.map((element, index) => {
      return (
        element.hidden ? (
          element.content
        ) : (
          <Grid
            item={true}
            sm={12}
            lg={
              element.type === "object" ||
              element.type === "array" ||
              element.__additional_property ||
              additional
                ? 12
                : 6
            }
            key={index}
          >
            {element.content}
          </Grid>
        )
      );
    })}
  </Grid>)

  const fieldTitle = uiSchema['ui:title'] || title;

  return (
    <>
      {fieldTitle ? (
        <>
          {schema.p_type !== "array" ? (
            <CustomTitleField
              id={`${idSchema.$id}-title`}
              title={additional ? "Value" : fieldTitle}
              description={description}
              properties={properties}
            />
          ) : null
          }
          {Object.keys(properties).length > 0 && show && Properties}
        </>
      ) : Properties}
    </>
  );
};

export default ObjectFieldTemplate;
