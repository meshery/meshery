import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { TextField, Grid, Button, Chip, MenuItem, styled, NoSsr, Alert } from '@layer5/sistent';
import { Controlled as CodeMirror } from 'react-codemirror2';
import { trueRandom } from '../../../lib/trueRandom';
import { usePostBoardImportMutation, useLazyQueryTemplateVarsQuery } from '@/rtk-query/connection';
import CodeIcon from '@mui/icons-material/Code';
import { updateProgress } from '@/store/slices/mesheryUi';

const PrometheusContainer = styled('div')(({ theme }) => ({
  padding: theme.spacing(5),
}));

const ButtonContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  justifyContent: 'flex-end',
  '& .actionButton': {
    marginTop: theme.spacing(3),
  },
}));

const AlignRight = styled('div')(({ theme }) => ({
  textAlign: 'right',
  marginBottom: theme.spacing(2),
}));

const PanelChipsContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  flexWrap: 'wrap',
  '& .chip': {
    margin: theme.spacing(0.25),
  },
}));

const StyledIcon = styled('img')(({ theme }) => ({
  width: theme.spacing(2.5),
}));

const FormControlWrapper = styled('div')(({ theme }) => ({
  marginTop: theme.spacing(2),
  minWidth: window.innerWidth * 0.25,
}));

const dummyBoard = `
{
  "id": null,
  "uid": "1234-exam-ple",
  "title": "Dashboard",
  "tags": [],
  "timezone": "browser",
  "editable": true,
  "hideControls": false,
  "graphTooltip": 1,
  "panels": [],
  "time": {
    "from": "now-3h",
    "to": "now"
  },
  "timepicker": {
    "time_options": [],
    "refresh_intervals": []
  },
  "templating": {
    "list": []
  },
  "annotations": {
    "list": []
  },
  "refresh": "10s",
  "schemaVersion": 17,
  "version": 0,
  "links": []
}
`;

const PrometheusSelectionComponent = (props) => {
  const {
    prometheusURL,
    handlePrometheusChipDelete,
    handlePrometheusClick,
    addSelectedBoardPanelConfig,
    handleError,
    connectionID,
  } = props;

  const [grafanaBoard, setGrafanaBoard] = useState(dummyBoard);
  const [grafanaBoardObject, setGrafanaBoardObject] = useState({});
  const [templateVars, setTemplateVars] = useState([]);
  const [templateVarOptions, setTemplateVarOptions] = useState([]);
  const [panels, setPanels] = useState([]);
  const [selectedPanels, setSelectedPanels] = useState([]);
  const [selectedTemplateVars, setSelectedTemplateVars] = useState([]);

  const cmEditorRef = useRef(null);
  const boardTimeoutRef = useRef(null);

  useEffect(() => {
    return () => {
      if (boardTimeoutRef.current) {
        clearTimeout(boardTimeoutRef.current);
      }
    };
  }, []);

  const getSelectedTemplateVar = (ind) => {
    return selectedTemplateVars[ind] ?? '';
  };

  const setSelectedTemplateVar = (ind, val) => {
    setSelectedTemplateVars((prev) => {
      const newVars = [...prev];
      newVars[ind] = val;
      for (let i = ind + 1; i < newVars.length; i++) {
        newVars[i] = '';
      }

      if (templateVars.length > ind + 1) {
        queryTemplateVars(ind + 1, templateVars, templateVarOptions, newVars);
      }

      return newVars;
    });
  };

  const [postBoardImport] = usePostBoardImportMutation();
  const [triggerTemplateQuery] = useLazyQueryTemplateVarsQuery();

  const boardChange = async (newVal) => {
    updateProgress({ showProgress: true });
    try {
      const result = await postBoardImport({
        connectionID,
        body: newVal,
      }).unwrap();

      const filteredPanels =
        result.panels?.filter((panel) =>
          panel.targets?.some((t) => t.datasource?.type?.toLowerCase() === 'prometheus'),
        ) || [];

      if (filteredPanels.length === 0) {
        return handleError('No panels found with target datasource as prometheus.');
      }

      setGrafanaBoardObject(result);
      setPanels(filteredPanels);
      setSelectedPanels(filteredPanels.map((p) => p.id));
      setTemplateVars(result.template_vars || []);
      setTemplateVarOptions([]);
      setSelectedTemplateVars([]);

      if (result.template_vars?.length > 0) {
        queryTemplateVars(0, result.template_vars, [], []);
      }
    } catch (error) {
      handleError(error);
    } finally {
      updateProgress({ showProgress: false });
    }
  };

  const queryTemplateVars = async (ind, vars, options, selectedVars) => {
    if (!vars?.[ind]) return;

    let queryStr = `query=${encodeURIComponent(vars[ind].query)}`;
    for (let i = ind; i > 0; i--) {
      queryStr += `&${vars[i - 1].name}=${selectedVars[i - 1]}`;
    }

    if (vars[ind].query.startsWith('label_values') && vars[ind].query.includes(',')) {
      const ed = new Date();
      const sd = new Date();
      sd.setDate(sd.getDate() - 1);
      queryStr += `&start=${Math.floor(sd.getTime() / 1000)}&end=${Math.floor(ed.getTime() / 1000)}`;
    }

    updateProgress({ showProgress: true });
    try {
      const result = await triggerTemplateQuery({ connectionID, query: queryStr }).unwrap();
      let tmpVarOpts = [];

      if (Array.isArray(result?.data)) {
        if (
          vars[ind].query.startsWith('label_values') &&
          vars[ind].query.includes(',') &&
          typeof result.data[0] === 'object'
        ) {
          const q = vars[ind].query.replace('label_values(', '').slice(0, -1);
          const qInd = q.split(',').pop().trim();
          result.data.forEach((d) => {
            const v = d[qInd];
            if (v && !tmpVarOpts.includes(v)) tmpVarOpts.push(v);
          });
        } else {
          tmpVarOpts = result.data;
        }
      } else if (result?.data?.result) {
        tmpVarOpts = result.data.result.map(({ metric }) => Object.values(metric)?.[0]);
      }

      setTemplateVarOptions((prev) => {
        const newOptions = [...prev];
        newOptions[ind] = tmpVarOpts;
        return newOptions;
      });
    } catch (error) {
      setTemplateVarOptions((prev) => {
        const newOptions = [...prev];
        newOptions[ind] = [vars[ind].Value];
        return newOptions;
      });
      handleError(error);
    } finally {
      updateProgress({ showProgress: false });
    }
  };

  const handleTemplateVarChange = (ind) => (e) => {
    setSelectedTemplateVar(ind, e.target.value);
  };

  const handlePanelSelection = (e) => {
    setSelectedPanels(e.target.value);
  };

  const addBoardConfig = () => {
    const boardConfig = {
      board: grafanaBoardObject,
      panels: panels.filter((p) => selectedPanels.includes(p.id)),
      templateVars: templateVars.map(({ name }, i) =>
        selectedTemplateVars[i] ? `${name}=${selectedTemplateVars[i]}` : '',
      ),
    };
    addSelectedBoardPanelConfig(boardConfig);
  };

  const genRandomNumberForKey = () => Math.floor(trueRandom() * 1000 + 1);

  const handleCodeChange = (editor, data, value) => {
    setGrafanaBoard(value);
    setGrafanaBoardObject({});
    setPanels([]);
    setSelectedPanels([]);
    setTemplateVars([]);
    setTemplateVarOptions([]);
    setSelectedTemplateVars([]);

    if (boardTimeoutRef.current) {
      clearTimeout(boardTimeoutRef.current);
    }

    boardTimeoutRef.current = setTimeout(() => {
      if (cmEditorRef.current?.state.lint.marked.length === 0) {
        boardChange(value);
      }
    }, 1000);
  };

  return (
    <NoSsr>
      <PrometheusContainer>
        <AlignRight>
          <Chip
            label={prometheusURL}
            onDelete={handlePrometheusChipDelete}
            onClick={handlePrometheusClick}
            icon={
              <StyledIcon src="/static/img/prometheus_logo_orange_circle.svg" alt="Prometheus" />
            }
            variant="outlined"
          />
        </AlignRight>

        <Grid container spacing={1}>
          <Grid item xs={12}>
            <div style={{ padding: '20px', display: 'flex' }}>
              <CodeIcon style={{ marginRight: '6px' }} />
              Paste your custom board JSON below.
            </div>

            <CodeMirror
              editorDidMount={(editor) => {
                cmEditorRef.current = editor;
              }}
              value={grafanaBoard}
              options={{
                theme: 'material',
                lineNumbers: true,
                lineWrapping: true,
                gutters: ['CodeMirror-lint-markers'],
                lint: true,
                mode: 'application/json',
              }}
              onBeforeChange={handleCodeChange}
              onChange={() => {}}
            />
          </Grid>

          {templateVars.map(
            ({ name }, ind) =>
              (ind === 0 || selectedTemplateVars[ind - 1]) && (
                <Grid item xs={12} sm={4} key={ind}>
                  <FormControlWrapper>
                    <TextField
                      select
                      id={`template_var_${ind}`}
                      name={`template_var_${ind}`}
                      label={`Template variable: ${name}`}
                      fullWidth
                      value={getSelectedTemplateVar(ind)}
                      margin="normal"
                      variant="outlined"
                      onChange={handleTemplateVarChange(ind)}
                    >
                      <MenuItem value="">
                        <em>None</em>
                      </MenuItem>
                      {templateVarOptions[ind]?.map((opt) => (
                        <MenuItem
                          key={`tmplVarOpt_${name}_${opt}_${ind}_${genRandomNumberForKey()}`}
                          value={opt}
                        >
                          {opt}
                        </MenuItem>
                      ))}
                    </TextField>
                  </FormControlWrapper>
                </Grid>
              ),
          )}

          {panels.length === 0 && (
            <Grid item xs={12} style={{ marginTop: '10px' }}>
              <Alert severity="error">
                Please load a valid Board JSON to be able to view the panels
              </Alert>
            </Grid>
          )}

          {panels.length > 0 && (
            <Grid item xs={12}>
              <TextField
                select
                id="panels"
                name="panels"
                label="Panels"
                fullWidth
                value={selectedPanels}
                margin="normal"
                variant="outlined"
                onChange={handlePanelSelection}
                SelectProps={{
                  multiple: true,
                  renderValue: (selected) => (
                    <PanelChipsContainer>
                      {selected.map((value) => {
                        const panel = panels.find((p) => p.id === value);
                        return (
                          <Chip key={`pl_${value}`} label={panel?.title || ''} className="chip" />
                        );
                      })}
                    </PanelChipsContainer>
                  ),
                }}
              >
                {panels.map((panel) => (
                  <MenuItem key={`panel_${panel.id}`} value={panel.id}>
                    {panel.title}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
          )}
        </Grid>

        {selectedPanels.length > 0 && (
          <ButtonContainer>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              size="large"
              onClick={addBoardConfig}
              className="actionButton"
            >
              Add
            </Button>
          </ButtonContainer>
        )}
      </PrometheusContainer>
    </NoSsr>
  );
};

PrometheusSelectionComponent.propTypes = {
  prometheusURL: PropTypes.string.isRequired,
  handlePrometheusClick: PropTypes.func.isRequired,
  handlePrometheusChipDelete: PropTypes.func.isRequired,
  addSelectedBoardPanelConfig: PropTypes.func.isRequired,
  handleError: PropTypes.func.isRequired,
};

export default PrometheusSelectionComponent;
