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
  return (
    <Box>
      <StyledTitle>
        <span> {title}</span> <span> ({components?.length || 0})</span>
      </StyledTitle>

      {components.map((component) => (
        <Box display={'flex'} alignItems={'center'} gap={2}>
          <ComponentIcon iconSrc={'/' + component.styles.svgColor} />
          {component.component.kind} &lsquo;{component.displayName}&rsquo; from model &lsquo;
          {component.model.name}&rsquo;
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
  );
};

export const RelationshipEvaluationTraceFormatter = ({ value: trace }) => {
  console.log('trace ', trace);
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
