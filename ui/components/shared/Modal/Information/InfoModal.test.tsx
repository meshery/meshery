import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import InfoModal from './InfoModal';

const notify = vi.fn();
const enqueueSnackbar = vi.fn();
const closeSnackbar = vi.fn();
const updatePatternMock = vi.fn().mockReturnValue({
  then: (cb: any) => {
    cb();
    return { catch: vi.fn() };
  },
});
const publishMock = vi.fn().mockReturnValue({
  unwrap: () =>
    Promise.resolve({
      then: (cb: any) => cb(),
      catch: vi.fn(),
    }),
});
const mockGetUserById = vi.fn();
const mockGetMeshModels = vi.fn();

let mockState: any = { ui: { user: { id: 'u1', roleNames: ['admin'] } } };

vi.mock('react-redux', () => ({
  useSelector: (fn: any) => fn(mockState),
}));

vi.mock('@/rtk-query/design', () => ({
  usePublishPatternMutation: () => [publishMock],
  useUpdatePatternFileMutation: () => [updatePatternMock],
}));

vi.mock('@/rtk-query/user', () => ({
  useGetUserByIdQuery: (id: string) => mockGetUserById(id),
}));

vi.mock('@/api/meshmodel', () => ({
  getMeshModels: () => mockGetMeshModels(),
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('@/store/ProviderStoreWrapper', () => ({
  default: ({ children }: any) => <div data-testid="provider-wrapper">{children}</div>,
}));

vi.mock('@/store/slices/mesheryUi', () => ({
  updateProgress: vi.fn(),
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('@/utils/objects', () => ({
  filterEmptyFields: (obj: any) => obj || {},
}));

vi.mock('@/utils/TooltipButton', () => ({
  default: ({ children, title, onClick }: any) => (
    <span data-testid="tooltip-btn" data-title={title} onClick={onClick}>
      {children}
    </span>
  ),
}));

vi.mock('../../../../utils/utils', () => ({
  getDesignVersion: (r: any) => r?.version || '',
  getSharableCommonHostAndprotocolLink: (r: any) => `https://meshery.io/${r?.id}`,
  modifyRJSFSchema: (s: any) => s,
}));

vi.mock('../../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('../../../../rtk-query/user', () => ({
  useGetUserByIdQuery: (id: string) => mockGetUserById(id),
}));

vi.mock('../../../../lib/event-types', () => ({
  EVENT_TYPES: { SUCCESS: 'success', ERROR: 'error' },
}));

vi.mock('../../../../constants/endpoints', () => ({
  MESHERY_CLOUD_PROD: 'https://cloud.meshery.io',
}));

vi.mock('../../../../css/icons.styles', () => ({
  iconMedium: {},
  iconSmall: {},
}));

vi.mock('../../../../assets/icons/Pattern', () => ({
  default: () => <svg data-testid="pattern-icon" />,
}));

vi.mock('@/assets/icons/ServiceMesheryIcon', () => ({
  default: () => <svg data-testid="service-icon" />,
}));

vi.mock('@/assets/icons', () => ({
  Close: () => <svg data-testid="close-icon" />,
  Lock: () => <svg data-testid="lock-icon" />,
  Public: () => <svg data-testid="public-icon" />,
}));

vi.mock('@/components/meshery-mesh-interface/PatternService/RJSF_wrapper', () => ({
  default: ({ formData, onChange }: any) => (
    <div data-testid="rjsf">
      <button
        onClick={() => onChange?.({ ...(formData || {}), type: 'Pattern' })}
        aria-label="form-change"
      >
        change
      </button>
    </div>
  ),
}));

vi.mock('js-yaml', () => ({
  default: { load: (s: string) => ({ loaded: s }) },
}));

vi.mock('notistack', () => ({
  useSnackbar: () => ({ enqueueSnackbar, closeSnackbar }),
}));

vi.mock('./styles', () => ({
  ActionContainer: ({ children }: any) => <div data-testid="action-container">{children}</div>,
  CopyLinkButton: ({ children }: any) => <span>{children}</span>,
  CreatAtContainer: ({ children }: any) => <span>{children}</span>,
  ResourceName: ({ children }: any) => <span data-testid="resource-name">{children}</span>,
}));

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = (props: any) => <Component {...props}>{props.children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    Avatar: ({ src }: any) => <img data-testid="avatar" alt="avatar" src={src} />,
    Box: ({ children, style }: any) => <div style={style}>{children}</div>,
    Button: ({ children, onClick, disabled, ...rest }: any) => (
      <button onClick={onClick} disabled={disabled} {...rest}>
        {children}
      </button>
    ),
    CircularProgress: () => <div data-testid="progress" />,
    CustomTooltip: ({ children, title }: any) => (
      <span data-tooltip={String(title)}>{children}</span>
    ),
    getFormatDate: (d: string) => `date:${d}`,
    Grid: ({ children }: any) => <div>{children}</div>,
    IconButton: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
    Link: ({ children, href }: any) => <a href={href}>{children}</a>,
    Modal: ({ open, closeModal, title, children }: any) =>
      open ? (
        <div data-testid="modal" data-title={title}>
          <button onClick={closeModal} aria-label="modal-close">
            close
          </button>
          {children}
        </div>
      ) : null,
    ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
    ModalButtonPrimary: ({ children, onClick, disabled }: any) => (
      <button onClick={onClick} disabled={disabled} data-testid="save-btn">
        {children}
      </button>
    ),
    ModalButtonSecondary: ({ children, onClick, disabled }: any) => (
      <button onClick={onClick} disabled={disabled} data-testid="publish-btn">
        {children}
      </button>
    ),
    ModalFooter: ({ children }: any) => <div data-testid="modal-footer">{children}</div>,
    publishCatalogItemSchema: { type: 'object' },
    publishCatalogItemUiSchema: {},
    Skeleton: () => <div data-testid="skeleton" />,
    Typography: ({ children, variant }: any) => <span data-variant={variant}>{children}</span>,
    VisibilityChipMenu: ({ value, onChange, enabled }: any) => (
      <button
        data-testid="visibility-menu"
        data-value={value}
        data-enabled={String(enabled)}
        onClick={() => onChange('private')}
      >
        {value}
      </button>
    ),
    useTheme: () => ({
      palette: {
        common: { white: '#fff' },
        divider: '#ccc',
      },
      shadows: ['none', 's1', 's2', 's3', 's4', 's5', 's6', 's7', 's8'],
    }),
    styled,
  };
});

const baseSelectedResource = {
  id: 'p1',
  name: 'My Design',
  userId: 'u1',
  visibility: 'public',
  catalogData: {
    pattern_caveats: 'foo',
    pattern_info: 'bar',
    type: 'pattern',
    compatibility: ['kubernetes'],
  },
  patternFile: 'name: x\n',
  createdAt: '2024-01-01',
  updatedAt: '2024-02-01',
};

describe('InfoModal', () => {
  beforeEach(() => {
    notify.mockReset();
    enqueueSnackbar.mockReset();
    closeSnackbar.mockReset();
    mockGetUserById.mockReturnValue({
      data: { id: 'user-1', firstName: 'Bob', lastName: 'Jones', avatarUrl: '/b.png' },
    });
    mockGetMeshModels.mockResolvedValue({
      models: [{ name: 'kubernetes', displayName: 'Kubernetes' }],
    });
    mockState = { ui: { user: { id: 'u1', roleNames: ['admin'] } } };
  });

  it('renders the InfoModal with the resource name as title', () => {
    render(
      <InfoModal
        infoModalOpen={true}
        handleInfoModalClose={vi.fn()}
        resourceOwnerID="u1"
        selectedResource={baseSelectedResource}
        patternFetcher={vi.fn()}
      />,
    );

    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'My Design');
    expect(screen.getByTestId('provider-wrapper')).toBeInTheDocument();
  });

  it('renders nothing when modal is closed', () => {
    render(
      <InfoModal
        infoModalOpen={false}
        handleInfoModalClose={vi.fn()}
        resourceOwnerID="u1"
        selectedResource={baseSelectedResource}
        patternFetcher={vi.fn()}
      />,
    );

    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('shows the resource name in the content', () => {
    render(
      <InfoModal
        infoModalOpen={true}
        handleInfoModalClose={vi.fn()}
        resourceOwnerID="u1"
        selectedResource={baseSelectedResource}
        patternFetcher={vi.fn()}
      />,
    );

    expect(screen.getByTestId('resource-name')).toHaveTextContent('My Design');
  });

  it('renders the publish to catalog button', () => {
    render(
      <InfoModal
        infoModalOpen={true}
        handleInfoModalClose={vi.fn()}
        resourceOwnerID="u1"
        selectedResource={baseSelectedResource}
        patternFetcher={vi.fn()}
      />,
    );

    expect(screen.getByTestId('publish-btn')).toHaveTextContent('Publish to Catalog');
  });

  it('shows Published when the resource is already published', () => {
    render(
      <InfoModal
        infoModalOpen={true}
        handleInfoModalClose={vi.fn()}
        resourceOwnerID="u1"
        selectedResource={{ ...baseSelectedResource, visibility: 'published' }}
        patternFetcher={vi.fn()}
      />,
    );

    expect(screen.getByTestId('publish-btn')).toHaveTextContent('Published');
  });

  it('renders the visibility menu', () => {
    render(
      <InfoModal
        infoModalOpen={true}
        handleInfoModalClose={vi.fn()}
        resourceOwnerID="u1"
        selectedResource={baseSelectedResource}
        patternFetcher={vi.fn()}
      />,
    );

    expect(screen.getByTestId('visibility-menu')).toHaveAttribute('data-value', 'public');
  });

  it('copies the link via the tooltip button', async () => {
    const user = userEvent.setup();
    const writeText = vi.fn();
    Object.defineProperty(navigator, 'clipboard', {
      value: { writeText },
      writable: true,
      configurable: true,
    });

    render(
      <InfoModal
        infoModalOpen={true}
        handleInfoModalClose={vi.fn()}
        resourceOwnerID="u1"
        selectedResource={baseSelectedResource}
        patternFetcher={vi.fn()}
      />,
    );

    await user.click(screen.getByTestId('tooltip-btn'));
    expect(writeText).toHaveBeenCalledWith('https://meshery.io/p1');
    expect(enqueueSnackbar).toHaveBeenCalled();
  });

  it('shows the design image fallback when imageURL is not present', () => {
    render(
      <InfoModal
        infoModalOpen={true}
        handleInfoModalClose={vi.fn()}
        resourceOwnerID="u1"
        selectedResource={{
          ...baseSelectedResource,
          catalogData: { ...baseSelectedResource.catalogData, imageURL: undefined },
        }}
        patternFetcher={vi.fn()}
      />,
    );

    expect(screen.getByTestId('service-icon')).toBeInTheDocument();
  });

  it('renders the design image when imageURL is provided', () => {
    render(
      <InfoModal
        infoModalOpen={true}
        handleInfoModalClose={vi.fn()}
        resourceOwnerID="u1"
        selectedResource={{
          ...baseSelectedResource,
          catalogData: {
            ...baseSelectedResource.catalogData,
            imageURL: ['/design-1.png'],
          },
        }}
        patternFetcher={vi.fn()}
      />,
    );

    expect(screen.getByAltText('My Design')).toHaveAttribute('src', '/design-1.png');
  });
});
