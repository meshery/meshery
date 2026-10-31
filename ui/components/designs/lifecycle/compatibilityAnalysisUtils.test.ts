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
      const postgresProd = analysis.resources
        .find((r) => r.resourceName === 'postgres')
        ?.clusterResults['cluster-3'];
      expect(postgresProd?.status).toBe('failed');
      expect(postgresProd?.possibleFix).toContain('Custom Resource Definition');
      expect(postgresProd?.dependencyImpact).toContain('database connectivity');

      // Check staging-cluster status on ingress
      const ingressStaging = analysis.resources
        .find((r) => r.resourceName === 'ingress')
        ?.clusterResults['cluster-2'];
      expect(ingressStaging?.status).toBe('warning');

      // Check cluster summaries
      expect(analysis.clusterSummaries).toHaveLength(3);
      const prodSummary = analysis.clusterSummaries.find((c) => c.clusterId === 'cluster-3');
      expect(prodSummary?.failed).toBe(1);
      expect(prodSummary?.passed).toBe(2);

      // Check next steps
      expect(analysis.nextSteps.some((s) => s.includes('postgres') && s.includes('prod-cluster'))).toBe(true);
    });

    it('gracefully handles empty dry-run responses without crashing', () => {
      const analysis = buildCompatibilityAnalysis({
        dryRunResponse: null,
        selectedK8sConnections: [],
      });
      expect(analysis.clusters).toHaveLength(1);
      expect(analysis.resources).toHaveLength(0);
      expect(analysis.metrics.totalResources).toBe(0);
      expect(analysis.nextSteps).toContain('All resources successfully validated across all targeted clusters.');
    });
  });
});
