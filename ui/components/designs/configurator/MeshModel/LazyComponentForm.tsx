import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  CircularProgress,
  ExpandMoreIcon,
  styled,
  Typography,
} from '@sistent/sistent';
import { isEmpty } from 'lodash';
import React from 'react';
import { getMeshModelComponent } from '../../../../api/meshmodel';
import { iconMedium } from '../../../../css/icons.styles';
import PatternServiceForm from '../../../meshery-mesh-interface/PatternServiceForm';
import { useNotification } from '../../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../../lib/event-types';

export const AccordionRoot = styled('div')({
  width: '100%',
  marginBottom: 8,
});

export const Heading = styled(Typography)(({ theme }) => ({
  fontSize: theme.typography.pxToRem(15),
  fontWeight: theme.typography.fontWeightRegular,
}));

/**
 * Accordion for one component in the Design Configurator. The component's schema is
 * fetched the first time it is expanded, and its settings form is rendered with its
 * own reference so edits update only this component in the design.
 * @param {{
 * component: Types.ComponentDefinition
 * disabled?: boolean
 * createSettingsChangeHandler: (reference: React.MutableRefObject<any>) => Function
 * }} param0
 * @returns
 */
export default function LazyComponentForm({
  component,
  disabled,
  createSettingsChangeHandler,
  ...otherprops
}) {
  const [expanded, setExpanded] = React.useState(false);
  // Once opened, the form stays mounted so its reference key (the component id in
  // the design) and its form state survive collapsing and re-expanding the accordion.
  const [hasBeenExpanded, setHasBeenExpanded] = React.useState(false);
  const [schemaSet, setSchemaSet] = React.useState({});
  // Each form gets its own reference so edits are written to the matching component.
  const formReference = React.useRef();
  const { notify } = useNotification();

  /**
   * Opens or closes the accordion. On first open it loads the component's schema.
   * @param {boolean} state true to expand, false to collapse
   */
  async function expand(state) {
    if (!state) {
      setExpanded(false);
      return;
    }

    setExpanded(true);
    setHasBeenExpanded(true);
    // useMeshModelComponents sets a top-level apiVersion to the latest version of the
    // kind; fall back to the definition's own version when it is not present.
    const {
      component: { version: componentApiVersion, kind },
      apiVersion = componentApiVersion,
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
    <AccordionRoot>
      <Accordion elevation={0} expanded={expanded} onChange={() => !disabled && expand(!expanded)}>
        <AccordionSummary expandIcon={<ExpandMoreIcon style={iconMedium} />}>
          <Heading>
            {component.displayName}{' '}
            {disabled && <em style={{ opacity: 0.5 }}>(contains invalid schema)</em>}
          </Heading>
        </AccordionSummary>
        <LazyAccordionDetails mounted={hasBeenExpanded}>
          {isEmpty(schemaSet) ? (
            <CircularProgress />
          ) : (
            <PatternServiceForm
              formData={{}}
              color={component?.styles?.primaryColor}
              {...otherprops}
              onSettingsChange={createSettingsChangeHandler?.(formReference)}
              reference={formReference}
              schemaSet={schemaSet}
            />
          )}
        </LazyAccordionDetails>
      </Accordion>
    </AccordionRoot>
  );
}

/**
 * Renders the accordion body only after it has been mounted once, then keeps it
 * mounted so collapsing does not reset the form or its component id.
 * @param {{ mounted: boolean, children?: React.ReactNode }} props
 */
function LazyAccordionDetails(props) {
  if (!props.mounted) return <AccordionDetails />;

  // @ts-ignore // LEE: This behavior is more like what we need - https://codesandbox.io/s/upbeat-tesla-uchsb?file=/src/MyAccordion.js
  return <AccordionDetails>{props.children}</AccordionDetails>;
}
