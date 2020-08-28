import React from 'react';
import NoSsr from '@material-ui/core/NoSsr';
import CytoscapeComponent from 'react-cytoscapejs';
import GraphStyle from './styles/styleContainer';
import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import { withStyles } from '@material-ui/core/styles';
import cytoscape from 'cytoscape';
import cxtmenu from 'cytoscape-cxtmenu';

cytoscape.use(cxtmenu)

const style = (theme) => ({
  
  zoomButton: {
    position: 'absolute',
    top: 'auto',
    bottom: theme.spacing(2),
    left: theme.spacing(2),
    right: 'auto'
  },

  saveButton: {
    position: 'absolute',
    top:'auto',
    bottom: theme.spacing(2),
    left: 'auto',
    right:theme.spacing(2),
  },

  div: {
    width: '100%',
    height:'100%',
    borderRadius: '5px',
    background: '#fff',
  }
});

let cxtMenuSettings = {
  menuRadius: 100, // the radius of the circular menu in pixels
  selector: 'node', // elements matching this Cytoscape.js selector will trigger cxtmenus
  commands: [ // an array of commands to list in the menu or a function that returns the array
    {
      fillColor: 'rgba(200, 200, 200, 0.75)', // optional: custom background color for item
      content: 'cmd1', // html/text content to be displayed in the menu
      contentStyle: {}, // css key:value pairs to set the command's css in js if you want
      select: function(ele){ // a function to execute when the command is selected
        console.log( ele.id() ) // `ele` holds the reference to the active element
      },
      enabled: true // whether the command is selectable
    },
    {
      fillColor: 'rgba(200, 200, 200, 0.75)', // optional: custom background color for item
      content: 'cmd2', // html/text content to be displayed in the menu
      contentStyle: {}, // css key:value pairs to set the command's css in js if you want
      select: function(ele){ // a function to execute when the command is selected
        console.log( ele.id() ) // `ele` holds the reference to the active element
      },
      enabled: true // whether the command is selectable
    },
    {
      fillColor: 'rgba(200, 200, 200, 0.75)', // optional: custom background color for item
      content: 'cmd3', // html/text content to be displayed in the menu
      contentStyle: {}, // css key:value pairs to set the command's css in js if you want
      select: function(ele){ // a function to execute when the command is selected
        console.log( ele.id() ) // `ele` holds the reference to the active element
      },
      enabled: true // whether the command is selectable
    }
  ], // function( ele ){ return [ /*...*/ ] }, // a function that returns commands or a promise of commands
  fillColor: 'rgba(0, 0, 0, 0.75)', // the background colour of the menu
  activeFillColor: 'rgba(1, 105, 217, 0.75)', // the colour used to indicate the selected command
  activePadding: 20, // additional size in pixels for the active command
  indicatorSize: 24, // the size in pixels of the pointer to the active command
  separatorWidth: 3, // the empty spacing in pixels between successive commands
  spotlightPadding: 4, // extra spacing in pixels between the element and the spotlight
  minSpotlightRadius: 24, // the minimum radius in pixels of the spotlight
  maxSpotlightRadius: 38, // the maximum radius in pixels of the spotlight
  openMenuEvents: 'cxttapstart taphold', // space-separated cytoscape events that will open the menu; only `cxttapstart` and/or `taphold` work here
  itemColor: 'white', // the colour of text in the command's content
  itemTextShadowColor: 'transparent', // the text shadow colour of the command's content
  zIndex: 9999, // the z-index of the ui div
  atMouse: false // draw menu at mouse position
};

class MesheryVisualizeComponent extends React.Component {
  
  constructor(props) {
    super(props)
  }

  zoomIn() {
    this.cy.zoom(0.5 + this.cy.zoom())
    if(this.cy.$(':selected').size())
      this.cy.center(this.cy.$(':selected'))
    else
      this.cy.center()
  }

  zoomOut() {
    this.cy.zoom(-0.5 + this.cy.zoom())
    if(this.cy.$(':selected').size())
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
    lnk.download = 'MeshMap - '+date.toDateString()+'/' + date.toLocaleTimeString() +'.png'
    
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


  render() {
    const {classes} = this.props

    //Checkout the docs for JSON format https://js.cytoscape.org/#notation/elements-json
    const elements = {

      nodes: [
        {
          data: {
            id: 'ingress_gateway',
            label: 'Ingress Gateway'
          },
          classes: ['namespace']

        },
        {
          data: {
            id: 'gateway_app',
            label: 'Istio Ingress Gateway',
            parent: 'ingress_gateway',

          },
          classes: ['app']
        },
        {
          data: {
            parent: 'gateway_app'
          },
          classes: ['proxy']
        },
        {
          data: {
            id: 'product_namespace',
            label: 'Product',
          },
          classes: ['namespace']
        },
        {
          data: {
            id: 'app_2',
            label: 'Product',
            parent: 'product_namespace'
          },
          classes: ['app']
        },
        {
          data: {
            id: 'proxy_2',
            parent: 'app_2'
          },
          classes: ['proxy']
        },
        {
          data: {
            id: 'service_1',
            parent: 'app_2'
          },
          classes: ['service']
        },
        {
          data: {
            source: 'product_namespace',
            target: 'ingress_gateway',
            label: 'traffic'
          }
        }
      ],
    };

    return (
      <NoSsr>
        <div style={{position: 'relative', width:'100%', height:'80%'}}>
          <div className={classes.div}>
            <CytoscapeComponent 
              elements={CytoscapeComponent.normalizeElements(elements)}
              style={ {width: '100%', height: '100%'} }
              layout={GraphStyle.getLayout()}
              stylesheet={GraphStyle.getStylesheetContainer()}
              cy={cy => {
                this.cy = cy
                this.cy.cxtmenu( cxtMenuSettings );
              }}
            />
          </div>
          <ButtonGroup className={classes.zoomButton} color="primary" aria-label="outlined primary button group">
            <Button onClick={this.zoomIn.bind(this)}>+</Button>
            <Button onClick={this.zoomOut.bind(this)}>-</Button>
            <Button onClick={this.fit.bind(this)}>fit</Button>
          </ButtonGroup>
          <ButtonGroup className={classes.saveButton} color="primary" aria-label="outlined primary button group">
            <Button id="download" onClick={this.saveGraph.bind(this)} style={{textDecoration:'none'}}>Save</Button>
          </ButtonGroup>
        </div>
      </NoSsr>
    )
  }
}

export default withStyles(style)(MesheryVisualizeComponent);