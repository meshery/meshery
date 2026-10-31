import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// MesheryPatterns is a large container; its split-out pieces are tested in
// their own files. This file covers only the page access gate: what a viewer
// without the view permission sees, and that no design list is requested.
const hoisted = vi.hoisted(() => ({
  permissions: {} as Record<string, boolean>,
  patternsQuery: vi.fn(),
}));

vi.mock('@sistent/sistent', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@sistent/sistent')>()),
  useHasPermission: (key: { id: string }) => hoisted.permissions[key.id] ?? false,
}));
vi.mock('@/rtk-query/design', () => {
  const mutation = () => [vi.fn()];
  return {
    useGetPatternsQuery: (...args: unknown[]) => hoisted.patternsQuery(...args),
    useClonePatternMutation: mutation,
    useDeletePatternFileMutation: mutation,
    useDeletePatternMutation: mutation,
    useDeployPatternMutation: mutation,
    useImportDesignMutation: mutation,
    usePublishPatternMutation: mutation,
    useUndeployPatternMutation: mutation,
    useUnpublishPatternMutation: mutation,
    useUpdatePatternFileMutation: mutation,
    useUploadPatternFileMutation: mutation,
    useEvaluateRelationshipsMutation: mutation,
  };
});
vi.mock('@/rtk-query/user', () => ({ useGetProviderCapabilitiesQuery: () => ({}) }));
vi.mock('@/utils/hooks/useTableUrlState', () => ({
  useTableUrlState: () => ({
    tableState: { page: 0, pageSize: 10, sortOrder: '', search: '', filters: { vis: '' } },
    updateTableState: vi.fn(),
  }),
}));
vi.mock('@/utils/hooks/useColumnVisibilityPreference', () => ({
  useColumnVisibilityPreference: () => ({
    columnVisibility: {},
    setColumnVisibilityByUser: vi.fn(),
    setColumnVisibilityByResponsive: vi.fn(),
  }),
}));
vi.mock('../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));
vi.mock('../../../utils/dimension', () => ({ useWindowDimensions: () => ({ width: 1280 }) }));
vi.mock('../../../api/meshmodel', () => ({
  getMeshModels: () => Promise.resolve({ models: [] }),
}));
vi.mock('../../../utils/utils', () => ({ modifyRJSFSchema: (schema: unknown) => schema }));
vi.mock('next/router', () => ({
  useRouter: () => ({ query: {}, pathname: '/configuration/designs', push: vi.fn() }),
}));
vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) =>
    selector({ ui: { selectedK8sContexts: [], catalogVisibility: false, user: {} } }),
}));
vi.mock('@xstate/react', () => ({ useActorRef: () => ({ send: vi.fn() }) }));
vi.mock('machines/validator/designValidator', () => ({ designValidationMachine: {} }));
vi.mock('@/assets/icons/Pattern', () => ({ default: () => null }));
vi.mock('../../general/error-404/index', () => ({
  default: () => <div data-testid="default-error" />,
}));
// Pieces only used on the allowed path; their own tests cover them.
vi.mock('./MesheryPatterns.columns', () => ({
  buildPatternColumns: () => [],
  buildPatternsTableOptions: () => ({}),
  PATTERN_COL_VIEWS: [],
}));
vi.mock('./MesheryPatternGridView', () => ({ default: () => null }));
// Both only build click handlers for the allowed view.
vi.mock('./design-lifecycle-handlers', () => ({
  buildDesignLifecycleHandlers: () => new Proxy({}, { get: () => vi.fn() }),
}));
vi.mock('./patterns-actions', () => ({
  createPatternsActions: () => new Proxy({}, { get: () => vi.fn() }),
}));
vi.mock('./MesheryPatternsToolbar', () => ({ default: () => null }));
vi.mock('./YAMLEditor', () => ({ default: () => null }));
vi.mock('../ImportDesignModal', () => ({ ImportDesignModal: () => null }));
vi.mock('../PublishDesignModal', () => ({ PublishDesignModal: () => null }));
vi.mock('../export/ExportDesignModal', () => ({ default: () => null }));
vi.mock('../../general/PromptComponent', () => ({ default: () => null }));
vi.mock('../../shared/Modal/Information/InfoModal', () => ({ default: () => null }));

import { SistentThemeProvider } from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import MesheryPatterns from './MesheryPatterns';

const allow = (...keys: Array<{ id: string }>) => {
  hoisted.permissions = Object.fromEntries(keys.map((key) => [key.id, true]));
};

const renderPatterns = (props = {}) =>
  render(
    <SistentThemeProvider>
      <MesheryPatterns {...props} />
    </SistentThemeProvider>,
  );

describe('MesheryPatterns access gate', () => {
  beforeEach(() => {
    hoisted.patternsQuery.mockReset().mockReturnValue({
      data: undefined,
      isLoading: false,
      refetch: vi.fn(),
    });
  });

  it('shows DefaultError and does not request designs without the view permission', () => {
    allow();

    renderPatterns();

    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(hoisted.patternsQuery).toHaveBeenCalledWith(expect.anything(), { skip: true });
    expect(hoisted.patternsQuery).not.toHaveBeenCalledWith(expect.anything(), { skip: false });
  });

  it('does not let the catalog permission open the Designs page', () => {
    allow(Keys.CatalogManagementViewCatalog);

    renderPatterns();

    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(hoisted.patternsQuery).toHaveBeenCalledWith(expect.anything(), { skip: true });
  });

  it('requests designs with the designs view permission', () => {
    allow(Keys.CatalogManagementViewDesigns);

    renderPatterns();

    expect(screen.queryByTestId('default-error')).not.toBeInTheDocument();
    expect(hoisted.patternsQuery).toHaveBeenCalledWith(expect.anything(), { skip: false });
  });

  it('requests designs on the Catalog page with only the catalog view permission', () => {
    allow(Keys.CatalogManagementViewCatalog);

    renderPatterns({ pageTitle: 'Catalog' });

    expect(screen.queryByTestId('default-error')).not.toBeInTheDocument();
    expect(hoisted.patternsQuery).toHaveBeenCalledWith(expect.anything(), { skip: false });
  });
});
