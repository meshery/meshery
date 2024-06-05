export const KubernetesManifests = {
  HelloWorld: {
    source_type: 'Kubernetes Manifest',
    name: 'GCP Sample',
    uri: 'https://raw.githubusercontent.com/GoogleCloudPlatform/microservices-demo/main/release/kubernetes-manifests.yaml',
    expectations: {
      numberOfComponents: 35,
      numberOfEdges: 20,
    },
  },
};
