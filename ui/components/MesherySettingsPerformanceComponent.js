/* eslint-disable */
import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import SaveOutlinedIcon from '@material-ui/icons/SaveOutlined';
import { Autocomplete } from '@material-ui/lab';
import {
  NoSsr,
  Tooltip,
  IconButton,
  CircularProgress,
  FormControl,
  RadioGroup,
  FormControlLabel,
  Radio,
  Button,
  TextField,
  Grid,
} from '@material-ui/core';
import dataFetch from '../lib/data-fetch';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateLoadTestPref, updateProgress } from '../lib/store';
import { durationOptions } from '../lib/prePopulatedOptions';
import { ctxUrl } from '../utils/multi-ctx';
import { withNotify } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';

const loadGenerators = ['fortio', 'wrk2', 'nighthawk'];

const MesherySettingsPerformanceComponent = (props) => {
  const { classes, notify } = props;
  const { qps: initialQps, c: initialC, t: initialT, gen: initialGen } = props;

  const [qps, setQps] = useState(initialQps);
  const [c, setC] = useState(initialC);
  const [t, setT] = useState(initialT);
  const [tValue, setTValue] = useState(initialT);
  const [gen, setGen] = useState(initialGen);
  const [blockRunTest, setBlockRunTest] = useState(false);
  const [tError, setTError] = useState('');

  const handleChange = (name) => (event) => {
    const value = event.target.value;
    if (name === 'qps' || name === 'c') {
      const parsedValue = parseInt(value, 10);
      if (!isNaN(parsedValue)) {
        name === 'qps' ? setQps(parsedValue) : setC(parsedValue);
      }
    } else {
      name === 't' ? setT(value) : setGen(value);
    }
  };

  const handleDurationChange = (event, newValue) => {
    setTValue(newValue);
    if (newValue !== null) {
      setTError('');
    }
  };

  const handleInputDurationChange = (event, newValue) => {
    setT(newValue);
  };

  const handleSubmit = () => {
    try {
      const tNum = parseInt(t.substring(0, t.length - 1), 10);
      if (isNaN(tNum) || tNum <= 0 || !['h', 'm', 's'].includes(t.slice(-1).toLowerCase())) {
        setTError('error-autocomplete-value');
        return;
      }
    } catch {
      setTError('error-autocomplete-value');
      return;
    }

    submitPerfPreference();
  };

  const submitPerfPreference = () => {
    const loadTestPrefs = { qps, c, t, gen };
    const requestBody = JSON.stringify({ loadTestPrefs });

    setBlockRunTest(true);
    props.updateProgress({ showProgress: true });

    dataFetch(
      ctxUrl('/api/user/prefs', props.selectedK8sContexts),
      {
        credentials: 'same-origin',
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json;charset=UTF-8' },
        body: requestBody,
      },
      (result) => {
        props.updateProgress({ showProgress: false });
        if (result) {
          notify({ message: 'Preferences saved', event_type: EVENT_TYPES.SUCCESS });
          props.updateLoadTestPref({ loadTestPref: { qps, c, t, gen } });
          setBlockRunTest(false);
        }
      },
      handleError('There was an error saving your preferences'),
    );
  };

  useEffect(() => {
    getLoadTestPrefs();
  }, []);

  const getLoadTestPrefs = () => {
    dataFetch(
      ctxUrl('/api/user/prefs', props.selectedK8sContexts),
      {
        credentials: 'same-origin',
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (result) {
          setQps(result.loadTestPrefs.qps);
          setC(result.loadTestPrefs.c);
          setT(result.loadTestPrefs.t);
          setGen(result.loadTestPrefs.gen);
        }
      },
      handleError('There was an error fetching your preferences'),
    );
  };

  const handleError = (msg) => (error) => {
    setBlockRunTest(false);
    let finalMsg = msg;
    if (typeof error === 'string') {
      finalMsg = `${msg}: ${error}`;
    }
    notify({ message: finalMsg, event_type: EVENT_TYPES.ERROR, details: error.toString() });
  };

  // const { blockRunTest, qps, t, c, gen, tValue, tError } = state;

  return (
    <NoSsr>
      <React.Fragment>
        <div sx={{ padding: theme.spacing(10) }}>
          <label>
            <strong>Performance Load Test Defaults</strong>
          </label>
          <Grid container spacing={3}>
            <Grid item xs={12} lg={4}>
              <TextField
                required
                id="c"
                name="c"
                label="Concurrent requests"
                type="number"
                fullWidth
                value={c}
                inputProps={{ min: '0', step: '1' }}
                margin="normal"
                variant="outlined"
                onChange={handleChange('c')}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12} lg={4}>
              <TextField
                required
                id="qps"
                name="qps"
                label="Queries per second"
                type="number"
                fullWidth
                value={qps}
                inputProps={{ min: '0', step: '1' }}
                margin="normal"
                variant="outlined"
                onChange={handleChange('qps')}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12} lg={4}>
              <Tooltip
                title={"Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."}
              >
                <Autocomplete
                  required
                  id="t"
                  name="t"
                  freeSolo
                  label="Duration*"
                  fullWidth
                  variant="outlined"
                  className={classes.errorValue}
                  classes={{ root: tError }}
                  value={tValue}
                  inputValue={t}
                  onChange={handleDurationChange}
                  onInputChange={handleInputDurationChange}
                  options={durationOptions}
                  style={{ marginTop: '16px', marginBottom: '8px' }}
                  renderInput={(params) => (
                    <TextField {...params} label="Duration*" variant="outlined" />
                  )}
                />
              </Tooltip>
            </Grid>
            <Grid item xs={12} lg={4}>
              <FormControl component="loadGenerator" className={classes.formControl}>
                <label>
                  <strong>Default Load Generator</strong>
                </label>
                <RadioGroup
                  aria-label="loadGenerator"
                  name="loadGenerator"
                  value={gen}
                  onChange={handleChange('gen')}
                  row
                >
                  {loadGenerators.map((lg) => (
                    <FormControlLabel
                      value={lg}
                      control={
                        <Radio
                          color="primary"
                          disabled={lg === 'wrk2'}
                          sx={{
                            '&.Mui-checked': {
                              color:
                                theme.palette.mode === 'dark'
                                  ? '#00B39F'
                                  : theme.palette.primary.main,
                            },
                          }}
                        />
                      }
                      label={lg}
                    />
                  ))}
                </RadioGroup>
              </FormControl>
            </Grid>
          </Grid>
          <div sx={{ display: 'flex', justifyContent: 'flex-end' }}>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              size="large"
              onClick={handleSubmit}
              sx={{ marginTop: theme.spacing(3), marginLeft: theme.spacing(1) }}
              disabled={blockRunTest}
            >
              <SaveOutlinedIcon style={{ marginRight: '3px' }} />
              {blockRunTest ? <CircularProgress size={30} /> : 'Save'}
            </Button>
          </div>
        </div>
      </React.Fragment>
    </NoSsr>
  );
};

MesherySettingsPerformanceComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateLoadTestPref: bindActionCreators(updateLoadTestPref, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const loadTestPref = state.get('loadTestPref').toJS();
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    ...loadTestPref,
    selectedK8sContexts,
  };
};

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withNotify(MesherySettingsPerformanceComponent));
