import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => {
      if (typeof Component === 'string') {
        return React.createElement(Component, props, children);
      }
      return React.createElement(Component, props, children);
    };
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  };

  return {
    styled,
    Box: (props: any) => <div {...props}>{props.children}</div>,
    Typography: (props: any) => <span {...props}>{props.children}</span>,
    Link: (props: any) => (
      <a href={props.href} {...props}>
        {props.children}
      </a>
    ),
  };
});

import {
  StyledSummaryBox,
  StyledSummaryItem,
  SectionHeading,
  StyledColorBox,
  StyledDocsRedirectLink,
} from './style';

describe('Stepper/style', () => {
  it('exports all styled components', () => {
    expect(StyledSummaryBox).toBeDefined();
    expect(StyledSummaryItem).toBeDefined();
    expect(SectionHeading).toBeDefined();
    expect(StyledColorBox).toBeDefined();
    expect(StyledDocsRedirectLink).toBeDefined();
  });

  it('renders components without crashing', () => {
    const { container } = render(
      <StyledSummaryBox>
        <StyledSummaryItem>summary item</StyledSummaryItem>
        <SectionHeading>heading</SectionHeading>
        <StyledColorBox color="#ff0000" />
        <StyledDocsRedirectLink href="https://docs.meshery.io">docs</StyledDocsRedirectLink>
      </StyledSummaryBox>,
    );

    expect(container.textContent).toContain('summary item');
    expect(container.textContent).toContain('heading');
    expect(container.textContent).toContain('docs');
  });
});
