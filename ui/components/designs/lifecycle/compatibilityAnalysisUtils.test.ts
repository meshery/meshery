import { describe, expect, it } from 'vitest';
import {
  buildCompatibilityAnalysis,
  generateDependencyImpact,
  generatePossibleFix,
} from './compatibilityAnalysisUtils';

describe('compatibilityAnalysisUtils', () => {
  describe('generatePossibleFix', () => {
    it('detects missing CRD errors', () => {
      const fix = generatePossibleFix(
        'error: no matches for kind "PostgresCluster" in version "postgresql.cnpg.io/v1"',
        'StatefulSet',
        'prod-cluster',
      );
      expect(fix).toContain('Custom Resource Definition (CRD)');
      expect(fix).toContain('prod-cluster');
    });

    it('detects apiVersion deprecation errors', () => {
      const fix = generatePossibleFix(
        'unsupported apiVersion networking.k8s.io/v1beta1 for Ingress',
        'Ingress',
        'dev-cluster',
      );
      expect(fix).toContain('apiVersion');
      expect(fix).toContain('dev-cluster');
    });

    it('detects quota limit errors', () => {
      const fix = generatePossibleFix(
        'exceeded quota: pods-quota exceeded',
        'Pod',
        'staging-cluster',
      );
      expect(fix).toContain('quota');
    });
  });

  describe('generateDependencyImpact', () => {
    it('identifies dependent workloads for database failures', () => {
      const impact = generateDependencyImpact('postgres', 'StatefulSet', [
        'frontend',
        'backend',
        'postgres',
        'redis',
      ]);
      expect(impact).toContain('database connectivity');
      expect(impact).toContain('frontend, backend, redis');
    });

    it('identifies external ingress disruption for ingress failures', () => {
      const impact = generateDependencyImpact('app-ingress', 'Ingress', [
        'frontend',
        'app-ingress',
      ]);
      expect(impact).toContain('External ingress routing');
    });
  });

  describe('buildCompatibilityAnalysis', () => {
    it('builds a multi-cluster matrix from dryRunResponse', () => {
      const dryRunResponse = {
        frontend: {
          'cluster-1': { success: true },
          'cluster-2': { success: true },
          'cluster-3': { success: true },
        },
        postgres: {
          'cluster-1': { success: true },
          'cluster-2': { success: true },
          'cluster-3': {
            success: false,
            error: {
              Status:
                'error: no matches for kind "PostgresCluster" in version "postgresql.cnpg.io/v1"',
            },
          },
        },
        ingress: {
          'cluster-1': { success: true },
          'cluster-2': {
            success: false,
            error: {
              Status: 'Warning: Ingress class not specified',
            },
          },
          'cluster-3': { success: true },
        },
      };

      const selectedK8sConnections = [
        { id: 'cluster-1', metadata: { name: 'dev-cluster', version: 'v1.29.2' } },
        { id: 'cluster-2', metadata: { name: 'staging-cluster', version: 'v1.28.5' } },
        { id: 'cluster-3', metadata: { name: 'prod-cluster', version: 'v1.27.9' } },
      ];

      const analysis = buildCompatibilityAnalysis({
        dryRunResponse,
        selectedK8sConnections,
      });

      expect(analysis.clusters).toHaveLength(3);
      expect(analysis.resources).toHaveLength(3);
      expect(analysis.metrics.totalClusters).toBe(3);
      expect(analysis.metrics.totalResources).toBe(3);
      expect(analysis.metrics.succeededCount).toBe(1); // frontend
      expect(analysis.metrics.warningCount).toBe(1); // ingress
      expect(analysis.metrics.failedCount).toBe(1); // postgres

      // Check prod-cluster status on postgres
      const postgresProd = analysis.resources.find((r) => r.resourceName === 'postgres')
        ?.clusterResults['cluster-3'];
      expect(postgresProd?.status).toBe('failed');
      expect(postgresProd?.possibleFix).toContain('Custom Resource Definition');
      expect(postgresProd?.dependencyImpact).toContain('database connectivity');

      // Check staging-cluster status on ingress
      const ingressStaging = analysis.resources.find((r) => r.resourceName === 'ingress')
        ?.clusterResults['cluster-2'];
      expect(ingressStaging?.status).toBe('warning');

      // Check cluster summaries
      expect(analysis.clusterSummaries).toHaveLength(3);
      const prodSummary = analysis.clusterSummaries.find((c) => c.clusterId === 'cluster-3');
      expect(prodSummary?.failed).toBe(1);
      expect(prodSummary?.passed).toBe(2);

      // Check next steps
      expect(
        analysis.nextSteps.some((s) => s.includes('postgres') && s.includes('prod-cluster')),
      ).toBe(true);
    });

    it('gracefully handles empty dry-run responses without crashing', () => {
      const analysis = buildCompatibilityAnalysis({
        dryRunResponse: null,
        selectedK8sConnections: [],
      });
      expect(analysis.clusters).toHaveLength(1);
      expect(analysis.clusters[0].version).toBe('Unknown');
      expect(analysis.clusters[0].status).toBe('unknown');
      expect(analysis.resources).toHaveLength(0);
      expect(analysis.metrics.totalResources).toBe(0);
      expect(analysis.nextSteps).toContain(
        'All resources successfully validated across all targeted clusters.',
      );
    });

    it('does not treat "all" as a cluster ID when resolving selectedK8sContexts', () => {
      const selectedK8sConnections = [
        { id: 'cluster-1', metadata: { name: 'dev-cluster', version: 'v1.29.0' } },
      ];
      const analysis = buildCompatibilityAnalysis({
        dryRunResponse: {
          frontend: { 'cluster-1': { success: true } },
        },
        selectedK8sConnections,
        selectedK8sContexts: ['all'],
      });

      expect(analysis.clusters.some((c) => c.id === 'all')).toBe(false);
      expect(analysis.clusters).toHaveLength(1);
      expect(analysis.clusters[0].id).toBe('cluster-1');
    });

    it('marks missing cluster results as incomplete and does not declare deployment readiness', () => {
      const dryRunResponse = {
        frontend: {
          'cluster-1': { success: true },
          // cluster-2 is absent from response
        },
      };

      const selectedK8sConnections = [
        { id: 'cluster-1', metadata: { name: 'dev-cluster', version: 'v1.29.0' } },
        { id: 'cluster-2', metadata: { name: 'prod-cluster', version: 'v1.28.0' } },
      ];

      const analysis = buildCompatibilityAnalysis({
        dryRunResponse,
        selectedK8sConnections,
      });

      expect(analysis.clusters).toHaveLength(2);
      const frontend = analysis.resources.find((r) => r.resourceName === 'frontend');
      expect(frontend?.clusterResults['cluster-2']?.status).toBe('incomplete');
      expect(frontend?.clusterResults['cluster-2']?.errorMessage).toContain(
        'No dry-run result was returned for prod-cluster',
      );

      // Incomplete cells should not count as succeeded
      expect(analysis.metrics.succeededCount).toBe(0);
      expect(analysis.metrics.incompleteCount).toBe(1);

      // Cluster summary should track incomplete count
      const prodSummary = analysis.clusterSummaries.find((c) => c.clusterId === 'cluster-2');
      expect(prodSummary?.incomplete).toBe(1);
      expect(prodSummary?.passed).toBe(0);

      // Next steps must not declare readiness
      expect(analysis.nextSteps).not.toContain('Design is ready for deployment.');
      expect(
        analysis.nextSteps.some((s) =>
          s.includes('unavailable or incomplete cluster prod-cluster'),
        ),
      ).toBe(true);
    });

    it('uses Unknown and unknown fallback values for missing version and status', () => {
      const selectedK8sConnections = [{ id: 'cluster-anon' }];
      const analysis = buildCompatibilityAnalysis({
        dryRunResponse: null,
        selectedK8sConnections,
      });

      expect(analysis.clusters[0].version).toBe('Unknown');
      expect(analysis.clusters[0].status).toBe('unknown');
    });
  });
});
