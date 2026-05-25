import { styled, Accordion, AccordionSummary } from '@sistent/sistent';
import type { AccordionProps, AccordionSummaryProps } from '@sistent/sistent';

export const StyledAccordion = styled(Accordion)<AccordionProps>(({ theme }) => ({
  border: `1px solid ${theme.palette.divider}`,
  '&.Mui-expanded': {
    margin: 'auto',
  },
}));

export const StyledAccordionSummary = styled(AccordionSummary)<AccordionSummaryProps>(
  ({ theme }) => ({
    borderBottom: `1px solid ${theme.palette.divider}`,
    '& .MuiAccordionSummary-content': {
      '&.Mui-expanded': {
        margin: '12px 0',
      },
    },
  }),
);
