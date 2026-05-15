/**
 * Relay Configuration
 *
 */

module.exports = {
  // Configuration options accepted by the `relay-compiler` command-line tool and `babel-plugin-relay`.
  language: 'typescript',
  src: './',
  schema: '../server/internal/graphql/schema/schema.graphql',
  // Remove the "%future added value" catch-all from generated enum unions.
  // Without it the generated TS is fully type-safe and does not need @ts-nocheck.
  noFutureProofEnums: true,
};
