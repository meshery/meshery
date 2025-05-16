import React from 'react';
import { canExpand } from '@rjsf/utils';
import AddIcon from '../../../../assets/icons/AddIcon';
import { Grid, Box, IconButton, Typography, useTheme, CssBaseline } from '@layer5/sistent';
import { CustomTextTooltip } from '../CustomTextTooltip';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import ExpandMoreIcon from '../../../../assets/icons/ExpandMoreIcon';
import ExpandLessIcon from '../../../../assets/icons/ExpandLessIcon';
import ErrorOutlineIcon from '../../../../assets/icons/ErrorOutlineIcon';
import { ERROR_COLOR } from '../../../../constants/colors';
import { iconMedium, iconSmall } from '../../../../css/icons.styles';
import { calculateGrid } from '../helper';

/**
 * Get the raw errors from the error schema.
 * @param {Object} errorSchema error schema.
 * @returns {Array} raw errors.
 */

const getRawErrors = (errorSchema) => {
  if (!errorSchema) return [];
  const errors = [];
  Object.keys(errorSchema).forEach((key) => {
    if (errorSchema?.[key]?.__errors) {
      errors.push(...errorSchema[key].__errors);
    }
  });
  return errors;
};

const ObjectFieldTemplate = ({
  description,
  title,
  properties,
  disabled,
  readonly,
  uiSchema,
  idSchema,
  schema,
  formData,
  onAddClick,

  errorSchema,
}) => {
  const additional = schema?.__additional_property; // check if the object is additional
  const theme = useTheme();
  const rawErrors = getRawErrors(errorSchema);

  // If the parent type is an `array`, then expand the current object.
  const [show, setShow] = React.useState(uiSchema['ui:options']?.expand || false);
  properties.forEach((property, index) => {
    if (schema.properties[property.name].type) {
      properties[index].type = schema.properties[property.name].type;
      properties[index].__additional_property =
        schema.properties[property.name]?.__additional_property || false;
    }
  });
  const CustomTitleField = ({ title, id, description, properties }) => {
    return (
      <Box mb={1} mt={1} id={id}>
        <CssBaseline />
        <Grid container justify="flex-start" alignItems="center">
          {canExpand(schema, uiSchema, formData) ? (
            <Grid
              item={true}
              onClick={() => {
                if (!show) setShow(true);
              }}
            >
              <IconButton
                className="object-property-expand"
                onClick={onAddClick(schema)}
                disabled={disabled || readonly}
              >
                <AddIcon
                  style={{
                    backgroundColor: `${theme.palette.mode === 'dark' ? '#00b39F' : '#647881'}`,
                    width: '1rem',
                    height: '1rem',
                    color: '#ffffff',
                    borderRadius: '.2rem',
                  }}
                />
              </IconButton>
            </Grid>
          ) : (
            Object.keys(properties).length > 0 && (
              <Grid item={true}>
                <IconButton className="object-property-expand" onClick={() => setShow(!show)}>
                  {show ? (
                    <ExpandLessIcon style={iconMedium} fill="gray" />
                  ) : (
                    <ExpandMoreIcon style={iconMedium} fill="gray" />
                  )}
                </IconButton>
              </Grid>
            )
          )}

          <Grid item mb={1} mt={1}>
            <Typography
              variant="body1"
              style={{ fontWeight: 'bold', display: 'inline', fontFamily: 'inherit', fontSize: 13 }}
            >
              {title.charAt(0).toUpperCase() + title.slice(1)}{' '}
            </Typography>
            {description && (
              <CustomTextTooltip title={description}>
                <IconButton
                  disableTouchRipple="true"
                  disableRipple="true"
                  component="span"
                  size="small"
                >
                  <HelpOutlineIcon
                    width="1rem"
                    height="1rem"
                    fill={theme.palette.mode === 'dark' ? 'white' : 'black'}
                    style={{ marginLeft: '4px', verticalAlign: 'middle', ...iconSmall }}
                  />
                </IconButton>
              </CustomTextTooltip>
            )}
            {rawErrors.length !== 0 && (
              <CustomTextTooltip bgColor={ERROR_COLOR} title={rawErrors?.join('  ')}>
                <IconButton
                  disableTouchRipple="true"
                  disableRipple="true"
                  component="span"
                  size="small"
                >
                  <ErrorOutlineIcon
                    width="1rem"
                    height="1rem"
                    fill="#B32700"
                    style={{ marginLeft: '4px', verticalAlign: 'middle', ...iconSmall }}
                  />
                </IconButton>
              </CustomTextTooltip>
            )}
          </Grid>
        </Grid>
      </Box>
    );
  };

  const Properties = (
    <Grid
      container={true}
      spacing={2}
      style={{
        padding: '.5rem',
        paddingTop: '0.7rem',
        width: '100%',
        margin: '0px',
        ...(Object.keys(properties).length === 0 || schema['$schema']
          ? { border: 'none', ...(uiSchema['styles'] || {}) }
          : uiSchema['styles']),
      }}
    >
      {properties.map((element, index) => {
        return element.hidden ? (
          element.content
        ) : (
          <Grid item={true} {...calculateGrid(element)} key={index}>
            {element.content}
          </Grid>
        );
      })}
    </Grid>
  );

  const fieldTitle = uiSchema['ui:title'] || title;

  return (
    <>
      {fieldTitle ? (
        <>
          <CustomTitleField
            id={`${idSchema.$id}-title`}
            title={additional ? 'Value' : fieldTitle}
            description={description}
            properties={properties}
          />
          {Object.keys(properties).length > 0 && show && Properties}
        </>
      ) : (
        Properties
      )}
    </>
  );
};

export default ObjectFieldTemplate;
