import React from 'react';
import NoSsr from '@material-ui/core/NoSsr';
import CytoscapeComponent from 'react-cytoscapejs';

class MesheryVisualizeComponent extends React.Component {
  render() {

    const elements = [
      { data: { id: 'one', label: 'Node 1' }, position: { x: 200, y: 200 } },
      { data: { id: 'two', label: 'Node 2' }, position: { x: 400, y: 400 } },
      { data: { source: 'one', target: 'two', label: 'Edge from Node1 to Node2' } }
    ];

    return (
      <NoSsr>
        <CytoscapeComponent elements={elements} style={ { width: '600px', height: '600px' } } />
      </NoSsr>
    )
  }
}

export default MesheryVisualizeComponent;