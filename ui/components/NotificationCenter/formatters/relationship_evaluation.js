import React from 'react';
import { Box, Typography, styled, Chip, Badge, Tooltip } from '@layer5/sistent';
import { ComponentIcon } from '@/components/DesignLifeCycle/common';
// import { StyledAccordion } from '@/components/StyledAccordion';
import { Box, styled } from '@layer5/sistent';

const StyledTitle = styled(Box)(() => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  fontSize: '16px',
  fontWeight: 'bold',
  paddingBottom: '0.5rem',
}));

export const ComponentsTrace = ({ components, title }) => {
  if (components.length == 0) {
    return null;
  }
};

const TraceSection = ({ title, items, type, children, emptyMessage = 'No changes' }) => {
  if (!items || items.length === 0) return null;

  return (
    <Box>
      <StyledTitle>
        <span> {title}</span> <span> ({components?.length || 0})</span>
      </StyledTitle>

      {components.map((component) => (
        <Box display={'flex'} alignItems={'center'} gap={2}>
          <ComponentIcon iconSrc={'/' + component.styles.svgColor} />
          {component.component.kind} &lsquo;{component.displayName}&rsquo; from model &lsquo;{component.model.name}&rsquo;
        </Box>
      ))}
    </Box>
  );
};

export const RelationshipsTrace = ({ relationships, title, action }) => {
  if (relationships.length == 0) {
    return null;
  }
  return (
    <Box>
      <StyledTitle>
        <span> {title}</span> <span> ({relationships?.length || 0})</span>
      </StyledTitle>

      {relationships.map((relationship) =>
        relationship.selectors.map((selector) => (
          <Box display={'flex'} alignItems={'center'} gap={2}>
            {action} {relationship.kind}-{relationship.subType}-{relationship.type} relationship
            from {selector?.allow?.from?.[0].kind} to {selector?.allow?.to?.[0].kind}
          </Box>
        )),
      )}
    </Box>
  </ItemRow>
);

// Relationship Trace Item
const RelationshipItem = ({ relationship, action }) => (
  <>
    {relationship.selectors.map((selector, index) => (
      <ItemRow key={index}>
        <Box flex={1}

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

         <Tooltip title={`Model: ${relationship.model.name}`}>
           <ModelBadge size="small" label={relationship.model.name} variant="outlined" />
         </Tooltip>
        </Box>
      </ItemRow>
    ))}
  </>
);

// Component Trace List
export const ComponentsTrace = ({ components, title, type }) => (
  <TraceSection title={title} items={components} type={type}>
    {components.map((component, index) => (
      <ComponentItem key={index} component={component} />
    ))}
  </TraceSection>
);

// Relationship Trace List
export const RelationshipsTrace = ({ relationships, title, type }) => (
  <TraceSection title={title} items={relationships} type={type}>
    {relationships.map((relationship, index) => (
      <RelationshipItem key={index} relationship={relationship} action={type} />
    ))}
  </TraceSection>
);

// Main Formatter Component
export const RelationshipEvaluationTraceFormatter = ({ value: trace }) => {
  const hasChanges =
    trace.componentsAdded.length > 0 ||
    trace.componentsRemoved.length > 0 ||
    trace.componentsUpdated.length > 0 ||
    trace.relationshipsAdded.length > 0 ||
    trace.relationshipsUpdated.length > 0 ||
    trace.relationshipsRemoved.length > 0;

  return (
    <Box flexDirection={'column'} gap={2} display={'flex'}>
      <ComponentsTrace title="New Component(s)" components={trace.componentsAdded} />
      <ComponentsTrace title="Deleted Component(s)" components={trace.componentsRemoved} />
      <ComponentsTrace title="Updated Component(s)" components={trace.componentsUpdated} />

      <RelationshipsTrace
        title="Relationship(s) Added"
        action="Added"
        relationships={trace.relationshipsAdded}
      />
      <RelationshipsTrace
        title="Relationship(s) Updated"
        action="Updated"
        relationships={trace.relationshipsUpdated}
      />
      <RelationshipsTrace
        title="Relationship(s) Deleted"
        action="Deleted"
        relationships={trace.relationshipsRemoved}
      />
    </Box>
  );
};
