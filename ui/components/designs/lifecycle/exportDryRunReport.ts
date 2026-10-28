/**
 * exportDryRunReport.ts
 *
 * Serializes Dry Run validation results into a standardized JSON report
 * matching Meshery's Dry Run validation schema (v1.0.0) and triggers
 * an immediate file download in the browser.
 */

import { downloadFileFromContent } from '@/utils/fileDownloader';
import { processDesign } from './common';

export interface DryRunErrorDetail {
  Type?: string;
  FieldPath?: string;
  Message?: string;
  Causes?: unknown[];
  Status?: string;
  [key: string]: unknown;
}

export interface DryRunErrorEntry {
  type: 'ComponentError' | 'RequestError' | string;
  compName?: string;
  component?: {
    component?: {
      kind?: string;
      version?: string;
    };
    [key: string]: unknown;
  };
  contextId?: string;
  errors: Array<DryRunErrorDetail | string>;
}

export interface DesignData {
  id?: string;
  name?: string;
  components?: unknown[];
  schemaVersion?: string;
  [key: string]: unknown;
}

export interface DryRunReportRunInfo {
  startTime?: string;
  completedTime: string;
  status: 'success' | 'failed' | 'incomplete';
  durationSeconds?: number;
}

export interface DryRunReportSummary {
  resourcesChecked: number;
  resourcesPassed: number;
  resourcesFailed: number;
  warnings: number;
}

export interface DryRunReportErrorItem {
  message: string;
  errorType: string;
  fieldPath: string;
  resourceKind?: string;
  apiVersion?: string;
}

export interface DryRunReportResourceResult {
  name: string;
  type: string;
  status: 'valid' | 'failed';
  errors: DryRunReportErrorItem[];
}

export interface DryRunValidationReport {
  schemaVersion: '1.0.0';
  reportType: 'meshery.dry-run-validation';
  design: {
    id: string;
    name: string;
  };
  run: DryRunReportRunInfo;
  summary: DryRunReportSummary;
  results: DryRunReportResourceResult[];
}

const SENSITIVE_KEY_REGEX =
  /(password|token|secret|authorization|bearer|credential|apikey|private[_-]?key)/i;

/**
 * Strips sensitive values (tokens, credentials, authorization headers)
 * to comply with the acceptance criteria for secure reporting.
 */
export function sanitizeData<T>(data: T): T {
  if (data === null || data === undefined) {
    return data;
  }
  if (typeof data === 'string') {
    // Redact bearer tokens or authorization headers embedded in strings
    return data
      .replace(/(bearer\s+)[a-zA-Z0-9_\-\.]+/gi, '$1[REDACTED]')
      .replace(
        /(password|secret|token)\s*[:=]\s*["']?[^"'\s]+["']?/gi,
        '$1=[REDACTED]',
      ) as unknown as T;
  }
  if (Array.isArray(data)) {
    return data.map((item) => sanitizeData(item)) as unknown as T;
  }
  if (typeof data === 'object') {
    const cleanObj: Record<string, unknown> = {};
    for (const [key, val] of Object.entries(data)) {
      if (SENSITIVE_KEY_REGEX.test(key)) {
        cleanObj[key] = '[REDACTED]';
      } else {
        cleanObj[key] = sanitizeData(val);
      }
    }
    return cleanObj as unknown as T;
  }
  return data;
}

/**
 * Normalizes an individual error item into the standardized report shape.
 */
function normalizeErrorItem(
  error: DryRunErrorDetail | string,
  kind?: string,
  version?: string,
): DryRunReportErrorItem {
  if (typeof error === 'string') {
    return {
      message: sanitizeData(error),
      errorType: 'ComponentError',
      fieldPath: '',
      resourceKind: kind,
      apiVersion: version,
    };
  }

  const rawMessage =
    error.Message ||
    error.Status ||
    (error.FieldPath ? `Error at ${error.FieldPath}` : 'Validation error');
  return {
    message: sanitizeData(String(rawMessage)),
    errorType: sanitizeData(error.Type || 'ComponentError'),
    fieldPath: sanitizeData(error.FieldPath || ''),
    resourceKind: kind,
    apiVersion: version,
  };
}

/**
 * Builds the complete Dry Run JSON validation report.
 */
export function buildDryRunReport(
  dryRunErrors: DryRunErrorEntry[] = [],
  design: DesignData = {},
  timing?: { startTime?: string; completedTime?: string; durationSeconds?: number },
): DryRunValidationReport {
  const safeErrors = Array.isArray(dryRunErrors) ? dryRunErrors : [];
  const safeDesign = design || {};

  const { configurableComponents } = processDesign(safeDesign);

  const componentErrorsList = safeErrors.filter((e) => e && e.type === 'ComponentError');
  const requestErrorsList = safeErrors.filter((e) => e && e.type === 'RequestError');

  // Map of component errors by component name
  const errorsByCompName = new Map<string, DryRunErrorEntry>();
  for (const entry of componentErrorsList) {
    if (entry.compName) {
      errorsByCompName.set(entry.compName, entry);
    }
  }

  const results: DryRunReportResourceResult[] = [];
  const processedCompNames = new Set<string>();

  // Process all configurable components from the design
  if (Array.isArray(configurableComponents) && configurableComponents.length > 0) {
    for (const comp of configurableComponents) {
      const compName = comp.name || comp.displayName || 'unnamed-component';
      processedCompNames.add(compName);

      const errEntry = errorsByCompName.get(compName);
      const kind = comp.component?.kind;
      const version = comp.component?.version;

      if (errEntry && Array.isArray(errEntry.errors) && errEntry.errors.length > 0) {
        results.push({
          name: compName,
          type: kind || 'Component',
          status: 'failed',
          errors: errEntry.errors.map((err) => normalizeErrorItem(err, kind, version)),
        });
      } else {
        results.push({
          name: compName,
          type: kind || 'Component',
          status: 'valid',
          errors: [],
        });
      }
    }
  }

  // Also include any component errors that were not in configurableComponents
  for (const entry of componentErrorsList) {
    const compName = entry.compName || 'unknown-component';
    if (!processedCompNames.has(compName)) {
      processedCompNames.add(compName);
      const kind = entry.component?.component?.kind;
      const version = entry.component?.component?.version;
      const errs = Array.isArray(entry.errors) ? entry.errors : [];
      results.push({
        name: compName,
        type: kind || 'Component',
        status: errs.length > 0 ? 'failed' : 'valid',
        errors: errs.map((err) => normalizeErrorItem(err, kind, version)),
      });
    }
  }

  // If there were request-level errors and no results yet, report a general entry
  if (requestErrorsList.length > 0 && results.length === 0) {
    for (const reqErr of requestErrorsList) {
      const errs = Array.isArray(reqErr.errors) ? reqErr.errors : [];
      results.push({
        name: 'DryRunRequest',
        type: 'Request',
        status: 'failed',
        errors: errs.map((err) => ({
          message: sanitizeData(typeof err === 'string' ? err : JSON.stringify(err)),
          errorType: 'RequestError',
          fieldPath: '',
        })),
      });
    }
  }

  const failedCount = results.filter((r) => r.status === 'failed').length;
  const passedCount = results.filter((r) => r.status === 'valid').length;
  const totalChecked = results.length;

  let runStatus: 'success' | 'failed' | 'incomplete';
  if (totalChecked === 0 && safeErrors.length === 0) {
    runStatus = 'incomplete';
  } else if (failedCount > 0 || requestErrorsList.length > 0) {
    runStatus = 'failed';
  } else {
    runStatus = 'success';
  }

  const completedTime = timing?.completedTime || new Date().toISOString();
  const startTime = timing?.startTime || completedTime;
  const durationSeconds = timing?.durationSeconds ?? 0;

  return {
    schemaVersion: '1.0.0',
    reportType: 'meshery.dry-run-validation',
    design: {
      id: sanitizeData(safeDesign.id || safeDesign.name || ''),
      name: sanitizeData(safeDesign.name || 'Unnamed Design'),
    },
    run: {
      startTime,
      completedTime,
      status: runStatus,
      durationSeconds,
    },
    summary: {
      resourcesChecked: totalChecked,
      resourcesPassed: passedCount,
      resourcesFailed: failedCount,
      warnings: 0,
    },
    results,
  };
}

/**
 * Serializes the report to a formatted JSON string and triggers
 * an immediate file download in the user's browser.
 */
export function exportDryRunReport(
  dryRunErrors: DryRunErrorEntry[] = [],
  design: DesignData = {},
  timing?: { startTime?: string; completedTime?: string; durationSeconds?: number },
): void {
  const report = buildDryRunReport(dryRunErrors, design, timing);
  const designSlug = (report.design.name || 'dry-run-report')
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');

  const fileName = `${designSlug || 'design'}-dry-run-report.json`;
  downloadFileFromContent(JSON.stringify(report, null, 2), fileName, 'application/json');
}
