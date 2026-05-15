export const grafanaSchema = {
  type: 'object',
  title: 'Grafana',
  properties: {
    credentialName: {
      title: 'Credential Name',
      type: 'string',
      description: 'Name of your credential',
    },
    secret: {
      type: 'object',
      title: 'Credential Secret',
      description: 'Credential secret for the Grafana instance',
      properties: {
        grafanaURL: {
          type: 'string',
          title: 'URL',
          description: 'URL of the Grafana instance',
        },
        grafanaAPIKey: {
          type: 'string',
          title: 'API Key',
          description: 'API Key for the Grafana instance',
        },
      },
      required: ['grafanaURL', 'grafanaAPIKey'],
    },
  },
  required: ['credentialName'],
};
