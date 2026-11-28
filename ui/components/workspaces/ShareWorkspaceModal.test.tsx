import React from 'react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ShareWorkspaceModal from './ShareWorkspaceModal';

const updatePatterns = vi.fn();
const updateView = vi.fn();
let shareModalProps: any;

// The visibility helpers in ./SpacesSwitcher/hooks are deliberately not mocked: the
// test follows the modal's handler down to the argument the mutation receives.
vi.mock('@sistent/sistent', () => ({
  ShareModal: (props: any) => {
    shareModalProps = props;
    return null;
  },
  DesignIcon: () => null,
  ViewIcon: () => null,
  PROMPT_VARIANTS: { DANGER: 'danger' },
  useTheme: () => ({}),
}));

vi.mock('@/rtk-query/design', () => ({
  useUpdatePatternFileMutation: () => [updatePatterns],
  useDeletePatternFileMutation: () => [vi.fn()],
}));

vi.mock('@/rtk-query/view', () => ({
  useUpdateViewVisibilityMutation: () => [updateView],
  useDeleteViewMutation: () => [vi.fn()],
}));

vi.mock('@/rtk-query/resource', () => ({
  useCreateAndRevokeResourceAccessRecordMutation: () => [vi.fn()],
  useGetAccessActorsInfoOfResourceQuery: () => ({ data: { users: [] }, isSuccess: true }),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetAllUsersQuery: vi.fn(),
  useGetLoggedInUserQuery: () => ({ data: { id: 'user-1' } }),
  useGetUserProfileSummaryByIdQuery: () => ({ data: { id: 'user-1' }, isSuccess: true }),
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('@/constants/endpoints', () => ({ MESHERY_CLOUD_PROD: 'https://cloud.example.com' }));

vi.mock('@/utils/Enum', () => ({
  APP_MODE: { DESIGN: 'design', OPERATOR: 'operator' },
  RESOURCE_TYPE: { DESIGN: 'design', VIEW: 'view', CATALOG: 'catalog' },
}));

vi.mock('@/utils/utils', () => ({
  isInOperatorMode: () => false,
  JsonParse: (value: string) => JSON.parse(value),
}));

vi.mock('lib/event-types', () => ({ EVENT_TYPES: {} }));

vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

vi.mock('@/store/slices/mesheryUi', () => ({ updateProgress: vi.fn() }));

vi.mock('@/utils/fileDownloader', () => ({ default: vi.fn(), downloadFileFromContent: vi.fn() }));

const design = {
  id: 'p-1',
  name: 'my design',
  userId: 'user-1',
  catalogData: { type: 'deployment' },
  patternFile: JSON.stringify({ name: 'my design' }),
};

describe('ShareWorkspaceModal visibility updates', () => {
  beforeEach(() => {
    updatePatterns.mockReset().mockResolvedValue({});
    updateView.mockReset().mockResolvedValue({});
    shareModalProps = undefined;
  });

  it('sends the design and its new visibility as updateBody, the key the mutation posts', async () => {
    render(<ShareWorkspaceModal resource={design} type="design" handleClose={vi.fn()} />);

    await shareModalProps.handleUpdateVisibility('public');

    expect(updatePatterns).toHaveBeenCalledWith({
      updateBody: {
        id: 'p-1',
        name: 'my design',
        catalogData: { type: 'deployment' },
        designFile: { name: 'my design' },
        visibility: 'public',
      },
    });
    expect(updateView).not.toHaveBeenCalled();
  });

  it('sends a view visibility change through the view mutation', async () => {
    render(<ShareWorkspaceModal resource={{ id: 'v-1' }} type="view" handleClose={vi.fn()} />);

    await shareModalProps.handleUpdateVisibility('private');

    expect(updateView).toHaveBeenCalledWith({ id: 'v-1', body: { visibility: 'private' } });
    expect(updatePatterns).not.toHaveBeenCalled();
  });

  it('uses the first resource when given a list', async () => {
    render(<ShareWorkspaceModal resource={[design]} type="design" handleClose={vi.fn()} />);

    await shareModalProps.handleUpdateVisibility('private');

    expect(updatePatterns.mock.calls[0][0].updateBody).toMatchObject({
      id: 'p-1',
      visibility: 'private',
    });
  });

  it('returns an error and calls no mutation when the list is empty', async () => {
    render(<ShareWorkspaceModal resource={[]} type="design" handleClose={vi.fn()} />);

    await expect(shareModalProps.handleUpdateVisibility('public')).resolves.toEqual({
      error: 'No resource selected',
    });
    expect(updatePatterns).not.toHaveBeenCalled();
    expect(updateView).not.toHaveBeenCalled();
  });

  it('propagates 4xx/5xx/network errors and success for designs', async () => {
    render(<ShareWorkspaceModal resource={design} type="design" handleClose={vi.fn()} />);

    updatePatterns.mockResolvedValueOnce({ error: { status: 403, data: 'Forbidden' } });
    expect(await shareModalProps.handleUpdateVisibility('public')).toEqual({ error: 'Forbidden' });

    updatePatterns.mockResolvedValueOnce({
      error: { status: 500, data: { message: 'Internal Server Error' } },
    });
    expect(await shareModalProps.handleUpdateVisibility('private')).toEqual({
      error: { message: 'Internal Server Error' },
    });

    updatePatterns.mockResolvedValueOnce({
      error: { status: 'FETCH_ERROR', error: 'Failed to fetch' },
    });
    expect(await shareModalProps.handleUpdateVisibility('private')).toEqual({
      error: 'Failed to fetch',
    });

    updatePatterns.mockResolvedValueOnce({ data: { success: true } });
    expect(await shareModalProps.handleUpdateVisibility('public')).toEqual({ error: undefined });
  });

  it('propagates 4xx/5xx/network errors and success for views', async () => {
    render(
      <ShareWorkspaceModal
        resource={{ id: 'v-1', name: 'my-view', userId: 'user-1' }}
        type="view"
        handleClose={vi.fn()}
      />,
    );

    updateView.mockResolvedValueOnce({ error: { status: 400, data: 'Bad Request' } });
    expect(await shareModalProps.handleUpdateVisibility('private')).toEqual({
      error: 'Bad Request',
    });

    updateView.mockResolvedValueOnce({
      error: { status: 500, data: { message: 'Internal Server Error' } },
    });
    expect(await shareModalProps.handleUpdateVisibility('public')).toEqual({
      error: { message: 'Internal Server Error' },
    });

    updateView.mockResolvedValueOnce({
      error: { status: 'FETCH_ERROR', error: 'Failed to fetch' },
    });
    expect(await shareModalProps.handleUpdateVisibility('private')).toEqual({
      error: 'Failed to fetch',
    });

    updateView.mockResolvedValueOnce({ data: { success: true } });
    expect(await shareModalProps.handleUpdateVisibility('public')).toEqual({ error: undefined });
  });

  it('returns undefined error without calling mutations for unsupported resource types', async () => {
    render(
      <ShareWorkspaceModal
        resource={{ id: 'c-1', name: 'catalog-item', userId: 'user-1' }}
        type="catalog"
        handleClose={vi.fn()}
      />,
    );

    expect(await shareModalProps.handleUpdateVisibility('public')).toEqual({ error: undefined });
    expect(updatePatterns).not.toHaveBeenCalled();
    expect(updateView).not.toHaveBeenCalled();
  });
});
