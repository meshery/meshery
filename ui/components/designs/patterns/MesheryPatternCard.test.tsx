import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('@/rtk-query/user', () => ({
  useGetUserByIdQuery: () => ({ data: {} }),
  useGetProviderCapabilitiesQuery: () => ({ data: {} }),
}));

vi.mock('@/utils/Enum', () => ({
  VISIBILITY: { PUBLISHED: 'published', PUBLIC: 'public', PRIVATE: 'private' },
}));

vi.mock('../../../constants/endpoints', () => ({
  MESHERY_CLOUD_PROD: 'https://meshery.io',
}));

vi.mock('../../shared/Modal/Information/InfoModal', () => ({
  VIEW_VISIBILITY: { PUBLIC: 'public', PRIVATE: 'private', PUBLISHED: 'published' },
}));

vi.mock('css/icons.styles', () => ({ iconMedium: {}, iconLarge: {} }));

vi.mock('@sistent/sistent', () => {
  const styled = (_Component: any) => (_factory?: any) => {
    const Styled = ({ children, ...props }: any) => <div {...props}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    Divider: () => <hr />,
    Grid2: ({ children }: any) => <div>{children}</div>,
    IconButton: ({ children, onClick, ...props }: any) => (
      <button type="button" onClick={onClick} {...props}>
        {children}
      </button>
    ),
    Typography: ({ children }: any) => <span>{children}</span>,
    CustomTooltip: ({ children }: any) => <div>{children}</div>,
    Link: ({ children }: any) => <a>{children}</a>,
    Avatar: () => <div data-testid="avatar" />,
    useTheme: () => ({
      palette: {
        icon: { default: 'icon-default' },
        background: { constant: { white: 'white' } },
      },
    }),
    useHasPermission: () => true,
    VisibilityChipMenu: ({ value }: any) => <div data-testid="visibility-chip">{value}</div>,
    InfoOutlinedIcon: () => <svg data-testid="info-outlined-icon" />,
    FullScreenIcon: () => <svg data-testid="fullscreen-icon" />,
    FullScreenExitIcon: () => <svg data-testid="fullscreen-exit-icon" />,
    AccountTreeIcon: () => <svg data-testid="account-tree-icon" />,
    crimson: { 40: '#F91313' },
    styled,
  };
});

vi.mock('@/assets/icons', () => ({
  Delete: () => <svg data-testid="delete-icon" />,
  Save: () => <svg data-testid="save-icon" />,
  Public: () => <svg data-testid="public-icon" />,
  GetApp: () => <svg data-testid="get-app-icon" />,
  DoneAll: () => <svg data-testid="done-all-icon" />,
  Lock: () => <svg data-testid="lock-icon" />,
  Edit: () => <svg data-testid="edit-icon" />,
}));

vi.mock('@/assets/icons/DryRunIcon', () => ({
  default: () => <svg data-testid="dryrun-icon" />,
}));

vi.mock('@/assets/icons/CheckIcon', () => ({
  default: () => <svg data-testid="check-icon" />,
}));

vi.mock('@/assets/icons/Pattern', () => ({
  default: () => <svg data-testid="pattern-icon" />,
}));

vi.mock('../../../public/static/img/UndeployIcon', () => ({
  default: () => <svg data-testid="undeploy-icon" />,
}));

vi.mock('../../../public/static/img/CloneIcon', () => ({
  default: () => <svg data-testid="clone-icon" />,
}));

vi.mock('react-moment', () => ({
  default: ({ children }: any) => <span data-testid="moment">{String(children)}</span>,
}));

vi.mock('../../general/FlipCard', () => ({
  default: ({ children }: any) => (
    <div data-testid="flip-card">
      {Array.isArray(children)
        ? children.map((c: any, i: number) => <div key={i}>{c}</div>)
        : children}
    </div>
  ),
}));

vi.mock('../../general/CodeMirror', () => ({
  UnControlled: ({ value }: any) => <pre data-testid="codemirror">{value}</pre>,
}));

vi.mock('../../general/YamlDialog', () => ({
  default: () => null,
}));

vi.mock('./Cards.styles', () => ({
  BottomContainer: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  CardBackGrid: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  CatalogCardButtons: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  UpdateDeleteButtons: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  YamlDialogTitleGrid: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  CardHeaderRight: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  GridBtnText: ({ children, ...props }: any) => <span>{children}</span>,
  GridCloneBtnText: ({ children, ...props }: any) => <span>{children}</span>,
  StyledCodeMirrorWrapper: ({ children, ...props }: any) => <div {...props}>{children}</div>,
}));

vi.mock('@/utils/TooltipButton', () => ({
  default: ({ children, onClick, 'data-testid': testId }: any) => (
    <button data-testid={testId} onClick={onClick}>
      {children}
    </button>
  ),
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: vi.fn(), query: {} }),
}));

vi.mock('@/utils/hooks/useResourceOwner', () => ({
  useResourceOwner: () => ({ owner: { name: 'Test User' }, hasCloudProfile: false }),
}));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: {
    CatalogManagementDeleteADesign: 'del-design',
    CatalogManagementEditDesign: 'edit-design',
    CatalogManagementValidateDesign: 'val-design',
    CatalogManagementDeployDesign: 'dep-design',
    CatalogManagementUndeployDesign: 'undep-design',
    CatalogManagementEvaluateRelationships: 'eval-rel',
    CatalogManagementUnpublishDesign: 'unpub-design',
  },
}));

vi.mock('./design-permissions', () => ({
  canEditDesign: () => true,
}));

vi.mock('./ActionButton', () => ({
  default: ({ options, 'data-testid': testId }: any) => (
    <div data-testid={testId || 'action-btn-mock'}>
      {options?.map((opt: any) => (
        <button
          key={opt.label}
          data-testid={opt['data-testid'] || `action-btn-option-${opt.label}`}
          onClick={opt.onClick}
        >
          {opt.label}
        </button>
      ))}
    </div>
  ),
}));

import MesheryPatternCard from './MesheryPatternCard';

describe('MesheryPatternCard', () => {
  const defaultProps = {
    id: 'pattern-1',
    name: 'Sample Pattern',
    updatedAt: '2024-01-01',
    createdAt: '2024-01-01',
    patternFile: '{"services":{}}',
    pattern: { id: 'pattern-1', userId: 'user-1', name: 'Sample Pattern' },
    deleteHandler: vi.fn(),
    handleVerify: vi.fn(),
    handleDryRun: vi.fn(),
    handleUnpublishModal: vi.fn(),
    handleDeploy: vi.fn(),
    handleUnDeploy: vi.fn(),
    handleDownload: vi.fn(),
    updateHandler: vi.fn(),
    handleClone: vi.fn(),
    setSelectedPatterns: vi.fn(),
    setYaml: vi.fn(),
    handleInfoModal: vi.fn(),
    handleEvaluate: vi.fn(),
  };

  it('renders Delete action in front ActionButton dropdown when not read-only and invokes deleteHandler', async () => {
    const user = userEvent.setup();
    const deleteHandler = vi.fn();

    render(
      <MesheryPatternCard {...defaultProps} isReadOnly={false} deleteHandler={deleteHandler} />,
    );

    const deleteOption = screen.getByTestId('pattern-btn-action-delete');
    expect(deleteOption).toBeInTheDocument();
    expect(deleteOption).toHaveTextContent('Delete');

    await user.click(deleteOption);
    expect(deleteHandler).toHaveBeenCalledTimes(1);
  });

  it('omits Delete action from front ActionButton dropdown when isReadOnly is true', () => {
    render(<MesheryPatternCard {...defaultProps} isReadOnly={true} />);

    expect(screen.queryByTestId('pattern-btn-action-delete')).not.toBeInTheDocument();
  });
});
