import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import NoSsr from '@mui/material/NoSsr';
import { Button, Divider, MenuItem, TextField, Grid, Typography, styled } from '@layer5/sistent';
import { blue } from '@mui/material/colors';
import PropTypes from 'prop-types';
import { useRouter } from 'next/router';
import SettingsIcon from '@mui/icons-material/Settings';
import MesheryAdapterPlayComponent from './MesheryAdapterPlayComponent';
import { bindActionCreators } from 'redux';
import { setAdapter } from '../lib/store';
import { UsesSistent } from './SistentWrapper';

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

const PaneSection = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.secondary.elevatedComponents,
  padding: theme.spacing(2.5),
  borderRadius: 4,
}));

const AlreadyConfigured = styled('div')(({ theme }) => ({
  textAlign: 'center',
  padding: theme.spacing(20),
}));

const MesheryPlayComponent = (props) => {
  const { meshAdapters, adapter: selectedAdapterProp, setAdapter } = props;
  const router = useRouter();

  // Initialize state
  const [adapter, setAdapterState] = useState(() => {
    if (meshAdapters && meshAdapters.length > 0) {
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
    } else if (meshAdapters.length > 0) {
      setAdapterState(meshAdapters[0]);
    }
  };

  useEffect(() => {
    router.events.on('routeChangeComplete', handleRouteChange);

    return () => {
      router.events.off('routeChangeComplete', handleRouteChange);
    };
  }, [router.events, meshAdapters]);

  useEffect(() => {
    if (meshAdapters?.length > 0) {
      handleRouteChange();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [meshAdapters?.length]);

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

  const handleAdapterChange = (event) => {
    if (event.target.value !== '') {
      const selectedAdapter = meshAdapters.find(
        ({ adapter_location }) => adapter_location === event.target.value,
      );
      if (selectedAdapter) {
        setAdapterState(selectedAdapter);
        setAdapter({ selectedAdapter: selectedAdapter.name });
      }
    }
  };

  const renderIndividualAdapter = () => {
    let adapCount = 0;
    let adapterToRender;
    meshAdapters.forEach((adap) => {
      if (adap.adapter_location === props.adapter) {
        adapterToRender = adap;
        meshAdapters.forEach((ad) => {
          if (ad.name === adap.name) adapCount += 1;
        });
      }
    });
    if (adapterToRender) {
      const imageIcon = pickImage(adapterToRender);
      return (
        <MesheryAdapterPlayComponent
          adapter={adapterToRender}
          adapCount={adapCount}
          adapter_icon={imageIcon}
        />
      );
    }
    return null;
  };

  if (meshAdapters.length === 0) {
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

  if (selectedAdapterProp && selectedAdapterProp !== '') {
    const indContent = renderIndividualAdapter();
    if (indContent) {
      return indContent;
    }
  }

  const imageIcon = pickImage(adapter);

  return (
    <UsesSistent>
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
                  value={adapter?.adapter_location || ''}
                  margin="normal"
                  variant="outlined"
                  onChange={handleAdapterChange}
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
                      // Removed getContentAnchorEl as it's deprecated in MUI v5
                    },
                  }}
                >
                  {meshAdapters.map((ada) => (
                    <MenuItem key={ada.adapter_location} value={ada.adapter_location}>
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
        {adapter?.adapter_location && (
          <MesheryAdapterPlayComponent adapter={adapter} adapter_icon={imageIcon} />
        )}
      </NoSsr>
    </UsesSistent>
  );
};

MesheryPlayComponent.propTypes = {
  meshAdapters: PropTypes.array.isRequired,
  setAdapter: PropTypes.func.isRequired,
  adapter: PropTypes.string,
};

// Mapping state to props
const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const meshAdapters = state.get('meshAdapters');
  const meshAdaptersts = state.get('meshAdaptersts');
  const selectedAdapter = state.get('selectedAdapter');
  return { k8sconfig, meshAdapters, meshAdaptersts, adapter: selectedAdapter };
};

// Mapping dispatch to props
const mapDispatchToProps = (dispatch) => ({
  setAdapter: bindActionCreators(setAdapter, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(MesheryPlayComponent);
