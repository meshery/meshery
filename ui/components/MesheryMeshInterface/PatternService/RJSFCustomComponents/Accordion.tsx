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
} from '@sistent/sistent';
import { iconSmall } from '../../../../css/icons.styles';

export default function SimpleAccordion(props) {
  const accordionDetailsRef = React.useRef<HTMLDivElement>(null);

  useEffect(() => {
    // for managing focus
    if (accordionDetailsRef.current) {
      accordionDetailsRef.current.scrollIntoView({
        behavior: 'smooth',
        block: 'nearest',
      });
    }
  }, []);
  const theme = useTheme();
  return (
    <Box width={'100%'} marginBottom={'0rem'}>
      <MuiAccordion defaultExpanded elevation={0}>
        <MuiAccordionSummary
          expandIcon={
            <ExpandMoreIcon
              width={iconSmall.width}
              height={iconSmall.height}
              style={iconSmall}
              fill="gray"
            />
          }
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
              style={{ padding: '0', ...iconSmall }}
              disabled={props.childProps.disabled || props.childProps.readonly}
              onClick={props.childProps.onDropIndexClick(props.childProps.index)}
            >
              <DeleteIcon
                width={iconSmall.width}
                height={iconSmall.height}
                fill="gray"
                style={iconSmall}
              />
            </IconButton>
          )}
        </MuiAccordionSummary>
        <MuiAccordionDetails ref={accordionDetailsRef}>{props.children}</MuiAccordionDetails>
      </MuiAccordion>
    </Box>
  );
}
