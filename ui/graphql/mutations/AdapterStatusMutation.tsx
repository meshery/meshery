import { graphql, commitMutation } from 'react-relay';
import { createRelayEnvironment } from 'lib/relayEnvironment';

interface AdapterStatusVariables {
  status: string;
  targetPort: string | number;
  adapter: string;
}

export default function changeAdapterState(
  onComplete: (_response: unknown, _errors: ReadonlyArray<unknown> | null) => void,
  variables: AdapterStatusVariables,
  onError?: (_error: Error) => void,
) {
  const environment = createRelayEnvironment({});
  const vars = {
    input: {
      targetStatus: variables.status,
      targetPort: String(variables.targetPort),
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
    onCompleted: (response, errors) => onComplete(response, errors ?? null),
    onError: (error) => {
      if (onError) {
        onError(error);
        return;
      }

      // Fall back to console.error so failures aren't silent.
      console.error('Error in AdapterStatusMutation:', error);
    },
  });
}
