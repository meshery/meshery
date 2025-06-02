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

const MeshMapSignUpcard = ({ hasAccessToMeshMap = false }) => {
  const handleSignUp = (e) => {
    window.open('https://docs.layer5.io/kanvas', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid2 size={LARGE_6_MED_12_GRID_STYLE}>
      <CardContainer>
        <Typography data-testid="kanvas-signup-heading" variant="h5" component="div">
          Kanvas
        </Typography>

        <FrontSideDescription variant="body">
          <ImageWrapper src="/static/img/kanvas-icon-color.svg" />
          Collaboratively design and manage your Kubernetes clusters, service mesh deployments, and
          cloud native apps. Kanvas is now publicly available.{' '}
          {!hasAccessToMeshMap && 'Sign-up today to for access!'}
        </FrontSideDescription>
        {
          <div style={{ textAlign: 'right' }}>
            <Button
              variant="contained"
              data-testid="kanvas-signup-btn"
              disabled={hasAccessToMeshMap}
              onClick={(e) => handleSignUp(e)}
            >
              {hasAccessToMeshMap ? 'Enabled' : 'Sign Up'}
            </Button>
          </div>
        }
      </CardContainer>
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
        maxWidth: '220px',
        maxHeight: '150px',
      }}
      data-testid="kanvas-snapshot-image"
      src="/static/img/meshmap-snapshot-logo.svg"
    />
  );
};

const MeshMapSnapShotCard = ({ githubActionEnabled = false }) => {
  const handleEnable = (e) => {
    window.open('https://cloud.layer5.io/connect/github/new/', '_blank');
    e.stopPropagation();
  };

  return (
    <>
      <Grid2 size={LARGE_6_MED_12_GRID_STYLE}>
        <CardContainer>
          <Typography data-testid="kanvas-snapshot-heading" variant="h5" component="div">
            GitHub Action: Kanvas Snapshot
          </Typography>

          <FrontSideDescription data-testid="kanvas-snapshot-description" variant="body">
            <MeshMapSnapShotLogo />
            Connect Kanvas to your GitHub repo and see changes pull request-to-pull request. Get
            snapshots of your infrastructure directly in your PRs.
          </FrontSideDescription>

          <div style={{ textAlign: 'right' }}>
            <Button
              variant="contained"
              color="primary"
              data-testid="kanvas-snapshot-enable-btn"
              disabled={githubActionEnabled}
              onClick={(e) => handleEnable(e)}
            >
              {githubActionEnabled ? 'Remove' : 'Enable'}
            </Button>
          </div>
        </CardContainer>
      </Grid2>
    </>
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
  const handleEnable = (e) => {
    window.open(
      'https://github.com/marketplace/actions/performance-testing-with-meshery',
      '_blank',
    );
    e.stopPropagation();
  };

  return (
    <>
      <Grid2 size={LARGE_6_MED_12_GRID_STYLE}>
        <CardContainer>
          <Typography data-testid="performance-analysis-heading" variant="h5" component="div">
            GitHub Action: Performance Analysis
          </Typography>

          <FrontSideDescription variant="body">
            <MesheryPerformanceActionLogo />
            Characterize the performance of your services using Meshery&apos;s performance analysis
            GitHub Action to benchmark and visually compare percentiles (e.g. P99) over time.
          </FrontSideDescription>

          <div style={{ textAlign: 'right' }}>
            <Button
              variant="contained"
              data-testid="performance-analysis-enable-btn"
              disabled={githubActionEnabled}
              onClick={(e) => handleEnable(e)}
            >
              {githubActionEnabled ? 'Remove' : 'Enable'}
            </Button>
          </div>
        </CardContainer>
      </Grid2>
    </>
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
  const handleDownload = (e) => {
    window.open('https://hub.docker.com/extensions/meshery/docker-extension-meshery', '_blank');
    e.stopPropagation();
  };

  return (
    <>
      <Grid2 size={LARGE_6_MED_12_GRID_STYLE}>
        <CardContainer>
          <Typography data-testid="docker-extension-heading" variant="h5" component="div">
            Meshery Docker Extension
          </Typography>

          <FrontSideDescription variant="body">
            <MesheryDockerExtensionLogo />
            Connect Meshery to your Kubernetes cluster via Docker Desktop and let MeshSync discover
            your clusters. Use Kanvas&apos;s no-code designer to collaboratively design and manage
            your infrastructure with ready-made patterns from Meshery Catalog.
          </FrontSideDescription>
          {
            <div style={{ textAlign: 'right' }}>
              <Button
                variant="contained"
                color="primary"
                data-testid="docker-extension-download-btn"
                onClick={(e) => handleDownload(e)}
              >
                Download
              </Button>
            </div>
          }
        </CardContainer>
      </Grid2>
    </>
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
  const handleLearnMore = (e) => {
    window.open('https://docs.meshery.io/extensions/helm-kanvas-snapshot', '_blank');
    e.stopPropagation();
  };

  return (
    <>
      <Grid2 size={LARGE_6_MED_12_GRID_STYLE}>
        <CardContainer>
          <Typography variant="h5" component="div">
            Kanvas Snapshot Helm Plugin
          </Typography>

          <FrontSideDescription variant="body">
            <MesheryHelmKanvasLogo />
            The Kanvas Snapshot Helm Plugin allows you to generate a visual snapshot of your Helm
            charts directly from the command line. It simplifies the process of creating Meshery
            Snapshots, providing a visual representation of packaged Helm charts.
          </FrontSideDescription>
          {
            <div style={{ textAlign: 'right' }}>
              <Button
                variant="contained"
                color="primary"
                data-testid="helm-kanvas-learn-more-btn"
                onClick={(e) => handleLearnMore(e)}
              >
                Learn More
              </Button>
            </div>
          }
        </CardContainer>
      </Grid2>
    </>
  );
};

const MesheryDesignEmbedExtension = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.layer5.io/kanvas/designer/embedding-designs/', '_blank');
    e.stopPropagation();
  };

  return (
    <>
      <Grid2 size={LARGE_6_MED_12_GRID_STYLE}>
        <CardContainer>
          <Typography variant="h5" component="div">
            Meshery Design Embed
          </Typography>

          <FrontSideDescription variant="body">
            <MesheryDesignEmbedLogo />
            Meshery Design Embedding lets you export designs in an interactive format that
            seamlessly integrates with websites, blogs, and platforms using HTML, CSS, and
            JavaScript, making it easy to share with stakeholders.
          </FrontSideDescription>
          {
            <div style={{ textAlign: 'right' }}>
              <Button
                variant="contained"
                color="primary"
                data-testid="design-embed-learn-more-btn"
                onClick={(e) => handleLearnMore(e)}
              >
                Learn More
              </Button>
            </div>
          }
        </CardContainer>
      </Grid2>
    </>
  );
};

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
            <Grid2 size={LARGE_6_MED_12_GRID_STYLE}>
              <CardContainer>
                <Typography data-testid="catalog-section-heading" variant="h5" component="div">
                  {'Meshery Catalog'}
                </Typography>

                <FrontSideDescription variant="body">
                  <CatalogIcon
                    data-testid="catalog-toggle-switch"
                    style={{
                      paddingRight: '1rem',
                      height: '80px',
                      width: '80px',
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
                </FrontSideDescription>

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
              </CardContainer>
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
