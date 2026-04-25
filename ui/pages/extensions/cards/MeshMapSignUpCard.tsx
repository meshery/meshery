import React from 'react';
import { Button, Grid2, Typography } from '@sistent/sistent';
import {
  UnifiedCardContainer,
  UnifiedDescription,
  UnifiedButtonContainer,
  ResponsiveImage,
} from '../shared';
import { GRID_SIZE } from '../constants';

const MeshMapSignUpCard = ({ hasAccessToMeshMap = false }) => {
  const handleSignUp = (e) => {
    window.open('https://docs.layer5.io/kanvas', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="kanvas-signup-heading" variant="h5" component="div">
          Kanvas
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/extensions/kanvas-icon-color.svg"
            alt="Kanvas Icon"
            testId="kanvas-signup-icon"
          />
          <div>
            Collaboratively design and manage your infra and apps. Kanvas is now publicly available.{' '}
            {!hasAccessToMeshMap && 'Sign-up today for access!'}
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            data-testid="kanvas-signup-btn"
            disabled={hasAccessToMeshMap}
            onClick={(e) => handleSignUp(e)}
          >
            {hasAccessToMeshMap ? 'Enabled' : 'Sign Up'}
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default MeshMapSignUpCard;
