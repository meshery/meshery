export const prometheusSchema = {
  type: 'object',
  title: 'Prometheus',
  properties: {
    credentialName: {
      title: 'Credential Name',
      type: 'string',
      description: 'Name of your credential',
    },
    secret: {
      type: 'object',
      title: 'Credential Secret',
      description: 'Credential secret for the Prometheus instance',
      properties: {
        prometheusURL: {
          type: 'string',
          title: 'URL',
          description: 'URL of the Prometheus instance',
        },
      },
      required: ['prometheusURL'],
    },
  },
  required: ['credentialName'],
};
