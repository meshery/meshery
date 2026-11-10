/**
 * Access Meshery (Playground) modal.
 *
 * Promo dialog shown when a user lands on a sandboxed Meshery Playground
 * instance, pointing them at the install / learning path resources. Composed
 * on the shared `Modal` primitive so the playground chrome matches the rest
 * of the app's modals.
 */
import React, { FC } from 'react';
import { Typography, Button } from '@sistent/sistent';
import { styled } from '@/theme';
import { Modal } from '@/components/shared/Modal';

const ImgWrapper = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  padding: '2rem',
});

const InsideImgWrapper = styled('img')(({ theme }) => ({
  padding: '0rem 0.5rem',
  content:
    theme.palette.mode === 'dark'
      ? "url('/static/img/meshery-logo/meshery-logo-text.svg')"
      : "url('/static/img/meshery-logo/meshery-logo-light-text.svg')",
}));

const InsideImgWrapperLogo = styled('img')({
  padding: '0rem 0.5rem',
});

const InstallButton = styled(Button)({
  marginBottom: '1rem',
});

const TutorialLink = styled('a')(({ theme }) => ({
  color: theme.palette.text.brand,
}));

const PromoBody = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  gap: '1rem',
});

const FooterRow = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  width: '100%',
});

export interface PlaygroundMeshDeployProps {
  isOpen: boolean;
  closeForm: () => void;
}

const PlaygroundMeshDeploy: FC<PlaygroundMeshDeployProps> = ({ isOpen, closeForm }) => {
  const handlePage = (e: React.MouseEvent<HTMLButtonElement>) => {
    window.open('https://meshery.io/#getting-started', '_blank', 'noopener,noreferrer');
    e.stopPropagation();
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={closeForm}
      title="The Cloud Native Playground"
      size="sm"
      aria-labelledby="access-meshery-modal-title"
      actions={
        <FooterRow>
          <InstallButton
            size="large"
            variant="contained"
            color="primary"
            onClick={handlePage}
            data-testid="install-meshery-button"
          >
            Install Meshery
          </InstallButton>
        </FooterRow>
      }
    >
      <PromoBody>
        <ImgWrapper>
          <InsideImgWrapperLogo
            width="20%"
            height="20%"
            src="/static/img/meshery-logo/meshery-logo.png"
          />
          <InsideImgWrapper width="50%" height="50%" />
        </ImgWrapper>
        <Typography gutterBottom>
          Meshery Playground gives you hands-on experience with designing cloud native systems -
          from your browser - using every CNCF project. Choose a{' '}
          <TutorialLink
            href="https://docs.meshery.io/guides/tutorials"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learning Path
          </TutorialLink>{' '}
          and work through labs as you visually and collaboratively learn-by-doing without having to
          install a single thing.
        </Typography>
        <Typography gutterBottom>
          To ensure that Meshery Playground remains a clean sandbox for all to use, many of
          Meshery&apos;s features are disabled. For full access to all of Meshery&apos;s features,
          deploy your own instance of Meshery.
        </Typography>
      </PromoBody>
    </Modal>
  );
};

export default PlaygroundMeshDeploy;
