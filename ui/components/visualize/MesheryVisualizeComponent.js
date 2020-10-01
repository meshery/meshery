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
import Drawer from './Drawer'

cytoscape.use(popper)
cytoscape.use(cxtmenu)

const DummyTooltip = (props) => {
  return <p>{props.data.data('app')}</p>;
};

const removeTooltip = () => {
  var el = document.getElementById('ccfc');
  if ( el ) {
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
    top:'auto',
    bottom: theme.spacing(2),
    left: '45%',
    right: '55%',
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
    super(props);
    this.cyPopperRef = React.createRef();
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

  changeLayout(type){
    var layout = this.cy.layout({
      name: type
    })
    layout.run();
  }

  render() {
    const {classes} = this.props

    //Checkout the docs for JSON format https://js.cytoscape.org/#notation/elements-json
    const elements = {
      
      nodes: [
        {
          data: {
            app: "details",
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            isGroup: "app",
            tcpIn: null,
            tcpOut: null,
            id: "bcbbc26d201909101e45d8edb0b617ae",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 301,
            y: 69
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "productpage",
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            isGroup: "app",
            tcpIn: null,
            tcpOut: null,
            id: "1cd79b3dd828bdcf35a57fdfc4e4f505",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 89.5,
            y: 307
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "ratings",
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            isGroup: "app",
            tcpIn: null,
            tcpOut: null,
            id: "4aa670e0682d4be3aaf24d8f589feb4b",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 493.5,
            y: 338.75
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "reviews",
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            isGroup: "app",
            tcpIn: null,
            tcpOut: null,
            id: "6519157be154675342fb76c41edc731c",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 301,
            y: 307
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            destServices: [
              {
                namespace: "unknown",
                name: "kubernetes.default.svc.cluster.local"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            service: "kubernetes.default.svc.cluster.local",
            tcpIn: 2130.18,
            tcpOut: null,
            id: "d40327e4c9dd917578c6c51cb641dbf4",
            nodeType: "service",
            namespace: "unknown"
          },
          position: {
            x: 301,
            y: 521
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "prometheus",
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            isRoot: true,
            tcpIn: null,
            tcpOut: 2130.18,
            version: "latest",
            workload: "prometheus",
            id: "c9e43c28f29e6debbbd908c816bc0a4e",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 89.5,
            y: 521
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "details",
            destServices: [
              {
                namespace: "default",
                name: "details"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            service: "details",
            tcpIn: null,
            tcpOut: null,
            id: "16d04f68bd507ca9b0707c2a576d1fc2",
            parent: "bcbbc26d201909101e45d8edb0b617ae",
            nodeType: "service",
            namespace: "default"
          },
          position: {
            x: 261,
            y: 58
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "details",
            destServices: [
              {
                namespace: "default",
                name: "details"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            tcpIn: null,
            tcpOut: null,
            version: "v1",
            workload: "details-v1",
            id: "721ef0b8cfba57d153213a5d659ae9da",
            parent: "bcbbc26d201909101e45d8edb0b617ae",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 342,
            y: 58
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "productpage",
            destServices: [
              {
                namespace: "default",
                name: "productpage"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            hasVS: true,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            service: "productpage",
            tcpIn: null,
            tcpOut: null,
            id: "c2efd356d9a25fd009efe2a323e12361",
            parent: "1cd79b3dd828bdcf35a57fdfc4e4f505",
            nodeType: "service",
            namespace: "default"
          },
          position: {
            x: 49.5,
            y: 295
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "productpage",
            destServices: [
              {
                namespace: "default",
                name: "productpage"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            tcpIn: null,
            tcpOut: null,
            version: "v1",
            workload: "productpage-v1",
            id: "06e488a37fc9aa5b0e0805db4f16ae69",
            parent: "1cd79b3dd828bdcf35a57fdfc4e4f505",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 150,
            y: 295
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "ratings",
            destServices: [
              {
                namespace: "default",
                name: "ratings"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            service: "ratings",
            tcpIn: null,
            tcpOut: null,
            id: "906196769ac4714aae43f4f789a36d9c",
            parent: "4aa670e0682d4be3aaf24d8f589feb4b",
            nodeType: "service",
            namespace: "default"
          },
          position: {
            x: 453.5,
            y: 327.75
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "ratings",
            destServices: [
              {
                namespace: "default",
                name: "ratings"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            tcpIn: null,
            tcpOut: null,
            version: "v1",
            workload: "ratings-v1",
            id: "4b64bda48e5a3c7e50ab1c63836c9469",
            parent: "4aa670e0682d4be3aaf24d8f589feb4b",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 535,
            y: 327.75
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "reviews",
            destServices: [
              {
                namespace: "default",
                name: "reviews"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            service: "reviews",
            tcpIn: null,
            tcpOut: null,
            id: "adbe9380f23fcbefe5fedd6beb0597ef",
            parent: "6519157be154675342fb76c41edc731c",
            nodeType: "service",
            namespace: "default"
          },
          position: {
            x: 261,
            y: 296
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "reviews",
            destServices: [
              {
                namespace: "default",
                name: "reviews"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            tcpIn: null,
            tcpOut: null,
            version: "v1",
            workload: "reviews-v1",
            id: "ce6eb1c1255b2c90e76a8f1a803cdb24",
            parent: "6519157be154675342fb76c41edc731c",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 344,
            y: 246
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        },
        {
          data: {
            app: "reviews",
            destServices: [
              {
                namespace: "default",
                name: "reviews"
              }
            ],
            grpcIn: null,
            grpcInErr: null,
            grpcOut: null,
            httpIn: null,
            httpIn3xx: null,
            httpIn4xx: null,
            httpIn5xx: null,
            httpOut: null,
            tcpIn: null,
            tcpOut: null,
            version: "v2",
            workload: "reviews-v2",
            id: "31150e7e5adf85b63f22fbd8255803d7",
            parent: "6519157be154675342fb76c41edc731c",
            nodeType: "app",
            namespace: "default"
          },
          position: {
            x: 344,
            y: 346
          },
          group: "nodes",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: false,
          classes: ""
        }
      ],
      edges: [
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            hasTraffic: true,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "tcp",
            responses: {
              "-": {
                flags: {
                  "-": "100.0"
                },
                hosts: {
                  "kubernetes.default.svc.cluster.local": "100.0"
                }
              }
            },
            responseTime: null,
            tcp: 2130.18,
            id: "ec73f6bac5146601e57af53f400f4b86",
            source: "c9e43c28f29e6debbbd908c816bc0a4e",
            target: "d40327e4c9dd917578c6c51cb641dbf4",
            isMtls: null
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        },
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "http",
            responseTime: null,
            tcp: null,
            id: "603b3196d57ed4275cde825c9b550cf8",
            source: "adbe9380f23fcbefe5fedd6beb0597ef",
            target: "31150e7e5adf85b63f22fbd8255803d7",
            traffic: {
              protocol: "http"
            }
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        },
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "http",
            responseTime: null,
            tcp: null,
            id: "5bfd2136556bf57f55046dee357f11ba",
            source: "adbe9380f23fcbefe5fedd6beb0597ef",
            target: "ce6eb1c1255b2c90e76a8f1a803cdb24",
            traffic: {
              protocol: "http"
            }
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        },
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "http",
            responseTime: null,
            tcp: null,
            id: "085e267222568260e02b64d1d47b48fb",
            source: "31150e7e5adf85b63f22fbd8255803d7",
            target: "906196769ac4714aae43f4f789a36d9c",
            traffic: {
              protocol: "http"
            }
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        },
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "http",
            responseTime: null,
            tcp: null,
            id: "a1dfe2e1b0907e1cb16670f1abe2e41e",
            source: "906196769ac4714aae43f4f789a36d9c",
            target: "4b64bda48e5a3c7e50ab1c63836c9469",
            traffic: {
              protocol: "http"
            }
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        },
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "http",
            responseTime: null,
            tcp: null,
            id: "49e7da37f0733e10454cdc64887377dc",
            source: "06e488a37fc9aa5b0e0805db4f16ae69",
            target: "adbe9380f23fcbefe5fedd6beb0597ef",
            traffic: {
              protocol: "http"
            }
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        },
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "http",
            responseTime: null,
            tcp: null,
            id: "3d6c17c6a0ba2aecb9fc3b4466502fb9",
            source: "4b64bda48e5a3c7e50ab1c63836c9469",
            target: "c2efd356d9a25fd009efe2a323e12361",
            traffic: {
              protocol: "http"
            }
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        },
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "http",
            responseTime: null,
            tcp: null,
            id: "86725a525e9ddca333f585f328c423e9",
            source: "c2efd356d9a25fd009efe2a323e12361",
            target: "06e488a37fc9aa5b0e0805db4f16ae69",
            traffic: {
              protocol: "http"
            }
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        },
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "http",
            responseTime: null,
            tcp: null,
            id: "531fdd02d44612b9a7a9f393dcbf6940",
            source: "06e488a37fc9aa5b0e0805db4f16ae69",
            target: "16d04f68bd507ca9b0707c2a576d1fc2",
            traffic: {
              protocol: "http"
            }
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        },
        {
          data: {
            grpc: null,
            grpcErr: null,
            grpcPercentErr: null,
            grpcPercentReq: null,
            http: null,
            http3xx: null,
            http4xx: null,
            http5xx: null,
            httpPercentErr: null,
            httpPercentReq: null,
            isMTLS: -1,
            protocol: "http",
            responseTime: null,
            tcp: null,
            id: "f774bb2441e6cc6e34630c04c150ca35",
            source: "16d04f68bd507ca9b0707c2a576d1fc2",
            target: "721ef0b8cfba57d153213a5d659ae9da",
            traffic: {
              protocol: "http"
            }
          },
          position: {
            x: 0,
            y: 0
          },
          group: "edges",
          removed: false,
          selected: false,
          selectable: false,
          locked: false,
          grabbable: true,
          pannable: true,
          classes: ""
        }
      ]
    };

    return (
      <NoSsr>
        <div style={{position: 'relative', width:'100%', height:'90%'}}>
          <div className={classes.div}>
            <CytoscapeComponent 
              elements={CytoscapeComponent.normalizeElements(elements)}
              style={ {width: '100%', height: '100%'} }
              layout={GraphStyle.getLayout()}
              stylesheet={GraphStyle.getStylesheetContainer()}
              cy={cy => {
                this.cy = cy;
                this.cy.cxtmenu( cxtMenuSettings );
                this.cy.elements().on('mouseover', (event) => {
                  this.cyPopperRef.current = event.target.popper({
                    content: createContentFromComponent(<DummyTooltip data={event.target}/>),
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
                this.cy.nodes().on('click', (event) => {
                  var el = document.getElementById('ndc');
                  if (el) {
                    el.remove();
                  }
                  el = document.getElementById('edc');
                  if (el) {
                    el.remove();
                  }
                  var dummyDomEle = document.createElement('div');
                  dummyDomEle.id = 'ndc';
                  ReactDOM.render(<Drawer data={event.target}/>, dummyDomEle);
                  document.body.appendChild(dummyDomEle);
                })
                this.cy.edges().on('click', (event) => {
                  var el = document.getElementById('edc');
                  if (el) {
                    el.remove();
                  }
                  el = document.getElementById('ndc');
                  if (el) {
                    el.remove();
                  }
                  var dummyDomEle = document.createElement('div');
                  dummyDomEle.id = 'edc';
                  ReactDOM.render(<Drawer data={event.target}/>, dummyDomEle);
                  document.body.appendChild(dummyDomEle);
                })
              }}
            />
          </div>
          <ButtonGroup className={classes.zoomButton} color="primary" aria-label="outlined primary button group">
            <Button onClick={this.zoomIn.bind(this)}>+</Button>
            <Button onClick={this.zoomOut.bind(this)}>-</Button>
            <Button onClick={this.fit.bind(this)}>fit</Button>
          </ButtonGroup>
          <ButtonGroup className={classes.layoutButton} color="primary" aria-label="outlined primary button group">
            <Button onClick={this.changeLayout.bind(this, 'cose')}>Cose</Button>
            <Button onClick={this.changeLayout.bind(this, 'breadthfirst')}>BFS</Button>
            <Button onClick={this.changeLayout.bind(this, 'circle')}>Circle</Button>
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