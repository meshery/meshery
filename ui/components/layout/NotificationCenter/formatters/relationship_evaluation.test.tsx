import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = ({ children, ...props }: any) => {
      if (typeof Component === 'string') {
        return React.createElement(Component, props, children);
      }
      return <div {...props}>{children}</div>;
    };
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
    Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
    styled,
    CustomTooltip: ({ children, title }: any) => (
      <div data-testid="custom-tooltip" data-title={title}>
        {children}
      </div>
    ),
    Collapse: ({ in: isIn, children }: any) => (
      <div data-testid="collapse" data-open={String(Boolean(isIn))}>
        {isIn ? children : null}
      </div>
    ),
    ErrorBoundary: ({ children }: any) => <div data-testid="error-boundary">{children}</div>,
    InfoIcon: () => <svg data-testid="info-icon" />,
    ExpandMoreIcon: () => <button data-testid="expand-more" type="button" />,
  };
});

vi.mock('@/components/designs/lifecycle/common', () => ({
  ComponentIcon: ({ iconSrc, label }: any) => (
    <span data-testid="component-icon" data-icon={iconSrc} data-label={label}>
      icon
    </span>
  ),
}));

vi.mock('@/assets/icons/ExpandLessIcon', () => ({
  default: ({ onClick }: any) => (
    <button data-testid="expand-less" onClick={onClick} type="button" />
  ),
}));

import {
  RelationshipEvaluationTraceFormatter,
  RelationshipEvaluationEventFormatter,
} from './relationship_evaluation';

describe('RelationshipEvaluationTraceFormatter', () => {
  it('returns null when actions is undefined', () => {
    const { container } = render(
      <RelationshipEvaluationTraceFormatter actions={undefined as any} design={{}} />,
    );
    expect(container).toBeEmptyDOMElement();
  });

  it('renders the empty state when there are no actions', () => {
    render(<RelationshipEvaluationTraceFormatter actions={[]} design={{}} />);
    expect(screen.getByText(/No changes detected in this evaluation/i)).toBeInTheDocument();
    expect(screen.getByTestId('info-icon')).toBeInTheDocument();
  });

  it('renders only sections with matching actions', () => {
    const design = {
      components: [
        {
          id: 'c1',
          displayName: 'Comp1',
          component: { kind: 'Pod' },
          styles: { svgColor: 'icon.svg' },
          model: { name: 'core', version: '0.1' },
        },
      ],
      relationships: [],
    };

    render(
      <RelationshipEvaluationTraceFormatter
        actions={[{ op: 'add_component', value: { item: { id: 'c1' } } }]}
        design={design as any}
      />,
    );

    expect(screen.getByText('Components Added')).toBeInTheDocument();
    expect(screen.queryByText('Components Deleted')).not.toBeInTheDocument();
    expect(screen.queryByText('Relationships Added')).not.toBeInTheDocument();
  });

  it('toggles expand/collapse on a section', () => {
    const design = {
      components: [
        {
          id: 'c1',
          displayName: 'Comp1',
          component: { kind: 'Pod' },
          styles: { svgColor: 'icon.svg' },
          model: { name: 'core', version: '0.1' },
        },
      ],
      relationships: [],
    };

    render(
      <RelationshipEvaluationTraceFormatter
        actions={[{ op: 'add_component', value: { item: { id: 'c1' } } }]}
        design={design as any}
      />,
    );

    // Initially collapsed — expand-more icon is visible, contents hidden
    expect(screen.getByTestId('expand-more')).toBeInTheDocument();
    expect(screen.queryByText(/Comp1/)).not.toBeInTheDocument();

    fireEvent.click(screen.getByText('Components Added'));
    expect(screen.getByTestId('expand-less')).toBeInTheDocument();
    expect(screen.getByText(/Comp1/)).toBeInTheDocument();
  });

  it('renders update message for update_component op with path and value', () => {
    const design = {
      components: [
        {
          id: 'c1',
          displayName: 'C',
          component: { kind: 'Pod' },
          styles: { svgColor: 'i.svg' },
          model: { name: 'core', version: '0.2' },
        },
      ],
      relationships: [],
    };

    render(
      <RelationshipEvaluationTraceFormatter
        actions={[
          {
            op: 'update_component',
            value: { id: 'c1', path: ['spec', 'replicas'], value: 3 },
          },
        ]}
        design={design as any}
      />,
    );

    fireEvent.click(screen.getByText('Components Updated'));
    expect(screen.getByText(/spec\.replicas/)).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  it('renders delete_component path: pulls component from action.value.component', () => {
    render(
      <RelationshipEvaluationTraceFormatter
        actions={[
          {
            op: 'delete_component',
            value: {
              component: {
                displayName: 'Gone',
                component: { kind: 'Service' },
                styles: { svgColor: 'g.svg' },
                model: { name: 'core', version: '0.3' },
              },
            },
          },
        ]}
        design={{} as any}
      />,
    );

    fireEvent.click(screen.getByText('Components Deleted'));
    expect(screen.getByText(/Gone/)).toBeInTheDocument();
  });

  it('renders relationships sections when relationship actions are present', () => {
    const design = {
      components: [],
      relationships: [
        {
          id: 'r1',
          kind: 'edge',
          subType: 'network',
          type: 'tcp',
          model: { name: 'core', version: '1.0' },
          selectors: [
            {
              allow: { from: [{ kind: 'pod' }], to: [{ kind: 'svc' }] },
            },
          ],
        },
      ],
    };

    render(
      <RelationshipEvaluationTraceFormatter
        actions={[{ op: 'add_relationship', value: { item: { id: 'r1' } } }]}
        design={design as any}
      />,
    );

    expect(screen.getByText('Relationships Added')).toBeInTheDocument();
    fireEvent.click(screen.getByText('Relationships Added'));
    expect(screen.getByText('pod')).toBeInTheDocument();
    expect(screen.getByText('svc')).toBeInTheDocument();
    expect(screen.getByText(/relationship from/i)).toBeInTheDocument();
  });
});

describe('RelationshipEvaluationEventFormatter', () => {
  it('wraps content in ErrorBoundary and shows the event description with legacy evaluation_response', () => {
    render(
      <RelationshipEvaluationEventFormatter
        event={{
          description: 'Evaluation finished',
          metadata: { evaluation_response: { actions: [], design: {} } },
        }}
      />,
    );

    expect(screen.getByTestId('error-boundary')).toBeInTheDocument();
    expect(screen.getByText('Evaluation finished')).toBeInTheDocument();
  });

  it('wraps content in ErrorBoundary and shows the event description with camelCase evaluationResponse', () => {
    render(
      <RelationshipEvaluationEventFormatter
        event={{
          description: 'Evaluation finished',
          metadata: { evaluationResponse: { actions: [], design: {} } },
        }}
      />,
    );

    expect(screen.getByTestId('error-boundary')).toBeInTheDocument();
    expect(screen.getByText('Evaluation finished')).toBeInTheDocument();
  });

  it('renders the empty state when no actions are present', () => {
    render(<RelationshipEvaluationEventFormatter event={{ description: 'no-op', metadata: {} }} />);

    expect(screen.getByText(/No changes detected/i)).toBeInTheDocument();
  });
});
