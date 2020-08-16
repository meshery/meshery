 export default class DefaultGraph {
  static getLayout() {
    return {
      name: 'cola',
      animate: false,
      fit: false,
      flow: { axis: 'x' },
      nodeDimensionsIncludeLabels: true,
      randomize: false
    }
  }

  //Styling the <div> wrrapping cydoscape
  static getStyleContainer() {
    return {
      width: '100%',
      height:'80%',
      borderRadius: '5px',
      background: '#fff'
    }
  }

  //Stling nodes and edges
  static getStylesheetContainer() {
    return [
      {
        selector: 'node',
        style: {
          'background-color': '#4caf50',
          'label': 'data(label)'
        }
      },
    ]
  }
}