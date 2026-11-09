import React from 'react';
import { render, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Covers the notification the Designs page raises when the publish schema
// (built from the mesh model list) can't be loaded. The page's children,
// RTK Query hooks and the validation actor are stubbed; only the container's
// own mount effects run.

// Hook results must be referentially stable across renders: the container
// keys effects on them.
const { notify, getMeshModels, actionDeps, PATTERNS_QUERY_RESULT, TABLE_STATE, REDUX_STATE } =
  vi.hoisted(() => ({
    notify: vi.fn(),
    getMeshModels: vi.fn(),
    actionDeps: vi.fn(),
    PATTERNS_QUERY_RESULT: {
      data: { patterns: [], totalCount: 0 },
      isLoading: false,
      refetch: () => {},
    },
    TABLE_STATE: {
      tableState: {
        page: 0,
        pageSize: 10,
        sortOrder: 'updated_at desc',
        search: '',
        filters: { vis: '' },
      },
      updateTableState: () => {},
    },
    REDUX_STATE: { ui: { selectedK8sContexts: [], catalogVisibility: true, user: {} } },
  }));

vi.mock('../../../api/meshmodel', () => ({ getMeshModels }));
// Record what the container passes to its actions so a test can tell when the
// mesh models have been stored in state.
vi.mock('./patterns-actions', async (importOriginal) => {
  const actual = await importOriginal<typeof import('./patterns-actions')>();
  return {
    ...actual,
    createPatternsActions: (deps) => {
      actionDeps(deps);
      return actual.createPatternsActions(deps);
    },
  };
});
vi.mock('../../../utils/hooks/useNotification', () => ({ useNotification: () => ({ notify }) }));
vi.mock('@/rtk-query/design', () => {
  const mutation = () => [vi.fn()];
  return {
    useClonePatternMutation: mutation,
    useDeletePatternFileMutation: mutation,
    useDeletePatternMutation: mutation,
    useDeployPatternMutation: mutation,
    useGetPatternsQuery: () => PATTERNS_QUERY_RESULT,
    useImportDesignMutation: mutation,
    usePublishPatternMutation: mutation,
    useUndeployPatternMutation: mutation,
    useUnpublishPatternMutation: mutation,
    useUpdatePatternFileMutation: mutation,
    useUploadPatternFileMutation: mutation,
    useEvaluateRelationshipsMutation: mutation,
  };
});
vi.mock('@/rtk-query/user', () => ({ useGetProviderCapabilitiesQuery: () => ({ data: {} }) }));
vi.mock('@/utils/hooks/useTableUrlState', () => ({ useTableUrlState: () => TABLE_STATE }));
vi.mock('@/utils/hooks/useColumnVisibilityPreference', () => {
  const visibility = {
    columnVisibility: {},
    setColumnVisibilityByUser: () => {},
    setColumnVisibilityByResponsive: () => {},
  };
  return { useColumnVisibilityPreference: () => visibility };
});
vi.mock('next/router', () => ({
  useRouter: () => ({ query: {}, pathname: '/', push: vi.fn(), replace: vi.fn() }),
}));
vi.mock('react-redux', () => ({
  useSelector: (selector: (state: unknown) => unknown) => selector(REDUX_STATE),
  useDispatch: () => vi.fn(),
}));
vi.mock('@xstate/react', () => ({ useActorRef: () => ({}) }));
vi.mock('machines/validator/designValidator', () => ({ designValidationMachine: {} }));
vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));

vi.mock('./MesheryPatternGridView', () => ({ default: () => null }));
vi.mock('./MesheryPatternsToolbar', () => ({ default: () => null }));
vi.mock('./YAMLEditor', () => ({ default: () => null }));
vi.mock('../../general/PromptComponent', () => ({ default: React.forwardRef(() => null) }));
vi.mock('../../shared/LoadingState/LoadingComponent', () => ({ default: () => null }));
vi.mock('../../shared/Modal/Information/InfoModal', () => ({ default: () => null }));
vi.mock('../../general/error-404/index', () => ({ default: () => <div>no access</div> }));
vi.mock('../export/ExportDesignModal', () => ({ default: () => null }));
vi.mock('../ImportDesignModal', () => ({ ImportDesignModal: () => null }));
vi.mock('../PublishDesignModal', () => ({ PublishDesignModal: () => null }));
vi.mock('@/assets/icons/Pattern', () => ({ default: () => null }));
// The column and lifecycle builders import JSX icon files from public/ that the
// test transform can't parse; they're covered by their own tests.
vi.mock('./MesheryPatterns.columns', () => ({
  buildPatternColumns: () => [],
  buildPatternsTableOptions: () => ({}),
  PATTERN_COL_VIEWS: [],
}));
vi.mock('./design-lifecycle-handlers', () => ({ buildDesignLifecycleHandlers: () => ({}) }));
// utils/utils imports pages/_app, which pulls in the JSX drawer icons.
vi.mock('../../../pages/_app', () => ({ mesheryExtensionRoute: '/' }));

vi.mock('@sistent/sistent', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@sistent/sistent')>();
  return {
    ...actual,
    ResponsiveDataTable: () => null,
    useHasPermission: () => true,
  };
});

import MesheryPatterns from './MesheryPatterns';

const notifiedEventTypes = () =>
  notify.mock.calls.map(([arg]) => arg?.event_type?.type ?? arg?.event_type);

describe('MesheryPatterns publish schema notifications', () => {
  beforeEach(() => {
    notify.mockReset();
    getMeshModels.mockReset();
    actionDeps.mockReset();
  });

  it('shows an error notification when the mesh models fail to load', async () => {
    getMeshModels.mockRejectedValue(new Error('registry unavailable'));

    render(<MesheryPatterns />);

    await waitFor(() => expect(notify).toHaveBeenCalledTimes(1));
    expect(notifiedEventTypes()).toEqual(['error']);
    expect(notify.mock.calls[0][0].message).toContain('registry unavailable');
  });

  it('does not notify when the mesh models load', async () => {
    getMeshModels.mockResolvedValue({ models: [{ displayName: 'Kubernetes' }] });

    render(<MesheryPatterns />);

    // Wait until the loaded models reach state, i.e. the fetch has fully settled.
    await waitFor(() =>
      expect(actionDeps).toHaveBeenLastCalledWith(
        expect.objectContaining({ meshModels: [{ displayName: 'Kubernetes' }] }),
      ),
    );
    expect(notify).not.toHaveBeenCalled();
  });
});
