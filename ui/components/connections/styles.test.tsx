import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  type AnyProps = React.HTMLAttributes<HTMLElement> & {
    children?: React.ReactNode;
    'data-testid'?: string;
  };
  const passthrough =
    (tag: keyof JSX.IntrinsicElements, defaultTestId?: string) =>
    ({ children, ...rest }: AnyProps) => {
      const elementProps = { ...rest } as Record<string, unknown>;
      if (!elementProps['data-testid'] && defaultTestId)
        elementProps['data-testid'] = defaultTestId;
      return React.createElement(tag, elementProps, children);
    };

  // styled(Component)(styleArg) -> returns a passthrough that delegates to Component
  // while preserving caller-provided data-testid attributes for assertions.
  const styledFactory = <P extends object>(Component: React.ComponentType<P> | string) => {
    return (styleArg: unknown) => {
      const RemoteStyle = ({ children, ...rest }: AnyProps) =>
        typeof Component === 'string'
          ? React.createElement(Component, rest as never, children)
          : React.createElement(
              Component as React.ComponentType<AnyProps>,
              rest as never,
              children,
            );
      (RemoteStyle as unknown as { __styleArg: unknown }).__styleArg = styleArg;
      RemoteStyle.displayName = 'StyledMock';
      return RemoteStyle;
    };
  };

  return {
    Box: passthrough('div'),
    Button: passthrough('button'),
    Chip: passthrough('div'),
    Grid: passthrough('div'),
    Grid2: passthrough('div'),
    MenuItem: passthrough('li'),
    Select: passthrough('div'),
    Tab: passthrough('div'),
    Tabs: passthrough('div'),
    StepConnector: passthrough('div'),
    StepLabel: passthrough('div'),
    Stepper: passthrough('div'),
    TableContainer: passthrough('div'),
    styled: styledFactory,
  };
});

vi.mock('@/theme', () => ({
  alpha: (color: string, value: number) => `alpha(${color}, ${value})`,
}));

vi.mock('../../utils/Enum', () => ({
  CONNECTION_STATES: {
    DISCOVERED: 'discovered',
    REGISTERED: 'registered',
    CONNECTED: 'connected',
    IGNORED: 'ignored',
    MAINTENANCE: 'maintenance',
    DISCONNECTED: 'disconnected',
    DELETED: 'deleted',
    NOTFOUND: 'not found',
  },
}));

import {
  CreateButton,
  InnerTableContainer,
  ActionListItem,
  ActionButton,
  ConnectionTab,
  ConnectionTabs,
  ConnectionStyledSelect,
  ConnectionStyledMenuItem,
  ConnectionIconText,
  ChipWrapper,
  DiscoveredChip,
  DeletedChip,
  RegisteredChip,
  ConnectedChip,
  IgnoredChip,
  DisconnectedChip,
  NotFoundChip,
  MaintainanceChip,
  ColumnWrapper,
  OperationButton,
  FormatterWrapper,
  ContentContainer,
  StepperContainer,
  CustomLabelStyle,
  ColorlibConnector,
  getConnectionStateColors,
} from './styles';

describe('connections/styles - styled components', () => {
  it('CreateButton renders its children as a Grid passthrough', () => {
    render(<CreateButton data-testid="create-button">create</CreateButton>);
    expect(screen.getByTestId('create-button')).toHaveTextContent('create');
  });

  it('InnerTableContainer renders as a TableContainer passthrough', () => {
    render(<InnerTableContainer data-testid="inner-container">row</InnerTableContainer>);
    expect(screen.getByTestId('inner-container')).toHaveTextContent('row');
  });

  it('ActionListItem renders children inside the Box passthrough', () => {
    render(<ActionListItem data-testid="action-list">item</ActionListItem>);
    expect(screen.getByTestId('action-list')).toHaveTextContent('item');
  });

  it('ActionButton renders as a styled Button', () => {
    render(<ActionButton data-testid="action-btn">act</ActionButton>);
    expect(screen.getByTestId('action-btn')).toHaveTextContent('act');
  });

  it('ConnectionTab and ConnectionTabs render children', () => {
    render(
      <ConnectionTabs data-testid="conn-tabs">
        <ConnectionTab data-testid="conn-tab">tab-label</ConnectionTab>
      </ConnectionTabs>,
    );
    expect(screen.getByTestId('conn-tabs')).toBeInTheDocument();
    expect(screen.getByTestId('conn-tab')).toHaveTextContent('tab-label');
  });

  it('ConnectionStyledSelect and ConnectionStyledMenuItem render', () => {
    render(
      <ConnectionStyledSelect data-testid="select">
        <ConnectionStyledMenuItem data-testid="menu">one</ConnectionStyledMenuItem>
      </ConnectionStyledSelect>,
    );
    expect(screen.getByTestId('select')).toBeInTheDocument();
    expect(screen.getByTestId('menu')).toHaveTextContent('one');
  });

  it('ConnectionIconText renders as a styled div', () => {
    render(<ConnectionIconText data-testid="icon-text">label</ConnectionIconText>);
    expect(screen.getByTestId('icon-text')).toHaveTextContent('label');
  });

  it('ChipWrapper renders as a styled Chip', () => {
    render(<ChipWrapper data-testid="chip-wrapper">label</ChipWrapper>);
    expect(screen.getByTestId('chip-wrapper')).toHaveTextContent('label');
  });

  it.each([
    ['DiscoveredChip', DiscoveredChip],
    ['DeletedChip', DeletedChip],
    ['RegisteredChip', RegisteredChip],
    ['ConnectedChip', ConnectedChip],
    ['IgnoredChip', IgnoredChip],
    ['DisconnectedChip', DisconnectedChip],
    ['NotFoundChip', NotFoundChip],
    ['MaintainanceChip', MaintainanceChip],
  ])('%s renders as a styled Chip', (label, Component) => {
    const testId = `chip-${label}`;
    render(<Component data-testid={testId}>{label}</Component>);
    expect(screen.getByTestId(testId)).toHaveTextContent(label);
  });

  it('ColumnWrapper renders as a styled div', () => {
    render(<ColumnWrapper data-testid="column-wrapper">column</ColumnWrapper>);
    expect(screen.getByTestId('column-wrapper')).toHaveTextContent('column');
  });

  it('OperationButton renders as a styled Grid2', () => {
    render(<OperationButton data-testid="op-btn">op</OperationButton>);
    expect(screen.getByTestId('op-btn')).toHaveTextContent('op');
  });

  it('FormatterWrapper renders as a styled Box', () => {
    render(<FormatterWrapper data-testid="formatter">f</FormatterWrapper>);
    expect(screen.getByTestId('formatter')).toHaveTextContent('f');
  });

  it('ContentContainer renders as a styled Grid2', () => {
    render(<ContentContainer data-testid="content">c</ContentContainer>);
    expect(screen.getByTestId('content')).toHaveTextContent('c');
  });

  it('StepperContainer, CustomLabelStyle, and ColorlibConnector all render', () => {
    render(
      <StepperContainer data-testid="stepper">
        <CustomLabelStyle data-testid="label">label</CustomLabelStyle>
        <ColorlibConnector data-testid="connector" />
      </StepperContainer>,
    );
    expect(screen.getByTestId('stepper')).toBeInTheDocument();
    expect(screen.getByTestId('label')).toHaveTextContent('label');
    expect(screen.getByTestId('connector')).toBeInTheDocument();
  });
});

describe('getConnectionStateColors', () => {
  const baseTheme = {
    palette: {
      primary: { main: '#1976d2' },
      warning: { main: '#ed6c02' },
      error: { main: '#d32f2f' },
    },
  } as never;

  it('maps each documented CONNECTION_STATES key to a palette color', () => {
    const colors = getConnectionStateColors(baseTheme);
    expect(colors.connected).toBe('#1976d2');
    expect(colors.registered).toBe('#1976d2');
    expect(colors.discovered).toBe('#ed6c02');
    expect(colors.ignored).toBe('#ed6c02');
    expect(colors.deleted).toBe('#d32f2f');
    expect(colors.maintenance).toBe('#ed6c02');
    expect(colors.disconnected).toBe('#ed6c02');
    // CONNECTION_STATES.NOTFOUND === 'not found'
    expect(colors['not found']).toBe('#ed6c02');
  });

  it('always exposes a key for every CONNECTION_STATE', () => {
    const colors = getConnectionStateColors(baseTheme);
    expect(Object.keys(colors).length).toBe(8);
  });
});
