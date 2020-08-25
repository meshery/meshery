import React from 'react';
import NoSsr from '@material-ui/core/NoSsr';
import CytoscapeComponent from 'react-cytoscapejs';
import GraphStyle from './styles/styleContainer';
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