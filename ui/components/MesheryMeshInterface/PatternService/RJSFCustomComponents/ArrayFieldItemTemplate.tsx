import React from 'react';
import { Box, Grid2, Paper, IconButton, useTheme } from '@sistent/sistent';
import DeleteIcon from '../../../../assets/icons/DeleteIcon';
import ExpandMoreIcon from '../../../../assets/icons/ExpandMoreIcon';
import {
  Accordion as MuiAccordion,
  AccordionSummary as MuiAccordionSummary,
  AccordionDetails as MuiAccordionDetails,
  Typography,
} from '@sistent/sistent';
import { iconSmall } from '../../../../css/icons.styles';
import { safeStringTitle } from '../helper';

const MesheryArrayFieldItemTemplate = (props) => {
  const { children, buttonsProps, index, itemKey } = props;
  const theme = useTheme();
  const accordionDetailsRef = React.useRef(null);

  const heading = `Item (${index})`;
  const safeHeading = safeStringTitle(heading);
  const displayHeading = safeHeading
    ? safeHeading.charAt(0).toUpperCase() + safeHeading.slice(1)
    : '';

  return (
    <Box width={'100%'} marginBottom={'0rem'} data-rjsf-itemkey={itemKey}>
      <MuiAccordion defaultExpanded elevation={0}>
        <MuiAccordionSummary
          expandIcon={<ExpandMoreIcon style={iconSmall} fill="gray" />}
          aria-controls={`panel-${itemKey}-content`}
          id={`panel-${itemKey}-header`}
          style={{ backgroundColor: 'F7F7F7' }}
        >
          <Typography
            style={{
              fontSize: theme.typography.pxToRem(15),
              fontWeight: theme.typography.fontWeightRegular,
            }}
          >
            {displayHeading}
          </Typography>

          {buttonsProps?.hasRemove && (
            <IconButton
              style={{ padding: '0' }}
              disabled={buttonsProps.disabled || buttonsProps.readonly}
              onClick={buttonsProps.onRemoveItem}
            >
              <DeleteIcon fill="gray" style={iconSmall} />
            </IconButton>
          )}
        </MuiAccordionSummary>
        <MuiAccordionDetails ref={accordionDetailsRef}>
          <Grid2 container={true} alignItems="center" size="grow">
            <Grid2 size={{ xs: 12 }}>
              <Box mb={2} style={{ border: '0.5px solid black' }}>
                <Paper elevation={0}>
                  <Box p={2}>{children}</Box>
                </Paper>
              </Box>
            </Grid2>
          </Grid2>
        </MuiAccordionDetails>
      </MuiAccordion>
    </Box>
  );
};

export default MesheryArrayFieldItemTemplate;
