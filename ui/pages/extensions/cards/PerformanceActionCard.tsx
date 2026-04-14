import React from 'react';
import { Button, Grid2, Typography } from '@sistent/sistent';
import {
  UnifiedCardContainer,
  UnifiedDescription,
  UnifiedButtonContainer,
  ResponsiveImage,
} from '../shared';
import { GRID_SIZE } from '../constants';

const PerformanceActionCard = ({ githubActionEnabled = false }) => {
  const handleEnable = (e) => {
    window.open(
      'https://github.com/marketplace/actions/performance-testing-with-meshery',
      '_blank',
    );
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="performance-analysis-heading" variant="h5" component="div">
          GitHub Action: Performance Analysis
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/extensions/smp-dark.svg"
            alt="Meshery Performance Logo"
            testId="performance-analysis-icon"
          />
          <div>
            Characterize the performance of your services using Meshery&apos;s performance analysis
            GitHub Action to benchmark and visually compare percentiles (e.g. P99) over time.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            data-testid="performance-analysis-enable-btn"
            disabled={githubActionEnabled}
            onClick={(e) => handleEnable(e)}
          >
            {githubActionEnabled ? 'Remove' : 'Enable'}
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default PerformanceActionCard;
