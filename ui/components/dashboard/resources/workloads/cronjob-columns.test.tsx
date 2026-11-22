import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';

vi.mock('../../../../utils/k8s-utils', () => ({ timeAgo: (v: string) => `AGO:${v}` }));
vi.mock('../../../../utils/multi-ctx', () => ({
  getK8sContextFromClusterId: () => ({ name: 'ctx' }),
}));
vi.mock('../../../connections/ConnectionChip', () => ({
  TooltipWrappedConnectionChip: () => null,
}));
vi.mock('../sortable-table-cell', () => ({
  DefaultTableCell: () => null,
  SortableTableCell: () => null,
}));
vi.mock('../../../../utils/Enum', () => ({
  CoreConnectionKinds: {
    meshery: 'meshery',
    kubernetes: 'kubernetes',
    prometheus: 'prometheus',
    grafana: 'grafana',
    github: 'github',
  },
}));
vi.mock('@/components/data-formatter', () => ({ FormatId: () => null }));
vi.mock('../../view', () => ({ Title: () => null }));
vi.mock('../config', () => ({ SINGLE_VIEW: 'single' }));

import { buildCronJobColumns } from './cronjob-columns';

describe('buildCronJobColumns', () => {
  it('returns the CronJob name and non-empty columns', () => {
    const cfg = buildCronJobColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'CronJob',
      ping: () => {},
    });
    expect(cfg.name).toBe('CronJob');
    expect(cfg.columns.length).toBeGreaterThan(0);
  });

  const renderSuspend = (attribute: unknown) => {
    const cfg = buildCronJobColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'CronJob',
      ping: () => {},
    });
    const column = cfg.columns.find((c) => c.label === 'Suspend');
    return render(<>{column?.options.customBodyRender(JSON.stringify(attribute))}</>);
  };

  it('shows True for a suspended CronJob', () => {
    renderSuspend({ schedule: '*/5 * * * *', suspend: true });
    expect(screen.getByText('True')).toBeInTheDocument();
  });

  it('shows False when suspend is false or omitted', () => {
    const { unmount } = renderSuspend({ schedule: '*/5 * * * *', suspend: false });
    expect(screen.getByText('False')).toBeInTheDocument();
    unmount();
    renderSuspend({ schedule: '*/5 * * * *' });
    expect(screen.getByText('False')).toBeInTheDocument();
  });
});
