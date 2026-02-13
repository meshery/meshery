/* eslint-disable no-unused-vars */
import React from 'react';
import { Box, Grid2, Paper, Button, IconButton, Typography, useTheme } from '@sistent/sistent';
import AddIcon from '@mui/icons-material/Add';
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward';
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';
import SimpleAccordion from './Accordion';
import { CustomTextTooltip } from '../CustomTextTooltip';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import { isMultiSelect } from '@rjsf/utils';
import ErrorOutlineIcon from '../../../../assets/icons/ErrorOutlineIcon';
import { ERROR_COLOR } from '../../../../constants/colors';
import { iconSmall } from '../../../../css/icons.styles';
import pluralize from 'pluralize';

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
  const { schema, registry, classes } = props;
  // TODO: update types so we don't have to cast registry as any
  if (isMultiSelect(schema, registry?.rootSchema)) {
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
      variant="body1"
      style={{ fontWeight: 'bold', display: 'inline', fontSize: '0.8rem' }}
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
      <Grid2 container={true} key={props.key} alignItems="center" size="grow">
        <Grid2 size={{ xs: 12 }}>
          <Box mb={2} style={{ border: '0.5px solid black' }}>
            <Paper elevation={0}>
              <Box p={2}>{props.children}</Box>
            </Paper>
          </Box>
        </Grid2>

        {props.hasToolbar && (
          <Grid2>
            {(props.hasMoveUp || props.hasMoveDown) && (
              <IconButton
                className="array-item-move-up"
                tabIndex={-1}
                style={btnStyle}
                disabled={props.disabled || props.readonly || !props.hasMoveUp}
                onClick={props.onReorderClick(props.index, props.index - 1)}
              >
                <ArrowUpwardIcon fontSize="small" />
              </IconButton>
            )}

            {(props.hasMoveUp || props.hasMoveDown) && (
              <IconButton
                tabIndex={-1}
                style={btnStyle}
                disabled={props.disabled || props.readonly || !props.hasMoveDown}
                onClick={props.onReorderClick(props.index, props.index + 1)}
              >
                <ArrowDownwardIcon fontSize="small" />
              </IconButton>
            )}
          </Grid2>
        )}
      </Grid2>
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
        title={getTitle(props)}
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
  const theme = useTheme();
  const { classes } = props;
  return (
    <Paper elevation={0}>
      <Box p={1}>
        <Grid2
          alignItems="center"
          justifyContent="space-between"
          style={{ marginBottom: '0.3rem' }}
          size="grow"
        >
          <Grid2 size={{ xs: 4 }}>
            <ArrayFieldTitle
              key={`array-field-title-${props.idSchema.$id}`}
              title={getTitle(props)}
              classes={classes}
            />

            {(props.uiSchema['ui:description'] || props.schema.description) && (
              <CustomTextTooltip title={props.schema.description}>
                <IconButton disableTouchRipple={true} disableRipple={true}>
                  <HelpOutlineIcon
                    width="14px"
                    height="14px"
                    fill={theme.palette.mode === 'dark' ? 'white' : 'gray'}
                    style={{ marginLeft: '4px', ...iconSmall }}
                  />
                </IconButton>
              </CustomTextTooltip>
            )}
            {props.rawErrors?.length > 0 && (
              <CustomTextTooltip
                bgColor={ERROR_COLOR}
                interactive={true}
                title={props.rawErrors?.join('  ')}
              >
                <IconButton
                  component="span"
                  size="small"
                  disableTouchRipple={true}
                  disableRipple={true}
                >
                  <ErrorOutlineIcon
                    width="16px"
                    height="16px"
                    fill={theme.palette.mode === 'dark' ? '#F91313' : '#B32700'}
                    style={{ marginLeft: '4px', verticalAlign: 'middle', ...iconSmall }}
                  />
                </IconButton>
              </CustomTextTooltip>
            )}
          </Grid2>
          <Grid2 size={{ xs: 4 }}>
            {props.canAdd && (
              <Grid2 container justifyContent="flex-end">
                <Grid2>
                  <Box mt={2}>
                    <IconButton
                      className="array-item-add"
                      onClick={props.onAddClick}
                      disabled={props.disabled || props.readonly}
                    >
                      <AddIcon />
                    </IconButton>
                  </Box>
                </Grid2>
              </Grid2>
            )}
          </Grid2>
        </Grid2>

        <Grid2 container={true} key={`array-item-list-${props.idSchema.$id}`} size={'grow'}>
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
        </Grid2>
      </Box>
    </Paper>
  );
};

export default ArrayFieldTemplate;
