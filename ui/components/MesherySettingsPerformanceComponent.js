/* eslint-disable */
import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import {
  NoSsr,
  Tooltip,
  CircularProgress,
  FormControl,
  RadioGroup,
  FormControlLabel,
  Radio,
  Button,
  Grid,
  TextField,
} from '@material-ui/core';
import SaveOutlinedIcon from '@material-ui/icons/SaveOutlined';
import { Autocomplete } from '@material-ui/lab';
import { withStyles } from '@material-ui/core/styles';
import dataFetch from '../lib/data-fetch';
import { updateLoadTestPref, updateProgress } from '../lib/store';
import { durationOptions } from '../lib/prePopulatedOptions';
import { EVENT_TYPES } from '../lib/event-types';
import { withNotify } from '../utils/hooks/useNotification';
import { ctxUrl } from '../utils/multi-ctx';

const loadGenerators = ['fortio', 'wrk2', 'nighthawk'];

const styles = (theme) => ({
  root: {
    padding: theme.spacing(10),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  margin: {
    margin: theme.spacing(1),
  },
  radio: {
    '&.Mui-checked': {
      color: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
});

const MesherySettingsPerformanceComponent = ({
  qps,
  c,
  t,
  updateProgress,
  selectedK8sContexts,
  notify,
  updateLoadTestPref,
  classes,
  gen,
}) => {
  const [tValue, settValue] = useState('');
  const [tError, settError] = useState('');
  const [blockRunTest, setBlockRunTest] = useState('');
  const [queriesPerSecond, setQueriesPerSecond] = useState(qps);
  const [concurrentRequests, setConcurrentRequests] = useState(c);
  const [timeDuration, setTimeDuration] = useState(t);
  const [generator, setGenerator] = useState(gen);

  const handleChange = (name) => (event) => {
    const value = event.target.value;
    switch (name) {
      case 'qps':
        setQueriesPerSecond(parseInt(value));
        break;
      case 'c':
        setConcurrentRequests(parseInt(value));
        break;
      case 'gen':
        setGenerator(value);
        break;
      default:
        break;
    }
  };

  const handleDurationChange = (event, newValue) => {
    settValue(newValue);
    if (newValue !== null) {
      settError('');
    }
  };

  const handleInputDurationChange = (event, newValue) => {
    setTimeDuration(newValue);
  };

  const handleSubmit = () => {
    let err = false;
    let tNum = 0;
    try {
      tNum = parseInt(timeDuration.substring(0, timeDuration.length - 1));
    } catch (ex) {
      err = true;
    }

    if (
      timeDuration === '' ||
      !(
        timeDuration.toLowerCase().endsWith('h') ||
        timeDuration.toLowerCase().endsWith('m') ||
        timeDuration.toLowerCase().endsWith('s')
      ) ||
      err ||
      tNum <= 0
    ) {
      settError('error-autocomplete-value');
      return;
    }

    submitPerfPreference();
  };

  const submitPerfPreference = () => {
    const loadTestPrefs = {
      qps: queriesPerSecond,
      c: concurrentRequests,
      t: timeDuration,
      gen: generator,
    };

    const requestBody = JSON.stringify({ loadTestPrefs });
    setBlockRunTest(true); // to block the button
    updateProgress({ showProgress: true });
    dataFetch(
      ctxUrl('/api/user/prefs', selectedK8sContexts),
      {
        credentials: 'same-origin',
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json;charset=UTF-8',
        },
        body: requestBody,
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          notify({ message: 'Preferences saved', event_type: EVENT_TYPES.SUCCESS });
          updateLoadTestPref({
            loadTestPref: {
              qps: queriesPerSecond,
              c: concurrentRequests,
              t: timeDuration,
              gen: generator,
            },
          });
          setBlockRunTest(false);
        }
      },
      handleError('There was an error saving your preferences'),
    );
  };

  const getLoadTestPrefs = () => {
    dataFetch(
      ctxUrl('/api/user/prefs', selectedK8sContexts),
      {
        credentials: 'same-origin',
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (result) {
          setQueriesPerSecond(result.loadTestPrefs.qps);
          setConcurrentRequests(result.loadTestPrefs.c);
          setTimeDuration(result.loadTestPrefs.t);
          setGen(result.loadTestPrefs.gen);
        }
      },
      () => {
        !qps || !t || !c ? handleError('There was an error fetching your preferences') : {};
      },
    );
  };

  const handleError = (msg) => {
    return (error) => {
      setBlockRunTest(false);
      let finalMsg = msg;
      if (typeof error === 'string') {
        finalMsg = `${msg}: ${error}`;
      }
      notify({ message: finalMsg, event_type: EVENT_TYPES.ERROR, details: error.toString() });
    };
  };

  useEffect(() => {
    getLoadTestPrefs();
  }, []);

  return (
    <NoSsr>
      <React.Fragment>
        <div className={classes.root}>
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
                value={concurrentRequests}
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
                value={queriesPerSecond}
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
                  inputValue={timeDuration}
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
                  value={generator}
                  onChange={handleChange('gen')}
                  row
                >
                  {loadGenerators.map((lg) => (
                    <FormControlLabel
                      value={lg}
                      control={
                        <Radio color="primary" disabled={lg === 'wrk2'} className={classes.radio} />
                      }
                      label={lg}
                    />
                  ))}
                </RadioGroup>
              </FormControl>
            </Grid>
          </Grid>
          <div className={classes.buttons}>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              size="large"
              onClick={handleSubmit}
              className={classes.button}
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

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withNotify(MesherySettingsPerformanceComponent)),
);
