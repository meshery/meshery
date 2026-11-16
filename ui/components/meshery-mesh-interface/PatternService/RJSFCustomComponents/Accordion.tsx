import React from 'react';
import {
  IconButton,
  Accordion as MuiAccordion,
  AccordionSummary as MuiAccordionSummary,
  AccordionDetails as MuiAccordionDetails,
  Typography,
  Box,
  useTheme,
  ExpandMoreIcon,
  DeleteIcon,
} from '@sistent/sistent';
import { iconSmall } from '../../../../css/icons.styles';
import { safeStringTitle } from '../helper';

export default function SimpleAccordion(props) {
  const theme = useTheme();
  return (
    <Box sx={{ width: '100%', marginBottom: '0rem' }}>
      <MuiAccordion defaultExpanded elevation={0}>
        <MuiAccordionSummary
          expandIcon={<ExpandMoreIcon style={iconSmall} fill="gray" />}
          aria-controls="panel1a-content"
          id="panel1a-header"
          style={{ backgroundColor: 'F7F7F7' }}
        >
          <Typography
            style={{
              fontSize: theme.typography.pxToRem(15),
              fontWeight: theme.typography.fontWeightRegular,
            }}
          >
            {(() => {
              const h = safeStringTitle(props.heading);
              return h ? h.charAt(0).toUpperCase() + h.slice(1) : '';
            })()}{' '}
          </Typography>

          {props.childProps.hasRemove && (
            <IconButton
              style={{ padding: '0', ...iconSmall }}
              disabled={props.childProps.disabled || props.childProps.readonly}
              onClick={props.childProps.onRemoveItem}
            >
              <DeleteIcon fill="gray" style={iconSmall} />
            </IconButton>
          )}
        </MuiAccordionSummary>
        <MuiAccordionDetails>{props.children}</MuiAccordionDetails>
      </MuiAccordion>
    </Box>
  );
}
