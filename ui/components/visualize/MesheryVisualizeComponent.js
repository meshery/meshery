import React from 'react';
import NoSsr from '@material-ui/core/NoSsr';
import CytoscapeComponent from 'react-cytoscapejs';
import GraphStyle from './styles/styleContainer';
import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import { withStyles } from '@material-ui/core/styles';
import cytoscape from 'cytoscape';
import cxtmenu from 'cytoscape-cxtmenu';
import popper from 'cytoscape-popper';
import ReactDOM from 'react-dom';
import ToggleButton from '@material-ui/lab/ToggleButton';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import elementsJson from './Elements';
import dagre from 'cytoscape-dagre';
import clsx from 'clsx';
import {
  TopologyIcon,
} from '@patternfly/react-icons';
import { Paper } from '@material-ui/core';
import logsJson from './logs';
import PrimaryDrawer from './drawer/PrimaryDrawer';
import SecondaryDrawer from './drawer/SecondaryDrawer'
import PerformanceModal from './PerformanceModal';
import Terminal from './terminal';
import GraphQL from './GraphqlData';

cytoscape.use(dagre)
cytoscape.use(popper)
cytoscape.use(cxtmenu)

const DummyTooltip = (props) => {
  return <p>{props.data.data('app')}</p>;
};

const removeTooltip = () => {
  var el = document.getElementById('ccfc');
  if (el) {
    el.remove();
  }
}

const createContentFromComponent = (component) => {
  removeTooltip();
  var dummyDomEle = document.createElement('div');
  dummyDomEle.id = 'ccfc';
  ReactDOM.render(component, dummyDomEle);
  document.body.appendChild(dummyDomEle);
  return dummyDomEle;
};

const formalities = () => {
  var el = document.getElementById('edc');
  if (el) {
    el.remove();
  }
  el = document.getElementById('ndc');
  if (el) {
    el.remove();
  }
}

const style = (theme) => ({

  zoomButton: {
    position: 'absolute',
    top: 'auto',
    bottom: theme.spacing(2),
    left: theme.spacing(2),
    right: 'auto'
  },

  layoutButton: {
    position: 'absolute',
    top: 'auto',
    bottom: theme.spacing(2),
    left: '45%',
    right: '55%',
  },

  saveButton: {
    position: 'absolute',
    top: 'auto',
    bottom: theme.spacing(2),
    left: 'auto',
    right: theme.spacing(2),
  },

  div: {
    width: '100%',
    height: '100%',
    borderRadius: '5px',
    background: '#fff',
  },

  wrapper: {
    position: 'relative',
    width: '95%',
    height: '60vh',
  },

  wrapper2: {
    position: 'relative',
    width: '70%',
    height: '60vh',
  },

  logsContainer: {
    marginTop: '20px',
  },

  logs: {
    width: '100%',
    minHeight: '13vh',
    overflow: 'auto',
    resize: 'vertical',
    color: '#fff',
    backgroundColor: '#253137'
  }
});

class MesheryVisualizeComponent extends React.Component {

  constructor(props) {
    super(props);
    this.cyPopperRef = React.createRef();
    this.state = {
      layout: 'cose',
      open: false,
      data: null,
      tab: 0,
      logs: [],
      showModal: false,
      urlForModal: '',
      elements: {},
    }
    this.prev = null;
  }

  toggleChildMenu(data, val) {
    if(!val) this.prev = null;
    this.setState({open: val, data: data});
  }

  togglePrimaryDrawer(tab){
    this.setState({open: true, tab: tab});
  }

  zoomIn() {
    this.cy.zoom(0.5 + this.cy.zoom())
    if (this.cy.$(':selected').size())
      this.cy.center(this.cy.$(':selected'))
    else
      this.cy.center()
  }

  zoomOut() {
    this.cy.zoom(-0.5 + this.cy.zoom())
    if (this.cy.$(':selected').size())
      this.cy.center(this.cy.$(':selected'))
    else
      this.cy.center()
  }

  fit() {
    this.cy.fit()
    

  }

  saveGraph() {
    let image = this.cy.png()
    var lnk = document.createElement('a'), date = new Date(), e
    lnk.href = image
    lnk.download = 'MeshMap - ' + date.toDateString() + '/' + date.toLocaleTimeString() + '.png'

    if (document.createEvent) {
      e = document.createEvent("MouseEvents");
      e.initMouseEvent("click", true, true, window,
        0, 0, 0, 0, 0, false, false, false,
        false, 0, null);

      lnk.dispatchEvent(e);
    } else if (lnk.fireEvent) {
      lnk.fireEvent("onclick");
    }

  }

  handleLayoutChange = (e, value) => {
    if (value) {
      var layout = this.cy.layout({
        name: value
      });
      layout.run();
      this.setState({ layout: value });
    }
  }

  togglePeformanceModal(ele) {
    if(ele){
      this.setState({urlForModal: ele.data.URL?ele.data.URL:''})
    }

    this.setState(prevState => ({showModal: !prevState.showModal}))
  }

  componentDidMount() {
    this.setState({logs: logsJson.logs.join('\n')});
    this.setState({elements: GraphQL.getData()}, () => console.log(this.state.elements));
  }

  render() {
    const { classes } = this.props
    const { layout, open, data, tab , elements} = this.state;
    //Checkout the docs for JSON format https://js.cytoscape.org/#notation/elements-json

    //Checkout the docs at https://github.com/cytoscape/cytoscape.js-cxtmenu/blob/master/demo-adaptative.html
    let cxtMenuSettings = {
      menuRadius: 100,
      selector: 'node',
      commands: [
        {
          content: '<svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="tachometer-alt" class="svg-inline--fa fa-tachometer-alt fa-w-18 fa-fw " role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 576 512" style="transform-origin: 0.5625em 0.5em;"><g transform="translate(0, 0)  scale(1.7, 1.7)  rotate(0 0 0)"><path fill="currentColor" d="M288 32C128.94 32 0 160.94 0 320c0 52.8 14.25 102.26 39.06 144.8 5.61 9.62 16.3 15.2 27.44 15.2h443c11.14 0 21.83-5.58 27.44-15.2C561.75 422.26 576 372.8 576 320c0-159.06-128.94-288-288-288zm0 64c14.71 0 26.58 10.13 30.32 23.65-1.11 2.26-2.64 4.23-3.45 6.67l-9.22 27.67c-5.13 3.49-10.97 6.01-17.64 6.01-17.67 0-32-14.33-32-32S270.33 96 288 96zM96 384c-17.67 0-32-14.33-32-32s14.33-32 32-32 32 14.33 32 32-14.33 32-32 32zm48-160c-17.67 0-32-14.33-32-32s14.33-32 32-32 32 14.33 32 32-14.33 32-32 32zm246.77-72.41l-61.33 184C343.13 347.33 352 364.54 352 384c0 11.72-3.38 22.55-8.88 32H232.88c-5.5-9.45-8.88-20.28-8.88-32 0-33.94 26.5-61.43 59.9-63.59l61.34-184.01c4.17-12.56 17.73-19.45 30.36-15.17 12.57 4.19 19.35 17.79 15.17 30.36zm14.66 57.2l15.52-46.55c3.47-1.29 7.13-2.23 11.05-2.23 17.67 0 32 14.33 32 32s-14.33 32-32 32c-11.38-.01-20.89-6.28-26.57-15.22zM480 384c-17.67 0-32-14.33-32-32s14.33-32 32-32 32 14.33 32 32-14.33 32-32 32z"></path></g></svg>',
          select: this.togglePeformanceModal.bind(this),
          enabled: true // whether the command is selectable
        },
        {
          content: 'cmd2', // html/text content to be displayed in the menu
          contentStyle: {}, // css key:value pairs to set the command's css in js if you want
          select: function (ele) { // a function to execute when the command is selected
            console.log(ele.id()) // `ele` holds the reference to the active element
          },
          enabled: false // whether the command is selectable
        },
        {
          content: 'cmd3', // html/text content to be displayed in the menu
          contentStyle: {}, // css key:value pairs to set the command's css in js if you want
          select: function (ele) { // a function to execute when the command is selected
            console.log(ele.id()) // `ele` holds the reference to the active element
          },
          enabled: false // whether the command is selectable
        }
      ],
      fillColor: '#607d8b77',
      activeFillColor: '#263237',
      activePadding: 10, // additional size in pixels for the active command
      indicatorSize: 24, // the size in pixels of the pointer to the active command
      separatorWidth: 2, // the empty spacing in pixels between successive commands
      spotlightPadding: 1,  // extra spacing in pixels between the element and the spotlight
      adaptativeNodeSpotlightRadius: true,
      minSpotlightRadius: 24, // the minimum radius in pixels of the spotlight
      maxSpotlightRadius: 38, // the maximum radius in pixels of the spotlight
      openMenuEvents: 'cxttapstart taphold', // space-separated cytoscape events that will open the menu; only `cxttapstart` and/or `taphold` work here
      itemColor: 'white', // the colour of text in the command's content
      zIndex: 9999, // the z-index of the ui div
      atMouse: false // draw menu at mouse position
    };

    return (
      <NoSsr>
        <PerformanceModal
          open={this.state.showModal}
          handleClose={() => this.setState({showModal: false})}
          urlForModal={this.state.urlForModal}
        />
        <PrimaryDrawer 
          toggle={(tab) => {
            this.togglePrimaryDrawer(tab);
          }}
        />
        <div className={clsx({ [classes.wrapper]: !open }, {
          [classes.wrapper2]: open,
        })} >
          <div className={classes.div}>
            <CytoscapeComponent
              elements={CytoscapeComponent.normalizeElements(elements)}
              style={{ width: '100%', height: '100%' }}
              layout={GraphStyle.getLayout()}
              stylesheet={GraphStyle.getStylesheetContainer()}
              cy={cy => {
                this.cy = cy;
                this.cy.cxtmenu(cxtMenuSettings);
                this.cy.elements().on('mouseover', (event) => {
                  this.cyPopperRef.current = event.target.popper({
                    content: createContentFromComponent(<DummyTooltip data={event.target} />),
                    popper: {
                      placement: 'right',
                      removeOnDestroy: true,
                    },
                  });
                });
                this.cy.elements().on('mouseout', () => {
                  if (this.cyPopperRef) {
                    removeTooltip();
                  }
                });
                this.cy.elements().on('click', (event) => { 
                  var sel = event.target;
                  if(!this.prev){
                    this.prev = event.target;
                  } else {
                    if(event.target === this.prev) return;
                    else this.prev = event.target;
                  }
                  formalities();
                  this.toggleChildMenu(sel, true);
                  this.cy.startBatch();
                  {
                    this.cy.elements().removeClass('semitransp','highlight');
                    this.cy.elements().difference(sel.outgoers().union(sel.incomers())).not(sel).addClass('semitransp');
                    sel.addClass('highlight').outgoers().union(sel.incomers()).addClass('highlight');
                    this.cy.resize();
                    this.cy.animate({
                      fit: {
                        eles: sel.outgoers().union(sel.incomers()),
                        padding: 200
                      }
                    },
                    {
                      duration: 500
                    }
                    )
                  }
                  this.cy.endBatch();
                });
              }}
            />
          </div>
          {
            open && <SecondaryDrawer data={data} open={open} tab={tab}
              toggle={(data, val) => {
                this.toggleChildMenu(data, val);
              }}
              togglePeformanceModal = {() => {
                this.togglePeformanceModal() 
              }}
            />
          }
          <ButtonGroup className={classes.zoomButton} color="primary" aria-label="outlined primary button group">
            <Button onClick={this.zoomIn.bind(this)}>+</Button>
            <Button onClick={this.zoomOut.bind(this)}>-</Button>
            <Button onClick={this.fit.bind(this)}>fit</Button>
          </ButtonGroup>
          <ToggleButtonGroup
            value={layout}
            exclusive
            onChange={this.handleLayoutChange}
            className={classes.layoutButton}
            size="small"
          >
            <ToggleButton value="cose" >
              <TopologyIcon /> 1
            </ToggleButton>
            <ToggleButton value="breadthfirst" >
              <TopologyIcon /> 2
            </ToggleButton>
            <ToggleButton value="circle" >
              <TopologyIcon /> 3
            </ToggleButton>
            <ToggleButton value="dagre" >
              <TopologyIcon /> 4
            </ToggleButton>
            <ToggleButton value="concentric" >
              <TopologyIcon /> 5
            </ToggleButton>
          </ToggleButtonGroup>
          <ButtonGroup className={classes.saveButton} color="primary" aria-label="outlined primary button group">
            <Button id="download" onClick={this.saveGraph.bind(this)} style={{ textDecoration: 'none' }}>Save</Button>
          </ButtonGroup>
          <Paper className={classes.logsContainer}>
            <Terminal className={classes.logs} />   
          </Paper>
        </div>
      </NoSsr>
    )
  }
}

export default withStyles(style)(MesheryVisualizeComponent);