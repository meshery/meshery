import React, { useEffect, useState } from 'react';
import Head from 'next/head';
import { Button, CatalogIcon, Grid2, Switch, Typography, useTheme } from '@sistent/sistent';
import { useGetUserPrefQuery, useUpdateUserPrefMutation } from '@/rtk-query/user';
import { Adapters } from '../components/extensions';
import DefaultError from '@/components/General/error-404';
import { EVENT_TYPES } from '../lib/event-types';
import { EXTENSION_NAMES } from '../utils/Enum';
import { useNotification } from '../utils/hooks/useNotification';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { LARGE_6_MED_12_GRID_STYLE } from '../css/grid.style';
import { CardContainer, FrontSideDescription, ImageWrapper } from '../css/icons.styles';
import { useDispatch, useSelector } from 'react-redux';
import { toggleCatalogContent, updatePage } from '@/store/slices/mesheryUi';
import { getPath } from 'lib/path';

// Hook to detect screen size
const useMediaQuery = (query) => {
  const [matches, setMatches] = useState(false);

  useEffect(() => {
    const media = window.matchMedia(query);
    if (media.matches !== matches) {
      setMatches(media.matches);
    }
    const listener = () => setMatches(media.matches);
    media.addListener(listener);
    return () => media.removeListener(listener);
  }, [matches, query]);

  return matches;
};

const ResponsiveCardContainer = ({ children }) => {
  const isMobile = useMediaQuery('(max-width: 768px)');

  return (
    <CardContainer
      style={{
        height: isMobile ? 'auto' : '280px',
        minHeight: isMobile ? '200px' : '280px',
        display: 'flex',
        flexDirection: 'column',
        position: 'relative',
      }}
    >
      {children}
    </CardContainer>
  );
};

const ResponsiveFrontSideDescription = ({ children, ...props }) => {
  const isMobile = useMediaQuery('(max-width: 768px)');

  return (
    <FrontSideDescription
      {...props}
      style={{
        flex: 1,
        display: 'flex',
        alignItems: 'flex-start',
        marginBottom: isMobile ? '50px' : '60px',
      }}
    >
      {children}
    </FrontSideDescription>
  );
};

// Button container that adapts for mobile
const ResponsiveButtonContainer = ({ children }) => {
  //const isMobile = useMediaQuery('(max-width: 768px)');

  return (
    <div
      style={{
        position: 'absolute',
        bottom: '16px',
        right: '16px',
        left: '16px',
        textAlign: 'right',
      }}
    >
      {children}
    </div>
  );
};

// Responsive grid size prop
const useResponsiveGridSize = () => {
  const isMobile = useMediaQuery('(max-width: 768px)');
  return isMobile ? { xs: 12 } : LARGE_6_MED_12_GRID_STYLE;
};

const MeshMapSignUpcard = ({ hasAccessToMeshMap = false }) => {
  const gridSize = useResponsiveGridSize();

  const handleSignUp = (e) => {
    window.open('https://docs.layer5.io/kanvas', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={gridSize}>
      <ResponsiveCardContainer>
        <Typography data-testid="kanvas-signup-heading" variant="h5" component="div">
          Kanvas
        </Typography>

        <ResponsiveFrontSideDescription variant="body">
          <ImageWrapper src="/static/img/kanvas-icon-color.svg" />
          Collaboratively design and manage your infra and apps. Kanvas is now publicly available.{' '}
          {!hasAccessToMeshMap && 'Sign-up today to for access!'}
        </ResponsiveFrontSideDescription>

        <ResponsiveButtonContainer>
          <Button
            variant="contained"
            data-testid="kanvas-signup-btn"
            disabled={hasAccessToMeshMap}
            onClick={(e) => handleSignUp(e)}
          >
            {hasAccessToMeshMap ? 'Enabled' : 'Sign Up'}
          </Button>
        </ResponsiveButtonContainer>
      </ResponsiveCardContainer>
    </Grid2>
  );
};

const MeshMapSnapShotLogo = () => {
  return (
    <img
      style={{
        paddingRight: '1rem',
        height: 'auto',
        width: 'auto',
        maxWidth: '120px',
        maxHeight: '75px',
      }}
      data-testid="kanvas-snapshot-image"
      src="/static/img/meshmap-snapshot-logo.svg"
    />
  );
};

const MeshMapSnapShotCard = ({ githubActionEnabled = false }) => {
  const gridSize = useResponsiveGridSize();

  const handleEnable = (e) => {
    window.open('https://cloud.layer5.io/connect/github/new/', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={gridSize}>
      <ResponsiveCardContainer>
        <Typography data-testid="kanvas-snapshot-heading" variant="h5" component="div">
          GitHub Action: Kanvas Snapshot
        </Typography>

        <ResponsiveFrontSideDescription data-testid="kanvas-snapshot-description" variant="body">
          <MeshMapSnapShotLogo />
          Connect Kanvas to your GitHub repo and see changes pull request-to-pull request. Get
          snapshots of your infrastructure directly in your PRs.
        </ResponsiveFrontSideDescription>

        <ResponsiveButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="kanvas-snapshot-enable-btn"
            disabled={githubActionEnabled}
            onClick={(e) => handleEnable(e)}
          >
            {githubActionEnabled ? 'Remove' : 'Enable'}
          </Button>
        </ResponsiveButtonContainer>
      </ResponsiveCardContainer>
    </Grid2>
  );
};

const MesheryPerformanceActionLogo = () => {
  return (
    <img
      style={{
        paddingRight: '1rem',
        height: 'auto',
        width: 'auto',
        maxWidth: '120px',
        maxHeight: '75px',
      }}
      src="/static/img/smp-dark.svg"
    />
  );
};

const MesheryPerformanceAction = ({ githubActionEnabled = false }) => {
  const gridSize = useResponsiveGridSize();

  const handleEnable = (e) => {
    window.open(
      'https://github.com/marketplace/actions/performance-testing-with-meshery',
      '_blank',
    );
    e.stopPropagation();
  };

  return (
    <Grid2 size={gridSize}>
      <ResponsiveCardContainer>
        <Typography data-testid="performance-analysis-heading" variant="h5" component="div">
          GitHub Action: Performance Analysis
        </Typography>

        <ResponsiveFrontSideDescription variant="body">
          <MesheryPerformanceActionLogo />
          Characterize the performance of your services using Meshery&apos;s performance analysis
          GitHub Action to benchmark and visually compare percentiles (e.g. P99) over time.
        </ResponsiveFrontSideDescription>

        <ResponsiveButtonContainer>
          <Button
            variant="contained"
            data-testid="performance-analysis-enable-btn"
            disabled={githubActionEnabled}
            onClick={(e) => handleEnable(e)}
          >
            {githubActionEnabled ? 'Remove' : 'Enable'}
          </Button>
        </ResponsiveButtonContainer>
      </ResponsiveCardContainer>
    </Grid2>
  );
};

const MesheryDockerExtensionLogo = () => {
  return (
    <img
      style={{
        paddingRight: '1rem',
        height: 'auto',
        width: 'auto',
        maxWidth: '120px',
        maxHeight: '75px',
      }}
      src="/static/img/docker.svg"
    />
  );
};

const MesheryDockerExtension = () => {
  const gridSize = useResponsiveGridSize();

  const handleDownload = (e) => {
    window.open('https://hub.docker.com/extensions/meshery/docker-extension-meshery', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={gridSize}>
      <ResponsiveCardContainer>
        <Typography data-testid="docker-extension-heading" variant="h5" component="div">
          Meshery Docker Extension
        </Typography>

        <ResponsiveFrontSideDescription variant="body">
          <MesheryDockerExtensionLogo />
          Connect Meshery to your Kubernetes cluster via Docker Desktop and let MeshSync discover
          your clusters. Use Kanvas&apos;s no-code designer to collaboratively design and manage
          your infrastructure with ready-made patterns from Meshery Catalog.
        </ResponsiveFrontSideDescription>

        <ResponsiveButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="docker-extension-download-btn"
            onClick={(e) => handleDownload(e)}
          >
            Download
          </Button>
        </ResponsiveButtonContainer>
      </ResponsiveCardContainer>
    </Grid2>
  );
};

const MesheryDesignEmbedLogo = () => {
  return (
    <img
      style={{
        paddingRight: '1rem',
        height: 'auto',
        width: 'auto',
        maxWidth: '120px',
        maxHeight: '75px',
      }}
      src="/static/img/meshmap.svg"
    />
  );
};

const MesheryHelmKanvasLogo = () => {
  return (
    <img
      style={{
        paddingRight: '1rem',
        height: 'auto',
        width: 'auto',
        maxWidth: '120px',
        maxHeight: '75px',
      }}
      src="/static/img/helm_chart.svg"
    />
  );
};

const MesheryHelmKanvasExtension = () => {
  const gridSize = useResponsiveGridSize();

  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/helm-kanvas-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={gridSize}>
      <ResponsiveCardContainer>
        <Typography variant="h5" component="div">
          Kanvas Snapshot Helm Plugin
        </Typography>

        <ResponsiveFrontSideDescription variant="body">
          <MesheryHelmKanvasLogo />
          The Kanvas Snapshot Helm Plugin allows you to generate a visual snapshot of your Helm
          charts directly from the command line. It simplifies the process of creating Meshery
          Snapshots, providing a visual representation of packaged Helm charts.
        </ResponsiveFrontSideDescription>

        <ResponsiveButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="helm-kanvas-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </ResponsiveButtonContainer>
      </ResponsiveCardContainer>
    </Grid2>
  );
};

const MesheryDesignEmbedExtension = () => {
  const gridSize = useResponsiveGridSize();

  const handleLearnMore = (e) => {
    window.open('https://docs.layer5.io/kanvas/designer/embedding-designs/', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={gridSize}>
      <ResponsiveCardContainer>
        <Typography variant="h5" component="div">
          Meshery Design Embed
        </Typography>

        <ResponsiveFrontSideDescription variant="body">
          <MesheryDesignEmbedLogo />
          Meshery Design Embedding lets you export designs in an interactive format that seamlessly
          integrates with websites, blogs, and platforms using HTML, CSS, and JavaScript, making it
          easy to share with stakeholders.
        </ResponsiveFrontSideDescription>

        <ResponsiveButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="design-embed-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </ResponsiveButtonContainer>
      </ResponsiveCardContainer>
    </Grid2>
  );
};

const Layer5CloudLogo = () => {
  return (
    <img
      style={{
        paddingRight: '1rem',
        height: 'auto',
        width: 'auto',
        maxWidth: '120px',
        maxHeight: '75px',
      }}
      src="/static/img/layer5.svg"
      data-testid="layer5-cloud-image"
    />
  );
};

const Layer5CloudExtension = () => {
  const gridSize = useResponsiveGridSize();

  const handleLearnMore = (e) => {
    window.open('https://meshery.io/extensions/layer5-cloud', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={gridSize}>
      <ResponsiveCardContainer>
        <Typography data-testid="layer5-cloud-heading" variant="h5" component="div">
          Layer5 Cloud
        </Typography>

        <ResponsiveFrontSideDescription variant="body">
          <Layer5CloudLogo />A comprehensive platform offering identity and collaboration services,
          private catalogs, GitOps, and multi-Meshery management. Leverage its extensible
          authorization framework and organizational hierarchy for streamlined cloud infrastructure
          management.
        </ResponsiveFrontSideDescription>

        <ResponsiveButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="layer5-cloud-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </ResponsiveButtonContainer>
      </ResponsiveCardContainer>
    </Grid2>
  );
};

const KubectlPluginLogo = () => {
  return (
    <img
      style={{
        paddingRight: '1rem',
        height: 'auto',
        width: 'auto',
        maxWidth: '120px',
        maxHeight: '75px',
      }}
      src="/static/img/kubectl.svg"
      data-testid="kubectl-plugin-image"
    />
  );
};

const KubectlPluginExtension = () => {
  const gridSize = useResponsiveGridSize();

  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/kubectl-kanvas-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={gridSize}>
      <ResponsiveCardContainer>
        <Typography data-testid="kubectl-plugin-heading" variant="h5" component="div">
          Kubectl Plugin for Kanvas Snapshot
        </Typography>

        <ResponsiveFrontSideDescription variant="body">
          <KubectlPluginLogo />
          Generate visual snapshots of your Kubernetes manifests directly from kubectl. cluster
          configurations and workflows with Kanvas Snapshots. Receive snapshots via email or get
          instant terminal URL display.
        </ResponsiveFrontSideDescription>

        <ResponsiveButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="kubectl-plugin-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </ResponsiveButtonContainer>
      </ResponsiveCardContainer>
    </Grid2>
  );
};

const KubectlMeshSyncLogo = () => {
  return (
    <img
      style={{
        paddingRight: '1rem',
        height: 'auto',
        width: 'auto',
        maxWidth: '120px',
        maxHeight: '75px',
      }}
      src="/static/img/meshsync.svg"
      data-testid="kubectl-meshsync-image"
    />
  );
};

const KubectlMeshSyncExtension = () => {
  const gridSize = useResponsiveGridSize();

  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/kubectl-meshsync-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={gridSize}>
      <ResponsiveCardContainer>
        <Typography data-testid="kubectl-meshsync-heading" variant="h5" component="div">
          Kubectl Plugin for MeshSync Snapshot
        </Typography>

        <ResponsiveFrontSideDescription variant="body">
          <KubectlMeshSyncLogo />
          Capture cluster state directly from kubectl with simplified networking and access
          requirements. Generate MeshSync snapshots for offline management and visualization in
          Meshery Server, without requiring full Meshery Operator deployment.
        </ResponsiveFrontSideDescription>

        <ResponsiveButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="kubectl-meshsync-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </ResponsiveButtonContainer>
      </ResponsiveCardContainer>
    </Grid2>
  );
};

export const WrappedLayer5CloudExtension = Layer5CloudExtension;
export const WrappedKubectlPluginExtension = KubectlPluginExtension;
export const WrappedKubectlMeshSyncExtension = KubectlMeshSyncExtension;
export const WrappedMeshMapSignupCard = MeshMapSignUpcard;
export const WrappedMeshMapSnapShopCard = MeshMapSnapShotCard;
export const WrappedMesheryPerformanceAction = MesheryPerformanceAction;
export const WrappedMesheryDockerExtension = MesheryDockerExtension;
export const WrappedMesheryEmbedDesignExtension = MesheryDesignEmbedExtension;
export const WrappedMesheryHelmKanvasExtension = MesheryHelmKanvasExtension;

const Extensions = () => {
  const [catalogContent, setCatalogContent] = useState(true);
  const [extensionPreferences, setExtensionPreferences] = useState({});
  const [hasAccessToMeshMap, setHasAccessToMeshMap] = useState(false);
  const { notify } = useNotification();
  const [updateUserPref] = useUpdateUserPrefMutation();
  const dispatch = useDispatch();
  const gridSize = useResponsiveGridSize();

  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Extensions' }));
  }, []);

  const { capabilitiesRegistry } = useSelector((state) => state.ui);

  const {
    data: userData,
    isSuccess: userDataFetched,
    isError: isUserError,
    error: userError,
  } = useGetUserPrefQuery();

  const handleToggle = () => {
    dispatch(toggleCatalogContent({ catalogVisibility: !catalogContent }));
    setCatalogContent(!catalogContent);
    handleCatalogPreference(!catalogContent);
  };

  const fetchUser = () => {
    if (userDataFetched && userData) {
      setExtensionPreferences(userData?.usersExtensionPreferences);
      setCatalogContent(userData?.usersExtensionPreferences?.catalogContent);
    } else if (isUserError) {
      console.log(userError);
    }
  };
  useEffect(() => {
    fetchUser();
  }, [userData]);

  const handleCatalogPreference = (catalogPref) => {
    let body = Object.assign({}, extensionPreferences);
    body['catalogContent'] = catalogPref;
    updateUserPref({ usersExtensionPreferences: body })
      .unwrap()
      .then(() => {
        notify({
          message: `Catalog Content was ${catalogPref ? 'enab' : 'disab'}led`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    const meshMapExtensionExists = capabilitiesRegistry?.extensions?.navigator?.filter(
      (val) => val.title.toLowerCase() === EXTENSION_NAMES.KANVAS,
    );
    if (typeof meshMapExtensionExists === 'object' && meshMapExtensionExists.length)
      setHasAccessToMeshMap(true);
  }, []);

  const theme = useTheme();

  return (
    <>
      <React.Fragment>
        <Head>
          <title>Extensions | Meshery</title>
        </Head>
        {CAN(keys.VIEW_EXTENSIONS.action, keys.VIEW_EXTENSIONS.subject) ? (
          <Grid2 container spacing={1} size="grow">
            <WrappedMeshMapSnapShopCard githubActionEnabled={false} />
            <WrappedMesheryPerformanceAction githubActionEnabled={false} />
            <WrappedMeshMapSignupCard hasAccessToMeshMap={hasAccessToMeshMap} />
            <WrappedMesheryHelmKanvasExtension />
            <WrappedMesheryDockerExtension />
            <WrappedMesheryEmbedDesignExtension />
            <WrappedLayer5CloudExtension />
            <WrappedKubectlPluginExtension />
            <WrappedKubectlMeshSyncExtension />
            <Grid2 size={gridSize}>
              <ResponsiveCardContainer>
                <Typography data-testid="catalog-section-heading" variant="h5" component="div">
                  {'Meshery Catalog'}
                </Typography>

                <ResponsiveFrontSideDescription variant="body">
                  <CatalogIcon
                    data-testid="catalog-toggle-switch"
                    style={{
                      paddingRight: '1rem',
                      height: '75px',
                      width: '120px',
                      flexShrink: 0,
                    }}
                  />

                  <div
                    style={{
                      display: 'inline',
                      position: 'relative',
                    }}
                  >
                    Enable access to the cloud native catalog, supporting design patterns,
                    WebAssembly filters (<span style={{ fontStyle: 'italic' }}>soon</span>), and OPA
                    policies (<span style={{ fontStyle: 'italic' }}>soon</span>). Import any catalog
                    item and customize.
                  </div>
                </ResponsiveFrontSideDescription>

                <Grid2
                  container
                  spacing={2}
                  direction="row"
                  justifyContent="space-between"
                  alignItems="baseline"
                  size="grow"
                  style={{
                    position: 'absolute',
                    paddingRight: '3rem',
                    paddingLeft: '.5rem',
                    bottom: '1.5rem',
                  }}
                >
                  <Typography variant="subtitle2" style={{ fontStyle: 'italic' }}>
                    Explore the{' '}
                    <a
                      href="https://meshery.io/catalog"
                      target="_blank"
                      rel="noreferrer"
                      style={{
                        textDecoration: 'none',
                        color: theme.palette.text.brand,
                      }}
                    >
                      Meshery Catalog
                    </a>
                  </Typography>

                  <div style={{ textAlign: 'right' }}>
                    <Switch
                      checked={catalogContent}
                      onChange={handleToggle}
                      name="OperatorSwitch"
                      color="primary"
                    />
                  </div>
                </Grid2>
              </ResponsiveCardContainer>
            </Grid2>
            <Adapters />
          </Grid2>
        ) : (
          <DefaultError />
        )}
      </React.Fragment>
    </>
  );
};

export default Extensions;
