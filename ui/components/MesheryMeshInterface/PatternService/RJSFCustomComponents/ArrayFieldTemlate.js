/* eslint-disable no-unused-vars */
import React from "react";

import { utils } from "@rjsf/core";

import Box from "@material-ui/core/Box";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";

import { Button, IconButton } from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add"
import SimpleAccordion from "./Accordion";

const { isMultiSelect, getDefaultRegistry } = utils;

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
  return <h3>{title}</h3>;
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
    paddingLeft : 6,
    paddingRight : 6,
    fontWeight : "bold",
    minWidth : 0
  };
  return (
    <SimpleAccordion childProps={props}>
      <Grid container={true} key={props.key} alignItems="center">
        <Grid item={true} xs style={{ overflow : "auto" }}>
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
      <ArrayFieldTitle
        key={`array-field-title-${props.idSchema.$id}`}
        TitleField={props.TitleField}
        idSchema={props.idSchema}
        title={props.uiSchema["ui:title"] || props.title}
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
        {props.items && props.items.map(DefaultArrayItem)}
      </div>

      {props.canAdd && (
        <Button
          className="array-item-add"
          onClick={props.onAddClick}
          disabled={props.disabled || props.readonly}
        >
          Add
        </Button>
      )}
    </fieldset>
  );
};

const DefaultNormalArrayFieldTemplate = (props) => {
  return (
    <Paper elevation={0}>
      <Box p={2}>
        <Grid item container alignItems="center" xs={12} justify="space-between" style={{ marginBottom : "0.3rem" }}>
          <Grid item xs={4}>
            <ArrayFieldTitle
              key={`array-field-title-${props.idSchema.$id}`}
              TitleField={props.TitleField}
              idSchema={props.idSchema}
              title={props.uiSchema["ui:title"] || props.title}
              required={props.required}
            />
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

        {(props.uiSchema["ui:description"] || props.schema.description) && (
          <ArrayFieldDescription
            key={`array-field-description-${props.idSchema.$id}`}
            DescriptionField={props.DescriptionField}
            idSchema={props.idSchema}
            description={
              props.uiSchema["ui:description"] || props.schema.description
            }
          />
        )}

        <Grid container={true} key={`array-item-list-${props.idSchema.$id}`}>
          {props.items && props.items.map((p) => DefaultArrayItem(p))}


        </Grid>
      </Box>
    </Paper>
  );
};

export default ArrayFieldTemplate;
