import React from 'react';
import Grid from '@material-ui/core/Grid';
import { canExpand } from '@rjsf/utils';
import { CssBaseline,useTheme, withStyles } from '@material-ui/core';
import AddIcon from '../../../../assets/icons/AddIcon';
import { Box, IconButton, Typography } from '@material-ui/core';
import { CustomTextTooltip } from '../CustomTextTooltip';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import ExpandMoreIcon from '../../../../assets/icons/ExpandMoreIcon';
import ExpandLessIcon from '../../../../assets/icons/ExpandLessIcon'
import ErrorOutlineIcon from '../../../../assets/icons/ErrorOutlineIcon';
import { ERROR_COLOR } from '../../../../constants/colors';
import { iconMedium, iconSmall } from '../../../../css/icons.styles';
import { calculateGrid, getHyperLinkDiv } from '../helper';

const styles = (theme) => ({
  objectFieldGrid : {
    // paddingLeft: "0.6rem",
    padding : ".5rem",
    paddingTop : "0.7rem",
    // margin : ".5rem",
    backgroundColor : theme.palette.type === 'dark' ? "#303030" : "#f4f4f4",
    border : `1px solid  ${theme.palette.type === 'dark' ? "rgba(255, 255, 255, .45)" : "rgba(0, 0, 0, .125)"}`,
    width : "100%",
    margin : "0px"
  },
  typography : {
    fontFamily : "inherit",
    fontSize : 13,
  },

});
/**
 * Get the raw errors from the error schema.
 * @param {Object} errorSchema error schema.
 * @returns {Array} raw errors.
*/

const getRawErrors = (errorSchema) => {
  if (!errorSchema) return [];
  const errors = [];
  Object.keys(errorSchema).forEach((key) => {
    if (errorSchema[key].__errors) {
      errors.push(...errorSchema[key].__errors);
    }
  });
  return errors;
};

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
  classes,
  errorSchema
}) => {
  const additional = schema?.__additional_property; // check if the object is additional
  const theme = useTheme();
  const rawErrors = getRawErrors(errorSchema)

  // If the parent type is an `array`, then expand the current object.
  const [show, setShow] = React.useState(false);
  properties.forEach((property, index) => {
    if (schema.properties[property.name].type) {
      properties[index].type = schema.properties[property.name].type;
      properties[index].__additional_property =
        schema.properties[property.name]?.__additional_property || false;
    }
  });
  const CustomTitleField = ({ title, id, description, properties }) => {
    return <Box mb={1} mt={1} id={id} >
      <CssBaseline />
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
                {show ? <ExpandLessIcon style={iconMedium} fill="gray" /> : <ExpandMoreIcon style={iconMedium} fill="gray"  />}
              </IconButton>
            </Grid>
          )
        )}

        <Grid item mb={1} mt={1}>
          <Typography variant="body1" className={classes.typography} style={{ fontWeight : "bold", display : "inline" }}>{title.charAt(0).toUpperCase() + title.slice(1)}{" "}
          </Typography>
          {description &&
            <CustomTextTooltip backgroundColor="#3C494F" title={getHyperLinkDiv(description)}>
              <IconButton disableTouchRipple="true" disableRipple="true" component="span" size="small">
                <HelpOutlineIcon width="14px" height="14px"  fill={theme.palette.type === 'dark' ? "white" : "black"}   style={{ marginLeft : "4px", verticalAlign : "middle", ...iconSmall }}/>
              </IconButton>
            </CustomTextTooltip>}
          {rawErrors.length !==0 &&
            <CustomTextTooltip backgroundColor={ERROR_COLOR} title={rawErrors?.map((error, index) => (
              <div key={index}>{error}</div>
            ))}>
              <IconButton disableTouchRipple="true" disableRipple="true" component="span" size="small">
                <ErrorOutlineIcon width="14px" height="14px" fill="red" style={{ marginLeft : "4px", verticalAlign : "middle", ...iconSmall }} />
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
            {...calculateGrid(element)}
            key={index}
          >
            {element.content}
          </Grid>
        )
      );
    })}
  </Grid>
  )

  const fieldTitle = uiSchema['ui:title'] || title;

  return (
    <>
      {fieldTitle ? (
        <>
          <CustomTitleField
            id={`${idSchema.$id}-title`}
            title={additional ? "Value" : fieldTitle}
            description={description}
            properties={properties}
          />
          {Object.keys(properties).length > 0 && show && Properties}
        </>
      ) : Properties}
    </>
  );
};

export default withStyles(styles)(ObjectFieldTemplate);
