import fetchView1 from './View1'
import fetchView2 from './View2'

export default class GraphQL {
  static getData(layout, filters) {

    var cyData = {
      "nodes": [],
      "edges": [],
    }
    let parentid = ""
    console.log("variables", filters, layout)

    if(layout=="view1"){
      fetchView1(filters).then(
        data => data.cluster.map(
          (cluster) => {
            cluster.id && cyData.nodes.push(
              {
                "data": {
                  "parent":parentid,
                  "label": cluster.id,
                  "id": cluster.id
                }
              }
            )
            parentid = cluster.id?cluster.id:parentid
            cluster.clusterNodes.map(
              node => {
                node.id && cyData.nodes.push(
                  {
                    "data": {
                      "parent":parentid,
                      "label": node.id,
                      "id": node.id
                    }
                  }
                )
              }
            )

            cluster.clusterNodes.map( node => {
              parentid = node.id?node.id:parentid
              node.pods.map(
                pod => {
                  pod.id && cyData.nodes.push(
                    {
                      "data": {
                        "parent": parentid,
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
    } else if(layout=="view2"){

      fetchView2(filters).then(
        data => data.cluster.map(
          (cluster) => {
            cluster.id && cyData.nodes.push(
              {
                "data": {
                  "parent":parentid,
                  "label": cluster.id,
                  "id": cluster.id
                }
              }
            )

            parentid = cluster.id?cluster.id:parentid
            cluster.namespaces.map(
              ns => {
                ns.id && cyData.nodes.push(
                  {
                    "data": {
                      "parent":parentid,
                      "label": ns.id,
                      "id": ns.id
                    }
                  }
                )
              }
            )

            cluster.namespaces.map( ns => {

              parentid = ns.id?ns.id:parentid
              if(ns.deployments){
                ns.deployments.map(
                  deploy => {
                    deploy.id && cyData.nodes.push(
                      {
                        "data": {
                          "parent": parentid,
                          "id": deploy.id,
                          "label": deploy.id,
                        }
                      }
                    )
                  }
                )
              }
            }
            )

            cluster.namespaces.map( ns => {
              parentid = ns.id?ns.id:parentid
              if(ns.deployments){
                ns.deployments.map( deploy => {
                  parentid = deploy.id?deploy.id:parentid
                  deploy.pods.map(
                    pod => {
                      pod.id && cyData.nodes.push(
                        {
                          "data": {
                            "parent": parentid,
                            "id": pod.id,
                            "label": pod.id,
                          }
                        }
                      )
                    }
                  )
                }
                )
              }
            }
            )
          }
        )
      )
    }
    
    return cyData
  }
}