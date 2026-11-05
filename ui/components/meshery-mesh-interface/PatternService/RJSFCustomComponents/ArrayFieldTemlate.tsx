import React from 'react';
import { Box, Grid2, Paper, Button, IconButton, Typography, useTheme } from '@sistent/sistent';
import AddIcon from '../../../../assets/icons/AddIcon';
import SimpleAccordion from './Accordion';
import { CustomTextTooltip } from '../CustomTextTooltip';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import { isFixedItems, type ArrayFieldTemplateProps } from '@rjsf/utils';
import ErrorOutlineIcon from '../../../../assets/icons/ErrorOutlineIcon';
import { iconSmall } from '../../../../css/icons.styles';
import pluralize from 'pluralize';
import { safeDisplayValue, safeStringTitle } from '../helper';

function getTitleForItem(props: ArrayFieldTemplateProps) {
  return pluralize.singular(getTitle(props));
}

function getTitle(props: ArrayFieldTemplateProps): string {
  return safeStringTitle(props.uiSchema?.['ui:title'] ?? props.title) || 'Unknown';
}

const ArrayFieldTemplate = (props: ArrayFieldTemplateProps) => {
  // RJSF's own ArrayField routes a multi-select array to its multi-select widget
  // before ever calling this template, so `isMultiSelect` would always be false
  // here. A tuple `items` is the only shape that reaches the template as a fixed
  // array, and `isFixedItems` is the very test RJSF's ArrayField uses to pick it.
  // `registry` is required by RJSFBaseProps, so there is no fallback to default
  // and no cast to make — and `getDefaultRegistry()` could not serve as one
  // anyway: it lives in @rjsf/core and omits `schemaUtils`.
  if (isFixedItems(props.schema)) {
    return <DefaultFixedArrayFieldTemplate {...props} />;
  } else {
    return <DefaultNormalArrayFieldTemplate {...props} />;
  }
};

const ArrayFieldTitle = ({ title }: { title?: string }) => {
  const safeTitle = safeStringTitle(title);
  if (!safeTitle) return null;

  return (
    <Typography
      variant="body1"
      style={{ fontWeight: 'bold', display: 'inline', fontSize: '0.8rem' }}
    >
      {safeTitle.charAt(0).toUpperCase() + safeTitle.slice(1)}
    </Typography>
  );
};

// Used in the two templates
// RJSF hands each entry of `items` as a ready-to-render element for one array
// item (keyed by a stable item key), so the item is rendered as-is and only
// wrapped in the Meshery accordion chrome; its fields and its own toolbar come
// from rendering the element. Spreading the element into props used to pass
// `{ key, type, props }` instead of the item's props, so its content and
// toolbar were lost.
const DefaultArrayItem = ({ heading, item }: { heading: string; item: React.ReactElement }) => {
  return (
    <SimpleAccordion heading={heading} childProps={item.props}>
      <Grid2 container={true} alignItems="center" size="grow">
        <Grid2 size={{ xs: 12 }}>
          <Box mb={2} style={{ border: '0.5px solid black' }}>
            <Paper elevation={0}>
              <Box p={2}>{item}</Box>
            </Paper>
          </Box>
        </Grid2>
      </Grid2>
    </SimpleAccordion>
  );
};

const DefaultFixedArrayFieldTemplate = (props: ArrayFieldTemplateProps) => {
  const safeId = props.fieldPathId?.$id ?? 'array-field';

  return (
    <fieldset className={props.className}>
      {props.canAdd && (
        <Button
          className="array-item-add"
          onClick={typeof props.onAddClick === 'function' ? props.onAddClick : undefined}
          disabled={props.disabled || props.readonly}
        >
          Add
        </Button>
      )}

      <ArrayFieldTitle title={getTitle(props)} />

      {(props.uiSchema?.['ui:description'] ?? props.schema?.description) != null && (
        <div className="field-description" key={`field-description-${safeId}`}>
          {safeDisplayValue(props.uiSchema?.['ui:description'] ?? props.schema?.description)}
        </div>
      )}

      <div className="row array-item-list" key={`array-item-list-${safeId}`}>
        {props.items &&
          props.items.map((item, idx) => {
            return (
              <DefaultArrayItem
                key={item.key ?? `${getTitle(props)}-${idx}`}
                heading={`${getTitleForItem(props)} (${idx})`}
                item={item}
              />
            );
          })}
      </div>
    </fieldset>
  );
};

const DefaultNormalArrayFieldTemplate = (props: ArrayFieldTemplateProps) => {
  const theme = useTheme();
  const safeId = props.fieldPathId?.$id ?? 'array-field';

  return (
    <Paper elevation={0}>
      <Box p={1}>
        <Grid2
          alignItems="center"
          justify="space-between"
          style={{ marginBottom: '0.3rem' }}
          size="grow"
        >
          <Grid2 size={{ xs: 4 }}>
            <ArrayFieldTitle title={getTitle(props)} />

            {(props.uiSchema?.['ui:description'] ?? props.schema?.description) != null && (
              <CustomTextTooltip
                title={safeStringTitle(
                  props.uiSchema?.['ui:description'] ?? props.schema?.description,
                )}
              >
                <IconButton disableTouchRipple="true" disableRipple="true">
                  <HelpOutlineIcon
                    width="14px"
                    height="14px"
                    fill={theme.palette.mode === 'dark' ? 'white' : 'gray'}
                    style={{ marginLeft: '4px', ...iconSmall }}
                  />
                </IconButton>
              </CustomTextTooltip>
            )}
            {(props.rawErrors?.length ?? 0) > 0 && (
              <CustomTextTooltip
                bgColor={theme.palette.error.main}
                interactive={true}
                title={safeStringTitle(
                  Array.isArray(props.rawErrors) ? props.rawErrors.join('  ') : props.rawErrors,
                )}
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
                    fill={theme.palette.mode === 'dark' ? '#F91313' : '#B32700'}
                    style={{ marginLeft: '4px', verticalAlign: 'middle', ...iconSmall }}
                  />
                </IconButton>
              </CustomTextTooltip>
            )}
          </Grid2>
          <Grid2 size={{ xs: 4 }}>
            {props.canAdd && (
              <Grid2 container justify="flex-end">
                <Grid2 item={true}>
                  <Box mt={2}>
                    <IconButton
                      className="array-item-add"
                      onClick={
                        typeof props.onAddClick === 'function' ? props.onAddClick : undefined
                      }
                      disabled={props.disabled || props.readonly}
                    >
                      <AddIcon width="18px" height="18px" fill="gray" />
                    </IconButton>
                  </Box>
                </Grid2>
              </Grid2>
            )}
          </Grid2>
        </Grid2>

        <Grid2 container={true} key={`array-item-list-${safeId}`} size={'grow'}>
          {props.items &&
            props.items.map((item, idx) => {
              return (
                <DefaultArrayItem
                  key={item.key ?? `${getTitle(props)}-${idx}`}
                  heading={`${getTitleForItem(props)} (${idx})`}
                  item={item}
                />
              );
            })}
        </Grid2>
      </Box>
    </Paper>
  );
};

export default ArrayFieldTemplate;
