/**
 * Unit tests for exportDryRunReport utility
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  buildDryRunReport,
  exportDryRunReport,
  sanitizeData,
  DryRunErrorEntry,
} from './exportDryRunReport';

// Mock file downloader
vi.mock('@/utils/fileDownloader', () => ({
  downloadFileFromContent: vi.fn(),
}));
vi.mock('./common', () => ({
  processDesign: (design: any = {}) => {
    const components = Array.isArray(design?.components) ? design.components : [];
    return {
      configurableComponents: components.filter(
        (component: any) => !component?.metadata?.isAnnotation,
      ),
    };
  },
}));

import { downloadFileFromContent } from '@/utils/fileDownloader';

const mockDesign = {
  id: 'ecommerce-app-id',
  name: 'ecommerce-app',
  schemaVersion: 'designs.meshery.io/v1beta1',
  components: [
    {
      name: 'frontend',
      component: { kind: 'Deployment', version: 'apps/v1' },
      metadata: { isAnnotation: false },
    },
    {
      name: 'service',
      component: { kind: 'Service', version: 'v1' },
      metadata: { isAnnotation: false },
    },
    {
      name: 'redis',
      component: { kind: 'StatefulSet', version: 'apps/v1' },
      metadata: { isAnnotation: false },
    },
    {
      name: 'postgres',
      component: { kind: 'StatefulSet', version: 'apps/v1' },
      metadata: { isAnnotation: false },
    },
  ],
};

const mockDryRunErrorsWithFailures: DryRunErrorEntry[] = [
  {
    type: 'ComponentError',
    compName: 'postgres',
    contextId: 'ctx-1',
    component: {
      component: {
        kind: 'StorageClass',
        version: 'storage.k8s.io/v1',
      },
    },
    errors: [
      {
        Type: 'ComponentError',
        FieldPath: 'spec.storageClassName',
        Message: 'StorageClass "fast-ssd" not found',
      },
    ],
  },
];

describe('sanitizeData', () => {
  it('redacts sensitive keys from objects', () => {
    const input = {
      name: 'cluster-config',
      password: 'super-secret-password',
      authToken: 'eyJhbGciOi...',
      secret_key: 'my-key',
      normalField: 'all-good',
    };
    const sanitized = sanitizeData(input);
    expect(sanitized.password).toBe('[REDACTED]');
    expect(sanitized.authToken).toBe('[REDACTED]');
    expect(sanitized.secret_key).toBe('[REDACTED]');
    expect(sanitized.normalField).toBe('all-good');
  });

  it('redacts authorization and bearer tokens in strings', () => {
    const message = 'Failed with Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9 and token="secret123"';
    const sanitized = sanitizeData(message);
    expect(sanitized).not.toContain('eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9');
    expect(sanitized).toContain('Bearer [REDACTED]');
    expect(sanitized).toContain('token=[REDACTED]');
  });
});

describe('buildDryRunReport', () => {
  it('follows the standardized schema version and reportType', () => {
    const report = buildDryRunReport([], mockDesign);
    expect(report.schemaVersion).toBe('1.0.0');
    expect(report.reportType).toBe('meshery.dry-run-validation');
  });

  it('correctly maps design details', () => {
    const report = buildDryRunReport([], mockDesign);
    expect(report.design.id).toBe('ecommerce-app-id');
    expect(report.design.name).toBe('ecommerce-app');
  });

  it('handles empty or missing design gracefully', () => {
    const report = buildDryRunReport([], null as any);
    expect(report.design.id).toBe('');
    expect(report.design.name).toBe('Unnamed Design');
  });

  it('generates successful report when there are no errors', () => {
    const report = buildDryRunReport([], mockDesign);
    expect(report.run.status).toBe('success');
    expect(report.summary.resourcesChecked).toBe(4);
    expect(report.summary.resourcesPassed).toBe(4);
    expect(report.summary.resourcesFailed).toBe(0);
    expect(report.summary.warnings).toBe(0);

    expect(report.results).toHaveLength(4);
    expect(report.results.every((r) => r.status === 'valid')).toBe(true);
  });

  it('generates failed report when component validation errors exist', () => {
    const report = buildDryRunReport(mockDryRunErrorsWithFailures, mockDesign);
    expect(report.run.status).toBe('failed');
    expect(report.summary.resourcesChecked).toBe(4);
    expect(report.summary.resourcesPassed).toBe(3);
    expect(report.summary.resourcesFailed).toBe(1);

    const postgresResult = report.results.find((r) => r.name === 'postgres');
    expect(postgresResult).toBeDefined();
    expect(postgresResult?.status).toBe('failed');
    expect(postgresResult?.errors).toHaveLength(1);
    expect(postgresResult?.errors[0]).toEqual({
      message: 'StorageClass "fast-ssd" not found',
      errorType: 'ComponentError',
      fieldPath: 'spec.storageClassName',
      resourceKind: 'StatefulSet',
      apiVersion: 'apps/v1',
    });

    const frontendResult = report.results.find((r) => r.name === 'frontend');
    expect(frontendResult?.status).toBe('valid');
    expect(frontendResult?.errors).toHaveLength(0);
  });

  it('marks run as incomplete when no resources were checked and no errors are present', () => {
    const report = buildDryRunReport([], { name: 'empty-design' });
    expect(report.run.status).toBe('incomplete');
    expect(report.summary.resourcesChecked).toBe(0);
    expect(report.summary.resourcesPassed).toBe(0);
    expect(report.summary.resourcesFailed).toBe(0);
  });

  it('captures RequestError entries properly', () => {
    const requestErrors: DryRunErrorEntry[] = [
      {
        type: 'RequestError',
        errors: [{ message: 'Cluster API unreachable' } as any],
      },
    ];
    const report = buildDryRunReport(requestErrors, {});
    expect(report.run.status).toBe('failed');
    expect(report.results[0].name).toBe('DryRunRequest');
    expect(report.results[0].status).toBe('failed');
    expect(report.results[0].errors[0].errorType).toBe('RequestError');
  });

  it('supports custom timing parameters', () => {
    const timing = {
      startTime: '2026-10-28T10:42:21Z',
      completedTime: '2026-10-28T10:42:25Z',
      durationSeconds: 4.2,
    };
    const report = buildDryRunReport([], mockDesign, timing);
    expect(report.run.startTime).toBe('2026-10-28T10:42:21Z');
    expect(report.run.completedTime).toBe('2026-10-28T10:42:25Z');
    expect(report.run.durationSeconds).toBe(4.2);
  });

  it('normalizes string-only errors', () => {
    const stringErrors: DryRunErrorEntry[] = [
      {
        type: 'ComponentError',
        compName: 'custom-comp',
        errors: ['Plain string error occurred'],
      },
    ];
    const report = buildDryRunReport(stringErrors, {});
    expect(report.results[0].errors[0]).toEqual({
      message: 'Plain string error occurred',
      errorType: 'ComponentError',
      fieldPath: '',
      resourceKind: undefined,
      apiVersion: undefined,
    });
  });
});

describe('exportDryRunReport', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('triggers browser download with formatted JSON and sanitized filename', () => {
    exportDryRunReport(mockDryRunErrorsWithFailures, mockDesign);

    expect(downloadFileFromContent).toHaveBeenCalledTimes(1);
    const [jsonContent, filename, mimeType] = (downloadFileFromContent as any).mock.calls[0];

    const parsed = JSON.parse(jsonContent);
    expect(parsed.schemaVersion).toBe('1.0.0');
    expect(parsed.reportType).toBe('meshery.dry-run-validation');
    expect(parsed.design.name).toBe('ecommerce-app');

    expect(filename).toBe('ecommerce-app-dry-run-report.json');
    expect(mimeType).toBe('application/json');
  });

  it('handles design names with spaces and symbols in download filename', () => {
    exportDryRunReport([], { name: 'My App / Production [v2]' });
    const [, filename] = (downloadFileFromContent as any).mock.calls[0];
    expect(filename).toBe('my-app-production-v2-dry-run-report.json');
  });

  it('falls back to "design-dry-run-report.json" when design has no name', () => {
    exportDryRunReport([], {});
    const [, filename] = (downloadFileFromContent as any).mock.calls[0];
    expect(filename).toBe('unnamed-design-dry-run-report.json');
  });
});
