import React, { useState } from 'react';
import { Box, Typography, styled, CustomTooltip, Collapse, ErrorBoundary } from '@sistent/sistent';
import { ComponentIcon } from '@/components/DesignLifeCycle/common';
import { InfoIcon } from '@sistent/sistent'; // Assuming MUI icons are available
import ExpandLessIcon from '@/assets/icons/ExpandLessIcon';
import ExpandMoreIcon from '@/assets/icons/ExpandMoreIcon';

// Styled components
const SectionContainer = styled(Box)(({ theme }) => ({
  padding: '16px',
  marginBottom: '16px',
  borderRadius: '3px',
  backgroundColor: theme.palette.background.blur.heavy,
  boxShadow: '0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24)',
}));

const SectionHeader = styled(Box)(({ theme, expanded }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: expanded ? '0 0 12px 0' : '0px',
  borderBottom: expanded ? `1px solid ${theme.palette.divider}` : 'none',
  marginBottom: expanded ? '12px' : '0px',
  cursor: 'pointer',
}));

const SectionTitle = styled(Box)(() => ({
  display: 'flex',
  alignItems: 'center',
  gap: '8px',
  fontSize: '16px',
  fontWeight: 600,
}));

const ItemRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  padding: '8px 0',
  gap: '8px',
  borderBottom: `1px solid ${theme.palette.divider}`,
  '&:last-child': {
    borderBottom: 'none',
  },
}));

const EmptyState = styled(Box)(({ theme }) => ({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  padding: '24px',
  background: theme.palette.background.blur.light,
  color: theme.palette.text.primary,
  fontStyle: 'italic',
}));

const TraceSection = ({ title, items, children, emptyMessage = 'No changes' }) => {
  const [expanded, setExpanded] = useState(false);
  const toggleExpanded = () => {
    setExpanded((prev) => !prev);
  };
  if (!items || items.length === 0) return null;

  return (
    <SectionContainer>
      <SectionHeader onClick={toggleExpanded} expanded={expanded}>
        <SectionTitle>
          <Typography variant="subtitle1">{title}</Typography>
        </SectionTitle>

        <Box display={'flex'} alignItems={'center'} gap={2}>
          <span> ( {items.length} )</span>

          {expanded ? (
            <ExpandLessIcon height={24} width={24} fill="white" onClick={toggleExpanded} />
          ) : (
            <ExpandMoreIcon onClick={toggleExpanded} height={24} width={24} fill="white" />
          )}
        </Box>
      </SectionHeader>
      <Collapse in={expanded}>
        {items.length > 0 ? children : <EmptyState>{emptyMessage}</EmptyState>}
      </Collapse>
    </SectionContainer>
  );
};

const getComponentFromAction = (action, componentMapping) => {
  if (action.op == 'delete_component') {
    return action?.value?.component;
  }

  if (action.op == 'add_component') {
    return componentMapping?.[action?.value?.item?.id];
  }

  return componentMapping?.[action?.value?.id];
};

const getRelationshipFromAction = (action, relationshipMapping) => {
  if (action.op == 'delete_relationship') {
    return action?.value?.relationship;
  }
  if (action.op == 'add_relationship') {
    return relationshipMapping?.[action?.value?.item?.id];
  }
  return relationshipMapping?.[action?.value?.id];
};

// Component Trace Item
const ComponentAction = ({ action, componentMapping }) => {
  console.log('ComponentAction', action, componentMapping);
  const component = getComponentFromAction(action, componentMapping);

  if (!component) {
    return null;
  }

  return (
    <ItemRow>
      <Box display="flex" alignItems="center" gap={2} flex={1}>
        <ComponentIcon iconSrc={'/' + component?.styles?.svgColor} />
        <Box
          display={'flex'}
          justifyItems={'space-between'}
          justifyContent={'space-between'}
          width={'100%'}
          alignItems={'center'}
        >
          <Typography variant="body2" fontWeight={500}>
            {component?.component?.kind} <strong> &quot;{component.displayName}&quot; </strong>
          </Typography>
          <CustomTooltip
            title={`Model: ${component?.model?.name}  Version: ${component?.model?.version}`}
          >
            <div>
              <ComponentIcon
                iconSrc={`/static/img/${component?.model?.name}.svg`}
                label={component?.model?.name}
              />
            </div>
          </CustomTooltip>
        </Box>
      </Box>
      {action.op === 'update_component' || action.op === 'update_component_configuration' ? (
        <Typography variant="body2" color="text.secondary">
          Updated <code>{formatPath(action?.value?.path)}</code> to{' '}
          <code>{JSON.stringify(action?.value?.value)}</code>
        </Typography>
      ) : (
        ''
      )}
    </ItemRow>
  );
};

const formatPath = (path) => {
  if (Array.isArray(path)) return path.join('.');

  if (typeof path === 'string') return path;

  return JSON.stringify(path);
};

// Relationship Trace Item
const RelationshipAction = ({ action, relationshipMapping }) => {
  const relationship = getRelationshipFromAction(action, relationshipMapping);
  console.log('RelationshipAction', action, relationshipMapping, relationship);

  if (!relationship) {
    return null;
  }
  return (
    <>
      {relationship.selectors.map((selector, index) => (
        <ItemRow key={index}>
          <Box
            flex={1}
            display={'flex'}
            justifyItems={'space-between'}
            justifyContent={'space-between'}
            width={'100%'}
            alignItems={'center'}
          >
            <Typography variant="body2">
              <span style={{ fontWeight: 500 }}>
                {relationship.kind}-{relationship.subType}-{relationship.type}
              </span>{' '}
              relationship from <strong>{selector?.allow?.from?.[0]?.kind || 'Unknown'}</strong> to{' '}
              <strong>{selector?.allow?.to?.[0]?.kind || 'Unknown'}</strong>
            </Typography>

            <CustomTooltip
              title={`Model: ${relationship.model.name} Version: ${relationship?.model?.version}`}
            >
              <div>
                <ComponentIcon
                  iconSrc={`/static/img/${relationship.model.name}.svg`}
                  label={relationship.model.name}
                />
              </div>
            </CustomTooltip>
          </Box>
          {action.op === 'update_relationship' ? (
            <Typography variant="body2" color="text.secondary">
              Updated <code>{formatPath(action?.value?.path)}</code> to{' '}
              <code>{JSON.stringify(action?.value?.value)}</code>
            </Typography>
          ) : (
            ''
          )}
        </ItemRow>
      ))}
    </>
  );
};

// Main Formatter Component
export const RelationshipEvaluationTraceFormatter = ({ actions, design }) => {
  console.log('RelationshipEvaluationTraceFormatter', actions, design);
  const componentMapping = React.useMemo(
    () =>
      design?.components?.reduce?.(
        (acc, component) => ({
          ...acc,
          [component.id]: component,
        }),
        {},
      ),
    [design?.components],
  );

  const relationshipMapping = React.useMemo(
    () =>
      design?.relationships?.reduce?.(
        (acc, relationship) => ({
          ...acc,
          [relationship.id]: relationship,
        }),
        {},
      ),
    [design?.relationships],
  );

  if (!actions) return null;

  console.log('Trace ', actions);

  const hasChanges = actions.length > 0;

  const componentsAdded = actions.filter?.((a) => a.op == 'add_component');
  const componentsRemoved = actions?.filter?.((a) => a.op == 'delete_component');
  const componentsUpdated = actions?.filter?.(
    (a) => a.op == 'update_component' || a.op == 'update_component_configuration',
  );

  const relationshipsAdded = actions?.filter?.((a) => a.op == 'add_relationship');
  const relationshipsRemoved = actions.filter?.((a) => a.op == 'delete_relationship');
  const relationshipsUpdated = actions?.filter?.((a) => a.op == 'update_relationship');

  return (
    <Box mt={2}>
      {!hasChanges ? (
        <EmptyState>
          <InfoIcon />
          <Typography ml={1}>No changes detected in this evaluation</Typography>
        </EmptyState>
      ) : (
        <Box flexDirection="column" display="flex">
          <TraceSection title="Components Added" items={componentsAdded}>
            {componentsAdded.map((action, index) => (
              <ComponentAction
                action={action}
                key={index}
                design={design}
                componentMapping={componentMapping}
              />
            ))}
          </TraceSection>

          <TraceSection title="Components Deleted" items={componentsRemoved}>
            {componentsRemoved.map((action, index) => (
              <ComponentAction
                action={action}
                key={index}
                design={design}
                componentMapping={componentMapping}
              />
            ))}
          </TraceSection>

          <TraceSection title="Components Updated" items={componentsUpdated}>
            {componentsUpdated.map((action, index) => (
              <ComponentAction
                action={action}
                key={index}
                design={design}
                componentMapping={componentMapping}
              />
            ))}
          </TraceSection>

          <TraceSection title="Relationships Added" items={relationshipsAdded}>
            {relationshipsAdded.map((action, index) => (
              <RelationshipAction
                action={action}
                key={index}
                design={design}
                relationshipMapping={relationshipMapping}
              />
            ))}
          </TraceSection>

          <TraceSection title="Relationships Deleted" items={relationshipsRemoved}>
            {relationshipsRemoved.map((action, index) => (
              <RelationshipAction
                action={action}
                key={index}
                design={design}
                relationshipMapping={relationshipMapping}
              />
            ))}
          </TraceSection>

          <TraceSection title="Relationships Updated" items={relationshipsUpdated}>
            {relationshipsUpdated.map((action, index) => (
              <RelationshipAction
                action={action}
                key={index}
                design={design}
                relationshipMapping={relationshipMapping}
              />
            ))}
          </TraceSection>
        </Box>
      )}
    </Box>
  );
};

export const RelationshipEvaluationEventFormatter = ({ event }) => {
  return (
    <ErrorBoundary>
      <Box mt={2}>
        <Typography variant="body1">{event.description}</Typography>
        <RelationshipEvaluationTraceFormatter
          actions={event?.metadata?.evaluation_response?.actions || []}
          design={event?.metadata?.evaluation_response?.design || {}}
        />
      </Box>
    </ErrorBoundary>
  );
};
