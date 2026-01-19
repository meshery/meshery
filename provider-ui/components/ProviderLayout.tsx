import React, { ReactNode } from 'react';
import { NoSsr } from '@sistent/sistent';
import { accentGrey, styled } from '@sistent/sistent';

const RootContainer = styled('div')(() => ({
  padding: '4vh 12vw',
  borderRadius: '.5rem',
  textAlign: 'center',
  backgroundColor: accentGrey[20],
  margin: 'auto',
}));

interface ProviderLayoutProps {
  children: ReactNode;
}

export default function ProviderLayout({ children }: ProviderLayoutProps): React.ReactElement {
  return (
    <>
      <NoSsr>
        <RootContainer data-cy="root">{children}</RootContainer>
      </NoSsr>
    </>
  );
}
