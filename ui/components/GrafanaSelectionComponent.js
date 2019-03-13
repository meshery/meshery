import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, TextField, Grid, Button, Chip, FormControl, InputLabel, Select, MenuItem, OutlinedInput } from '@material-ui/core';
import ReactDOM from 'react-dom';

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

            panels: [],
            selectedPanels: [],

            boardLabelWidth: 0,
            panelLabelWidth: 0,
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
      }
    
      handleChange = name => event => {
        if (name === "grafanaBoard"){
            this.props.grafanaBoards.forEach((board) => {
                if (board.uri === this.state.grafanaBoard) {
                    this.setState({panels: board.panels, selectedPanels: board.panels.map(panel=>panel.id)}); // selecting all panels by default
                }
            });
        }

        this.setState({ [name]: event.target.value });
      };
    
      
    
    handleError = error => {
      // this.setState({timerDialogOpen: false });
      // this.setState({showSnackbar: true, snackbarVariant: 'error', snackbarMessage: `Grafana was not configured successfully with msg: ${error}`});
    }

    static getDerivedStateFromProps(props, state){
      if (state.grafanaBoards.sort().join() !== props.grafanaBoards.sort().join()) {
        return {
          grafanaBoards: props.grafanaBoards,
          grafanaBoard: '',
          selectedPanels: [],
        };
      }
      return null;
    }
    
    render = () => {
        const { classes, grafanaBoardSearch, grafanaURL, handleGrafanaBoardSearchChange, handleGrafanaChipDelete, addSelectedPanels } = this.props;
        const { panels, grafanaBoards, selectedPanels, panelLabelWidth, boardLabelWidth, grafanaBoard } = this.state;
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
                <Grid item xs={12} sm={4}>
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
                            {grafanaBoards.map((board) => (
                                <MenuItem key={board.uri} value={board.uri}>{board.title}</MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                </Grid>
                
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
                                <MenuItem key={panel.id} value={panel.id}>{panel.title}</MenuItem>
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
                        onClick={addSelectedPanels}
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
};

export default withStyles(grafanaStyles)(GrafanaSelectionComponent);