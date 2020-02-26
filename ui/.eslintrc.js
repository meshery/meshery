module.exports = {
    "env": {
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
        "plugin:react/recommended"
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
    "plugins": ["react"],
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
        "no-undef": "off",
        "react/prop-types": "off"
    }
};