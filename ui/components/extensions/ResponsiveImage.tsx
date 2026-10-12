import React from 'react';
import { styled } from '@sistent/sistent';

export type ResponsiveImageProps = {
  src: string;
  alt?: string;
  testId?: string;
};

const StyledResponsiveImage = styled('img')({
  height: 'auto',
  width: 'auto',
  maxWidth: '140px',
  maxHeight: '85px',
  flexShrink: 0,
});

export const ResponsiveImage = ({ src, alt, testId }: ResponsiveImageProps) => (
  <StyledResponsiveImage data-testid={testId} src={src} alt={alt || ''} />
);
