export default class GraphStyle {
  static getLayout() {
    //Right now using cose. Adding more options in future
    return {
      name: 'cose',
      animate: false,
      fit: false,
      nodeDimensionsIncludeLabels: true
    };
  }

  //Styling nodes and edges
  //Work needed to be done. Feel free to play around with it
  static getStylesheetContainer() {
    return [
      {
        selector: "node",
        style: {
          "background-color": "rgb(255,255,255)",
          "background-width": "80%",
          "background-height": "80%",
          "background-position-x": "1px",
          "background-position-y": "1px",
          "border-width": "1px",
          ghost: "yes",
          "label": "data(label)",
          "ghost-offset-x": "1px",
          "ghost-offset-y": "1px",
          "ghost-opacity": "0.4",
          height: "25px",
          width: "25px",
          "z-index": "10"
        }
      },
      {
        selector: "node[?isGroup]",
        style: {
          "background-color": "rgb(255,255,255)"
        }
      },
      {
        selector: "node:selected",
        style: {
          "border-width": "3px"
        }
      },
      {
        selector: "node.mousehighlight",
        style: {
          "font-size": "11px"
        }
      },
      {
        selector: "node.mousehighlight[^isGroup]",
        style: {
        }
      },
      {
        selector: "node.mousedim",
        style: {
          opacity: "0.6"
        }
      },
      {
        selector: "edge",
        style: {
          "curve-style": "bezier",
          "font-family": "Verdana,Arial,Helvetica,sans-serif,pficon",
          "font-size": "6px",
          label: "edge",
          "line-style": "solid",
          "target-arrow-shape": "vee",
          "text-events": "yes",
          "text-outline-color": "rgb(255,255,255)",
          "text-outline-width": "1px",
          "text-wrap": "wrap",
          width: "2px"
        }
      },
      {
        selector: "edge:selected",
        style: {
          width: "4px",
          label: "edge"
        }
      },
      {
        selector: "edge[protocol = \"tcp\"]",
        style: {
          "target-arrow-shape": "triangle-cross",
          "line-style": "solid"
        }
      },
      {
        selector: "edge.mousehighlight",
        style: {
          "font-size": "10px"
        }
      },
      {
        selector: "edge.mousehover",
        style: {
          label: "edge"
        }
      },
      {
        selector: "edge.mousedim",
        style: {
          opacity: "0.3"
        }
      },
      {
        selector: ".find[^isGroup]",
        style: {
          "overlay-color": "rgb(240,171,0)",
          "overlay-padding": "8px",
          "overlay-opacity": "0.5"
        }
      },
      {
        selector: 'node.highlight',
        style: {
          'border-color': '#123',
          'border-width': '2px'
        }
      },
      {
        selector: 'node.semitransp',
        style:{ 'opacity': '0.5' }
      },
      {
        selector: 'edge.highlight',
        style: { 'mid-target-arrow-color': '#FFF' }
      },
      {
        selector: 'edge.semitransp',
        style:{ 'opacity': '0.2' }
      }
    ]
  }
}