import React from 'react';
import NoSsr from '@material-ui/core/NoSsr';
import CytoscapeComponent from 'react-cytoscapejs';
import DefaultGraph from './styles/styleContainer'


class MesheryVisualizeComponent extends React.Component {
  
  constructor(props) {
    super(props)
    //cy is a ref to the core object of cytoscope
    this.myCy = React.createRef()
  }


  componentDidMount() {
    console.log(this.myCy)
    // this.myCy.style()
    // this.myCy.current.zoom(100)
  }

  render() {

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
        <div style={DefaultGraph.getStyleContainer()}>
          <CytoscapeComponent 
            elements={CytoscapeComponent.normalizeElements(elements)}
            style={ {width: '100%', height: '100%'} }
            layout={{ name: 'cose' }}
            cy={(cy) => {
              this.myCy = cy 
            }} />
        </div>
      </NoSsr>
    )
  }
}

export default MesheryVisualizeComponent;