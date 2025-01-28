import React, { useEffect } from 'react';
import DeleteIcon from '../../../../assets/icons/DeleteIcon';
import ExpandMoreIcon from '../../../../assets/icons/ExpandMoreIcon';
import { iconSmall } from '../../../../css/icons.styles';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  IconButton,
  useTheme,
} from '@layer5/sistent';
import { AccordionHeading, AccordionRoot } from '../style';
import { UsesSistent } from '@/components/SistentWrapper';

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

  const theme = useTheme();
  const accordionDetailsRef = React.useRef(null);

  return (
    <UsesSistent>
      <AccordionRoot>
        <Accordion defaultExpanded elevation={0}>
          <AccordionSummary
            expandIcon={<ExpandMoreIcon style={iconSmall} fill="gray" />}
            aria-controls="panel1a-content"
            id="panel1a-header"
            style={{ backgroundColor: 'F7F7F7' }}
          >
            <AccordionHeading
              sx={{
                fontSize: theme.typography.pxToRem(15),
                fontWeight: theme.typography.fontWeightRegular,
              }}
            >
              {props.heading?.charAt(0).toUpperCase() + props.heading?.slice(1)}{' '}
            </AccordionHeading>

            {props.childProps.hasRemove && (
              <IconButton
                style={{ padding: '0', iconSmall }}
                disabled={props.childProps.disabled || props.childProps.readonly}
                onClick={props.childProps.onDropIndexClick(props.childProps.index)}
              >
                <DeleteIcon fill="gray" style={iconSmall} />
              </IconButton>
            )}
          </AccordionSummary>
          <AccordionDetails ref={accordionDetailsRef}>{props.children}</AccordionDetails>
        </Accordion>
      </AccordionRoot>
    </UsesSistent>
  );
}
