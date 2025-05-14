import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { TextField, Grid, Button, Chip, MenuItem, useTheme, styled, Box } from '@layer5/sistent';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress } from '../../../lib/store';
import { trueRandom } from '../../../lib/trueRandom';
import { useLazyGetGrafanaTemplateVarsQuery } from '@/rtk-query/telemetry';

const GrafanaRoot = styled(Box)(() => {
  const theme = useTheme();
  return {
    padding: theme.spacing(5),
    backgroundColor: theme.palette.background.card,
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  };
});

const ButtonContainer = styled(Box)({
  display: 'flex',
  justifyContent: 'flex-end',
});

const StyledButton = styled(Button)(({ theme }) => ({
  marginTop: theme.spacing(3),
}));

const AlignRight = styled(Box)(({ theme }) => ({
  textAlign: 'right',
  marginBottom: theme.spacing(2),
}));

const GrafanaIcon = styled('img')(({ theme }) => ({
  width: theme.spacing(2.5),
}));

const PanelChips = styled(Box)({
  display: 'flex',
  flexWrap: 'wrap',
});

const StyledChip = styled(Chip)(({ theme }) => ({
  margin: theme.spacing(0.25),
}));

const StyledTextField = styled(TextField)(({ theme }) => ({
  margin: theme.spacing(1),
}));

function GrafanaSelectionComponent(props) {
  const {
    grafanaURL,
    connectionID,
    handleGrafanaBoardSearchChange,
    handleGrafanaChipDelete,
    handleGrafanaClick,
    addSelectedBoardPanelConfig: addSelectedBoardPanelConfigProp,
    handleError,
    updateProgress,
    grafanaBoards: propGrafanaBoards,
  } = props;

  const [grafanaBoards, setGrafanaBoards] = useState([]);
  const [grafanaBoard, setGrafanaBoard] = useState('');
  const [templateVars, setTemplateVars] = useState([]);
  const [templateVarOptions, setTemplateVarOptions] = useState([]); // will contain the array of options at the respective index, ex: [[v1], [v3, v4]]
  const [panels, setPanels] = useState([]);
  const [selectedPanels, setSelectedPanels] = useState([]);
  const [selectedTemplateVars, setSelectedTemplateVars] = useState([]); // will contain the selected val at the respective index: [v1, v3]

  const [getTemplateVars] = useLazyGetGrafanaTemplateVarsQuery();
  useEffect(() => {
    if (
      JSON.stringify([...grafanaBoards].sort()) !== JSON.stringify([...propGrafanaBoards].sort())
    ) {
      setGrafanaBoards(propGrafanaBoards);
      setGrafanaBoard('');
      setSelectedPanels([]);
    }
  }, [propGrafanaBoards]);

  const handleChange = (name) => (event) => {
    if (name === 'grafanaBoard') {
      boardChange(event.target.value);
    } else if (name.startsWith('template_var_')) {
      setSelectedTemplateVar(parseInt(name.replace('template_var_', '')), event.target.value);
    } else if (name === 'selectedPanels') {
      setSelectedPanels(event.target.value);
    }
  };

  const getSelectedTemplateVar = (ind) =>
    selectedTemplateVars[ind] !== undefined ? selectedTemplateVars[ind] : undefined;

  const setSelectedTemplateVar = (ind, val) => {
    const newSelectedTemplateVars = [...selectedTemplateVars];
    newSelectedTemplateVars[ind] = val;
    for (let i = ind + 1; i < newSelectedTemplateVars.length; i++) {
      newSelectedTemplateVars[i] = '';
    }
    setSelectedTemplateVars(newSelectedTemplateVars);
    if (templateVars.length > ind + 1) {
      queryTemplateVars(ind + 1);
    }
  };

  const boardChange = (newVal) => {
    propGrafanaBoards.forEach((board) => {
      if (board.uri === newVal) {
        setGrafanaBoard(newVal);
        setPanels(board.panels);
        setSelectedPanels(board.panels?.map((panel) => panel.id) || []);
        setTemplateVars(
          board.template_vars && board.template_vars.length > 0 ? board.template_vars : [],
        );
        setTemplateVarOptions([]);
        setSelectedTemplateVars([]);
        if (board.template_vars && board.template_vars.length > 0) {
          queryTemplateVars(0);
        }
      }
    });
  };

  const queryTemplateVars = (ind) => {
    if (templateVars.length > 0) {
      const selectedVars = [];
      for (let i = 0; i < ind; i++) {
        if (selectedTemplateVars[i]) {
          selectedVars.push({
            name: templateVars[i].name,
            value: selectedTemplateVars[i],
          });
        }
      }

      let startTime, endTime;
      if (
        templateVars[ind].query.startsWith('label_values') &&
        templateVars[ind].query.indexOf(',') > -1
      ) {
        // series query needs a start and end time or else it will take way longer to return. . .
        // but at this point this component does not have the time range selection bcoz the time range selection comes after this component makes its selections
        // hence for now just limiting the time period to the last 24hrs

        const ed = new Date();
        const sd = new Date();
        sd.setDate(sd.getDate() - 1);
        startTime = Math.floor(sd.getTime() / 1000);
        endTime = Math.floor(ed.getTime() / 1000);
      }
      updateProgress({ showProgress: true });
      getTemplateVars({
        connectionID,
        templateVarQuery: templateVars[ind].query,
        datasourceId: templateVars[ind].datasource.id,
        selectedVars,
        startTime,
        endTime,
      })
        .unwrap()
        .then((result) => {
          updateProgress({ showProgress: false });
          if (typeof result !== 'undefined') {
            let tmpVarOpts = [];
            if (Array.isArray(result.data)) {
              if (result.data.length > 0) {
                if (
                  templateVars[ind].query.startsWith('label_values') &&
                  templateVars[ind].query.indexOf(',') > -1 &&
                  typeof result.data[0] === 'object'
                ) {
                  let q = templateVars[ind].query.replace('label_values(', '');
                  q = q.substr(0, q.length - 1);
                  const qInd = q.substr(q.lastIndexOf(',')).replace(',', '').trim();
                  result.data.forEach((d) => {
                    const v = d[qInd];
                    if (typeof v !== 'undefined' && v !== null && tmpVarOpts.indexOf(v) === -1) {
                      tmpVarOpts.push(v);
                    }
                  });
                } else {
                  tmpVarOpts = result.data;
                }
              }
            } else {
              tmpVarOpts = result.data.result.map(({ metric }) => {
                const tmpArr = Object.keys(metric);
                if (tmpArr.length > 0) {
                  return metric[tmpArr[0]];
                }
              });
            }
            setTemplateVarOptions((prev) => {
              const newOptions = [...prev];
              newOptions[ind] = tmpVarOpts;
              return newOptions;
            });
          }
        })
        .catch((error) => {
          updateProgress({ showProgress: false });
          setTemplateVarOptions((prev) => {
            const newOptions = [...prev];
            newOptions[ind] = [templateVars[ind].Value];
            return newOptions;
          });
          handleError(error);
        });
    }
  };

  const addSelectedBoardPanelConfig = () => {
    let boardConfig = {};
    grafanaBoards.forEach((board) => {
      if (grafanaBoard === board.uri) {
        boardConfig.board = board;
      }
    });

    boardConfig.panels = panels.filter(({ id }) => selectedPanels.indexOf(id) > -1);

    boardConfig.templateVars = templateVars.map(({ name }, index) =>
      typeof selectedTemplateVars[index] !== 'undefined'
        ? `${name}=${selectedTemplateVars[index]}`
        : '',
    );
    addSelectedBoardPanelConfigProp(boardConfig);
  };

  const genRandomNumberForKey = () => Math.floor(trueRandom() * 1000 + 1);

  return (
    <NoSsr>
      <GrafanaRoot>
        <AlignRight>
          <StyledChip
            label={grafanaURL}
            onDelete={handleGrafanaChipDelete}
            onClick={handleGrafanaClick}
            icon={<GrafanaIcon src="/static/img/grafana_icon.svg" />}
            key="graf-key"
            variant="outlined"
          />
        </AlignRight>
        <Grid container spacing={1}>
          <Grid item xs={12} sm={6}>
            <StyledTextField
              id="grafanaBoardSearch"
              name="grafanaBoardSearch"
              label="Board Search"
              fullWidth
              value={props.grafanaBoardSearch}
              variant="outlined"
              onChange={handleGrafanaBoardSearchChange('grafanaBoardSearch')}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <StyledTextField
              select
              id="grafanaBoard"
              name="grafanaBoard"
              label="Board"
              fullWidth
              value={grafanaBoard}
              variant="outlined"
              onChange={handleChange('grafanaBoard')}
            >
              {grafanaBoards?.map((board) => (
                <MenuItem key={`bd_---_${board.uri}`} value={board.uri}>
                  {board.title}
                </MenuItem>
              ))}
            </StyledTextField>
          </Grid>

          {templateVars.length > 0 &&
            templateVars.map(({ name }, ind) => {
              if (ind === 0 || typeof getSelectedTemplateVar(ind - 1) !== 'undefined') {
                return (
                  <Grid item xs={12} sm={4} key={ind}>
                    <StyledTextField
                      select
                      id={`template_var_${ind}`}
                      name={`template_var_${ind}`}
                      label={`Template variable: ${name}`}
                      fullWidth
                      value={getSelectedTemplateVar(ind)}
                      variant="outlined"
                      onChange={handleChange(`template_var_${ind}`)}
                    >
                      <MenuItem
                        key={`tmplVarOpt__-___${ind}_${genRandomNumberForKey()}`}
                        value=""
                      />
                      {templateVarOptions[ind]?.map((opt) => (
                        <MenuItem
                          key={`tmplVarOpt__-__${name}_${opt}_${ind}_${genRandomNumberForKey()}`}
                          value={opt}
                        >
                          {opt}
                        </MenuItem>
                      ))}
                    </StyledTextField>
                  </Grid>
                );
              }
              return null;
            })}

          <Grid item xs={12}>
            <StyledTextField
              select
              id="panels"
              name="panels"
              label="Panels"
              fullWidth
              value={selectedPanels}
              variant="outlined"
              onChange={handleChange('selectedPanels')}
              SelectProps={{
                multiple: true,
                renderValue: (selected) => (
                  <PanelChips>
                    {selected.map((value) => {
                      let selVal = '';
                      let panelId = '';
                      panels.forEach((panel) => {
                        if (panel.id === value) {
                          selVal = panel.title;
                          panelId = panel.id;
                        }
                      });
                      return <StyledChip key={`pl_--_${panelId}`} label={selVal} />;
                    })}
                  </PanelChips>
                ),
              }}
            >
              {panels?.map((panel) => (
                <MenuItem key={`panel_-__-${panel.id}`} value={panel.id}>
                  {panel.title}
                </MenuItem>
              ))}
            </StyledTextField>
          </Grid>
        </Grid>

        <ButtonContainer>
          <StyledButton
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            onClick={addSelectedBoardPanelConfig}
          >
            Add
          </StyledButton>
        </ButtonContainer>
      </GrafanaRoot>
    </NoSsr>
  );
}

GrafanaSelectionComponent.propTypes = {
  grafanaURL: PropTypes.string.isRequired,
  handleGrafanaBoardSearchChange: PropTypes.func.isRequired,
  handleGrafanaChipDelete: PropTypes.func.isRequired,
  handleGrafanaClick: PropTypes.func.isRequired,
  addSelectedBoardPanelConfig: PropTypes.func.isRequired,
  handleError: PropTypes.func.isRequired,
  updateProgress: PropTypes.func.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});
const mapStateToProps = () => ({});

export default connect(mapStateToProps, mapDispatchToProps)(GrafanaSelectionComponent);
