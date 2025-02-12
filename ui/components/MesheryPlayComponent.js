import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import NoSsr from '@mui/material/NoSsr';
import {
  Button,
  Divider,
  MenuItem,
  TextField,
  Grid,
  Typography,
  styled,
  useTheme,
  gray,
  charcoal,
} from '@layer5/sistent';
import PropTypes from 'prop-types';
import { useRouter } from 'next/router';
import SettingsIcon from '@mui/icons-material/Settings';
import MesheryAdapterPlayComponent from './MesheryAdapterPlayComponent';
import { bindActionCreators } from 'redux';
import { setAdapter } from '../lib/store';

const StyledButton = styled(Button)(({ theme }) => ({
  marginTop: theme.spacing(3),
  marginLeft: theme.spacing(1),
}));

const StyledIcon = styled(SettingsIcon)(({ theme }) => ({
  fontSize: 23,
  width: theme.spacing(2.5),
  marginRight: theme.spacing(0.5),
  alignSelf: 'flex-start',
}));

const PlayRoot = styled('div')(({ theme }) => ({
  padding: theme.spacing(0),
  marginBottom: theme.spacing(2),
}));

export const PaneSection = styled('div')(() => {
  const theme = useTheme();
  return {
    backgroundColor: theme.palette.mode === 'dark' ? gray[20] : charcoal[90],
    padding: theme.spacing(2.5),
    borderRadius: 4,
  };
});

const AlreadyConfigured = styled('div')({
  textAlign: 'center',
  height: '100%',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  flexDirection: 'column',
});

const MesheryPlayComponent = (props) => {
  const { meshAdapters } = props;
  const router = useRouter();

  // Initialize state
  const [adapter, setAdapterState] = useState(() => {
    if (meshAdapters && meshAdapters.size > 0) {
      return meshAdapters[0];
    }
    return {};
  });

  const handleRouteChange = () => {
    const queryParam = router?.query?.adapter;
    if (queryParam) {
      const selectedAdapter = meshAdapters.find(
        ({ adapter_location }) => adapter_location === queryParam,
      );
      if (selectedAdapter) {
        setAdapterState(selectedAdapter);
      }
    } else if (meshAdapters.size > 0) {
      setAdapterState(meshAdapters.get(0));
    }
  };

  useEffect(() => {
    router.events.on('routeChangeComplete', handleRouteChange);

    return () => {
      router.events.off('routeChangeComplete', handleRouteChange);
    };
  }, [router.events]);

  useEffect(() => {
    if (meshAdapters?.size > 0) {
      handleRouteChange();
    }
  }, [meshAdapters?.size]);

  const handleConfigure = () => {
    router.push('/settings?settingsCategory=Adapters');
  };

  const pickImage = (adapter) => {
    let image = '/static/img/meshery-logo.png';
    let imageIcon = <img src={image} style={{ width: '24px' }} alt="Meshery Logo" />;
    if (adapter && adapter.name) {
      image = `/static/img/${adapter.name.toLowerCase()}.svg`;
      imageIcon = <img src={image} style={{ width: '24px' }} alt={`${adapter.name} Logo`} />;
    }
    return imageIcon;
  };

  const handleAdapterChange = () => {
    return (event) => {
      const { setAdapter } = props;
      if (event.target.value !== '') {
        const selectedAdapter = meshAdapters.filter(
          ({ adapter_location }) => adapter_location === event.target.value,
        );
        if (selectedAdapter && selectedAdapter.size === 1) {
          setAdapterState(selectedAdapter.get(0));
          setAdapter({ selectedAdapter: selectedAdapter.get(0).name });
        }
      }
    };
  };

  const renderIndividualAdapter = () => {
    let adapCount = 0;
    let adapter;
    meshAdapters.forEach((adap) => {
      if (adap.adapter_location === props.adapter) {
        adapter = adap;
        meshAdapters.forEach((ad) => {
          if (ad.name == adap.name) adapCount += 1;
        });
      }
    });
    if (adapter) {
      const imageIcon = pickImage(adapter);
      return (
        <>
          <MesheryAdapterPlayComponent
            adapter={adapter}
            adapCount={adapCount}
            adapter_icon={imageIcon}
          />
        </>
      );
    }
    return '';
  };

  if (meshAdapters.size === 0) {
    return (
      <NoSsr>
        <AlreadyConfigured>
          <Typography variant="h6">Adapter Unavailable</Typography>
          <Typography variant="subtitle1">Connect Meshery Adapter(s) in Settings</Typography>
          <StyledButton
            variant="contained"
            color="primary"
            size="large"
            onClick={handleConfigure}
            startIcon={<StyledIcon />}
          >
            Configure Settings
          </StyledButton>
        </AlreadyConfigured>
      </NoSsr>
    );
  }

  if (props.adapter && props.adapter !== '') {
    const indContent = renderIndividualAdapter();
    if (indContent !== '') {
      return indContent;
    }
  }

  const imageIcon = pickImage(adapter);

  return (
    <NoSsr>
      <PlayRoot>
        <Grid container>
          <Grid item xs={12}>
            <PaneSection>
              <TextField
                select
                id="adapter_id"
                name="adapter_name"
                label="Select Meshery Adapter"
                data-cy="lifecycle-service-mesh-type"
                fullWidth
                value={adapter && adapter.adapter_location ? adapter.adapter_location : ''}
                margin="normal"
                variant="outlined"
                sx={{
                  display: 'flex',
                }}
                onChange={handleAdapterChange()}
                SelectProps={{
                  MenuProps: {
                    anchorOrigin: {
                      vertical: 'bottom',
                      horizontal: 'left',
                    },
                    transformOrigin: {
                      vertical: 'top',
                      horizontal: 'left',
                    },
                  },
                }}
              >
                {meshAdapters.map((ada) => (
                  <MenuItem
                    key={`${ada.adapter_location}_${new Date().getTime()}`}
                    value={ada.adapter_location}
                    sx={{
                      display: 'flex',
                    }}
                  >
                    {pickImage(ada)}
                    <Typography variant="body1" sx={{ ml: 1 }}>
                      {ada.adapter_location}
                    </Typography>
                  </MenuItem>
                ))}
              </TextField>
            </PaneSection>
          </Grid>
        </Grid>
      </PlayRoot>
      <Divider variant="fullWidth" light />
      {adapter && adapter.adapter_location && (
        <MesheryAdapterPlayComponent adapter={adapter} adapter_icon={imageIcon} />
      )}
    </NoSsr>
  );
};

MesheryPlayComponent.propTypes = {
  meshAdapters: PropTypes.object.isRequired,
  setAdapter: PropTypes.func.isRequired,
  adapter: PropTypes.string,
};

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const meshAdapters = state.get('meshAdapters');
  const meshAdaptersts = state.get('meshAdaptersts');
  const selectedAdapter = state.get('selectedAdapter');
  return { k8sconfig, meshAdapters, meshAdaptersts, selectedAdapter };
};

const mapDispatchToProps = (dispatch) => ({
  setAdapter: bindActionCreators(setAdapter, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(MesheryPlayComponent);
