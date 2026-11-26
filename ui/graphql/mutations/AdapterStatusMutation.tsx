import { graphql, commitMutation } from 'react-relay';
import { createRelayEnvironment } from 'lib/relayEnvironment';

interface AdapterStatusVariables {
  status: string;
  targetPort: string;
  adapter: string;
}

export default function changeAdapterState(
  onComplete: (_response: unknown, _errors?: unknown) => void,
  variables: AdapterStatusVariables,
) {
  const environment = createRelayEnvironment({});
  const vars = {
    input: {
      targetStatus: variables.status,
      targetPort: variables.targetPort,
      adapter: variables.adapter,
    },
  };

  const adapterStatusMutation = graphql`
    mutation AdapterStatusMutation($input: AdapterStatusInput) {
      adapterStatus: changeAdapterStatus(input: $input)
    }
  `;

  commitMutation(environment, {
    mutation: adapterStatusMutation,
    variables: vars,
    onCompleted: onComplete,
    // Hand network/server failures to the caller too, otherwise it never hears
    // back and can't reset its spinner or show an error.
    onError: (error) => onComplete(null, [error]),
  });
}
