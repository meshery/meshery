import React, { useEffect, useState } from 'react';
import Head from 'next/head';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Button, Grid, Switch, Typography } from '@material-ui/core';
import { withStyles } from '@material-ui/core/styles';
import { CatalogIcon } from '@layer5/sistent';
import { useGetUserPrefQuery, useUpdateUserPrefMutation } from '@/rtk-query/user';
import { UsesSistent } from '@/components/SistentWrapper';
import { Adapters } from '../components/extensions';
import DefaultError from '@/components/General/error-404';
import { toggleCatalogContent } from '../lib/store';
import { EVENT_TYPES } from '../lib/event-types';
import { EXTENSIONS } from '../utils/Enum';
import { useNotification } from '../utils/hooks/useNotification';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { LARGE_6_MED_12_GRID_STYLE } from '../css/grid.style';
import { extensionStyles as styles } from '../css/icons.styles';

const INITIAL_GRID_SIZE = { lg: 6, md: 12, xs: 12 };

const MeshMapSignUpcard = ({ classes, hasAccessToMeshMap = false }) => {
  const handleSignUp = (e) => {
    window.open('https://docs.layer5.io/kanvas', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
      <div className={classes.card}>
        <Typography className={classes.frontContent} variant="h5" component="div">
          Kanvas
        </Typography>

        <Typography className={classes.frontSideDescription} variant="body">
          <img className={classes.img} src="/static/img/meshmap.svg" />
          Collaboratively design and manage your Kubernetes clusters, service mesh deployments, and
          cloud native apps. Kanvas is now publicly available.{' '}
          {!hasAccessToMeshMap && 'Sign-up today to for access!'}
        </Typography>
        {
          <div style={{ textAlign: 'right' }}>
            <Button
              variant="contained"
              color="primary"
              disabled={hasAccessToMeshMap}
              className={classes.button}
              onClick={(e) => handleSignUp(e)}
            >
              {hasAccessToMeshMap ? 'Enabled' : 'Sign Up'}
            </Button>
          </div>
        }
      </div>
    </Grid>
  );
};

const LogoStyles = () => ({
  img: {
    paddingRight: '1rem',
    height: 'auto',
    width: 'auto',
    maxWidth: '220px',
    maxHeight: '150px',
  },
});

const MeshMapSnapShotLogo = withStyles(LogoStyles)(({ classes }) => {
  return <img className={classes.img} src="/static/img/meshmap-snapshot-logo.svg" />;
});

const MeshMapSnapShotCard = ({ classes, githubActionEnabled = false }) => {
  const handleEnable = (e) => {
    window.open('https://meshery.layer5.io/connect/github/new/', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
      <div className={classes.card}>
        <Typography className={classes.frontContent} variant="h5" component="div">
          GitHub Action: Kanvas Snapshot
        </Typography>

        <Typography className={classes.frontSideDescription} variant="body">
          <MeshMapSnapShotLogo />
          Connect Kanvas to your GitHub repo and see changes pull request-to-pull request. Get
          snapshots of your infrastructure directly in your PRs.
        </Typography>
        {
          <div style={{ textAlign: 'right' }}>
            <Button
              variant="contained"
              color="primary"
              disabled={githubActionEnabled}
              className={classes.button}
              onClick={(e) => handleEnable(e)}
            >
              {githubActionEnabled ? 'Remove' : 'Enable'}
            </Button>
          </div>
        }
      </div>
    </Grid>
  );
};

const MesheryPerformacneLogoStyles = () => ({
  img: {
    paddingRight: '1rem',
    height: 'auto',
    width: 'auto',
    maxWidth: '120px',
    maxHeight: '75px',
  },
});

const MesheryPerformanceActionLogo = withStyles(MesheryPerformacneLogoStyles)(({ classes }) => {
  return <img className={classes.img} src="/static/img/smp-dark.svg" />;
});

const MesheryPerformanceAction = ({ classes, githubActionEnabled = false }) => {
  const handleEnable = (e) => {
    window.open(
      'https://github.com/marketplace/actions/performance-testing-with-meshery',
      '_blank',
    );
    e.stopPropagation();
  };

  return (
    <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
      <div className={classes.card}>
        <Typography className={classes.frontContent} variant="h5" component="div">
          GitHub Action: Performance Analysis
        </Typography>

        <Typography className={classes.frontSideDescription} variant="body">
          <MesheryPerformanceActionLogo />
          Characterize the performance of your services using Meshery&apos;s performance analysis
          GitHub Action to benchmark and visually compare percentiles (e.g. P99) over time.
        </Typography>
        {
          <div style={{ textAlign: 'right' }}>
            <Button
              variant="contained"
              color="primary"
              disabled={githubActionEnabled}
              className={classes.button}
              onClick={(e) => handleEnable(e)}
            >
              {githubActionEnabled ? 'Remove' : 'Enable'}
            </Button>
          </div>
        }
      </div>
    </Grid>
  );
};

const MesheryDockerExtensionLogo = withStyles(MesheryPerformacneLogoStyles)(({ classes }) => {
  return <img className={classes.img} src="/static/img/docker.svg" />;
});

const MesheryDockerExtension = ({ classes }) => {
  const handleDownload = (e) => {
    window.open('https://hub.docker.com/extensions/meshery/docker-extension-meshery', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
      <div className={classes.card}>
        <Typography className={classes.frontContent} variant="h5" component="div">
          Meshery Docker Extension
        </Typography>

        <Typography className={classes.frontSideDescription} variant="body">
          <MesheryDockerExtensionLogo />
          Connect Meshery to your Kubernetes cluster via Docker Desktop and let MeshSync discover
          your clusters. Use MeshMap&apos;s no-code designer to collaboratively design and manage
          your infrastructure with ready-made patterns from Meshery Catalog.
        </Typography>
        {
          <div style={{ textAlign: 'right' }}>
            <Button
              variant="contained"
              color="primary"
              className={classes.button}
              onClick={(e) => handleDownload(e)}
            >
              Download
            </Button>
          </div>
        }
      </div>
    </Grid>
  );
};

const MesheryDesignEmbedLogo = withStyles(MesheryPerformacneLogoStyles)(({ classes }) => {
  return <img className={classes.img} src="/static/img/meshmap.svg" />;
});

const MesheryDesignEmbedExtension = ({ classes }) => {
  const handleLearnMore = (e) => {
    window.open('https://docs.layer5.io/kanvas/designer/embedding-designs/', '_blank');
    e.stopPropagation();
  };

  return (
    <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
      <div className={classes.card}>
        <Typography className={classes.frontContent} variant="h5" component="div">
          Meshery Design Embed
        </Typography>

        <Typography className={classes.frontSideDescription} variant="body">
          <MesheryDesignEmbedLogo />
          Meshery Design Embedding lets you export designs in an interactive format that seamlessly
          integrates with websites, blogs, and platforms using HTML, CSS, and JavaScript, making it
          easy to share with stakeholders.
        </Typography>
        {
          <div style={{ textAlign: 'right' }}>
            <Button
              variant="contained"
              color="primary"
              className={classes.button}
              onClick={(e) => handleLearnMore(e)}
            >
              Learn More
            </Button>
          </div>
        }
      </div>
    </Grid>
  );
};

export const WrappedMeshMapSignupCard = withStyles(styles)(MeshMapSignUpcard);
export const WrappedMeshMapSnapShopCard = withStyles(styles)(MeshMapSnapShotCard);
export const WrappedMesheryPerformanceAction = withStyles(styles)(MesheryPerformanceAction);
export const WrappedMesheryDockerExtension = withStyles(styles)(MesheryDockerExtension);
export const WrappedMesheryEmbedDesignExtension = withStyles(styles)(MesheryDesignEmbedExtension);
const Extensions = ({ classes, toggleCatalogContent, capabilitiesRegistry }) => {
  const [catalogContent, setCatalogContent] = useState(true);
  const [extensionPreferences, setExtensionPreferences] = useState({});
  const [hasAccessToMeshMap, setHasAccessToMeshMap] = useState(false);
  const { notify } = useNotification();
  const [updateUserPref] = useUpdateUserPrefMutation();

  const {
    data: userData,
    isSuccess: userDataFetched,
    isError: isUserError,
    error: userError,
  } = useGetUserPrefQuery();

  const handleToggle = () => {
    toggleCatalogContent({ catalogVisibility: !catalogContent });
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
      (val) => val.title.toLowerCase() === EXTENSIONS.MESHMAP,
    );
    if (typeof meshMapExtensionExists === 'object' && meshMapExtensionExists.length)
      setHasAccessToMeshMap(true);
  }, []);

  return (
    <React.Fragment>
      <Head>
        <title>Extensions | Meshery</title>
      </Head>
      {CAN(keys.VIEW_EXTENSIONS.action, keys.VIEW_EXTENSIONS.subject) ? (
        <Grid container spacing={1}>
          <WrappedMeshMapSnapShopCard githubActionEnabled={false} />
          <WrappedMesheryPerformanceAction githubActionEnabled={false} />
          <WrappedMeshMapSignupCard hasAccessToMeshMap={hasAccessToMeshMap} />
          <WrappedMesheryDockerExtension />
          <WrappedMesheryEmbedDesignExtension />
          <Grid item {...INITIAL_GRID_SIZE}>
            <div className={classes.card}>
              <Typography className={classes.frontContent} variant="h5" component="div">
                {'Meshery Catalog'}
              </Typography>

              <Typography className={classes.frontSideDescription} variant="body">
                <UsesSistent>
                  <CatalogIcon className={classes.img} />
                </UsesSistent>
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
              </Typography>

              <Grid
                container
                spacing={2}
                className={classes.grid}
                direction="row"
                justifyContent="space-between"
                alignItems="baseline"
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
                    className={classes.link}
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
                    classes={{
                      switchBase: classes.switchBase,
                      track: classes.track,
                      checked: classes.checked,
                    }}
                  />
                </div>
              </Grid>
            </div>
          </Grid>
          <Adapters />
        </Grid>
      ) : (
        <DefaultError />
      )}
    </React.Fragment>
  );
};

const mapStateToProps = (state) => ({
  catalogVisibility: state.get('catalogVisibility'),
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
});

const mapDispatchToProps = (dispatch) => ({
  toggleCatalogContent: bindActionCreators(toggleCatalogContent, dispatch),
});

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(Extensions));
