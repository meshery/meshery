import React, { useEffect } from 'react';
import DeleteIcon from '../../../../assets/icons/DeleteIcon';
import ExpandMoreIcon from '../../../../assets/icons/ExpandMoreIcon';
import {
  IconButton,
  Accordion as MuiAccordion,
  AccordionSummary as MuiAccordionSummary,
  AccordionDetails as MuiAccordionDetails,
  Typography,
  Box,
  useTheme,
} from '@layer5/sistent';
import { iconSmall } from '../../../../css/icons.styles';

export default function SimpleAccordion(props) {
  useEffect(() => {
    // for managing focus
    if (accordionDetailsRef.current) {
      accordionDetailsRef.current.scrollIntoView({
        behavior: 'smooth',
        block: 'nearest',
      });
    }
  }, []);
  // React Hook useEffect has an unnecessary dependency: 'accordionDetailsRef.current'. Either exclude it or remove the dependency array. Mutable values like 'accordionDetailsRef.current' aren't valid dependencies because mutating them doesn't re-render the component.

  const accordionDetailsRef = React.useRef(null);
  const theme = useTheme();
  return (
    <Box width={'100%'} marginBottom={'0rem'}>
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
            {props.heading?.charAt(0).toUpperCase() + props.heading?.slice(1)}{' '}
          </Typography>

          {props.childProps.hasRemove && (
            <IconButton
              style={{ padding: '0', iconSmall }}
              // style={btnStyle, iconSmall}
              disabled={props.childProps.disabled || props.childProps.readonly}
              onClick={props.childProps.onDropIndexClick(props.childProps.index)}
            >
              <DeleteIcon fill="gray" style={iconSmall} />
            </IconButton>
          )}
        </MuiAccordionSummary>
        <MuiAccordionDetails ref={accordionDetailsRef}>{props.children}</MuiAccordionDetails>
      </MuiAccordion>
    </Box>
  );
}
