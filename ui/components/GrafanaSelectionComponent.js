import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, TextField, Grid, Button, Chip, FormControl, InputLabel, Select, MenuItem, OutlinedInput } from '@material-ui/core';
import ReactDOM from 'react-dom';
import dataFetch from '../lib/data-fetch';

const grafanaStyles = theme => ({
    root: {
      padding: theme.spacing(10),
    },
    buttons: {
      display: 'flex',
    //   justifyContent: 'flex-end',
    },
    button: {
      marginTop: theme.spacing(3),
    //   marginLeft: theme.spacing(1),
    },
    margin: {
      margin: theme.spacing(1),
    },
    chartTitle: {
      textAlign: 'center',
    },
    icon: {
        width: theme.spacing(2.5),
    },
    alignRight: {
        textAlign: 'right',
    },
    formControl: {
        margin: theme.spacing(1),
        minWidth: 180,
    },
    panelChips: {
        display: 'flex',
        flexWrap: 'wrap',
    },
    panelChip: {
        margin: theme.spacing(0.25),
    }
  });

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
    PaperProps: {
        style: {
        maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
        width: 250,
        },
    },
};

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

            boardLabelWidth: 0,
            panelLabelWidth: 0,
            templateVarLabelWidth: 0,
          };
    }

    componentDidMount() {
        if (this.InputLabelRef) {
            console.log(`board label width: ${ReactDOM.findDOMNode(this.InputLabelRef).offsetWidth}`)
            this.setState({
                boardLabelWidth: ReactDOM.findDOMNode(this.InputLabelRef).offsetWidth,
            });
        }
        if (this.PanelInputLabelRef) {
            console.log(`panel label width: ${ReactDOM.findDOMNode(this.PanelInputLabelRef).offsetWidth}`)
            this.setState({
                panelLabelWidth: ReactDOM.findDOMNode(this.PanelInputLabelRef).offsetWidth,
            });
        }
        if (this.TemplateVarLabelRef) { // this is mostly not going to work bcoz this section wont exist when the component mounts
          console.log(`panel label width: ${ReactDOM.findDOMNode(this.TemplateVarLabelRef).offsetWidth}`)
          this.setState({
            templateVarLabelWidth: ReactDOM.findDOMNode(this.TemplateVarLabelRef).offsetWidth,
          });
        }
      }
    
      handleChange = name => event => {
        if (name === 'grafanaBoard'){
            this.boardChange(event.target.value);
        } else if (name.startsWith('template_var_')) {
          this.setSelectedTemplateVar(parseInt(name.replace('template_var_', '')), event.target.value);
        } else {
          this.setState({ [name]: event.target.value });
        }
      };

      getSelectedTemplateVar = (ind) => {
        const {selectedTemplateVars} = this.state;
        return selectedTemplateVars[ind]?selectedTemplateVars[ind]:'';
      }

      setSelectedTemplateVar = (ind, val) => {
        const {templateVars, templateVarOptions, selectedTemplateVars} = this.state;
        selectedTemplateVars[ind] = val;
        for (let i=ind+1; i<selectedTemplateVars.length;i++){
          selectedTemplateVars[i] = '';
        }
        this.setState({selectedTemplateVars});
        if (templateVars.length > ind+1){
          this.queryTemplateVars(ind+1, templateVars, templateVarOptions, selectedTemplateVars);
        }
      }


      boardChange = (newVal) => {
        this.props.grafanaBoards.forEach((board) => {
          // if (board.uri === this.state.grafanaBoard) {
          if (board.uri === newVal) {
              this.setState({
                grafanaBoard: newVal,
                panels: board.panels, 
                selectedPanels: board.panels.map(panel=>panel.id), // selecting all panels by default
                templateVars: board.template_vars && board.template_vars.length > 0?board.template_vars:[],
                templateVarOptions: [],
                selectedTemplateVars: [],
              });
              if(board.template_vars && board.template_vars.length > 0){
                this.queryTemplateVars(0, board.template_vars, [], []);
              }
          }
        });
      }
    
      
      queryTemplateVars = (ind, templateVars, templateVarOptions, selectedTemplateVars) => {
        // const {templateVars, templateVarOptions, selectedTemplateVars} = this.state;
        if (templateVars.length > 0) {
          let queryURL = `/api/grafana/query?query=${encodeURIComponent(templateVars[ind].query)}&dsid=${templateVars[ind].datasource.id}`;
          for(let i=ind;i>0;i--){ // assumption
            queryURL += `&${templateVars[i-1].name}=${selectedTemplateVars[i-1]}`;
          }
          
          let self = this;
          dataFetch(queryURL, { 
            credentials: 'same-origin',
            credentials: 'include',
          }, result => {
            if (typeof result !== 'undefined'){
              var tmpVarOpts = [];
              // result.data check if it is an array or object
              if (Array.isArray(result.data)){
                tmpVarOpts = result.data;
              } else {
                tmpVarOpts = result.data.result.map(({metric}) => {
                  const tmpArr = Object.keys(metric);
                  if (tmpArr.length > 0) {
                    return metric[tmpArr[0]];
                  }
                })
              }
              templateVarOptions[ind] = tmpVarOpts;
              this.setState({templateVarOptions});
            }
          }, self.props.handleError);
      }
    }

    static getDerivedStateFromProps(props, state){
      if (JSON.stringify(state.grafanaBoards.sort()) !== JSON.stringify(props.grafanaBoards.sort())) {
        return {
          grafanaBoards: props.grafanaBoards,
          grafanaBoard: '',
          selectedPanels: [],
        };
      }
      return null;
    }

    addSelectedBoardPanelConfig = () => {
      const {grafanaBoard, grafanaBoards, templateVars, selectedTemplateVars, selectedPanels, panels} = this.state;
      const boardConfig = {};
      grafanaBoards.forEach((board) => {
        if (grafanaBoard === board.uri){
          boardConfig.board = board;
          return;
        }
      })
      
      boardConfig.panels = panels.filter(({id}) => selectedPanels.indexOf(id) > -1);

      boardConfig.templateVars = templateVars.map(({name}, index) => `${name}=${selectedTemplateVars[index]}`);

      this.props.addSelectedBoardPanelConfig(boardConfig);
    }
    
    render = () => {
        const { classes, grafanaBoardSearch, grafanaURL, handleGrafanaBoardSearchChange, 
          handleGrafanaChipDelete } = this.props;
        const { panels, grafanaBoards, selectedPanels, panelLabelWidth, boardLabelWidth, 
          templateVarLabelWidth, grafanaBoard, templateVars, templateVarOptions, 
          selectedTemplateVars } = this.state;
        return (
          <NoSsr>
        <React.Fragment>
            <div className={classes.root}>
            <div className={classes.alignRight}>
                <Chip 
                    label={grafanaURL}
                    onDelete={handleGrafanaChipDelete} 
                    icon={<img src="/static/img/grafana_icon.svg" className={classes.icon} />} 
                    variant="outlined" />
            </div>
            <Grid container spacing={5}>
                <Grid item xs={12}>
                <TextField
                    id="grafanaBoardSearch"
                    name="grafanaBoardSearch"
                    label="Board Search"
                    fullWidth
                    value={grafanaBoardSearch}
                    margin="normal"
                    variant="outlined"
                    onChange={handleGrafanaBoardSearchChange('grafanaBoardSearch')} // this event will be sent to the parent
                />
                </Grid>
                <Grid item xs={12}>
                    <FormControl variant="outlined" className={classes.formControl}>
                        <InputLabel
                            ref={ref => {
                            this.InputLabelRef = ref;
                            }}
                            htmlFor="grafanaBoard"
                        >
                            Board
                        </InputLabel>
                        <Select
                            value={grafanaBoard}
                            onChange={this.handleChange('grafanaBoard')}
                            input={
                                <OutlinedInput
                                  name="grafanaBoard"
                                  labelWidth={boardLabelWidth}
                                //   fullWidth
                                  id="grafanaBoard"
                                />
                            }
                            // inputProps={{
                            //     name: 'grafanaBoard',
                            //     id: 'grafanaBoard',
                            // }}
                        >
                            {grafanaBoards && grafanaBoards.map((board) => (
                                <MenuItem key={board.uri} value={board.uri}>{board.title}</MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                </Grid>
                {templateVars.map(({name}, ind) => {
                    if (ind ===0 || this.getSelectedTemplateVar(ind-1) !== ''){
                    return (
                    <Grid item xs={12} sm={4}>
                        <FormControl variant="outlined" className={classes.formControl}>
                            <InputLabel
                                ref={ref => {
                                this.TemplateVarLabelRef = ref;
                                }}
                                htmlFor={'template_var_'+ind}
                            >
                                {name}
                            </InputLabel>
                            <Select
                                value={this.getSelectedTemplateVar(ind)}
                                onChange={this.handleChange('template_var_'+ind)}
                                input={
                                    <OutlinedInput
                                      name={'template_var_'+ind}
                                      labelWidth={boardLabelWidth}
                                      id={'template_var_'+ind}
                                    />
                                }
                            >   
                                <MenuItem key={'tmplVarOpt__'+ind} value={''}></MenuItem>
                                {templateVarOptions[ind] && templateVarOptions[ind].map((opt) => (
                                    <MenuItem key={'tmplVarOpt_'+name+'_'+opt+'_'+ind} value={opt}>{opt}</MenuItem>
                                ))}
                            </Select>
                        </FormControl>
                    </Grid>
                    );
                    } else {
                      return null;
                    }
                  })}

                <Grid item xs={12}>
                    <FormControl variant="outlined" className={classes.formControl}>
                        <InputLabel
                            ref={ref => {
                            this.PanelInputLabelRef = ref;
                            }}
                            htmlFor="panels"
                        >
                            Panels
                        </InputLabel>
                        <Select
                            multiple
                            value={selectedPanels}
                            onChange={this.handleChange('selectedPanels')}
                            input={
                                <OutlinedInput
                                  name="panels"
                                  labelWidth={panelLabelWidth}
                                  id="panels"
                                />
                            }
                            renderValue={selected => (
                                <div className={classes.panelChips}>
                                  {selected.map(value => {
                                      let selVal = '';
                                      let panelId = '';
                                      panels.forEach(panel=>{
                                        if (panel.id === value){
                                            selVal = panel.title;
                                            panelId = panel.id;
                                        }
                                      })
                                      return (
                                    <Chip key={panelId} label={selVal} className={classes.panelChip} />
                                      )
                                  })}
                                </div>
                            )}
                            // MenuProps={MenuProps}
                        >
                            {panels.map((panel) => (
                                <MenuItem key={'panel_'+panel.id} value={panel.id}>{panel.title}</MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                </Grid>
                <Grid item xs={12} sm={4}>
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
                </Grid>
            </Grid>
            </div>
        </React.Fragment>
        </NoSsr>
        );
    }
}

GrafanaSelectionComponent.propTypes = {
  classes: PropTypes.object.isRequired,
  grafanaURL: PropTypes.string.isRequired,
  grafanaBoards: PropTypes.array.isRequired,
  handleGrafanaBoardSearchChange: PropTypes.func.isRequired,
  handleGrafanaChipDelete: PropTypes.func.isRequired,
  addSelectedBoardPanelConfig: PropTypes.func.isRequired,
  handleError: PropTypes.func.isRequired,
};

export default withStyles(grafanaStyles)(GrafanaSelectionComponent);