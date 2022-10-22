import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {
  NoSsr, TextField, Grid, Button, Chip, MenuItem,
} from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Controlled as CodeMirror } from 'react-codemirror2';
import { updateProgress } from '../lib/store';
import { trueRandom } from '../lib/trueRandom';
import dataFetch from '../lib/data-fetch';
import CodeIcon from '@material-ui/icons/Code';
import Alert from '@material-ui/lab/Alert';

const promStyles = (theme) => ({
  prometheusWrapper : { padding : theme.spacing(5), },
  buttons : { display : 'flex',
    justifyContent : 'flex-end', },
  button : { marginTop : theme.spacing(3), },
  margin : { margin : theme.spacing(1), },
  chartTitle : { textAlign : 'center', },
  icon : { width : theme.spacing(2.5), },
  alignRight : { textAlign : 'right',
    marginBottom : theme.spacing(2), },
  formControl : { marginTop : theme.spacing(2),
    minWidth : window.innerWidth * 0.25, },
  panelChips : { display : 'flex',
    flexWrap : 'wrap', },
  panelChip : { margin : theme.spacing(0.25), },
});

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

class PrometheusSelectionComponent extends Component {
  constructor(props) {
    super(props);

    this.state = {
      grafanaBoard : dummyBoard,
      grafanaBoardError : false,
      grafanaBoardObject : {},
      templateVars : [],
      templateVarOptions : [], // will contain the array of options at the respective index, ex: [[v1], [v3, v4]]

      panels : [],
      selectedPanels : [],
      selectedTemplateVars : [], // will contain the selected val at the respective index: [v1, v3]
    };
  }

      handleChange = (name) => (event) => {
        if (name === 'grafanaBoard') {
          // if (this.cmEditor.state.lint.marked.length > 0) {
          //   this.setState({grafanaBoardError: true});
          //   return
          // }
          // if(this.boardTimeout) {
          //   clearnTimeout(this.boardTimeout);
          // }
          // const self = this;
          // this.boardTimeout = setTimeout(function(){
          //   self.boardChange(event.target.value);
          // }, 500);
        } else if (name.startsWith('template_var_')) {
          this.setSelectedTemplateVar(parseInt(name.replace('template_var_', '')), event.target.value);
        } else {
          this.setState({ [name] : event.target.value });
        }
      };

      getSelectedTemplateVar = (ind) => {
        const { selectedTemplateVars } = this.state;
        const retVal = typeof selectedTemplateVars[ind] !== 'undefined'
          ? selectedTemplateVars[ind]
          : undefined;
        return retVal;
      }

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
      }


      boardChange = (newVal) => {
        this.props.updateProgress({ showProgress : true });
        const self = this;
        dataFetch('/api/telemetry/metrics/board_import', {
          method : 'POST',
          credentials : 'include',
          headers : { 'Content-Type' : 'application/json', },
          body : newVal,
        },
        (result) => {
          this.props.updateProgress({ showProgress : false });
          var panels = result.panels.filter((panel) => panel.targets !== undefined && panel.datasource === 'Prometheus');
          if (typeof result !== 'undefined') {
            this.setState({
              grafanaBoardObject : result,
              panels : panels,
              selectedPanels : panels.map((panel) => panel.id), // selecting all panels by default
              templateVars : result.template_vars && result.template_vars.length > 0 ? result.template_vars : [],
              templateVarOptions : [],
              selectedTemplateVars : [],
            });
            if (result.template_vars && result.template_vars.length > 0) {
              this.queryTemplateVars(0, result.template_vars, [], []);
            }
          }
        }, self.props.handleError);

        // this.props.grafanaBoards.forEach((board) => {
        //   if (board.uri === newVal) {
        //       this.setState({
        //         grafanaBoard: newVal,
        //         panels: board.panels,
        //         selectedPanels: board.panels.map(panel=>panel.id), // selecting all panels by default
        //         templateVars: board.template_vars && board.template_vars.length > 0?board.template_vars:[],
        //         templateVarOptions: [],
        //         selectedTemplateVars: [],
        //       });
        //       if(board.template_vars && board.template_vars.length > 0){
        //         this.queryTemplateVars(0, board.template_vars, [], []);
        //       }
        //   }
        // });
      }


      queryTemplateVars = (ind, templateVars, templateVarOptions, selectedTemplateVars) => {
        if (templateVars.length > 0) {
          let queryURL = `/api/telemetry/metrics/query?query=${encodeURIComponent(templateVars[ind].query)}`;
          for (let i = ind; i > 0; i--) {
            queryURL += `&${templateVars[i - 1].name}=${selectedTemplateVars[i - 1]}`;
          }
          if (templateVars[ind].query.startsWith('label_values') && templateVars[ind].query.indexOf(',') > -1) {
            // series query needs a start and end time or else it will take way longer to return. . .
            // but at this point this component does not have the time range selection bcoz the time range selection comes after this component makes its selections
            // hence for now just limiting the time period to the last 24hrs
            const ed = new Date();
            const sd = new Date();
            sd.setDate(sd.getDate() - 1);
            queryURL += `&start=${Math.floor(sd.getTime() / 1000)}&end=${Math.floor(ed.getTime() / 1000)}`; // accounts for the last 24hrs
          }
          this.props.updateProgress({ showProgress : true });
          const self = this;
          dataFetch(queryURL, { credentials : 'include', }, (result) => {
            this.props.updateProgress({ showProgress : false });
            if (typeof result !== 'undefined') {
              let tmpVarOpts = [];
              // result.data check if it is an array or object
              if (Array.isArray(result.data)) {
                if (result.data.length > 0) {
                  if (templateVars[ind].query.startsWith('label_values') && templateVars[ind].query.indexOf(',') > -1
                    && typeof result.data[0] === 'object') {
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
          }, (error) => {
            templateVarOptions[ind] = [templateVars[ind].Value];
            this.setState({ templateVarOptions });
            self.props.handleError(error);
          });
        }
      }

      // static getDerivedStateFromProps(props, state){
      //   if (JSON.stringify(state.grafanaBoards.sort()) !== JSON.stringify(props.grafanaBoards.sort())) {
      //     return {
      //       grafanaBoards: props.grafanaBoards,
      //       grafanaBoard: '',
      //       selectedPanels: [],
      //     };
      //   }
      //   return null;
      // }

    addSelectedBoardPanelConfig = () => {
      const {
        grafanaBoardObject, templateVars, selectedTemplateVars, selectedPanels, panels,
      } = this.state;
      const boardConfig = {};
      boardConfig.board = grafanaBoardObject;
      boardConfig.panels = panels.filter(({ id }) => selectedPanels.indexOf(id) > -1);
      boardConfig.templateVars = templateVars.map(({ name }, index) => (typeof selectedTemplateVars[index] !== 'undefined'
        ? `${name}=${selectedTemplateVars[index]}`
        : ''));
      this.props.addSelectedBoardPanelConfig(boardConfig);
    }

    genRandomNumberForKey = () => Math.floor((trueRandom() * 1000) + 1)

    render = () => {
      const self = this;
      const {
        classes, prometheusURL,
        handlePrometheusChipDelete, handlePrometheusClick,
      } = this.props;
      const {
        panels, selectedPanels,
        grafanaBoard, templateVars, templateVarOptions,
      } = this.state;
      return (
        <NoSsr>
          <React.Fragment>
            <div className={classes.prometheusWrapper}>
              <div className={classes.alignRight}>
                <Chip
                  label={prometheusURL}
                  onDelete={handlePrometheusChipDelete}
                  onClick={handlePrometheusClick}
                  key="prometh-key"
                  icon={<img src="/static/img/prometheus_logo_orange_circle.svg" className={classes.icon} />}
                  variant="outlined"
                />
              </div>
              <Grid container spacing={1}>
                <Grid item xs={12}>
                  <div style={{ padding : '20px', display : 'flex' }}>
                    <CodeIcon style={{ marginRight : '6px' }} />
                  Paste your custom board JSON below.
                  </div>

                  <CodeMirror
                    editorDidMount={(editor) => {
                      this.cmEditor = editor;
                    }}
                    value={grafanaBoard}
                    options={{
                      theme : 'material',
                      lineNumbers : true,
                      lineWrapping : true,
                      gutters : ['CodeMirror-lint-markers'],
                      lint : true,
                      mode : 'application/json',
                    }}
                    onBeforeChange={(editor, data, value) => {
                      self.setState({
                        grafanaBoard : value,
                        grafanaBoardObject : {},
                        panels : [],
                        selectedPanels : [],
                        templateVars : [],
                        templateVarOptions : [],
                        selectedTemplateVars : [],
                      });
                      if (typeof self.boardTimeout !== 'undefined') {
                        clearTimeout(self.boardTimeout);
                      }
                      self.boardTimeout = setTimeout(() => {
                        console.log(`lint error count: ${self.cmEditor.state.lint.marked.length}`);
                        if (value !== '' && self.cmEditor.state.lint.marked.length === 0) {
                          self.setState({ grafanaBoardError : false });
                        } else {
                          self.setState({ grafanaBoardError : true });
                          return;
                        }
                        self.boardChange(value);
                      }, 1000);
                    }}
                    onChange={() => {
                    }}
                  />
                </Grid>
                {templateVars.length > 0
                  && templateVars.map(({ name }, ind) => {
                    // if (ind === 0 || this.getSelectedTemplateVar(ind-1) !== ''){
                    if (ind === 0 || typeof this.getSelectedTemplateVar(ind - 1) !== 'undefined') {
                      return (

                        <Grid item xs={12} sm={4} key={ind}>

                          <TextField
                            select
                            id={`template_var_${ind}`}
                            name={`template_var_${ind}`}
                            label={`Template variable: ${name}`}
                            fullWidth
                            value={this.getSelectedTemplateVar(ind)}
                            margin="normal"
                            variant="outlined"
                            onChange={this.handleChange(`template_var_${ind}`)}
                          >
                            <MenuItem key={`tmplVarOpt__-___${ind}_${self.genRandomNumberForKey()}`} value="" />
                            {templateVarOptions[ind] && templateVarOptions[ind].map((opt) => (
                              <MenuItem key={`tmplVarOpt__-__${name}_${opt}_${ind}_${self.genRandomNumberForKey()}`} value={opt}>{opt}</MenuItem>
                            ))}
                          </TextField>
                        </Grid>
                      );
                    }
                    return null;
                  })}

                {panels.length === 0 && (
                  <Grid item xs={12} style={{ marginTop : '10px' }}>
                    <Alert severity="error">Please load a valid Board JSON to be able to view the panels</Alert>
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
                      onChange={this.handleChange('selectedPanels')}
                      SelectProps={{ multiple : true,
                        renderValue : (selected) => (
                          <div className={classes.panelChips}>
                            {selected.map((value) => {
                              let selVal = '';
                              let panelId = '';
                              panels.forEach((panel) => {
                                if (panel.id === value) {
                                  selVal = panel.title;
                                  panelId = panel.id;
                                }
                              });
                              return (
                                <Chip key={`pl_--_${panelId}`} label={selVal} className={classes.panelChip} />
                              );
                            })}
                          </div>
                        ), }}
                    >
                      {panels.map((panel) => (
                        <MenuItem key={`panel_-__-${panel.id}`} value={panel.id}>{panel.title}</MenuItem>
                      ))}
                    </TextField>
                  </Grid>
                )}
              </Grid>
              {selectedPanels.length > 0 && (
                <div className={classes.buttons}>
                  <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    size="large"
                    onClick={this.addSelectedBoardPanelConfig}
                    className={classes.button}
                  >
                  Add
                  </Button>
                </div>
              )}
            </div>
          </React.Fragment>
        </NoSsr>
      );
    }
}

PrometheusSelectionComponent.propTypes = {
  classes : PropTypes.object.isRequired,
  prometheusURL : PropTypes.string.isRequired,
  handlePrometheusClick : PropTypes.func.isRequired,
  handlePrometheusChipDelete : PropTypes.func.isRequired,
  addSelectedBoardPanelConfig : PropTypes.func.isRequired,
  handleError : PropTypes.func.isRequired,
};

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });
const mapStateToProps = () => ({});

export default withStyles(promStyles)(connect(
  mapStateToProps,
  mapDispatchToProps,
)(PrometheusSelectionComponent));
