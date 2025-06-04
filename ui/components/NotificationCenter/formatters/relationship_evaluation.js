import React, { memo, useState } from 'react';
import { Box, Typography, styled, CustomTooltip, Collapse } from '@sistent/sistent';
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
  alignItems: 'center',
  padding: '8px 0',
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

// Component Trace Item
const ComponentItem = ({ component }) => (
  <ItemRow>
    <Box display="flex" alignItems="center" gap={2} flex={1}>
      <ComponentIcon iconSrc={'/' + component.styles.svgColor} />
      <Box
        display={'flex'}
        justifyItems={'space-between'}
        justifyContent={'space-between'}
        width={'100%'}
        alignItems={'center'}
      >
        <Typography variant="body2" fontWeight={500}>
          {component.component.kind} <strong> &quot;{component.displayName}&quot; </strong>
        </Typography>
        <CustomTooltip
          title={`Model: ${component.model.name}  Version: ${component?.model?.version}`}
        >
          <div>
            <ComponentIcon
              iconSrc={`/static/img/${component.model.name}.svg`}
              label={component.model.name}
            />
          </div>
        </CustomTooltip>
      </Box>
    </Box>
  </ItemRow>
);

// Relationship Trace Item
const RelationshipItem = ({ relationship }) => {
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
        </ItemRow>
      ))}
    </>
  );
};

// Component Trace List
export const ComponentsTrace = ({ components, title }) => (
  <TraceSection title={title} items={components}>
    {components.map((component, index) => (
      <ComponentItem key={index} component={component} />
    ))}
  </TraceSection>
);

// Relationship Trace List
export const RelationshipsTrace = ({ relationships, title }) => (
  <TraceSection title={title} items={relationships}>
    {relationships.map((relationship, index) => (
      <RelationshipItem key={index} relationship={relationship} />
    ))}
  </TraceSection>
);

// Main Formatter Component
export const RelationshipEvaluationTraceFormatter = memo(function RelationshipTraceFormatter({
  trace,
}) {
  if (!trace) return null;

  const hasChanges =
    trace.componentsAdded?.length > 0 ||
    trace.componentsRemoved?.length > 0 ||
    trace.componentsUpdated?.length > 0 ||
    trace.relationshipsAdded?.length > 0 ||
    trace.relationshipsUpdated?.length > 0 ||
    trace.relationshipsRemoved?.length > 0;

  return (
    <Box mt={2}>
      {!hasChanges ? (
        <EmptyState>
          <InfoIcon />
          <Typography ml={1}>No changes detected in this evaluation</Typography>
        </EmptyState>
      ) : (
        <Box flexDirection="column" display="flex">
          <ComponentsTrace
            title="Components Added"
            components={trace.componentsAdded || []}
            type="added"
          />
          <ComponentsTrace
            title="Components Updated"
            components={trace.componentsUpdated || []}
            type="updated"
          />
          <ComponentsTrace
            title="Components Removed"
            components={trace.componentsRemoved || []}
            type="deleted"
          />
          <RelationshipsTrace
            title="Relationships Added"
            type="added"
            relationships={trace.relationshipsAdded || []}
          />
          <RelationshipsTrace
            title="Relationships Updated"
            type="updated"
            relationships={trace.relationshipsUpdated || []}
          />
          <RelationshipsTrace
            title="Relationships Removed"
            type="deleted"
            relationships={trace.relationshipsRemoved || []}
          />
        </Box>
      )}
    </Box>
  );
});

export const RelationshipEvaluationEventFormatter = ({ event }) => {
  return (
    <Box mt={2}>
      <Typography variant="body1">{event.description}</Typography>
      <RelationshipEvaluationTraceFormatter trace={event?.metadata?.trace} />
    </Box>
  );
};
