import React from 'react';
import { Box, Grid2, Paper, Button, IconButton, Typography, useTheme } from '@sistent/sistent';
import AddIcon from '@mui/icons-material/Add';
import { CustomTextTooltip } from '../CustomTextTooltip';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import { isMultiSelect, getDefaultFormState } from '@rjsf/utils';
import ErrorOutlineIcon from '../../../../assets/icons/ErrorOutlineIcon';
import { ERROR_COLOR } from '../../../../constants/colors';
import { iconSmall } from '../../../../css/icons.styles';
import { safeDisplayValue, safeStringTitle } from '../helper';

function getTitle(props) {
  if (!props) return 'Unknown';
  return safeStringTitle(props.uiSchema['ui:title'] ?? props.title) || 'Unknown';
}

const ArrayFieldTemplate = (props) => {
  const { schema, registry = getDefaultFormState(), classes } = props;
  const safeId = props.fieldPathId?.$id ?? props.idSchema?.$id ?? 'array-field';
  const safeProps = { ...props, fieldPathId: props.fieldPathId ?? { $id: safeId, path: [] } };
  if (isMultiSelect(schema, registry.rootSchema)) {
    return <DefaultFixedArrayFieldTemplate {...safeProps} />;
  } else {
    return <DefaultNormalArrayFieldTemplate {...safeProps} />;
  }
};

const ArrayFieldTitle = ({ title, classes }) => {
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

const DefaultFixedArrayFieldTemplate = (props) => {
  const { classes } = props;
  const safeId = props.fieldPathId?.$id ?? props.idSchema?.$id ?? 'array-field';

  return (
    <fieldset className={props.className}>
      {props.canAdd && (
        <Button
          className="rjsf-array-item-add"
          onClick={typeof props.onAddClick === 'function' ? props.onAddClick : undefined}
          disabled={props.disabled || props.readonly}
        >
          Add
        </Button>
      )}

      <ArrayFieldTitle
        key={`array-field-title-${safeId}`}
        title={getTitle(props)}
        required={props.required}
        classes={classes}
      />

      {(props.uiSchema['ui:description'] ?? props.schema?.description) != null && (
        <div className="field-description" key={`field-description-${safeId}`}>
          {safeDisplayValue(props.uiSchema['ui:description'] ?? props.schema?.description)}
        </div>
      )}

      <div className="row rjsf-array-item-list" key={`array-item-list-${safeId}`}>
        {props.items}
      </div>
    </fieldset>
  );
};

const DefaultNormalArrayFieldTemplate = (props) => {
  const theme = useTheme();
  const { classes } = props;
  const safeId = props.fieldPathId?.$id ?? props.idSchema?.$id ?? 'array-field';

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
            <ArrayFieldTitle
              key={`array-field-title-${safeId}`}
              title={getTitle(props)}
              required={props.required}
              classes={classes}
            />

            {(props.uiSchema['ui:description'] ?? props.schema?.description) != null && (
              <CustomTextTooltip
                title={safeStringTitle(
                  props.uiSchema['ui:description'] ?? props.schema?.description,
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
            {props.rawErrors?.length > 0 && (
              <CustomTextTooltip
                bgColor={ERROR_COLOR}
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
                <Grid2>
                  <Box mt={2}>
                    <IconButton
                      className="rjsf-array-item-add"
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
          {props.items}
        </Grid2>
      </Box>
    </Paper>
  );
};

export default ArrayFieldTemplate;
