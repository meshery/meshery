import React, { useEffect, useState } from 'react';
import Head from 'next/head';
import { Button, CatalogIcon, Grid2, Switch, Typography, useTheme, Box } from '@sistent/sistent';
import { useGetUserPrefQuery, useUpdateUserPrefMutation } from '@/rtk-query/user';
import { Adapters } from '../components/extensions';
import DefaultError from '@/components/General/error-404';
import { EVENT_TYPES } from '../lib/event-types';
import { EXTENSION_NAMES } from '../utils/Enum';
import { useNotification } from '../utils/hooks/useNotification';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { CardContainer, FrontSideDescription } from '../css/icons.styles';
import { useDispatch, useSelector } from 'react-redux';
import { toggleCatalogContent, updatePage } from '@/store/slices/mesheryUi';
import { getPath } from 'lib/path';

const UnifiedCardContainer = ({ children, sx = {} }) => {
  return (
    <CardContainer
      sx={{
        minHeight: {
          xs: '280px',
          sm: '260px',
          lg: '280px',
        },
        display: 'flex',
        flexDirection: 'column',
        position: 'relative',
        ...sx,
      }}
    >
      {children}
    </CardContainer>
  );
};

const UnifiedDescription = ({ children, hasIcon = false, ...props }) => {
  return (
    <FrontSideDescription
      {...props}
      sx={{
        flex: 1,
        display: 'flex',
        flexDirection: {
          xs: hasIcon ? 'column' : 'row',
          sm: 'row',
        },
        alignItems: {
          xs: hasIcon ? 'center' : 'flex-start',
          sm: 'flex-start',
        },
        textAlign: {
          xs: hasIcon ? 'center' : 'left',
          sm: 'left',
        },
        gap: {
          xs: hasIcon ? '12px' : '8px',
          sm: '8px',
        },
        marginBottom: {
          xs: '50px',
          sm: '45px',
          lg: '40px',
        },
        overflow: 'hidden',
        wordWrap: 'break-word',
        hyphens: 'auto',
        '& > *:not(img)': {
          maxWidth: '100%',
          overflow: 'hidden',
          wordWrap: 'break-word',
          hyphens: 'auto',
        },
      }}
    >
      {children}
    </FrontSideDescription>
  );
};

const UnifiedButtonContainer = ({ children }) => {
  return (
    <Box
      sx={{
        position: 'absolute',
        bottom: 12,
        left: 12,
        right: 12,
        textAlign: 'right',
      }}
    >
      {children}
    </Box>
  );
};

// Responsive image component with slightly smaller sizes
const ResponsiveImage = ({ src, alt, testId }) => {
  return (
    <img
      style={{
        height: 'auto',
        width: 'auto',
        maxWidth: '140px',
        maxHeight: '85px',
        minWidth: '100px',
        minHeight: '60px',
        flexShrink: 0,
      }}
      data-testid={testId}
      src={src}
      alt={alt || ''}
    />
  );
};

const GRID_SIZE = {
  xs: 12,
  sm: 12,
  md: 12,
  lg: 6,
  xl: 4,
};

const MeshMapSignUpcard = ({ hasAccessToMeshMap = false }) => {
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
            src="/static/img/kanvas-icon-color.svg"
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

const MeshMapSnapShotCard = ({ githubActionEnabled = false }) => {
  const handleEnable = (e) => {
    window.open('https://cloud.layer5.io/connect/github/new/', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="kanvas-snapshot-heading" variant="h5" component="div">
          GitHub Action: Kanvas Snapshot
        </Typography>

        <UnifiedDescription data-testid="kanvas-snapshot-description" variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/meshmap-snapshot-logo.svg"
            alt="Kanvas Snapshot Logo"
            testId="kanvas-snapshot-image"
          />
          <div>
            Connect Kanvas to your GitHub repo and see changes pull request-to-pull request. Get
            snapshots of your infrastructure directly in your PRs.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="kanvas-snapshot-enable-btn"
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

const MesheryPerformanceAction = ({ githubActionEnabled = false }) => {
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
            src="/static/img/smp-dark.svg"
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

const MesheryDockerExtension = () => {
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
            src="/static/img/docker.svg"
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

const MesheryHelmKanvasExtension = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/helm-kanvas-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography variant="h5" component="div">
          Kanvas Snapshot Helm Plugin
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/helm_chart.svg"
            alt="Helm Chart Logo"
            testId="helm-kanvas-icon"
          />
          <div>
            The Kanvas Snapshot Helm Plugin allows you to generate a visual snapshot of your Helm
            charts directly from the command line. It simplifies the process of creating Meshery
            Snapshots, providing a visual representation of packaged Helm charts.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="helm-kanvas-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

const MesheryDesignEmbedExtension = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.layer5.io/kanvas/designer/embedding-designs/', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography variant="h5" component="div">
          Meshery Design Embed
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/meshmap.svg"
            alt="Meshery Design Logo"
            testId="design-embed-icon"
          />
          <div>
            Meshery Design Embedding lets you export designs in an interactive format that
            seamlessly integrates with websites, blogs, and platforms using HTML, CSS, and
            JavaScript, making it easy to share with stakeholders.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="design-embed-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

const Layer5CloudExtension = () => {
  const handleLearnMore = (e) => {
    window.open('https://meshery.io/extensions/layer5-cloud', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="layer5-cloud-heading" variant="h5" component="div">
          Layer5 Cloud
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/layer5.svg"
            alt="Layer5 Cloud Logo"
            testId="layer5-cloud-image"
          />
          <div>
            A comprehensive platform offering identity and collaboration services, private catalogs,
            GitOps, and multi-Meshery management. Leverage its extensible authorization framework
            and organizational hierarchy for streamlined cloud infrastructure management.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="layer5-cloud-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

const KubectlPluginExtension = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/kubectl-kanvas-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="kubectl-plugin-heading" variant="h5" component="div">
          Kubectl Plugin for Kanvas Snapshot
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/kubectl.svg"
            alt="Kubectl Plugin Logo"
            testId="kubectl-plugin-image"
          />
          <div>
            Generate visual snapshots of your Kubernetes manifests directly from kubectl. Cluster
            configurations and workflows with Kanvas Snapshots. Receive snapshots via email or get
            instant terminal URL display.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="kubectl-plugin-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

const KubectlMeshSyncExtension = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/kubectl-meshsync-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="kubectl-meshsync-heading" variant="h5" component="div">
          Kubectl Plugin for MeshSync Snapshot
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <ResponsiveImage
            src="/static/img/meshsync.svg"
            alt="MeshSync Logo"
            testId="kubectl-meshsync-image"
          />
          <div>
            Capture cluster state directly from kubectl with simplified networking and access
            requirements. Generate MeshSync snapshots for offline management and visualization in
            Meshery Server, without requiring full Meshery Operator deployment.
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          <Button
            variant="contained"
            color="primary"
            data-testid="kubectl-meshsync-learn-more-btn"
            onClick={(e) => handleLearnMore(e)}
          >
            Learn More
          </Button>
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

// Catalog card with special layout considerations
const CatalogCard = ({ catalogContent, handleToggle, theme }) => {
  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography data-testid="catalog-section-heading" variant="h5" component="div">
          {'Meshery Catalog'}
        </Typography>

        <UnifiedDescription variant="body" hasIcon={true}>
          <CatalogIcon
            data-testid="catalog-toggle-switch"
            style={{
              height: '85px', // Reduced from 100px
              width: '140px', // Reduced from 160px
              flexShrink: 0,
            }}
          />

          <div
            style={{
              display: 'inline',
              position: 'relative',
            }}
          >
            Enable access to the cloud native catalog, supporting design patterns, WebAssembly
            filters (<span style={{ fontStyle: 'italic' }}>soon</span>), and OPA policies (
            <span style={{ fontStyle: 'italic' }}>soon</span>). Import any catalog item and
            customize.
          </div>
        </UnifiedDescription>

        <Box
          sx={{
            position: 'absolute',
            bottom: 12, // Reduced from 16
            left: 12, // Reduced from 16
            right: 12, // Reduced from 16
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
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

          <Switch
            checked={catalogContent}
            onChange={handleToggle}
            name="OperatorSwitch"
            color="primary"
          />
        </Box>
      </UnifiedCardContainer>
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
  const theme = useTheme();

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

  return (
    <>
      <React.Fragment>
        <Head>
          <title>Extensions | Meshery</title>
        </Head>
        {CAN(keys.VIEW_EXTENSIONS.action, keys.VIEW_EXTENSIONS.subject) ? (
          <Grid2 container spacing={2} size="grow">
            <WrappedMeshMapSnapShopCard githubActionEnabled={false} />
            <WrappedMesheryPerformanceAction githubActionEnabled={false} />
            <WrappedMeshMapSignupCard hasAccessToMeshMap={hasAccessToMeshMap} />
            <WrappedMesheryHelmKanvasExtension />
            <WrappedMesheryDockerExtension />
            <WrappedMesheryEmbedDesignExtension />
            <WrappedLayer5CloudExtension />
            <WrappedKubectlPluginExtension />
            <WrappedKubectlMeshSyncExtension />
            <CatalogCard
              catalogContent={catalogContent}
              handleToggle={handleToggle}
              theme={theme}
            />
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
