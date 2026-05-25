import React from 'react';
import { Box, IconButton, Typography } from '@sistent/sistent';
import { ArrowBack, GetApp as GetAppIcon } from '@/assets/icons';
import MesheryChart from '../MesheryChart';
import { iconMedium } from '../../css/icons.styles';

interface PerformanceTestResultsProps {
  testResult: any;
  chartStyle: React.CSSProperties;
  onBack: () => void;
}

/**
 * PerformanceTestResults renders the post-run results view: header bar with
 * back / download buttons and the MesheryChart for the latest runner result.
 *
 * Extracted from `performance/index.tsx` (was the inline `Results` component)
 * in Phase 5.a so the entry point stays under the 600-line size budget.
 */
const PerformanceTestResults: React.FC<PerformanceTestResultsProps> = ({
  testResult,
  chartStyle,
  onBack,
}) => {
  if (!testResult || !testResult.runner_results) {
    return null;
  }

  return (
    <div>
      <Box display="flex" justifyContent="space-between" alignItems="center">
        <IconButton onClick={onBack}>
          <ArrowBack />
        </IconButton>
        <Typography variant="h6" gutterBottom sx={{ textAlign: 'center' }} id="timerAnchor">
          Test Results
        </Typography>
        <IconButton
          key="download"
          aria-label="download"
          color="inherit"
          href={`/api/perf/profile/result/${encodeURIComponent(testResult.meshery_id)}`}
        >
          <GetAppIcon style={iconMedium} />
        </IconButton>
      </Box>
      <div style={chartStyle}>
        <MesheryChart
          rawdata={[testResult && testResult.runner_results ? testResult : {}]}
          data={[testResult && testResult.runner_results ? testResult.runner_results : {}]}
        />
      </div>
    </div>
  );
};

export default PerformanceTestResults;
