import React, { useEffect, useState } from 'react';
import Head from 'next/head';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { CatalogIcon, Button, Grid, Switch, Typography, Box } from '@layer5/sistent';
import { styled } from '@mui/material/styles';
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

const StyledCard = styled('div')(() => ({
  padding: '32px',
  position: 'relative',
  boxShadow: '0 4px 8px 0 rgba(0, 0, 0, 0.2)',
  minHeight: '250px',
  borderRadius: '8px',
  transformStyle: 'preserve-3d',
  backgroundColor: '#FFF',
}));

const StyledButton = styled(Button)(({ theme }) => ({
  color: '#fff',
  minWidth: '100px',
  borderRadius: '5px',
  backgroundColor: '#607d8b',
  '&:hover': {
    backgroundColor: theme.palette.grey[700],
  },
}));

const StyledImg = styled('img')(() => ({
  paddingRight: '1rem',
  height: 'auto',
  width: 'auto',
  maxWidth: '220px',
  maxHeight: '150px',
}));

const MeshMapSignUpcard = ({ hasAccessToMeshMap = false }) => {
  const handleSignUp = (e) => {
    window.open('https://docs.layer5.io/kanvas', '_blank');
    e.stopPropagation();
  };

  return (
    <UsesSistent>
      <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
        <StyledCard>
          <Typography variant="h5" component="div" sx={{ fontWeight: 'bold' }}>
            Kanvas
          </Typography>

          <Box sx={{ display: 'flex', alignItems: 'center', padding: '16px 0' }}>
            <StyledImg src="/static/img/meshmap.svg" alt="MeshMap" style={{ maxWidth: '150px' }} />
            <Typography variant="body1" sx={{ fontSize: '14px' }}>
              Collaboratively design and manage your Kubernetes clusters, service mesh deployments,
              and cloud native apps. Kanvas is now publicly available.{' '}
              {!hasAccessToMeshMap && 'Sign-up today for access!'}
            </Typography>
          </Box>

          <div style={{ textAlign: 'right' }}>
            <StyledButton disabled={hasAccessToMeshMap} onClick={(e) => handleSignUp(e)}>
              {hasAccessToMeshMap ? 'Enabled' : 'Sign Up'}
            </StyledButton>
          </div>
        </StyledCard>
      </Grid>
    </UsesSistent>
  );
};

const MeshMapSnapShotLogo = () => {
  return <StyledImg src="/static/img/meshmap-snapshot-logo.svg" alt="MeshMap Snapshot Logo" />;
};

const MeshMapSnapShotCard = ({ githubActionEnabled = false }) => {
  const handleEnable = (e) => {
    window.open('https://meshery.layer5.io/connect/github/new/', '_blank');
    e.stopPropagation();
  };

  return (
    <UsesSistent>
      <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
        <StyledCard>
          <Typography variant="h5" component="div" sx={{ fontWeight: 'bold' }}>
            GitHub Action: Kanvas Snapshot
          </Typography>

          <Box sx={{ display: 'flex', alignItems: 'center', padding: '16px 0' }}>
            <MeshMapSnapShotLogo />
            <Typography variant="body1" sx={{ fontSize: '14px' }}>
              Connect Kanvas to your GitHub repo and see changes pull request-to-pull request. Get
              snapshots of your infrastructure directly in your PRs.
            </Typography>
          </Box>

          <div style={{ textAlign: 'right' }}>
            <StyledButton disabled={githubActionEnabled} onClick={(e) => handleEnable(e)}>
              {githubActionEnabled ? 'Remove' : 'Enable'}
            </StyledButton>
          </div>
        </StyledCard>
      </Grid>
    </UsesSistent>
  );
};

const MesheryPerformanceActionLogo = styled('img')({
  paddingRight: '1rem',
  height: 'auto',
  width: 'auto',
  maxWidth: '120px',
  maxHeight: '75px',
});

const MesheryPerformanceAction = ({ githubActionEnabled = false }) => {
  const handleEnable = (e) => {
    window.open(
      'https://github.com/marketplace/actions/performance-testing-with-meshery',
      '_blank',
    );
    e.stopPropagation();
  };

  return (
    <UsesSistent>
      <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
        <StyledCard>
          <Typography variant="h5" component="div" sx={{ fontWeight: 'bold' }}>
            GitHub Action: Performance Analysis
          </Typography>

          <Box sx={{ display: 'flex', alignItems: 'center', padding: '16px 0' }}>
            <MesheryPerformanceActionLogo
              src="/static/img/smp-dark.svg"
              alt="Performance Action Logo"
              sx={{ marginRight: '1rem', width: 'auto', height: 'auto', maxWidth: '150px' }} // Adjust size and spacing
            />
            <Typography variant="body1" sx={{ fontSize: '14px' }}>
              Characterize the performance of your services using Meshery&apos;s performance
              analysis GitHub Action to benchmark and visually compare percentiles (e.g. P99) over
              time.
            </Typography>
          </Box>

          <div style={{ textAlign: 'right' }}>
            <StyledButton disabled={githubActionEnabled} onClick={(e) => handleEnable(e)}>
              {githubActionEnabled ? 'Remove' : 'Enable'}
            </StyledButton>
          </div>
        </StyledCard>
      </Grid>
    </UsesSistent>
  );
};

const MesheryDockerExtensionLogo = styled('img')({
  paddingRight: '1rem',
  height: 'auto',
  width: 'auto',
  maxWidth: '120px',
  maxHeight: '75px',
});

const MesheryDockerExtension = () => {
  const handleDownload = (e) => {
    window.open('https://hub.docker.com/extensions/meshery/docker-extension-meshery', '_blank');
    e.stopPropagation();
  };

  return (
    <UsesSistent>
      <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
        <StyledCard>
          <Typography variant="h5" component="div" sx={{ fontWeight: 'bold' }}>
            Meshery Docker Extension
          </Typography>

          <Box sx={{ display: 'flex', alignItems: 'center', padding: '16px 0' }}>
            <MesheryDockerExtensionLogo
              src="/static/img/docker.svg"
              alt="Docker Extension Logo"
              sx={{ marginRight: '1rem', width: 'auto', height: 'auto', maxWidth: '150px' }} // Adjust size and spacing
            />
            <Typography variant="body1" sx={{ fontSize: '14px' }}>
              Connect Meshery to your Kubernetes cluster via Docker Desktop and let MeshSync
              discover your clusters. Use MeshMap&apos;s no-code designer to collaboratively design
              and manage your infrastructure with ready-made patterns from Meshery Catalog.
            </Typography>
          </Box>

          <div style={{ textAlign: 'right' }}>
            <StyledButton onClick={(e) => handleDownload(e)}>Download</StyledButton>
          </div>
        </StyledCard>
      </Grid>
    </UsesSistent>
  );
};

const MesheryDesignEmbedLogo = styled('img')({
  paddingRight: '1rem',
  height: 'auto',
  width: 'auto',
  maxWidth: '120px',
  maxHeight: '75px',
});

const MesheryDesignEmbedExtension = () => {
  const handleLearnMore = (e) => {
    window.open('https://docs.layer5.io/kanvas/designer/embedding-designs/', '_blank');
    e.stopPropagation();
  };

  return (
    <UsesSistent>
      <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
        <StyledCard>
          <Typography variant="h5" component="div" sx={{ fontWeight: 'bold' }}>
            Meshery Design Embed
          </Typography>

          <Box sx={{ display: 'flex', alignItems: 'center', padding: '16px 0' }}>
            <MesheryDesignEmbedLogo
              src="/static/img/meshmap.svg"
              alt="Design Embed Logo"
              sx={{ marginRight: '1rem', width: 'auto', height: 'auto', maxWidth: '150px' }} // Adjust size and spacing
            />
            <Typography variant="body1" sx={{ fontSize: '14px' }}>
              Meshery Design Embedding lets you export designs in an interactive format that
              seamlessly integrates with websites, blogs, and platforms using HTML, CSS, and
              JavaScript, making it easy to share with stakeholders.
            </Typography>
          </Box>

          <div style={{ textAlign: 'right' }}>
            <StyledButton onClick={(e) => handleLearnMore(e)}>Learn More</StyledButton>
          </div>
        </StyledCard>
      </Grid>
    </UsesSistent>
  );
};

const Extensions = ({ toggleCatalogContent, capabilitiesRegistry }) => {
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
        <UsesSistent>
          <Grid container spacing={1}>
            <MeshMapSnapShotCard githubActionEnabled={false} />
            <MesheryPerformanceAction githubActionEnabled={false} />
            <MeshMapSignUpcard hasAccessToMeshMap={hasAccessToMeshMap} />
            <MesheryDockerExtension />
            <MesheryDesignEmbedExtension />
            <Grid item {...LARGE_6_MED_12_GRID_STYLE}>
              <StyledCard>
                <Typography variant="h5" component="div" sx={{ fontWeight: 'bold' }}>
                  {'Meshery Catalog'}
                </Typography>

                <Box sx={{ display: 'flex', alignItems: 'center', padding: '16px 0' }}>
                  <CatalogIcon
                    style={{
                      paddingRight: '1rem',
                      height: 'auto',
                      width: 'auto',
                      minWidth: '100px',
                      minHeight: '80px',
                      maxWidth: '220px',
                      maxHeight: '150px',
                    }}
                  />
                  <div style={{ display: 'inline', position: 'relative' }}>
                    <Typography variant="body1" sx={{ fontSize: '14px' }}>
                      Enable access to the cloud-native catalog, supporting design patterns,
                      WebAssembly filters (<span style={{ fontStyle: 'italic' }}>soon</span>), and
                      OPA policies (<span style={{ fontStyle: 'italic' }}>soon</span>). Import any
                      catalog item and customize.
                    </Typography>
                  </div>
                </Box>

                <Grid
                  container
                  spacing={2}
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
                    <a href="https://meshery.io/catalog" target="_blank" rel="noreferrer">
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
                </Grid>
              </StyledCard>
            </Grid>
            <Adapters />
          </Grid>
        </UsesSistent>
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

export default connect(mapStateToProps, mapDispatchToProps)(Extensions);
