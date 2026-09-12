import React from 'react';
import { Button, CircularProgress } from '@sistent/sistent';
import { SaveOutlined as SaveOutlinedIcon } from '@/assets/icons';

import { Keys } from '@meshery/schemas/permissions';

interface PerformanceFormActionsProps {
  disableTest: boolean;
  blockRunTest: boolean;
  hasTestResult: boolean;
  onAbort: () => void;
  onShowResults: () => void;
  onSaveProfile: () => void;
  onRunTest: () => void;
}

/**
 * PerformanceFormActions renders the four buttons in the form footer:
 * Clear, Results (when a test has been run), Save Profile, and Run Test.
 *
 * Extracted from `performance/index.tsx` in Phase 5.a so the entry point
 * stays under the 600-line size budget.
 */
const PerformanceFormActions: React.FC<PerformanceFormActionsProps> = ({
  disableTest,
  blockRunTest,
  hasTestResult,
  onAbort,
  onShowResults,
  onSaveProfile,
  onRunTest,
}) => {
  return (
    <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
      <Button
        type="submit"
        variant="contained"
        color="primary"
        size="large"
        sx={{ marginLeft: '1rem' }}
        disabled={disableTest}
        onClick={onAbort}
      >
        Clear
      </Button>
      {hasTestResult && (
        <Button
          type="submit"
          variant="contained"
          color="primary"
          size="large"
          sx={{ marginLeft: '1rem' }}
          disabled={disableTest}
          onClick={onShowResults}
        >
          Results
        </Button>
      )}
      <Button
        type="submit"
        variant="contained"
        color="primary"
        size="large"
        onClick={onSaveProfile}
        sx={{ marginLeft: '1rem' }}
        disabled={disableTest}
        startIcon={<SaveOutlinedIcon />}
      >
        Save Profile
      </Button>
      <Button
        type="submit"
        data-testid="run-performance-test"
        variant="contained"
        color="primary"
        size="large"
        onClick={onRunTest}
        sx={{ marginLeft: '1rem' }}
        disabled={blockRunTest || disableTest}
        permissionKey={Keys.PerformanceManagementRunTest}
      >
        {blockRunTest ? <CircularProgress size={30} /> : 'Run Test'}
      </Button>
    </div>
  );
};

export default PerformanceFormActions;
