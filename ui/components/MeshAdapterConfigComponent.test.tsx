import React from 'react';
import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const showPrompt = vi.fn();
const manageAdapterMutator = vi.fn();
const pingAdapterTrigger = vi.fn();
const notify = vi.fn();
const updateProgress = vi.fn();
const dispatch = vi.fn();

const defaultAdapters = [
  {
    name: 'Istio',
    adapterLocation: 'localhost:10000',
    uniqueID: 'adapter-istio-1',
  },
  {
    name: 'Linkerd',
    adapterLocation: 'localhost:10001',
    uniqueID: 'adapter-linkerd-2',
  },
];

vi.mock('@/rtk-query/system', () => ({
  useGetAdaptersQuery: () => ({ data: defaultAdapters }),
  useGetAvailableAdaptersQuery: () => ({ data: [] }),
  useLazyPingAdapterQuery: () => [pingAdapterTrigger],
  useManageAdapterMutation: () => [manageAdapterMutator],
}));

vi.mock('react-redux', () => ({
  useDispatch: () => dispatch,
  useSelector: (selector: (state: unknown) => unknown) =>
    selector({
      adapter: {
        meshAdapters: defaultAdapters,
        meshAdaptersts: {},
      },
    }),
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: (...args: any[]) => updateProgress(...args),
}));

vi.mock('@/store/slices/adapter', () => ({
  updateAdaptersInfo: (payload: any) => ({ type: 'adapter/updateAdaptersInfo', payload }),
}));

vi.mock('../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('./general/PromptComponent', () => ({
  default: React.forwardRef(function Prompt(_props, ref) {
    React.useImperativeHandle(ref, () => ({ show: showPrompt }));
    return <div data-testid="prompt-component" />;
  }),
}));

vi.mock('./general/ReactSelectWrapper', () => ({
  default: ({ label, value }: any) => <div data-testid={`select-${label}`}>{value}</div>,
}));

vi.mock('./general/CustomAvatar', () => ({
  default: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({ id: 'id', function: 'fn' }) }),
}));

vi.mock('css/icons.styles', () => ({
  iconMedium: {},
}));

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = (props: any) => <Component {...props}>{props.children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };

  return {
    styled,
    NoSsr: ({ children }: any) => <>{children}</>,
    Grid2: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children }: any) => <span>{children}</span>,
    Button: ({ children, onClick, ...rest }: any) => (
      <button onClick={onClick} {...rest}>
        {children}
      </button>
    ),
    Chip: ({ label, onDelete, onClick }: any) => (
      <div data-testid={`adapter-chip-${label}`}>
        <span data-testid="chip-label" onClick={onClick}>
          {label}
        </span>
        {onDelete && (
          <button data-testid={`chip-delete-btn-${label}`} onClick={onDelete}>
            delete
          </button>
        )}
      </div>
    ),
    TextField: (props: any) => <input data-testid="text-field" {...props} />,
    Tooltip: ({ children }: any) => <>{children}</>,
    Avatar: () => <img data-testid="avatar" alt="" />,
    useTheme: () => ({
      spacing: (val: number) => `${val * 8}px`,
      palette: {
        background: { card: '#fff' },
        icon: { default: '#000' },
      },
    }),
    PROMPT_VARIANTS: { DANGER: 'danger' },
  };
});

import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';

describe('MeshAdapterConfigComponent - Adapter Deletion Confirmation', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    showPrompt.mockResolvedValue('CANCEL');
    manageAdapterMutator.mockReturnValue({
      unwrap: () => Promise.resolve([]),
    });
  });

  it('clicking the adapter delete icon opens the confirmation prompt without calling manageAdapter', async () => {
    const user = userEvent.setup();
    render(<MeshAdapterConfigComponent />);

    const deleteBtn = screen.getByTestId('chip-delete-btn-localhost:10000');
    await user.click(deleteBtn);

    expect(showPrompt).toHaveBeenCalledTimes(1);
    expect(showPrompt).toHaveBeenCalledWith(
      expect.objectContaining({
        title: 'Delete Adapter?',
        subtitle: expect.stringContaining('localhost:10000'),
        primaryOption: 'DELETE',
        variant: 'danger',
      }),
    );
    expect(manageAdapterMutator).not.toHaveBeenCalled();
  });

  it('clicking Cancel closes the dialog and does not send any DELETE mutation', async () => {
    showPrompt.mockResolvedValue('CANCEL');
    const user = userEvent.setup();
    render(<MeshAdapterConfigComponent />);

    const deleteBtn = screen.getByTestId('chip-delete-btn-localhost:10000');
    await user.click(deleteBtn);

    expect(showPrompt).toHaveBeenCalledTimes(1);
    expect(manageAdapterMutator).not.toHaveBeenCalled();
    expect(dispatch).not.toHaveBeenCalled();
  });

  it('confirming Delete calls manageAdapter with correct payload and updates adapter state', async () => {
    showPrompt.mockResolvedValue('DELETE');
    const updatedAdapters: any[] = [];
    manageAdapterMutator.mockReturnValue({
      unwrap: () => Promise.resolve(updatedAdapters),
    });

    const user = userEvent.setup();
    render(<MeshAdapterConfigComponent />);

    const deleteBtn = screen.getByTestId('chip-delete-btn-localhost:10000');
    await user.click(deleteBtn);

    expect(showPrompt).toHaveBeenCalledTimes(1);
    expect(manageAdapterMutator).toHaveBeenCalledTimes(1);
    expect(manageAdapterMutator).toHaveBeenCalledWith({
      method: 'DELETE',
      adapter: 'localhost:10000',
    });

    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Adapter removed.',
      }),
    );
    expect(dispatch).toHaveBeenCalledWith({
      type: 'adapter/updateAdaptersInfo',
      payload: { meshAdapters: updatedAdapters },
    });
  });

  it('prompt rejection displays prompt error and does not report adapter deletion error', async () => {
    const promptError = new Error('Dialog cancelled unexpectedly');
    showPrompt.mockRejectedValue(promptError);

    const user = userEvent.setup();
    render(<MeshAdapterConfigComponent />);

    const deleteBtn = screen.getByTestId('chip-delete-btn-localhost:10000');
    await user.click(deleteBtn);

    expect(showPrompt).toHaveBeenCalledTimes(1);
    expect(manageAdapterMutator).not.toHaveBeenCalled();
    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Failed to display confirmation prompt',
        details: promptError.toString(),
      }),
    );
    expect(notify).not.toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Adapter was not removed due to an error',
      }),
    );
  });

  it('failed deletion preserves existing error handling and does not update adapter state', async () => {
    showPrompt.mockResolvedValue('DELETE');
    const deleteError = new Error('Network failure');
    manageAdapterMutator.mockReturnValue({
      unwrap: () => Promise.reject(deleteError),
    });

    const user = userEvent.setup();
    render(<MeshAdapterConfigComponent />);

    const deleteBtn = screen.getByTestId('chip-delete-btn-localhost:10000');
    await user.click(deleteBtn);

    expect(showPrompt).toHaveBeenCalledTimes(1);
    expect(manageAdapterMutator).toHaveBeenCalledTimes(1);
    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Adapter was not removed due to an error',
        details: deleteError.toString(),
      }),
    );
    expect(dispatch).not.toHaveBeenCalled();
  });

  it('allows clicking delete on another adapter while first prompt is waiting, without locking', async () => {
    // Keep the prompt unresolved for the first adapter
    let resolveFirstPrompt: (val: string) => void;
    const firstPromptPromise = new Promise<string>((resolve) => {
      resolveFirstPrompt = resolve;
    });

    showPrompt.mockReturnValueOnce(firstPromptPromise).mockResolvedValueOnce('CANCEL');

    const user = userEvent.setup();
    render(<MeshAdapterConfigComponent />);

    const firstDeleteBtn = screen.getByTestId('chip-delete-btn-localhost:10000');
    const secondDeleteBtn = screen.getByTestId('chip-delete-btn-localhost:10001');

    // Click first adapter delete
    await user.click(firstDeleteBtn);
    expect(showPrompt).toHaveBeenCalledTimes(1);

    // Deletion lock is NOT held while prompt is waiting, so clicking second adapter's delete icon is not blocked
    await user.click(secondDeleteBtn);
    expect(showPrompt).toHaveBeenCalledTimes(2);

    // Resolve first prompt as CANCEL
    await act(async () => {
      resolveFirstPrompt('CANCEL');
    });

    expect(manageAdapterMutator).not.toHaveBeenCalled();
  });

  it('repeated clicks while deletion is in progress do not trigger duplicate mutations', async () => {
    let resolveDeletePromise: (val: any) => void;
    const deletePromise = new Promise((resolve) => {
      resolveDeletePromise = resolve;
    });

    showPrompt.mockResolvedValue('DELETE');
    manageAdapterMutator.mockReturnValue({
      unwrap: () => deletePromise,
    });

    const user = userEvent.setup();
    render(<MeshAdapterConfigComponent />);

    const deleteBtn = screen.getByTestId('chip-delete-btn-localhost:10000');

    // First click initiates delete flow and resolves prompt to DELETE
    const firstClick = user.click(deleteBtn);

    // Second click during active in-flight mutation
    await user.click(deleteBtn);

    // Finish first deletion
    await act(async () => {
      resolveDeletePromise!([]);
      await firstClick;
    });

    expect(manageAdapterMutator).toHaveBeenCalledTimes(1);
  });

  it('clicking the adapter chip label still invokes pingAdapter (preserves selection/ping)', async () => {
    pingAdapterTrigger.mockReturnValue({
      unwrap: () => Promise.resolve({}),
    });

    const user = userEvent.setup();
    render(<MeshAdapterConfigComponent />);

    const chipLabels = screen.getAllByTestId('chip-label');
    await user.click(chipLabels[0]);

    expect(pingAdapterTrigger).toHaveBeenCalledWith('localhost:10000');
    expect(manageAdapterMutator).not.toHaveBeenCalled();
    expect(showPrompt).not.toHaveBeenCalled();
  });
});
