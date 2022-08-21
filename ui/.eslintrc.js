module.exports = {
  "env" : {
    "browser" : true,
    "es6" : true,
    "node" : true // tells the parser that we are using nodejs
  },
  'settings' : { 'react' : { 'version' : require('./package.json').dependencies.react, }, },
  "extends" : [
    "eslint:recommended",
    "plugin:react/recommended",
    "plugin:cypress/recommended"
  ],
  "globals" : {
    "Atomics" : "readonly",
    "SharedArrayBuffer" : "readonly"
  },
  "parser" : "@babel/eslint-parser",
  "parserOptions" : {
    "ecmaFeatures" : { "jsx" : true },
    "ecmaVersion" : 2018,
    "sourceType" : "module"
  },
  "plugins" : [
    "react",
  ],
  "rules" : {
    "array-bracket-spacing" : ["error", "never"],
    "comma-style" : ["error"],
    "arrow-spacing" : [
      "error",
      {
        "after" : true,
        "before" : true
      }
    ],
    "jsx-quotes" : ["error", "prefer-double"],
    "block-scoped-var" : "error",
    "keyword-spacing" : "error",
    "no-trailing-spaces" : "error",
    "object-curly-spacing" : ["error", "always"],
    "arrow-spacing" : ["error", {
      "before" : true,
      "after" : true,
    }],
    "key-spacing" : ["error", {
      "beforeColon" : true,
      "afterColon" : true,
    }],
    "block-spacing" : "error",
    "brace-style" : [
      "error",
      "1tbs"
    ],
    'indent' : [
      'error', 2, {
        "FunctionExpression" : { "parameters" : "first" },
        "FunctionDeclaration" : { "parameters" : "first" },
        "MemberExpression" : 1,
        "SwitchCase" : 1,
        "outerIIFEBody" : 0,
        "VariableDeclarator" : { "var" : 2, "let" : 2, "const" : 3 },
        ignoredNodes : ['TemplateLiteral']
      }
    ],
    "react/react-in-jsx-scope" : "off",
    "no-undef" : "error",
    "react/prop-types" : 0,
    "react/jsx-uses-vars" : [
      2
    ],
    "react/jsx-no-undef" : "error",
    "no-console" : 0,
    "no-unused-vars" : "error",
    "react/jsx-key" : "off",
    "no-dupe-keys" : "off",
    "react/jsx-filename-extension" : [1, { "extensions" : [".js", ".jsx"] }],
    "react/prop-types" : "off"
  }
};
