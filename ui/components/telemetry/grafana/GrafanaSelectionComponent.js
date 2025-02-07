import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@mui/material';
import { TextField, Grid, Button, Chip, MenuItem, useTheme, styled, Box } from '@layer5/sistent';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import dataFetch from '../../../lib/data-fetch';
import { updateProgress } from '../../../lib/store';
import { trueRandom } from '../../../lib/trueRandom';
import { UsesSistent } from '@/components/SistentWrapper';

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

class GrafanaSelectionComponent extends Component {
  constructor(props) {
    super(props);

    this.state = {
      grafanaBoards: [],
      grafanaBoard: '',
      grafanaBoardObject: {},
      templateVars: [],
      templateVarOptions: [], // will contain the array of options at the respective index, ex: [[v1], [v3, v4]]

      panels: [],
      selectedPanels: [],
      selectedTemplateVars: [], // will contain the selected val at the respective index: [v1, v3]
    };
  }

  handleChange = (name) => (event) => {
    if (name === 'grafanaBoard') {
      this.boardChange(event.target.value);
    } else if (name.startsWith('template_var_')) {
      this.setSelectedTemplateVar(parseInt(name.replace('template_var_', '')), event.target.value);
    } else {
      this.setState({ [name]: event.target.value });
    }
  };

  getSelectedTemplateVar = (ind) => {
    const { selectedTemplateVars } = this.state;
    const retVal =
      typeof selectedTemplateVars[ind] !== 'undefined' ? selectedTemplateVars[ind] : undefined;
    return retVal;
  };

  setSelectedTemplateVar = (ind, val) => {
    const { templateVars, templateVarOptions, selectedTemplateVars } = this.state;
    selectedTemplateVars[ind] = val;
    for (let i = ind + 1; i < selectedTemplateVars.length; i++) {
      selectedTemplateVars[i] = '';
    }
    this.setState({ selectedTemplateVars });
    if (templateVars.length > ind + 1) {
      this.queryTemplateVars(ind + 1, templateVars, templateVarOptions, selectedTemplateVars);
    }
  };

  boardChange = (newVal) => {
    this.props.grafanaBoards.forEach((board) => {
      if (board.uri === newVal) {
        this.setState({
          grafanaBoard: newVal,
          panels: board.panels,
          selectedPanels: board.panels?.map((panel) => panel.id), // selecting all panels by default
          templateVars:
            board.template_vars && board.template_vars.length > 0 ? board.template_vars : [],
          templateVarOptions: [],
          selectedTemplateVars: [],
        });
        if (board.template_vars && board.template_vars.length > 0) {
          this.queryTemplateVars(0, board.template_vars, [], []);
        }
      }
    });
  };

  queryTemplateVars = (ind, templateVars, templateVarOptions, selectedTemplateVars) => {
    if (templateVars.length > 0) {
      let queryURL = `/api/telemetry/metrics/grafana/query/${
        this.props.connectionID
      }?query=${encodeURIComponent(templateVars[ind].query)}&dsid=${
        templateVars[ind].datasource.id
      }`;
      for (let i = ind; i > 0; i--) {
        queryURL += `&${templateVars[i - 1].name}=${selectedTemplateVars[i - 1]}`;
      }
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
        queryURL += `&start=${Math.floor(sd.getTime() / 1000)}&end=${Math.floor(
          ed.getTime() / 1000,
        )}`; // accounts for the last 24hrs
      }
      this.props.updateProgress({ showProgress: true });
      dataFetch(
        queryURL,
        { credentials: 'include' },
        (result) => {
          this.props.updateProgress({ showProgress: false });
          if (typeof result !== 'undefined') {
            let tmpVarOpts = [];
            // result.data check if it is an array or object
            if (Array.isArray(result.data)) {
              if (result.data.length > 0) {
                if (
                  templateVars[ind].query.startsWith('label_values') &&
                  templateVars[ind].query.indexOf(',') > -1 &&
                  typeof result.data[0] === 'object'
                ) {
                  let q = templateVars[ind].query.replace('label_values(', '');
                  q = q.substr(0, q.length - 1); // to remove ')'
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
            templateVarOptions[ind] = tmpVarOpts;
            this.setState({ templateVarOptions });
          }
        },
        (error) => {
          templateVarOptions[ind] = [templateVars[ind].Value];
          this.setState({ templateVarOptions });
          this.props.handleError(error);
        },
      );
    }
  };

  static getDerivedStateFromProps(props, state) {
    if (JSON.stringify(state.grafanaBoards.sort()) !== JSON.stringify(props.grafanaBoards.sort())) {
      return { grafanaBoards: props.grafanaBoards, grafanaBoard: '', selectedPanels: [] };
    }
    return null;
  }

  addSelectedBoardPanelConfig = () => {
    const {
      grafanaBoard,
      grafanaBoards,
      templateVars,
      selectedTemplateVars,
      selectedPanels,
      panels,
    } = this.state;
    const boardConfig = {};
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

    this.props.addSelectedBoardPanelConfig(boardConfig);
  };

  genRandomNumberForKey = () => Math.floor(trueRandom() * 1000 + 1);

  render() {
    const {
      grafanaBoardSearch,
      grafanaURL,
      handleGrafanaBoardSearchChange,
      handleGrafanaChipDelete,
      handleGrafanaClick,
    } = this.props;
    const {
      panels,
      grafanaBoards,
      selectedPanels,
      grafanaBoard,
      templateVars,
      templateVarOptions,
    } = this.state;
    return (
      <UsesSistent>
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
                  value={grafanaBoardSearch}
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
                  onChange={this.handleChange('grafanaBoard')}
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
                  if (ind === 0 || typeof this.getSelectedTemplateVar(ind - 1) !== 'undefined') {
                    return (
                      <Grid item xs={12} sm={4} key={ind}>
                        <StyledTextField
                          select
                          id={`template_var_${ind}`}
                          name={`template_var_${ind}`}
                          label={`Template variable: ${name}`}
                          fullWidth
                          value={this.getSelectedTemplateVar(ind)}
                          variant="outlined"
                          onChange={this.handleChange(`template_var_${ind}`)}
                        >
                          <MenuItem
                            key={`tmplVarOpt__-___${ind}_${this.genRandomNumberForKey()}`}
                            value=""
                          />
                          {templateVarOptions[ind]?.map((opt) => (
                            <MenuItem
                              key={`tmplVarOpt__-__${name}_${opt}_${ind}_${this.genRandomNumberForKey()}`}
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
                  onChange={this.handleChange('selectedPanels')}
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
                onClick={this.addSelectedBoardPanelConfig}
              >
                Add
              </StyledButton>
            </ButtonContainer>
          </GrafanaRoot>
        </NoSsr>
      </UsesSistent>
    );
  }
}

GrafanaSelectionComponent.propTypes = {
  grafanaURL: PropTypes.string.isRequired,
  //grafanaBoards: PropTypes.array.isRequired,
  handleGrafanaBoardSearchChange: PropTypes.func.isRequired,
  handleGrafanaChipDelete: PropTypes.func.isRequired,
  handleGrafanaClick: PropTypes.func.isRequired,
  addSelectedBoardPanelConfig: PropTypes.func.isRequired,
  handleError: PropTypes.func.isRequired,
  updateProgress: PropTypes.func.isRequired, // Added for completeness
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = () => ({});

export default connect(mapStateToProps, mapDispatchToProps)(GrafanaSelectionComponent);
