export const kubernetesSchema = {
  type: 'object',
  title: 'Kubernetes',
  properties: {
    credentialName: {
      title: 'Credential Name',
      type: 'string',
      description: 'Name of your credential',
    },
    secret: {
      type: 'object',
      title: 'Credential Secret',
      description: 'Credential secret for the Kubernetes cluster',
      properties: {
        clusterName: {
          type: 'string',
          title: 'Cluster Name',
          description: 'Name of the Kubernetes cluster',
        },
        clusterServerURL: {
          type: 'string',
          title: 'Server URL',
          description: 'URL of the Kubernetes cluster',
        },
        auth: {
          type: 'object',
          title: 'Auth',
          description: 'Kubernetes cluster authentication',
          properties: {
            clusterUserName: {
              type: 'string',
              title: 'User Name',
              description: 'Name of the Kubernetes cluster user',
            },
            clusterToken: {
              type: 'string',
              title: 'Token',
              description: 'Token of the Kubernetes cluster user',
            },
            clusterClientCertificateData: {
              type: 'string',
              title: 'Client Certificate Data',
              description: 'Certificate data of the Kubernetes cluster',
            },
            clusterClientKeyData: {
              type: 'string',
              title: 'Client Key Data',
              description: 'Client Key data of the Kubernetes cluster',
            },
            clusterCertificateAuthorityData: {
              type: 'string',
              title: 'Certificate Authority Data',
              description: 'Certificate Authority data of the Kubernetes cluster',
            },
          },
          required: [
            'clusterUserName',
            'clusterToken',
            'clusterClientCertificateData',
            'clusterClientKeyData',
            'clusterCertificateAuthorityData',
          ],
        },
      },
      required: ['clusterName', 'clusterServerURL'],
    },
  },
  required: ['credentialName'],
};
