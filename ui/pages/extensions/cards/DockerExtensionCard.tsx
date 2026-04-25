import React from 'react';
import { Button, Grid2, Typography } from '@sistent/sistent';
import {
  UnifiedCardContainer,
  UnifiedDescription,
  UnifiedButtonContainer,
  ResponsiveImage,
} from '../shared';
import { GRID_SIZE } from '../constants';

const DockerExtensionCard = () => {
  const handleDownload = (e) => {
    window.open('https://hub.docker.com/extensions/meshery/docker-extension-meshery', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="docker-extension-heading" variant="h5" component="div">
          Meshery Docker Extension
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/extensions/docker.svg"
            alt="Docker Logo"
            testId="docker-extension-icon"
          />
          <div>
            Connect Meshery to your Kubernetes cluster via Docker Desktop and let MeshSync discover
            your clusters. Use Kanvas&apos;s no-code designer to collaboratively design and manage
            your infrastructure with ready-made patterns from Meshery Catalog.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="docker-extension-download-btn"
            onClick={(e) => handleDownload(e)}
          >
            Download
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export default DockerExtensionCard;
