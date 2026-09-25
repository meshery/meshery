import { RuleTester } from 'eslint';
import rule from '../eslint-rules/no-unclosed-graphql-ws-connections';

const ruleTester = new RuleTester({
  languageOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
  },
});

ruleTester.run('no-unclosed-graphql-ws-connections', rule, {
  valid: [
    {
      // Closing via the effect's cleanup return is acceptable.
      code: `
        import { createClient } from 'graphql-ws';
        function Component() {
          useEffect(() => {
            const client = createClient({ url: 'ws://localhost/graphql' });
            client.subscribe({ query: 'subscription { events }' }, {
              next() {}, error() {}, complete() {},
            });
            return () => client.dispose();
          }, []);
          return null;
        }
      `,
    },
    {
      // Closing the client inline (before the effect completes) is acceptable.
      code: `
        import { createClient } from 'graphql-ws';
        useEffect(() => {
          const client = createClient({ url: 'ws://localhost/graphql' });
          client.dispose();
        }, []);
      `,
    },
    {
      // A module-scope client lives for the app lifetime — not our concern.
      code: `
        import { createClient } from 'graphql-ws';
        export const client = createClient({ url: 'ws://localhost/graphql' });
      `,
    },
    {
      // Same shape, but the import is not graphql-ws.
      code: `
        import { createClient } from 'other-ws-library';
        useEffect(() => {
          const client = createClient({ url: 'ws://localhost/graphql' });
        }, []);
      `,
    },
    {
      // `new Client` from a module that is not graphql-ws is untouched.
      code: `
        import { Client } from 'somewhere-else';
        useEffect(() => {
          const ws = new Client({ url: 'ws://localhost/graphql' });
        }, []);
      `,
    },
    {
      // The client is stored in an outer variable — outside this rule's scope.
      code: `
        import { createClient } from 'graphql-ws';
        let sharedClient;
        useEffect(() => {
          sharedClient = createClient({ url: 'ws://localhost/graphql' });
        }, []);
      `,
    },
    {
      // A function-declaration cleanup that disposes of the client and is
      // returned by reference is acceptable.
      code: `
        import { createClient } from 'graphql-ws';
        function Component() {
          useEffect(() => {
            const client = createClient({ url: 'ws://localhost/graphql' });
            function cleanup() {
              client.dispose();
            }
            return cleanup;
          }, []);
          return null;
        }
      `,
    },
    {
      // A function-valued variable cleanup returned by reference is acceptable.
      code: `
        import { createClient } from 'graphql-ws';
        function Component() {
          useEffect(() => {
            const client = createClient({ url: 'ws://localhost/graphql' });
            const cleanup = () => {
              client.dispose();
            };
            return cleanup;
          }, []);
          return null;
        }
      `,
    },
  ],
  invalid: [
    {
      // createClient() inside useEffect with no cleanup or close call.
      code: `
        import { createClient } from 'graphql-ws';
        function Component() {
          useEffect(() => {
            const client = createClient({ url: 'ws://localhost/graphql' });
            client.subscribe({ query: 'subscription { events }' }, {
              next() {}, error() {}, complete() {},
            });
          }, []);
          return null;
        }
      `,
      errors: [{ messageId: 'unclosedWsConnection' }],
    },
    {
      // Older-style `new Client(...)` is caught too.
      code: `
        import { Client } from 'graphql-ws';
        useEffect(() => {
          const ws = new Client({ url: 'ws://localhost/graphql' });
        }, []);
      `,
      errors: [{ messageId: 'unclosedWsConnection' }],
    },
    {
      // Namespace imports are resolved.
      code: `
        import * as gqlWs from 'graphql-ws';
        useEffect(() => {
          const ws = gqlWs.createClient({ url: 'ws://localhost/graphql' });
        }, []);
      `,
      errors: [{ messageId: 'unclosedWsConnection' }],
    },
    {
      // Default import aliases are resolved.
      code: `
        import createClient from 'graphql-ws';
        useEffect(() => {
          const client = createClient({ url: 'ws://localhost/graphql' });
        }, []);
      `,
      errors: [{ messageId: 'unclosedWsConnection' }],
    },
    {
      // An empty cleanup does not dispose of the client.
      code: `
        import { createClient } from 'graphql-ws';
        useEffect(() => {
          const client = createClient({ url: 'ws://localhost/graphql' });
          return () => {};
        }, []);
      `,
      errors: [{ messageId: 'unclosedWsConnection' }],
    },
    {
      // A cleanup that disposes of a different client does not count.
      code: `
        import { createClient } from 'graphql-ws';
        useEffect(() => {
          const client = createClient({ url: 'ws://localhost/graphql' });
          return () => otherClient.dispose();
        }, []);
      `,
      errors: [{ messageId: 'unclosedWsConnection' }],
    },
    {
      // client.terminate() is not a disposal — it only drops the current
      // socket and the connection may retry.
      code: `
        import { createClient } from 'graphql-ws';
        useEffect(() => {
          const client = createClient({ url: 'ws://localhost/graphql' });
          client.terminate();
        }, []);
      `,
      errors: [{ messageId: 'unclosedWsConnection' }],
    },
    {
      // A function-declaration cleanup that is declared but never returned
      // or called does not dispose of the client.
      code: `
        import { createClient } from 'graphql-ws';
        useEffect(() => {
          const client = createClient({ url: 'ws://localhost/graphql' });
          function cleanup() {
            client.dispose();
          }
        }, []);
      `,
      errors: [{ messageId: 'unclosedWsConnection' }],
    },
  ],
});
