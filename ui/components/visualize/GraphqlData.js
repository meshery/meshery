// import React from 'react';
// import {graphql, QueryRenderer} from 'react-relay';
import {fetchQuery, graphql} from 'relay-runtime';
import environment from '../../lib/relayEnvironment';

export default class GraphQL {
  static getData() {
    const query = graphql`
        query GraphqlDataQuery {
          cluster {
            clusterNodes {
              id
              parentid
            }
          }
        }
      `;
    
    var cyData = {
      "nodes": [],
      "edges": [],
    }

    fetchQuery(environment, query).then(data => data.cluster.map(data => data.clusterNodes.map(data => cyData.nodes.push({...data, "label": data.id, })) ))
    
    return cyData
  }
}