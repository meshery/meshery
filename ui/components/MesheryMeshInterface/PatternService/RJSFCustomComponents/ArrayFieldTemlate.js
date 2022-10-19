/* eslint-disable no-unused-vars */
import React, { useEffect } from "react";

import { utils } from "@rjsf/core";

import Box from "@material-ui/core/Box";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";

import { Button, IconButton, Typography } from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import SimpleAccordion from "./Accordion";
import EnlargedTextTooltip from "../EnlargedTextTooltip";
import HelpOutlineIcon from "../HelpOutlineIcon";
const { isMultiSelect, getDefaultRegistry } = utils;

function getTitleForItem(props) {
  const title = getTitle(props);

  // remove plurals
  if (title.endsWith("es")) {
    return title.substring(0, title.length - 2);
  }
  if (title.endsWith("s")) {
    return title.substring(0, title.length - 1);
  }

  return title;
}

function getTitle(props) {
  if (!props) {
    return "Unknown"
  }
  return props.uiSchema["ui:title"] || props.title
}

const ArrayFieldTemplate = (props) => {
  const { schema, registry = getDefaultRegistry() } = props;

  // TODO: update types so we don't have to cast registry as any
  if (isMultiSelect(schema, registry.rootSchema)) {
    return <DefaultFixedArrayFieldTemplate {...props} />;
  } else {
    return <DefaultNormalArrayFieldTemplate {...props} />;
  }
};

const ArrayFieldTitle = ({ TitleField, idSchema, title, required }) => {
  if (!title) {
    return null;
  }

  const id = `${idSchema.$id}__title`;
  // return <h3>{title?.charAt(0)?.toUpperCase() + title?.slice(1)}</h3>;
  return <Typography variant="body1" style={{ fontWeight : "bold", marginLeft : ".5rem", display : "inline" }}>{title.charAt(0).toUpperCase() + title.slice(1)}</Typography>;
  // return <TitleField id={id} title={title} required={required} />;
};

const ArrayFieldDescription = ({ DescriptionField, idSchema, description }) => {
  if (!description) {
    return null;
  }

  const id = `${idSchema.$id}__description`;
  return <DescriptionField id={id} description={description} />;
};

// Used in the two templates
const DefaultArrayItem = (props) => {
  const btnStyle = {
    flex : 1,
    paddingLeft : 0,
    paddingRight : 6,
    fontWeight : "bold",
  };

  return (
    <SimpleAccordion heading={props.heading} childProps={props}>
      <Grid container={true} key={props.key} alignItems="center">
        <Grid item={true} xs >
          <Box mb={2} style={{ border : "0.5px solid black" }}>
            <Paper elevation={0}>
              <Box p={2}>{props.children}</Box>
            </Paper>
          </Box>
        </Grid>

        {props.hasToolbar && (
          <Grid item={true}>
            {(props.hasMoveUp || props.hasMoveDown) && (
              <IconButton
                icon="arrow-up"
                className="array-item-move-up"
                tabIndex={-1}
                style={btnStyle}
                iconProps={{ fontSize : "small" }}
                disabled={props.disabled || props.readonly || !props.hasMoveUp}
                onClick={props.onReorderClick(props.index, props.index - 1)}
              />
            )}

            {(props.hasMoveUp || props.hasMoveDown) && (
              <IconButton
                icon="arrow-down"
                tabIndex={-1}
                style={btnStyle}
                iconProps={{ fontSize : "small" }}
                disabled={
                  props.disabled || props.readonly || !props.hasMoveDown
                }
                onClick={props.onReorderClick(props.index, props.index + 1)}
              />
            )}
          </Grid>
        )}
      </Grid>
    </SimpleAccordion>
  );
};

const DefaultFixedArrayFieldTemplate = (props) => {
  return (
    <fieldset className={props.className}>
      {props.canAdd && (
        <Button
          className="array-item-add"
          onClick={props.onAddClick}
          disabled={props.disabled || props.readonly}
        >
          Add
        </Button>
      )}

      <ArrayFieldTitle
        key={`array-field-title-${props.idSchema.$id}`}
        TitleField={props.TitleField}
        idSchema={props.idSchema}
        title={getTitle(props)}
        required={props.required}
      />

      {(props.uiSchema["ui:description"] || props.schema.description) && (
        <div
          className="field-description"
          key={`field-description-${props.idSchema.$id}`}
        >
          {props.uiSchema["ui:description"] || props.schema.description}
        </div>
      )}

      <div
        className="row array-item-list"
        key={`array-item-list-${props.idSchema.$id}`}
      >
        {props.items && props.items.map((item, idx) => {
          return <DefaultArrayItem key={`${getTitle(props)}-${idx}`} heading={`${getTitleForItem(props)} (${idx})`} {...item} />
        })}
      </div>


    </fieldset>
  );
};

const DefaultNormalArrayFieldTemplate = (props) => {
  return (
    <Paper elevation={0}>
      <Box p={1}>
        <Grid item container alignItems="center" xs={12} justify="space-between" style={{ marginBottom : "0.3rem" }}>
          <Grid item xs={4}>
            <ArrayFieldTitle
              key={`array-field-title-${props.idSchema.$id}`}
              TitleField={props.TitleField}
              idSchema={props.idSchema}
              title={getTitle(props)}
              required={props.required}
            />

            {
              (props.uiSchema["ui:description"] || props.schema.description) &&
              <EnlargedTextTooltip title={props.uiSchema["ui:description"] || props.schema.description}>
                <HelpOutlineIcon style={{ marginLeft : '4px' }} />
              </EnlargedTextTooltip>
            }

          </Grid>
          <Grid item xs={4}>
            {props.canAdd && (
              <Grid container justify="flex-end">
                <Grid item={true}>
                  <Box mt={2}>
                    <IconButton
                      className="array-item-add"
                      onClick={props.onAddClick}
                      disabled={props.disabled || props.readonly}
                    >
                      <AddIcon />
                    </IconButton>
                  </Box>
                </Grid>
              </Grid>
            )}
          </Grid>
        </Grid>

        {/* {(props.uiSchema["ui:description"] || props.schema.description) && (
          <ArrayFieldDescription
            key={`array-field-description-${props.idSchema.$id}`}
            DescriptionField={CustomDescriptionField}
            idSchema={props.idSchema}
            description={
              props.uiSchema["ui:description"] || props.schema.description
            }
          />
        )} */}

        <Grid container={true} key={`array-item-list-${props.idSchema.$id}`}>
          {props.items && props.items.map((item, idx) => {
            return <DefaultArrayItem key={`${getTitle(props)}-${idx}`} heading={`${getTitleForItem(props)} (${idx})`} {...item} />
          })}


        </Grid>
      </Box>
    </Paper>
  );
};

export default ArrayFieldTemplate;
