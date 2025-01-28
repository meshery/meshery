import { styled } from '@layer5/sistent';

export const CardButton = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'flex-end',
  alignItems: 'center',
}));

export const BottomPart = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
}));

export const ResultContainer = styled('div')(() => ({
  margin: '0 0 1rem',
  '& div': {
    display: 'flex',
    alignItems: 'center',
  },
}));

export const PaginationWrapper = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  marginTop: '2rem',
}));
