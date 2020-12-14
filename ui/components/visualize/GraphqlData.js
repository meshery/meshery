// import React from 'react';
// import {graphql, QueryRenderer} from 'react-relay';
import {fetchQuery, graphql} from 'relay-runtime';
import environment from '../../lib/relayEnvironment';

export default class GraphQL {
  static getData(layout, filters) {

    var cyData = {
      "nodes": [],
      "edges": [],
    }
    
    if(layout=="layout1"){
      let query = graphql`
        query GraphqlDataQuery($nodeID: ID ) {
          cluster {
            clusterNodes(nodeid: $nodeID) {
              id
              parentid
              pods {
                id
                parentid
                name
              }
            }
          }
        }
      `;
      let variables = filters
      console.log("filters", filters)
      fetchQuery(environment, query, variables).then(
        data => data.cluster.map(
          (cluster) => {
            console.log("data", data)
            cluster.clusterNodes.map(
              node => {
                cyData.nodes.push(
                  {
                    "data": {
                      "parent":node.parentid,
                      "label": node.id,
                      "id": node.id
                    }
                  }
                )
                node.pods.map(
                  pod => {
                    console.log("Pods", pod)
                    cyData.nodes.push(
                      {
                        "data": {
                          "parent": pod.parentid,
                          "id": pod.id,
                          "label": pod.name,
                        }
                      }
                    )
                  }
                )

              }
            )
          }
        )
      )
    } else if(layout=="layout2"){
      // let query = graphql`
      //   query GraphqlDataQuery {
      //     cluster {
      //       clusterNodes {
      //         id
      //         parentid
      //         pods {
      //           id
      //           parentid
      //           name
      //         }
      //       }
      //     }
      //   }
      // `;
      // fetchQuery(environment, query).then(
      //   data => data.cluster.map(
      //     (cluster) => {
      //       console.log("data", data)
      //       cluster.clusterNodes.map(
      //         node => {
      //           cyData.nodes.push(
      //             {
      //               "data": {
      //                 "parent":node.parentid,
      //                 "label": node.id,
      //                 "id": node.id
      //               }
      //             }
      //           )
      //           console.log("node", node)
      //           node.pods.map(
      //             pods => {
      //               cyData.nodes.push(
      //                 {
      //                   "data": {
      //                     "parent": pods.parentid,
      //                     "id": pods.id,
      //                     "label": pods.name,
      //                   }
      //                 }
      //               )
      //             }
      //           )

      //         }
      //       )
      //     }
      //   )
      // )
    }
    
    

    
    return cyData
  }
}