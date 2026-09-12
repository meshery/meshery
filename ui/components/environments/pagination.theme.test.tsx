import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

/**
 * Regression coverage for the "previous pagination button invisible in dark
 * mode" bug (#21345).
 *
 * Sistent's ChevronLeftIcon/ChevronRightIcon default their fill to a hardcoded
 * onyx black, so when they are used as PaginationItem slot icons the chevrons
 * render black-on-black on dark backgrounds. The Environments and
 * Workspaces pages pass an explicit fill="currentColor" through the slot so
 * the chevrons follow the button's text color (light in dark mode, dark in
 * light mode), matching how the library's own default icons behave.
 *
 * This test renders the real MUI Pagination components (the sistent module is
 * mocked at the boundary with the real MUI implementations underneath) and
 * pins that the slot chevrons carry fill="currentColor". Without the fix the
 * icons fall back to their hardcoded black fill and this test fails.
 *
 * @sistent/sistent is intentionally NOT loaded for real here: its bundle
 * pulls the full application asset graph (see the excluded registry tests in
 * vitest.config.ts), which does not resolve under vitest.
 */

vi.mock('@sistent/sistent', async () => {
  const mui = await import('@mui/material');
  const { default: ChevronLeft } = await import('@mui/icons-material/ChevronLeft');
  const { default: ChevronRight } = await import('@mui/icons-material/ChevronRight');
  const { default: Delete } = await import('@mui/icons-material/Delete');
  const { default: Edit } = await import('@mui/icons-material/Edit');
  const { default: SyncAlt } = await import('@mui/icons-material/SyncAlt');

  const defaultTheme = mui.createTheme();

  return {
    SistentThemeProvider: ({ children }: { children?: React.ReactNode }) => <>{children}</>,
    Pagination: mui.Pagination,
    PaginationItem: mui.PaginationItem,
    Button: mui.Button,
    Grid2: mui.Grid,
    Typography: mui.Typography,
    NoSsr: mui.NoSsr,
    IconButton: mui.IconButton,
    Checkbox: mui.Checkbox,
    Card: mui.Card,
    Box: mui.Box,
    styled: mui.styled,
    useTheme: () => ({
      ...defaultTheme,
      palette: {
        ...defaultTheme.palette,
        icon: { default: 'currentColor', secondary: 'currentColor' },
      },
    }),
    ChevronLeftIcon: ChevronLeft,
    ChevronRightIcon: ChevronRight,
    DeleteIcon: Delete,
    EditIcon: Edit,
    SyncAltIcon: SyncAlt,
    Modal: ({ open, children }: { open?: boolean; children?: React.ReactNode }) =>
      open ? <>{children}</> : null,
    ModalBody: ({ children }: { children?: React.ReactNode }) => <>{children}</>,
    ModalFooter: ({ children }: { children?: React.ReactNode }) => <>{children}</>,
    PrimaryActionButtons: () => null,
    TransferList: () => null,
    SearchBar: () => null,
    ErrorBoundary: ({ children }: { children?: React.ReactNode }) => <>{children}</>,
    PromptComponent: () => null,
    useHasPermission: () => true,
    createAndEditEnvironmentSchema: {},
    createAndEditEnvironmentUiSchema: {},
    PROMPT_VARIANTS: {},
  };
});

vi.mock('next/router', () => ({
  withRouter: (component: React.FC) => component,
}));

vi.mock('../../public/static/img/drawer-icons/lifecycle_mgmt_svg', () => ({}));

vi.mock('../shared/Modal/Modal', () => ({
  RJSFModalWrapper: () => null,
}));

vi.mock('../general/error-404/index', () => ({
  default: () => null,
}));

vi.mock('../lifecycle/general', () => ({
  EmptyState: () => null,
  FlipCard: () => null,
}));

vi.mock('./environment-card', () => ({
  default: () => null,
}));

vi.mock('react-redux', () => ({
  useSelector: () => ({ organization: { id: 'org-1' } }),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: vi.fn(),
  mesheryUiReducer: (state: unknown) => state,
  default: (state: unknown) => state,
}));

vi.mock('../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

const ENV_FIXTURE = {
  id: 'env-1',
  name: 'prod',
  description: 'production environment',
};

const ENVIRONMENTS_QUERY_RESULT = {
  data: { environments: [ENV_FIXTURE], totalCount: 23 },
  isLoading: false,
  isError: false,
  error: undefined,
};

const ENV_CONNECTIONS_QUERY_RESULT = {
  data: { connections: [], totalCount: 0 },
  isLoading: false,
};

vi.mock('../../rtk-query/environments', () => ({
  useGetEnvironmentsQuery: () => ENVIRONMENTS_QUERY_RESULT,
  useGetEnvironmentConnectionsQuery: () => ENV_CONNECTIONS_QUERY_RESULT,
  useAddConnectionToEnvironmentMutation: () => [vi.fn()],
  useRemoveConnectionFromEnvironmentMutation: () => [vi.fn()],
  useCreateEnvironmentMutation: () => [vi.fn()],
  useUpdateEnvironmentMutation: () => [vi.fn()],
  useDeleteEnvironmentMutation: () => [vi.fn()],
}));

import Environments from './index';

describe('Environments pagination dark-mode visibility (#21345)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders 3 pages and passes fill=currentColor to the slot chevrons', () => {
    const { container } = render(<Environments />);

    const pageButtons = container.querySelectorAll(
      '.MuiPaginationItem-root:not(.MuiPaginationItem-previousNext)',
    );
    expect(pageButtons.length).toBe(3);

    const chevronSvgs = container.querySelectorAll('.MuiPaginationItem-previousNext svg');
    expect(chevronSvgs.length).toBe(2);

    chevronSvgs.forEach((svg) => {
      expect(svg.getAttribute('fill')).toBe('currentColor');
    });
  });
});
