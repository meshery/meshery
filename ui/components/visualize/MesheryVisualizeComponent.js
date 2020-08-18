import React from 'react';
import NoSsr from '@material-ui/core/NoSsr';
import CytoscapeComponent from 'react-cytoscapejs';
import DefaultGraph from './styles/styleContainer';
import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import { withStyles } from '@material-ui/core/styles';

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
    window.location.href=image.replace("image/png", "image/octet-stream")
  }

  render() {
    const {classes} = this.props

    const elements = {// get node j and the edges coming out from it

      nodes: [
        // Namespace > App > Proxy and Service
        {
          data: {
            id: 'ingress_gateway',
            label: 'Ingress Gateway'
          },
          classes: ['namespace']

          // selected: false, // whether the element is selected (default false)

          // selectable: true, // whether the selection state is mutable (default true)

          // locked: false, // when locked a node's position is immutable (default false)

          // grabbable: true, // whether the node can be grabbed and moved by the user

          // pannable: false, // whether dragging the node causes panning instead of grabbing
        },
        {
          data: {
            id: 'gateway_app',
            label: 'Istio Ingress Gateway',
            parent: 'ingress_gateway',

          },
          classes: ['app'] // an array (or a space separated string) of class names that the element has
        },
        {
          data: {
            parent: 'gateway_app'
          },
          classes: ['proxy'] // an array (or a space separated string) of class names that the element has
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

      ],
    };

    return (
      <NoSsr>
        <div style={{position: 'relative', width:'100%', height:'80%'}}>
          <div className={classes.div}>
            <CytoscapeComponent 
              elements={CytoscapeComponent.normalizeElements(elements)}
              style={ {width: '100%', height: '100%'} }
              layout={{ name: 'cose' }}
              stylesheet={DefaultGraph.getStylesheetContainer()}
              cy={cy => this.cy = cy}
            />
          </div>
          <ButtonGroup className={classes.zoomButton} color="primary" aria-label="outlined primary button group">
            <Button onClick={this.zoomIn.bind(this)}>+</Button>
            <Button onClick={this.zoomOut.bind(this)}>-</Button>
            <Button onClick={this.fit.bind(this)}>fit</Button>
          </ButtonGroup>
          <ButtonGroup className={classes.saveButton} color="primary" aria-label="outlined primary button group">
            <Button onClick={this.saveGraph.bind(this)}>Save</Button>
          </ButtonGroup>
        </div>
      </NoSsr>
    )
  }
}

export default withStyles(style)(MesheryVisualizeComponent);