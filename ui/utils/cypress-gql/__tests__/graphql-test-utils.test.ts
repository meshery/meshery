import { describe, expect, it } from 'vitest';
import { aliasMutation, aliasQuery, hasOperationName } from '../graphql-test-utils';

describe('hasOperationName', () => {
  it('returns true when the request body contains the matching operationName', () => {
    const req = { body: { operationName: 'GetUser' } } as never;
    expect(hasOperationName(req, 'GetUser')).toBe(true);
  });

  it('returns false when the request body has a different operationName', () => {
    const req = { body: { operationName: 'OtherOp' } } as never;
    expect(hasOperationName(req, 'GetUser')).toBe(false);
  });

  it('returns false when the request body does not have operationName key', () => {
    const req = { body: {} } as never;
    expect(hasOperationName(req, 'GetUser')).toBe(false);
  });

  it('does not include inherited operationName properties', () => {
    // The function uses Object.prototype.hasOwnProperty.call to guard against
    // prototype lookups. Confirm an inherited operationName is ignored.
    const proto = { operationName: 'GetUser' };
    const body = Object.create(proto);
    const req = { body } as never;
    expect(hasOperationName(req, 'GetUser')).toBe(false);
  });
});

describe('aliasQuery', () => {
  it('assigns a gql<Name>Query alias when operationName matches', () => {
    const req = { body: { operationName: 'Foo' } } as { body: unknown; alias?: string };
    aliasQuery(req as never, 'Foo');
    expect(req.alias).toBe('gqlFooQuery');
  });

  it('does nothing when operationName does not match', () => {
    const req = { body: { operationName: 'Bar' } } as { body: unknown; alias?: string };
    aliasQuery(req as never, 'Foo');
    expect(req.alias).toBeUndefined();
  });
});

describe('aliasMutation', () => {
  it('assigns a gql<Name>Mutation alias when operationName matches', () => {
    const req = { body: { operationName: 'CreateThing' } } as {
      body: unknown;
      alias?: string;
    };
    aliasMutation(req as never, 'CreateThing');
    expect(req.alias).toBe('gqlCreateThingMutation');
  });

  it('does not touch the alias when operationName does not match', () => {
    const req = { body: { operationName: 'OtherThing' }, alias: 'preset' } as {
      body: unknown;
      alias: string;
    };
    aliasMutation(req as never, 'CreateThing');
    expect(req.alias).toBe('preset');
  });
});
