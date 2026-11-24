import React from 'react';
import { render, screen, fireEvent, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
  Card: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
  Chip: ({ label, ...rest }: any) => <span {...rest}>{label}</span>,
  DialogContent: ({ children }: any) => <div>{children}</div>,
  ModalButtonPrimary: ({ children, onClick, disabled }: any) => (
    <button type="button" onClick={onClick} disabled={disabled} data-testid="btn-proceed">
      {children}
    </button>
  ),
  ModalButtonSecondary: ({ children, onClick }: any) => (
    <button type="button" onClick={onClick} data-testid="btn-cancel">
      {children}
    </button>
  ),
  Typography: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('@/components/shared/Modal', () => ({
  Modal: ({ children, isOpen, title, actions }: any) =>
    isOpen ? (
      <div data-testid="impact-modal">
        <h2>{title}</h2>
        <div>{children}</div>
        <div>{actions}</div>
      </div>
    ) : null,
}));

vi.mock('@/assets/icons/AlertIcon', () => ({
  default: () => <svg data-testid="alert-icon" />,
}));

vi.mock('@/assets/icons/CheckIcon', () => ({
  default: () => <svg data-testid="check-icon" />,
}));

vi.mock('@/theme', () => ({
  useTheme: () => ({
    palette: {
      warning: { main: '#ffa000' },
      success: { main: '#4caf50' },
      background: { paper: '#fff', default: '#fafafa' },
      divider: '#eee',
    },
  }),
}));

import DesignImpactPreviewModal from './DesignImpactPreviewModal';
import { DesignImpactAnalysisResult } from './designImpactAnalysis';

describe('DesignImpactPreviewModal', () => {
  const dummyResult: DesignImpactAnalysisResult = {
    targetType: 'component',
    targetId: 'comp-1',
    targetName: 'frontend-svc',
    targetKind: 'Service',
    modificationType: 'update',
    directlyAffected: [
      {
        id: 'comp-2',
        name: 'frontend-deploy',
        kind: 'Deployment',
        model: 'kubernetes',
        explanation: 'Directly connected via Edge with Service/frontend-svc',
      },
    ],
    indirectlyAffected: [
      {
        id: 'comp-3',
        name: 'frontend-pvc',
        kind: 'PersistentVolumeClaim',
        model: 'kubernetes',
        explanation: 'Transitively dependent through Deployment/frontend-deploy',
      },
    ],
    affectedRelationships: [
      {
        id: 'rel-1',
        kind: 'Edge',
        type: 'network',
        subType: 'ServiceToDeployment',
      },
    ],
    totalAffectedCount: 2,
    summary:
      'Modifying Service/frontend-svc potentially affects 1 directly connected and 1 indirect component(s).',
  };

  it('renders modal with target details and affected components', () => {
    render(
      <DesignImpactPreviewModal
        isOpen
        onClose={vi.fn()}
        onProceed={vi.fn()}
        impactResult={dummyResult}
      />,
    );

    expect(screen.getByTestId('impact-modal')).toBeInTheDocument();
    expect(screen.getByText(/Design Impact Preview: Modify frontend-svc/)).toBeInTheDocument();
    expect(screen.getByText(/Service: frontend-svc/)).toBeInTheDocument();
    expect(screen.getByText('UPDATE')).toBeInTheDocument();
    // Scope to the directly-affected section: the PVC explanation in the
    // indirect section repeats the 'Deployment/frontend-deploy' label.
    const directlyAffectedSection = screen.getByText(/Directly Affected Components/)
      .parentElement as HTMLElement;
    expect(
      within(directlyAffectedSection).getByText(/Deployment\/frontend-deploy/),
    ).toBeInTheDocument();
    expect(screen.getByText(/PersistentVolumeClaim\/frontend-pvc/)).toBeInTheDocument();
  });

  it('invokes onProceed when Continue is clicked', () => {
    const onProceed = vi.fn();
    render(
      <DesignImpactPreviewModal
        isOpen
        onClose={vi.fn()}
        onProceed={onProceed}
        impactResult={dummyResult}
      />,
    );

    fireEvent.click(screen.getByTestId('btn-proceed'));
    expect(onProceed).toHaveBeenCalledTimes(1);
  });

  it('invokes onClose when Cancel is clicked', () => {
    const onClose = vi.fn();
    render(
      <DesignImpactPreviewModal
        isOpen
        onClose={onClose}
        onProceed={vi.fn()}
        impactResult={dummyResult}
      />,
    );

    fireEvent.click(screen.getByTestId('btn-cancel'));
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('renders message when no components or relationships are affected', () => {
    const emptyResult: DesignImpactAnalysisResult = {
      targetType: 'component',
      targetId: 'isolated',
      targetName: 'isolated',
      targetKind: 'Pod',
      modificationType: 'delete',
      directlyAffected: [],
      indirectlyAffected: [],
      affectedRelationships: [],
      totalAffectedCount: 0,
      summary: 'No dependent components or relationships found for Pod/isolated.',
    };

    render(
      <DesignImpactPreviewModal
        isOpen
        onClose={vi.fn()}
        onProceed={vi.fn()}
        impactResult={emptyResult}
      />,
    );

    expect(
      screen.getByText(
        'No dependencies or connected relationships will be disrupted by this change.',
      ),
    ).toBeInTheDocument();
  });
});
