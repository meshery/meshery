import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import ConfirmationMsg from '../../designs/lifecycle/DeployConfirmationModal';

const dispatch = vi.fn();
let selectorState: any = {
  ui: {
    selectedK8sContexts: [],
    k8sConfig: [],
  },
};

const triggerPing = vi.fn().mockReturnValue({ unwrap: () => Promise.resolve() });
const notify = vi.fn();

vi.mock('react-redux', () => ({
  useDispatch: () => dispatch,
  useSelector: (fn: any) => fn(selectorState),
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  useFilterK8sContexts: (ctxs: any) => ctxs,
}));

vi.mock('@/rtk-query/connection', () => ({
  useLazyPingKubernetesQuery: () => [triggerPing],
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  setK8sContexts: (payload: any) => ({ type: 'setK8sContexts', payload }),
  updateProgress: vi.fn(),
}));

vi.mock('../../../utils/helpers/common', () => ({
  successHandlerGenerator: () => vi.fn(),
  errorHandlerGenerator: () => vi.fn(),
}));

vi.mock('../../../utils/multi-ctx', () => ({
  getK8sConfigIdsFromK8sConfig: (ctxs: any) => ctxs?.map((c: any) => c.id) || [],
}));

vi.mock('../../../utils/Enum', () => ({
  ACTIONS: { VERIFY: 0, UNDEPLOY: 1, DEPLOY: 2 },
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('@/utils/permission_constants', () => ({
  keys: {
    VALIDATE_DESIGN: { action: 'validate', resource: 'design', subject: 'design' },
    UNDEPLOY_DESIGN: { action: 'undeploy', resource: 'design', subject: 'design' },
    DEPLOY_DESIGN: { action: 'deploy', resource: 'design', subject: 'design' },
  },
}));

vi.mock('../EmptyState/K8sContextEmptyState', () => ({
  K8sEmptyState: ({ message }: any) => <div data-testid="k8s-empty">{message || 'empty'}</div>,
}));

vi.mock('../../layout/Header/Header', () => ({
  K8sContextConnectionChip: ({ ctx, selected, onSelectChange }: any) => (
    <button
      onClick={onSelectChange}
      data-testid={`ctx-chip-${ctx.id}`}
      data-selected={String(selected)}
    >
      {ctx.name}
    </button>
  ),
}));

vi.mock('@/components/shared/Modal', () => ({
  Modal: ({ isOpen, onClose, title, children, actions }: any) =>
    isOpen ? (
      <div data-testid="sistent-modal" data-title={title}>
        <button onClick={onClose} aria-label="modal-close">
          close
        </button>
        <div data-testid="modal-body">{children}</div>
        {actions}
      </div>
    ) : null,
}));

vi.mock('../../connections/ConnectionChip', () => ({
  TooltipWrappedConnectionChip: ({ title, handlePing }: any) => (
    <button onClick={handlePing} data-testid={`conn-chip-${title}`}>
      {title}
    </button>
  ),
}));

vi.mock('../../../lib/event-types', () => ({
  EVENT_TYPES: { INFO: 'info', WARNING: 'warning' },
}));

vi.mock('../../../assets/icons/shapes/RoundedTriangle', () => ({
  RoundedTriangleShape: () => <svg data-testid="triangle" />,
}));

vi.mock('../../../assets/icons/shapes/Octagon', () => ({
  default: () => <svg data-testid="octagon" />,
}));

vi.mock('../../../assets/icons/Pattern', () => ({
  default: () => <svg data-testid="pattern-icon" />,
}));

vi.mock('../../../css/icons.styles', () => ({
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = (props: any) => <Component {...props}>{props.children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    Box: ({ children, display }: any) => (
      <div data-display={display} data-testid="box">
        {children}
      </div>
    ),
    Button: ({ children, onClick, disabled, ...rest }: any) => (
      <button onClick={onClick} disabled={disabled} {...rest}>
        {children}
      </button>
    ),
    Checkbox: ({ checked, onChange, color }: any) => (
      <input
        type="checkbox"
        checked={!!checked}
        onChange={onChange}
        data-color={color}
        data-testid="checkbox"
      />
    ),
    DialogActions: ({ children }: any) => <div data-testid="dialog-actions">{children}</div>,
    DialogContent: ({ children }: any) => <div data-testid="dialog-content">{children}</div>,
    DialogContentText: ({ children }: any) => <div data-testid="dialog-text">{children}</div>,
    TextField: ({ id, label, onChange }: any) => (
      <input id={id} aria-label={label} onChange={onChange} data-testid="text-field" />
    ),
    Typography: ({ children, variant }: any) => <p data-variant={variant}>{children}</p>,
    Tab: ({ label, onClick, disabled, ...rest }: any) => (
      <button onClick={onClick} disabled={disabled} {...rest}>
        {label}
      </button>
    ),
    Tabs: ({ value, children }: any) => (
      <div data-testid="tabs" data-value={value}>
        {children}
      </div>
    ),
    DoneAllIcon: () => <svg data-testid="done-all-icon" />,
    DoneIcon: () => <svg data-testid="done-icon" />,
    RemoveDoneIcon: () => <svg data-testid="remove-done-icon" />,
    Modal: ({ open, closeModal, title, children }: any) =>
      open ? (
        <div data-testid="sistent-modal" data-title={title}>
          <button onClick={closeModal} aria-label="modal-close">
            close
          </button>
          {children}
        </div>
      ) : null,
    ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
    useTheme: () => ({
      palette: {
        common: { white: '#fff' },
        warning: { main: '#f80' },
        error: { main: '#f00' },
        icon: { default: '#000' },
        action: { disabledBackground: '#ccc' },
      },
      spacing: (n: number) => `${n * 8}px`,
      shadows: ['none', '1px', '2px', '3px', '4px', '5px', '6px', '7px', '8px'],
    }),
    Search: () => <svg data-testid="search-icon" />,
    styled,
  };
});

describe('ConfirmationMsg', () => {
  beforeEach(() => {
    dispatch.mockReset();
    notify.mockReset();
    triggerPing.mockClear();
    selectorState = {
      ui: { selectedK8sContexts: [], k8sConfig: [] },
    };
  });

  const baseProps = {
    open: true,
    handleClose: vi.fn(),
    submit: { deploy: vi.fn(), unDeploy: vi.fn() },
    title: 'Are you sure?',
    componentCount: 3,
    tab: 2, // deploy
    errors: { validationError: 1, deploymentError: 2 },
  };

  it('renders the modal with the title when open', () => {
    render(<ConfirmationMsg {...baseProps} />);
    expect(screen.getByTestId('sistent-modal')).toHaveAttribute('data-title', 'Are you sure?');
  });

  it('uses default Confirmation title when none provided', () => {
    const props = { ...baseProps };
    delete (props as any).title;
    render(<ConfirmationMsg {...(props as any)} />);
    expect(screen.getByTestId('sistent-modal')).toHaveAttribute('data-title', 'Confirmation');
  });

  it('falls back to K8sEmptyState if there are no contexts', () => {
    selectorState = { ui: { selectedK8sContexts: [], k8sConfig: [] } };
    render(<ConfirmationMsg {...baseProps} />);

    expect(screen.getByTestId('k8s-empty')).toBeInTheDocument();
  });

  it('renders connection chips when k8s contexts exist', () => {
    selectorState = {
      ui: {
        selectedK8sContexts: [],
        k8sConfig: [
          { id: '1', name: 'cluster-1', connectionId: 'c1' },
          { id: '2', name: 'cluster-2', connectionId: 'c2' },
        ],
      },
    };
    render(<ConfirmationMsg {...baseProps} />);

    expect(screen.getByTestId('conn-chip-cluster-1')).toBeInTheDocument();
    expect(screen.getByTestId('conn-chip-cluster-2')).toBeInTheDocument();
  });

  it('handles the validate tab content', () => {
    selectorState = { ui: { selectedK8sContexts: [], k8sConfig: [] } };
    render(<ConfirmationMsg {...baseProps} tab={0} validationBody="Validation passed" />);

    expect(screen.getByTestId('dialog-text')).toHaveTextContent('Validation passed');
  });

  it('triggers ping when a connection chip is clicked', async () => {
    const user = userEvent.setup();
    selectorState = {
      ui: {
        selectedK8sContexts: [],
        k8sConfig: [{ id: '1', name: 'cluster-1', connectionId: 'c1' }],
      },
    };
    render(<ConfirmationMsg {...baseProps} />);

    await user.click(screen.getByTestId('conn-chip-cluster-1'));
    expect(triggerPing).toHaveBeenCalledWith('c1');
  });

  it('handles the "all" checkbox toggle', async () => {
    const user = userEvent.setup();
    selectorState = {
      ui: {
        selectedK8sContexts: [],
        k8sConfig: [{ id: '1', name: 'cluster-1', connectionId: 'c1' }],
      },
    };
    render(<ConfirmationMsg {...baseProps} />);

    const checkboxes = screen.getAllByTestId('checkbox');
    await user.click(checkboxes[0]);
    expect(dispatch).toHaveBeenCalledWith({
      type: 'setK8sContexts',
      payload: { selectedK8sContexts: ['all'] },
    });
  });
});
