import { styled, Accordion, AccordionSummary } from '@layer5/sistent';

export const ExpansionPanel = styled(Accordion)({
  border: '1px solid rgba(0,0,0,.125)',
  '&.Mui-expanded': {
    margin: 'auto',
  },
});

export const ExpansionPanelSummary = styled(AccordionSummary)({
  borderBottom: '1px solid rgba(0,0,0,.125)',
  '& .MuiAccordionSummary-content': {
    '&.Mui-expanded': {
      margin: '12px 0',
    },
  },
});
