import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import CrossClusterCompatibilityAnalysis from './CrossClusterCompatibilityAnalysis';

// Mock sistent components and theme
vi.mock('@sistent/sistent', async () => {
  const actual = await vi.importActual('@sistent/sistent');
  return {
    ...actual,
    useTheme: () => ({
      palette: {
        mode: 'dark',
        primary: { main: '#00B39F' },
        success: { main: '#4caf50' },
        warning: { main: '#ff9800' },
        error: { main: '#f44336' },
        text: { primary: '#fff', secondary: '#aaa' },
        divider: '#333',
        background: { default: '#121212', paper: '#1e1e1e' },
      },
      shape: { borderRadius: 4 },
      breakpoints: { down: () => '@media (max-width: 900px)' },
    }),
  };
});

vi.mock('@/theme', () => ({
  alpha: (color: string, opacity: number) => `rgba(${color}, ${opacity})`,
  useTheme: () => ({
    palette: {
      mode: 'dark',
      primary: { main: '#00B39F' },
      success: { main: '#4caf50' },
      warning: { main: '#ff9800' },
      error: { main: '#f44336' },
      text: { primary: '#fff', secondary: '#aaa' },
      divider: '#333',
      background: { default: '#121212', paper: '#1e1e1e' },
    },
    shape: { borderRadius: 4 },
    breakpoints: { down: () => '@media (max-width: 900px)' },
  }),
}));

describe('CrossClusterCompatibilityAnalysis', () => {
  const sampleDryRunResponse = {
    frontend: {
      'c-dev': { success: true, component: { kind: 'Deployment' } },
      'c-staging': { success: true, component: { kind: 'Deployment' } },
      'c-prod': { success: true, component: { kind: 'Deployment' } },
    },
    postgres: {
      'c-dev': { success: true, component: { kind: 'StatefulSet' } },
      'c-staging': { success: true, component: { kind: 'StatefulSet' } },
      'c-prod': {
        success: false,
        component: { kind: 'StatefulSet' },
        error: {
          Status: 'error: no matches for kind "PostgresCluster" in version "postgresql.cnpg.io/v1"',
        },
      },
    },
    ingress: {
      'c-dev': { success: true, component: { kind: 'Ingress' } },
      'c-staging': {
        success: false,
        component: { kind: 'Ingress' },
        error: {
          Status: 'Warning: Ingress class not specified',
        },
      },
      'c-prod': { success: true, component: { kind: 'Ingress' } },
    },
  };

  const sampleConnections = [
    { id: 'c-dev', metadata: { name: 'dev-cluster', version: 'v1.29.2' } },
    { id: 'c-staging', metadata: { name: 'staging-cluster', version: 'v1.28.5' } },
    { id: 'c-prod', metadata: { name: 'prod-cluster', version: 'v1.27.9' } },
  ];

  it('renders the header and metric cards correctly', () => {
    render(
      <CrossClusterCompatibilityAnalysis
        dryRunResponse={sampleDryRunResponse}
        selectedK8sConnections={sampleConnections}
      />,
    );

    expect(screen.getByText('Cross-Cluster Compatibility Analysis')).toBeInTheDocument();
    expect(screen.getByText('Clusters Analyzed')).toBeInTheDocument();
    expect(screen.getByText('Resources Checked')).toBeInTheDocument();
    expect(screen.getByText('Succeeded')).toBeInTheDocument();
    expect(screen.getByText('Warnings')).toBeInTheDocument();
    expect(screen.getAllByText('Failed').length).toBeGreaterThan(0);

    // Check metric numbers
    expect(screen.getAllByText('3').length).toBeGreaterThan(0);
  });

  it('renders the compatibility matrix table with resources and clusters', () => {
    render(
      <CrossClusterCompatibilityAnalysis
        dryRunResponse={sampleDryRunResponse}
        selectedK8sConnections={sampleConnections}
      />,
    );

    expect(screen.getByText('Compatibility Matrix')).toBeInTheDocument();
    expect(screen.getByText('frontend')).toBeInTheDocument();
    expect(screen.getByText('postgres')).toBeInTheDocument();
    expect(screen.getByText('ingress')).toBeInTheDocument();

    expect(screen.getAllByText('dev-cluster').length).toBeGreaterThan(0);
    expect(screen.getAllByText('staging-cluster').length).toBeGreaterThan(0);
    expect(screen.getAllByText('prod-cluster').length).toBeGreaterThan(0);
  });

  it('displays detailed error information, possible fix, and dependency impact when selecting an error', () => {
    render(
      <CrossClusterCompatibilityAnalysis
        dryRunResponse={sampleDryRunResponse}
        selectedK8sConnections={sampleConnections}
      />,
    );

    // Click on failed badge
    const failedBadges = screen.getAllByTestId('status-badge-failed');
    fireEvent.click(failedBadges[0]);

    // Validation Details panel should show error details
    const detailsPanel = screen.getByTestId('validation-details-panel');
    expect(detailsPanel).toBeInTheDocument();
    expect(screen.getByText(/Validation Details/)).toBeInTheDocument();
    expect(screen.getByText(/no matches for kind "PostgresCluster"/)).toBeInTheDocument();
    expect(screen.getByText('Possible Fix')).toBeInTheDocument();
    expect(screen.getByText('Dependency Impact')).toBeInTheDocument();
    expect(screen.getByText(/Custom Resource Definition \(CRD\)/)).toBeInTheDocument();
    expect(screen.getByText(/database connectivity/)).toBeInTheDocument();
  });

  it('triggers onReRunAnalysis when re-run button is clicked', () => {
    const handleReRun = vi.fn();
    render(
      <CrossClusterCompatibilityAnalysis
        dryRunResponse={sampleDryRunResponse}
        selectedK8sConnections={sampleConnections}
        onReRunAnalysis={handleReRun}
      />,
    );

    const reRunBtn = screen.getByTestId('rerun-analysis-btn');
    fireEvent.click(reRunBtn);

    expect(handleReRun).toHaveBeenCalledTimes(1);
  });
});
