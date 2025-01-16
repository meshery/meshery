import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  CircularProgress,
  Typography,
  styled,
} from '@layer5/sistent';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import { isEmpty } from 'lodash';
import React from 'react';
import { getMeshModelComponent } from '../../../api/meshmodel';
import { iconMedium } from '../../../css/icons.styles';
import PatternServiceForm from '../../MesheryMeshInterface/PatternServiceForm';
// eslint-disable-next-line no-unused-vars
import * as Types from '../MeshModel/hooks/types';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import { UsesSistent } from '@/components/SistentWrapper';

export const AccordionRoot = styled('div')({
  width: '100%',
  marginBottom: 8,
});

export const Heading = styled(Typography)(({ theme }) => ({
  fontSize: theme.typography.pxToRem(15),
  fontWeight: theme.typography.fontWeightRegular,
}));

/**
 *
 * @param {{
 * component: Types.ComponentDefinition
 * }} param0
 * @returns
 */
export default function LazyComponentForm({ component, disabled, ...otherprops }) {
  const [expanded, setExpanded] = React.useState(false);
  const [schemaSet, setSchemaSet] = React.useState({});
  const { notify } = useNotification();

  async function expand(state) {
    if (!state) {
      setExpanded(false);
      return;
    }

    setExpanded(true);
    const {
      component: { vesion: apiVersion, kind },
      model,
    } = component;
    const {
      name: modelName,
      model: { version },
    } = model;
    try {
      if (isEmpty(schemaSet)) {
        const res = await getMeshModelComponent(modelName, kind, version, apiVersion);
        if (res.components[0]) {
          setSchemaSet({
            workload: JSON.parse(res.components[0].component.schema), // has to be removed
          });
        } else {
          throw new Error('found null in component definition');
        }
      }
    } catch (error) {
      notify({
        message: `error getting schema: ${error?.message}`,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
      });
    }
  }
  return (
    <UsesSistent>
      <AccordionRoot>
        <Accordion
          elevation={0}
          expanded={expanded}
          onChange={() => !disabled && expand(!expanded)}
        >
          <AccordionSummary expandIcon={<ExpandMoreIcon style={iconMedium} />}>
            <Heading>
              {component.displayName}{' '}
              {disabled && <em style={{ opacity: 0.5 }}>(contains invalid schema)</em>}
            </Heading>
          </AccordionSummary>
          <LazyAccordionDetails expanded={expanded}>
            {isEmpty(schemaSet) ? (
              <CircularProgress />
            ) : (
              <PatternServiceForm
                formData={{}}
                color={component?.styles?.primaryColor}
                {...otherprops}
                schemaSet={schemaSet}
              />
            )}
          </LazyAccordionDetails>
        </Accordion>
      </AccordionRoot>
    </UsesSistent>
  );
}

function LazyAccordionDetails(props) {
  if (!props.expanded) return <AccordionDetails />;

  // @ts-ignore // LEE: This behavior is more like what we need - https://codesandbox.io/s/upbeat-tesla-uchsb?file=/src/MyAccordion.js
  return <AccordionDetails>{props.children}</AccordionDetails>;
}
