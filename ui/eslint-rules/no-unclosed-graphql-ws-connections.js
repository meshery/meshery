/**
 * @fileoverview Disallows GraphQL-WS WebSocket connections created inside a
 * React useEffect callback without a matching close.
 *
 * Meshery UI ships graphql-ws (see ui/remote-component.config.js). Creating a
 * client inside a useEffect and never disposing of it leaks the socket for the
 * lifetime of the SPA. This rule flags client creation in useEffect when the
 * effect neither (a) returns a cleanup function that closes the connection, nor
 * (b) calls a close method (dispose/close/terminate/unsubscribe/down) on the
 * client within the effect body.
 *
 * graphql-ws exposes both `createClient(...)` (v5+) and a `Client` class
 * (`new Client(...)`) for older versions; named, default and namespace imports
 * from 'graphql-ws' are covered.
 */
'use strict';

const GRAPHQL_WS_MODULE = 'graphql-ws';
const CLOSE_METHODS = new Set(['dispose', 'close', 'terminate', 'unsubscribe', 'down']);

function isGraphQLWsCreation(node, factoryNames, classNames, namespaceNames) {
  const callee = node.callee;
  if (callee.type === 'Identifier') {
    if (node.type === 'CallExpression') {
      return factoryNames.has(callee.name);
    }
    if (node.type === 'NewExpression') {
      return classNames.has(callee.name);
    }
    return false;
  }

  if (callee.type !== 'MemberExpression') return false;
  if (callee.object.type !== 'Identifier' || !namespaceNames.has(callee.object.name)) return false;
  if (callee.property.type !== 'Identifier') return false;

  if (node.type === 'CallExpression') {
    return callee.property.name === 'createClient';
  }
  if (node.type === 'NewExpression') {
    return callee.property.name === 'Client';
  }
  return false;
}

function findSurroundingUseEffect(node) {
  for (let current = node.parent; current; current = current.parent) {
    if (current.type !== 'CallExpression') continue;
    const callee = current.callee;
    const isUseEffect =
      (callee.type === 'Identifier' && callee.name === 'useEffect') ||
      (callee.type === 'MemberExpression' &&
        callee.property.type === 'Identifier' &&
        callee.property.name === 'useEffect');
    if (isUseEffect) return current;
  }
  return null;
}

function isFunction(node) {
  return node.type === 'ArrowFunctionExpression' || node.type === 'FunctionExpression';
}

/**
 * Walks a function's body and reports whether it returns a cleanup function
 * (a non-void return) or calls any close method, without descending into
 * nested functions (a nested function's call does not manage the effect's
 * own resources synchronously at cleanup time).
 */
function effectManagesResources(body) {
  if (!body || body.type !== 'BlockStatement') return false;

  let managed = false;
  const visit = (node) => {
    if (managed) return;
    if (node.type === 'ReturnStatement' && node.argument && isFunction(node.argument)) {
      managed = true;
      return;
    }
    if (
      node.type === 'CallExpression' &&
      node.callee.type === 'MemberExpression' &&
      node.callee.property.type === 'Identifier' &&
      CLOSE_METHODS.has(node.callee.property.name)
    ) {
      managed = true;
      return;
    }
    if (node.type === 'FunctionExpression' || node.type === 'ArrowFunctionExpression') return;
    for (const key of Object.keys(node)) {
      if (key === 'parent') continue;
      const child = node[key];
      if (Array.isArray(child)) {
        for (const item of child) {
          if (item && typeof item.type === 'string') visit(item);
        }
      } else if (child && typeof child.type === 'string') {
        visit(child);
      }
    }
  };

  for (const statement of body.body) visit(statement);
  return managed;
}

function findEnclosingVariableDeclarator(node) {
  for (let current = node.parent; current; current = current.parent) {
    if (
      current.type === 'VariableDeclarator' &&
      current.init === node &&
      current.id.type === 'Identifier'
    ) {
      return current;
    }
    // Do not look past function boundaries — a client handed to an outer
    // variable is outside this rule's scope.
    if (current.type === 'FunctionExpression' || current.type === 'ArrowFunctionExpression') break;
  }
  return null;
}

module.exports = {
  meta: {
    type: 'problem',
    docs: {
      description: 'Disallows unclosed GraphQL-WS WebSocket connections created in useEffect',
      category: 'Possible Errors',
      recommended: true,
    },
    schema: [],
    messages: {
      unclosedWsConnection:
        'Unclosed GraphQL-WS WebSocket connection in useEffect. Close the connection (call client.dispose()) or return a cleanup function before the effect completes.',
    },
  },

  create(context) {
    const factoryNames = new Set();
    const classNames = new Set();
    const namespaceNames = new Set();

    return {
      ImportDeclaration(node) {
        if (node.source.value !== GRAPHQL_WS_MODULE) return;
        for (const spec of node.specifiers) {
          if (spec.type === 'ImportSpecifier') {
            if (spec.imported.name === 'createClient') factoryNames.add(spec.local.name);
            else if (spec.imported.name === 'Client') classNames.add(spec.local.name);
          } else if (spec.type === 'ImportDefaultSpecifier') {
            factoryNames.add(spec.local.name);
          } else if (spec.type === 'ImportNamespaceSpecifier') {
            namespaceNames.add(spec.local.name);
          }
        }
      },

      CallExpression(node) {
        reportUnclosedCreation(context, node, factoryNames, classNames, namespaceNames);
      },

      NewExpression(node) {
        reportUnclosedCreation(context, node, factoryNames, classNames, namespaceNames);
      },
    };

    function reportUnclosedCreation(context, node, factoryNames, classNames, namespaceNames) {
      if (!isGraphQLWsCreation(node, factoryNames, classNames, namespaceNames)) return;

      const useEffectCall = findSurroundingUseEffect(node);
      if (!useEffectCall) return;

      const callback = useEffectCall.arguments[0];
      if (!callback || !isFunction(callback)) return;

      if (!findEnclosingVariableDeclarator(node)) return;

      // A created-but-never-disposed client is fine as long as the effect
      // returns a cleanup or calls a close method anywhere in its body.
      if (effectManagesResources(callback.body)) return;

      context.report({
        node,
        messageId: 'unclosedWsConnection',
      });
    }
  },
};
