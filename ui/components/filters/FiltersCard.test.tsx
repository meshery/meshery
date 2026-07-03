import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const can = vi.fn(() => true);
const getUserByIdQuery = vi.fn();

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => can(...args),
}));

vi.mock('@/utils/permission_constants', () => ({
  keys: {
    PUBLISH_WASM_FILTER: { action: 'publish', subject: 'wasm-filter' },
    UNPUBLISH_WASM_FILTER: { action: 'unpublish', subject: 'wasm-filter' },
    DOWNLOAD_A_WASM_FILTER: { action: 'download', subject: 'wasm-filter' },
    CLONE_WASM_FILTER: { action: 'clone', subject: 'wasm-filter' },
    DETAILS_OF_WASM_FILTER: { action: 'details', subject: 'wasm-filter' },
    EDIT_WASM_FILTER: { action: 'edit', subject: 'wasm-filter' },
    DELETE_WASM_FILTER: { action: 'delete', subject: 'wasm-filter' },
  },
}));

vi.mock('../../rtk-query/user', () => ({
  useGetUserByIdQuery: (...args: unknown[]) => getUserByIdQuery(...args),
}));

vi.mock('../../utils/Enum', () => ({
  VISIBILITY: { PUBLISHED: 'published', PUBLIC: 'public', PRIVATE: 'private' },
}));

vi.mock('../../constants/endpoints', () => ({
  MESHERY_CLOUD_PROD: 'https://meshery.io',
}));

vi.mock('../shared/Modal/Information/InfoModal', () => ({
  VIEW_VISIBILITY: { PUBLIC: 'public', PRIVATE: 'private', PUBLISHED: 'published' },
}));

vi.mock('css/icons.styles', () => ({ iconMedium: {} }));

vi.mock('@sistent/sistent', () => {
  const styled = (_Component: any) => (_factory?: any) => {
    const Styled = ({ children }: any) => <div>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    Divider: () => <hr />,
    Grid2: ({ children }: any) => <div>{children}</div>,
    IconButton: ({ children, onClick, disabled }: any) => (
      <button onClick={onClick} disabled={disabled}>
        {children}
      </button>
    ),
    Typography: ({ children }: any) => <span>{children}</span>,
    Tooltip: ({ children, title }: any) => <div data-tip={title}>{children}</div>,
    Link: ({ children }: any) => <a>{children}</a>,
    Avatar: ({ src }: any) => <div data-testid="avatar" data-src={src} />,
    useTheme: () => ({
      palette: {
        icon: { default: 'icon-default' },
        background: { constant: { white: 'white' } },
      },
    }),
    VisibilityChipMenu: ({ value }: any) => <div data-testid="visibility-chip">{value}</div>,
    InfoOutlinedIcon: () => <svg data-testid="info-outlined-icon" />,
    FullScreenIcon: () => <svg data-testid="fullscreen-icon" />,
    FullScreenExitIcon: () => <svg data-testid="fullscreen-exit-icon" />,
    styled,
  };
});

vi.mock('@/assets/icons', () => ({
  Delete: () => <svg data-testid="delete-icon" />,
  Save: () => <svg data-testid="save-icon" />,
  Public: () => <svg data-testid="public-icon" />,
  GetApp: () => <svg data-testid="get-app-icon" />,
  Lock: () => <svg data-testid="lock-icon" />,
}));

vi.mock('react-moment', () => ({
  default: ({ children }: any) => <span data-testid="moment">{String(children)}</span>,
}));

vi.mock('../FlipCard', () => ({
  default: ({ children }: any) => (
    <div data-testid="flip-card">
      {/* React's children for the FlipCard are an array of two parts (front, back); render both */}
      {Array.isArray(children)
        ? children.map((c: any, i: number) => <div key={i}>{c}</div>)
        : children}
    </div>
  ),
}));

vi.mock('../CodeMirror', () => ({
  UnControlled: ({ value }: any) => <pre data-testid="codemirror">{value}</pre>,
}));

vi.mock('../designs/patterns/Cards.styles', () => ({
  BottomContainer: ({ children }: any) => <div>{children}</div>,
  CardBackGrid: ({ children }: any) => <div>{children}</div>,
  CatalogCardButtons: ({ children }: any) => <div>{children}</div>,
  UpdateDeleteButtons: ({ children }: any) => <div>{children}</div>,
  YamlDialogTitleGrid: ({ children }: any) => <div>{children}</div>,
  GridBtnText: ({ children }: any) => <span>{children}</span>,
  GridCloneBtnText: ({ children }: any) => <span>{children}</span>,
  CardHeaderRight: ({ children }: any) => <div>{children}</div>,
  StyledCodeMirrorWrapper: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('../YamlDialog', () => ({
  default: ({ name }: any) => <div data-testid="yaml-dialog">{name}</div>,
}));

vi.mock('../../public/static/img/CloneIcon', () => ({
  default: () => <svg data-testid="clone-icon" />,
}));

vi.mock('../../utils/TooltipButton', () => ({
  default: ({ children, onClick, title, disabled }: any) => (
    <button onClick={onClick} disabled={disabled} title={title} data-testid={`btn-${title}`}>
      {children}
    </button>
  ),
}));

import FiltersCard from './FiltersCard';

describe('FiltersCard', () => {
  beforeEach(() => {
    can.mockReset();
    can.mockReturnValue(true);
    getUserByIdQuery.mockReset();
    getUserByIdQuery.mockReturnValue({ data: { avatarUrl: 'https://a.io/u.png' } });
  });

  const baseProps = {
    name: 'my filter',
    updatedAt: '2026-01-01T00:00:00Z',
    createdAt: '2025-12-01T00:00:00Z',
    filterResource: 'foo: bar',
    handleClone: vi.fn(),
    handleDownload: vi.fn(),
    deleteHandler: vi.fn(),
    setYaml: vi.fn(),
    visibility: 'public',
    handlePublishModal: vi.fn(),
    handleUnpublishModal: vi.fn(),
    updateHandler: vi.fn(),
    canPublishFilter: true,
    handleInfoModal: vi.fn(),
    ownerId: 'u-1',
  };

  it('renders the filter name and "Modified On" timestamp', () => {
    render(<FiltersCard {...baseProps} />);
    // Name renders on both front and back of the flip card
    expect(screen.getAllByText('my filter').length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText(/Modified On:/)).toBeInTheDocument();
  });

  it('shows Publish button for non-published visibility when canPublishFilter is true', () => {
    render(<FiltersCard {...baseProps} visibility="public" />);
    expect(screen.getByTestId('btn-Publish')).toBeInTheDocument();
  });

  it('shows Unpublish button for published visibility', () => {
    render(<FiltersCard {...baseProps} visibility="published" />);
    expect(screen.getByTestId('btn-Unpublish')).toBeInTheDocument();
  });

  it('shows the Clone button only for published filters', () => {
    const { rerender } = render(<FiltersCard {...baseProps} visibility="public" />);
    expect(screen.queryByTestId('btn-Clone')).not.toBeInTheDocument();

    rerender(<FiltersCard {...baseProps} visibility="published" />);
    expect(screen.getByTestId('btn-Clone')).toBeInTheDocument();
  });

  it('invokes handleDownload when Download button is clicked', async () => {
    const user = userEvent.setup();
    const handleDownload = vi.fn();
    render(<FiltersCard {...baseProps} handleDownload={handleDownload} />);

    await user.click(screen.getByTestId('btn-Download'));
    expect(handleDownload).toHaveBeenCalled();
  });

  it('invokes handlePublishModal via genericClickHandler (stopPropagation)', async () => {
    const user = userEvent.setup();
    const handlePublishModal = vi.fn();
    render(<FiltersCard {...baseProps} handlePublishModal={handlePublishModal} />);

    await user.click(screen.getByTestId('btn-Publish'));
    expect(handlePublishModal).toHaveBeenCalled();
  });

  it('invokes handleInfoModal when Info button is clicked', async () => {
    const user = userEvent.setup();
    const handleInfoModal = vi.fn();
    render(<FiltersCard {...baseProps} handleInfoModal={handleInfoModal} />);
    await user.click(screen.getByTestId('btn-Filter Information'));
    expect(handleInfoModal).toHaveBeenCalled();
  });

  it('disables actions when permission CAN returns false', () => {
    can.mockReturnValue(false);
    render(<FiltersCard {...baseProps} visibility="published" />);
    expect(screen.getByTestId('btn-Download')).toBeDisabled();
    expect(screen.getByTestId('btn-Clone')).toBeDisabled();
    expect(screen.getByTestId('btn-Filter Information')).toBeDisabled();
  });

  it('passes ownerId to useGetUserByIdQuery', () => {
    render(<FiltersCard {...baseProps} ownerId="user-42" />);
    expect(getUserByIdQuery).toHaveBeenCalledWith('user-42');
  });

  it('renders the avatar with the fetched owner avatarUrl', () => {
    render(<FiltersCard {...baseProps} />);
    expect(screen.getByTestId('avatar')).toHaveAttribute('data-src', 'https://a.io/u.png');
  });
});
