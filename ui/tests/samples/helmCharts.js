export const HelmCharts = {
  HelloWorld: {
    source_type: 'Helm Chart',
    name: 'HelloWorld',
    uri: 'https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz',
    expectations: {
      numberOfComponents: 3,
      numberOfEdges: 1,
    },
  },
};
