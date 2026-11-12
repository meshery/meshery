import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

// The top-level MeshModelDetails component is a simple type-dispatcher.
// We mock all of the internal data display sub-components and the heavy
// dependencies (RTK Query, ReactJson, theming) so we can test only the
// branch logic in the default export.

vi.mock('@/assets/styles/general/tool.styles', () => ({
  DetailsContainer: ({ children, isEmpty }: any) => (
    <div data-testid="details-container" data-is-empty={String(isEmpty)}>
      {children}
    </div>
  ),
  Segment: ({ children }: any) => <div>{children}</div>,
  FullWidth: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('../../constants/navigator', () => ({
  MODELS: 'Models',
  COMPONENTS: 'Components',
  RELATIONSHIPS: 'Relationships',
  REGISTRANTS: 'Registrants',
}));

vi.mock('@/components/data-formatter', () => ({
  FormatStructuredData: () => <div data-testid="format-data" />,
  reorderObjectProperties: (o: any) => o,
}));

vi.mock('@sistent/sistent', () => ({
  FormControl: ({ children }: any) => <div>{children}</div>,
  Select: ({ children }: any) => <select>{children}</select>,
  MenuItem: ({ children }: any) => <option>{children}</option>,
  CircularProgress: () => <div data-testid="loading" />,
  useTheme: () => ({ palette: { text: { default: 'black' } } }),
  Button: ({ children, onClick, ...rest }: any) => (
    <button onClick={onClick} {...rest}>
      {children}
    </button>
  ),
  DownloadIcon: () => <svg />,
  ExpandMoreIcon: () => <svg />,
  Accordion: ({ children }: any) => <div>{children}</div>,
  AccordionSummary: ({ children }: any) => <div>{children}</div>,
  AccordionDetails: ({ children }: any) => <div>{children}</div>,
  styled: (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) =>
      React.createElement(Component, props, children);
    return StyledComponent;
  },
}));

vi.mock('@/utils/Enum', () => ({
  REGISTRY_ITEM_STATES: { IGNORED: 'ignored', ENABLED: 'enabled' },
}));

vi.mock('@/utils/fallback', () => ({
  normalizeStaticImagePath: (p: string) => p,
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useUpdateEntityStatusMutation: () => [vi.fn(), { isLoading: false }],
  useGetComponentsQuery: () => ({
    data: { components: [] },
    isSuccess: true,
  }),
  useGetMeshModelsQuery: () => ({
    data: { models: [] },
    isSuccess: true,
  }),
}));

vi.mock('./MeshModel.style', () => ({
  JustifyAndAlignCenter: ({ children }: any) => <div>{children}</div>,
  StyledKeyValueFormattedValue: ({ children }: any) => <div>{children}</div>,
  StyledKeyValuePropertyDiv: ({ children }: any) => <div>{children}</div>,
  StyledKeyValueProperty: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('./helper', () => ({
  reactJsonTheme: () => ({}),
}));

vi.mock('css/icons.styles', () => ({
  iconSmall: {},
}));

vi.mock('next/dynamic', () => ({
  default: () => () => <div data-testid="react-json" />,
}));

import MeshModelDetails from './MeshModelDetails';

describe('MeshModelDetails', () => {
  it('renders "No <view> selected" when data is empty', () => {
    render(<MeshModelDetails view="Models" showDetailsData={{ type: '', data: {} }} />);
    expect(screen.getByText('No Models selected')).toBeInTheDocument();
    expect(screen.getByTestId('details-container')).toHaveAttribute('data-is-empty', 'true');
  });

  it('renders "No <view> selected" when type is "none"', () => {
    render(
      <MeshModelDetails
        view="Relationships"
        showDetailsData={{ type: 'none', data: { id: 'x' } }}
      />,
    );
    expect(screen.getByText('No Relationships selected')).toBeInTheDocument();
  });

  it('renders null content for unknown type with non-empty data', () => {
    render(
      <MeshModelDetails
        view="Models"
        showDetailsData={{ type: 'unknown-type', data: { id: 'foo' } }}
      />,
    );
    // No content rendered for unknown type
    expect(screen.queryByText('No Models selected')).not.toBeInTheDocument();
    expect(screen.getByTestId('details-container')).toHaveAttribute('data-is-empty', 'false');
  });
});
