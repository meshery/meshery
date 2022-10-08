import React from 'react';

import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/styles';

import { utils } from '@rjsf/core';

import AddButton from "@material-ui/icons/Add";
import { Box, IconButton, Typography } from '@material-ui/core';
import EnlargedTextTooltip from '../EnlargedTextTooltip';
import HelpOutlineIcon from '../HelpOutlineIcon';
import ArrowDown from '@material-ui/icons/KeyboardArrowDown';
import ArrowUp from '@material-ui/icons/KeyboardArrowUp';

const { canExpand } = utils;

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
  const classes = useStyles();
  const [show, setShow] = React.useState(false);

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
            xs={element.name === "name" || element.name === "namespace" ? 6 : 12}
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
          <CustomTitleField
            id={`${idSchema.$id}-title`}
            title={fieldTitle}
            description={description}
            properties={properties}
          />
          {Object.keys(properties).length > 0 && show && Properties}
        </>
      ) : Properties}
    </>
  );
};

export default ObjectFieldTemplate;
