import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  // Mirror the project pattern: styled(Comp) returns a function that returns
  // a thin pass-through component.
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledMock';
    return StyledComponent;
  };

  return {
    styled,
    Card: ({ children, ...props }: any) => (
      <div data-testid="card" {...props}>
        {children}
      </div>
    ),
    Chip: ({ label, ...props }: any) => (
      <div data-testid="chip" {...props}>
        {label}
      </div>
    ),
    Grid: ({ children, container, item, ...props }: any) => (
      <div
        data-testid="grid"
        data-container={String(!!container)}
        data-item={String(!!item)}
        {...props}
      >
        {children}
      </div>
    ),
    TableCell: ({ children, ...props }: any) => (
      <td data-testid="table-cell" {...props}>
        {children}
      </td>
    ),
  };
});

import {
  AdapterChip,
  AdapterTableHeader,
  AdapterSmWrapper,
  SecondaryTable,
  PaneSection,
  ChipNamespaceContainer,
  InputWrapper,
  AdapterCard,
} from './adapter-play-styled';

describe('adapter-play-styled', () => {
  it('exports an AdapterChip that wraps Chip', () => {
    render(<AdapterChip label="my-chip" />);
    expect(screen.getByTestId('chip')).toHaveTextContent('my-chip');
  });

  it('exports an AdapterTableHeader that wraps TableCell', () => {
    render(
      <table>
        <thead>
          <tr>
            <AdapterTableHeader>Header</AdapterTableHeader>
          </tr>
        </thead>
      </table>,
    );
    expect(screen.getByTestId('table-cell')).toHaveTextContent('Header');
  });

  it('exports an AdapterSmWrapper that renders a div', () => {
    render(<AdapterSmWrapper>contents</AdapterSmWrapper>);
    expect(screen.getByText('contents')).toBeInTheDocument();
  });

  it('exports a SecondaryTable that renders a div', () => {
    render(<SecondaryTable>secondary</SecondaryTable>);
    expect(screen.getByText('secondary')).toBeInTheDocument();
  });

  it('exports a PaneSection that renders a div', () => {
    render(<PaneSection>pane</PaneSection>);
    expect(screen.getByText('pane')).toBeInTheDocument();
  });

  it('exports a ChipNamespaceContainer that wraps Grid', () => {
    render(
      <ChipNamespaceContainer container>
        <span>inner</span>
      </ChipNamespaceContainer>,
    );
    expect(screen.getByTestId('grid')).toBeInTheDocument();
    expect(screen.getByText('inner')).toBeInTheDocument();
  });

  it('exports an InputWrapper that renders a div', () => {
    render(<InputWrapper>input</InputWrapper>);
    expect(screen.getByText('input')).toBeInTheDocument();
  });

  it('exports an AdapterCard that wraps Card', () => {
    render(<AdapterCard>card content</AdapterCard>);
    expect(screen.getByTestId('card')).toHaveTextContent('card content');
  });
});
