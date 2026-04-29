import { styled, Accordion, AccordionSummary } from '@sistent/sistent';

export const StyledAccordion = styled(Accordion)({
  border: '1px solid rgba(0,0,0,.125)',
  '&.Mui-expanded': {
    margin: 'auto',
  },
});

export const StyledAccordionSummary = styled(AccordionSummary)({
  borderBottom: '1px solid rgba(0,0,0,.125)',
  '& .MuiAccordionSummary-content': {
    '&.Mui-expanded': {
      margin: '12px 0',
    },
  },
});
