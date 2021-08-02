module.exports = {
    "env": {
        "node": true,
        "browser": true,
        "es6": true
    },
    'settings': {
        'react': {
            'version': require('./package.json').dependencies.react,
        },
    },
    "extends": [
        "eslint:recommended",
        "plugin:react/recommended",
        "plugin:cypress/recommended"
    ],
    "globals": {
        "Atomics": "readonly",
        "SharedArrayBuffer": "readonly"
    },
    'parser': 'babel-eslint',
    "parserOptions": {
        "ecmaFeatures": {
            "jsx": true
        },
        "ecmaVersion": 2018,
        "sourceType": "module"
    },
    "plugins": [
        "react",
        "cypress"
    ],
    "rules": {
        "arrow-spacing": [
            "error",
            {
                "after": true,
                "before": true
            }
        ],
        "block-scoped-var": "error",
        "block-spacing": "error",
        "brace-style": [
            "error",
            "1tbs"
        ],
        "react/react-in-jsx-scope": "off",
        "no-undef": "error",
        "react/prop-types": 0,
        "react/jsx-uses-vars": [
            2
        ],
        "react/jsx-no-undef": "error",
        "no-console": 0,
        "no-unused-vars": "error",
        "react/jsx-key": "off",
        "no-dupe-keys": "off",
        "react/jsx-filename-extension": [1, { "extensions": [".js", ".jsx"] }]
    }
};