export default class DefaultGraph {
  static getLayout() {
    return {
      name: 'cose-bilkent',
      animate: false,
      fit: false,
      nodeDimensionsIncludeLabels: true
    };
  }

  //Stling nodes and edges
  static getStylesheetContainer() {
    return [
      {
        selector: "node.proxy",
        style: {
          width: 50,
          height: 20,
          shape: 'rectangle',
        }
      },
      {
        selector: "node.service",
        style: {
          width: 50,
          height: 20,
          shape: 'rectangle',
        }
      },
      {
        selector: "node.app",
        style: {
          shape: 'round-rectangle',
          label: 'data(label)',
        }
      },
      {
        selector: "node.namespace",
        style: {
          shape: 'round-rectangle',
          label: 'data(label)',
        }
      },
      {
        selector: "edge",
        style: {
          width: 1,
          opacity: 0.6,
          "line-color": '#000',
          events: "no"
        }
      },
      {
        selector: "node:selected",
        style: {
          'border-width': 2,
          'border-color': '#F7F0F0'
        }
      },
      {
        selector: "node[label]",
        style: {
          'text-valign': 'bottom',
          'font-family': ['Helvetica', 'Arial', 'sans-serif'],
          'text-wrap': 'wrap',
          'text-max-width': '100',
        }
      },

    ]
  }
}