import React from 'react';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/styles';
import { canExpand } from '@rjsf/utils';
import AddIcon from '../../../../assets/icons/AddIcon';
import { Box, IconButton, Typography } from '@material-ui/core';
import { CustomTextTooltip } from '../CustomTextTooltip';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import ExpandMoreIcon from '../../../../assets/icons/ExpandMoreIcon';
import ExpandLessIcon from '../../../../assets/icons/ExpandLessIcon'
import ErrorOutlineIcon from '../../../../assets/icons/ErrorOutlineIcon';
import { ERROR_COLOR } from '../../../../constants/colors';

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
  rawErrors
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
              <AddIcon width="18px" height="18px" fill="white" style={{ backgroundColor : "#647881", width : "1.25rem", height : "1.25rem", color : "#ffffff", borderRadius : ".2rem" }} />
            </IconButton>
          </Grid>
        ) : (
          Object.keys(properties).length > 0 && (
            <Grid item={true}>
              <IconButton
                className="object-property-expand"
                onClick={() => setShow(!show)}
              >
                {show ? <ExpandLessIcon width="18px" height="18px" fill="gray" /> : <ExpandMoreIcon width="18px" height="18px" fill="gray"  />}
              </IconButton>
            </Grid>
          )
        )}

        <Grid item mb={1} mt={1}>
          <Typography variant="body1" style={{ fontWeight : "bold", display : "inline" }}>{title.charAt(0).toUpperCase() + title.slice(1)}{" "}
          </Typography>
          {description &&
            <CustomTextTooltip backgroundColor="#3C494F" title={description}>
              <IconButton disableTouchRipple="true" disableRipple="true" component="span" size="small">
                <HelpOutlineIcon width="14px" height="14px" fill="black" style={{ marginLeft : "4px", verticalAlign : "middle" }}/>
              </IconButton>
            </CustomTextTooltip>}
          {rawErrors?.length &&
            <CustomTextTooltip backgroundColor={ERROR_COLOR} title={rawErrors?.map((error, index) => (
              <div key={index}>{error}</div>
            ))}>
              <IconButton disableTouchRipple="true" disableRipple="true" component="span" size="small">
                <ErrorOutlineIcon width="14px" height="14px" fill="red" style={{ marginLeft : "4px", verticalAlign : "middle" }} />
              </IconButton>
            </CustomTextTooltip>}
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
