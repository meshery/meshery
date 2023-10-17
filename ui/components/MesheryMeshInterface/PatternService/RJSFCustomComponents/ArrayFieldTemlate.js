/* eslint-disable no-unused-vars */
import React from 'react';
import Box from '@material-ui/core/Box';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { Button, IconButton, Typography, withStyles } from '@material-ui/core';
import AddIcon from '@material-ui/icons/Add';
import SimpleAccordion from './Accordion';
import { CustomTextTooltip } from '../CustomTextTooltip';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import { isMultiSelect, getDefaultFormState } from '@rjsf/utils';
import ErrorOutlineIcon from '../../../../assets/icons/ErrorOutlineIcon';
import { ERROR_COLOR } from '../../../../constants/colors';
import { iconSmall } from '../../../../css/icons.styles';
import { getHyperLinkDiv } from '../helper';
import pluralize from 'pluralize';
const styles = (theme) => ({
  typography: {
    fontSize: '0.8rem',
  },
  root: {
    '& .MuiPaper-root': {
      backgroundColor: '#f4f4f4',
    },
  },
});
function getTitleForItem(props) {
  const title = getTitle(props);

  return pluralize.singular(title);
}

function getTitle(props) {
  if (!props) {
    return 'Unknown';
  }
  return props.uiSchema['ui:title'] || props.title;
}

const ArrayFieldTemplate = (props) => {
  const { schema, registry = getDefaultFormState(), classes } = props;
  // TODO: update types so we don't have to cast registry as any
  if (isMultiSelect(schema, registry.rootSchema)) {
    return <DefaultFixedArrayFieldTemplate {...props} />;
  } else {
    return <DefaultNormalArrayFieldTemplate {...props} />;
  }
};

const ArrayFieldTitle = ({ title, classes }) => {
  if (!title) {
    return null;
  }

  return (
    <Typography
      className={classes.typography}
      variant="body1"
      style={{ fontWeight: 'bold', display: 'inline' }}
    >
      {title.charAt(0).toUpperCase() + title.slice(1)}
    </Typography>
  );
};

// Used in the two templates
const DefaultArrayItem = (props) => {
  const btnStyle = {
    flex: 1,
    paddingLeft: 0,
    paddingRight: 6,
    fontWeight: 'bold',
  };

  return (
    <SimpleAccordion heading={props.heading} childProps={props}>
      <Grid container={true} key={props.key} alignItems="center">
        <Grid item={true} xs>
          <Box mb={2} style={{ border: '0.5px solid black' }}>
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
                iconProps={{ fontSize: 'small' }}
                disabled={props.disabled || props.readonly || !props.hasMoveUp}
                onClick={props.onReorderClick(props.index, props.index - 1)}
              />
            )}

            {(props.hasMoveUp || props.hasMoveDown) && (
              <IconButton
                icon="arrow-down"
                tabIndex={-1}
                style={btnStyle}
                iconProps={{ fontSize: 'small' }}
                disabled={props.disabled || props.readonly || !props.hasMoveDown}
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
  const { classes } = props;
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
        classes={classes}
      />

      {(props.uiSchema['ui:description'] || props.schema.description) && (
        <div className="field-description" key={`field-description-${props.idSchema.$id}`}>
          {props.uiSchema['ui:description'] || props.schema.description}
        </div>
      )}

      <div className="row array-item-list" key={`array-item-list-${props.idSchema.$id}`}>
        {props.items &&
          props.items.map((item, idx) => {
            return (
              <DefaultArrayItem
                key={`${getTitle(props)}-${idx}`}
                heading={`${getTitleForItem(props)} (${idx})`}
                {...item}
              />
            );
          })}
      </div>
    </fieldset>
  );
};

const DefaultNormalArrayFieldTemplate = (props) => {
  const { classes, theme } = props;
  return (
    <Paper className={classes.root} elevation={0}>
      <Box p={1}>
        <Grid
          item
          container
          alignItems="center"
          xs={12}
          justify="space-between"
          style={{ marginBottom: '0.3rem' }}
        >
          <Grid item xs={4}>
            <ArrayFieldTitle
              key={`array-field-title-${props.idSchema.$id}`}
              TitleField={props.TitleField}
              idSchema={props.idSchema}
              title={getTitle(props)}
              required={props.required}
              classes={classes}
            />

            {(props.uiSchema['ui:description'] || props.schema.description) && (
              <CustomTextTooltip
                backgroundColor="#3C494F"
                title={getHyperLinkDiv(props.schema.description)}
              >
                <IconButton disableTouchRipple="true" disableRipple="true">
                  <HelpOutlineIcon
                    width="14px"
                    height="14px"
                    fill={theme.palette.type === 'dark' ? 'white' : 'gray'}
                    style={{ marginLeft: '4px', ...iconSmall }}
                  />
                </IconButton>
              </CustomTextTooltip>
            )}
            {props.rawErrors?.length > 0 && (
              <CustomTextTooltip
                backgroundColor={ERROR_COLOR}
                interactive={true}
                title={props.rawErrors?.map((error, index) => (
                  <div key={index}>{error}</div>
                ))}
              >
                <IconButton
                  component="span"
                  size="small"
                  disableTouchRipple="true"
                  disableRipple="true"
                >
                  <ErrorOutlineIcon
                    width="16px"
                    height="16px"
                    fill={theme.palette.type === 'dark' ? '#F91313' : '#B32700'}
                    style={{ marginLeft: '4px', verticalAlign: 'middle', ...iconSmall }}
                  />
                </IconButton>
              </CustomTextTooltip>
            )}
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
                      <AddIcon width="18px" height="18px" fill="gray" />
                    </IconButton>
                  </Box>
                </Grid>
              </Grid>
            )}
          </Grid>
        </Grid>

        <Grid container={true} key={`array-item-list-${props.idSchema.$id}`}>
          {props.items &&
            props.items.map((item, idx) => {
              return (
                <DefaultArrayItem
                  key={`${getTitle(props)}-${idx}`}
                  heading={`${getTitleForItem(props)} (${idx})`}
                  {...item}
                />
              );
            })}
        </Grid>
      </Box>
    </Paper>
  );
};

export default withStyles(styles, { withTheme: true })(ArrayFieldTemplate);
