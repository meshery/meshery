const nextConfig = require('eslint-config-next');
const prettierRecommended = require('eslint-plugin-prettier/recommended');
const unusedImports = require('eslint-plugin-unused-imports');
const globals = require('globals');

// ESLint 10: eslint-config-next's babel-based parser returns a scope manager that
// doesn't implement addGlobals (new ESLint 10 API). Replace it with espree (ESLint's
// built-in parser) for JS/JSX files; the TS entry already uses @typescript-eslint/parser.
const patchedNextConfig = nextConfig.map((cfg) => {
  if (cfg.name === 'next') {
    const { parser: _babelParser, globals: _g, ...restLangOpts } = cfg.languageOptions ?? {};
    return {
      ...cfg,
      languageOptions: {
        ...restLangOpts,
        parserOptions: {
          ...restLangOpts.parserOptions,
          ecmaFeatures: { jsx: true },
          ecmaVersion: 'latest',
          sourceType: 'module',
        },
      },
    };
  }
  return cfg;
});

// Temporary allowlists for legacy files that still violate the new UI guardrails.
// Keep the rules active for new/clean files while letting incremental refactors
// remove entries from these lists over time.
const legacyRestrictedImportOffenders = [
  'components/Dashboard/charts/NodeStatusChart.tsx',
  'components/Dashboard/charts/PodStatusChart.tsx',
  'components/Dashboard/index.tsx',
  'components/Dashboard/resources/resources-sub-menu.tsx',
  'components/Dashboard/widgets/HoneyComb/HoneyCombComponent.tsx',
  'components/DataFormatter/index.tsx',
  'components/DesignLifeCycle/DeploymentSummary.tsx',
  'components/DesignLifeCycle/DryRun.tsx',
  'components/DesignLifeCycle/ValidateDesign.tsx',
  'components/DesignLifeCycle/styles.tsx',
  'components/ExportModal.tsx',
  'components/General/Modals/Information/InfoModal.tsx',
  'components/General/Modals/Modal.tsx',
  'components/Header.tsx',
  'components/Lifecycle/Environments/environment-card.tsx',
  'components/Lifecycle/Environments/index.tsx',
  'components/Lifecycle/Workspaces/WorkspaceGridView.tsx',
  'components/MesheryAdapterPlayComponent.tsx',
  'components/MesheryChart.tsx',
  'components/MesheryDateTimePicker.tsx',
  'components/MesheryFilters/Filters.tsx',
  'components/MesheryFilters/FiltersCard.tsx',
  'components/MesheryMeshInterface/PatternService/RJSF.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/ArrayFieldTemlate.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/CustomFileWidget.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/CustomSelectWidget.tsx',
  'components/MesheryMeshInterface/PatternService/helper.tsx',
  'components/MesheryMeshInterface/PatternServiceForm.tsx',
  'components/MesheryPatterns/MesheryPatternCard.tsx',
  'components/MesheryPatterns/MesheryPatterns.tsx',
  'components/MesheryPlayComponent.tsx',
  'components/MesherySettingsEnvButtons.tsx',
  'components/NotificationCenter/formatters/common.tsx',
  'components/NotificationCenter/index.tsx',
  'components/NotificationCenter/notificationCenter.style.tsx',
  'components/Performance/MesheryResults.tsx',
  'components/Performance/PerformanceCard.tsx',
  'components/Performance/PerformanceResults.tsx',
  'components/Performance/index.tsx',
  'components/Performance/style.tsx',
  'components/ReactSelectWrapper.tsx',
  'components/Registry/RegistryModal.tsx',
  'components/RelationshipBuilder/RelationshipFormStepper.tsx',
  'components/Settings/Registry/ComponentTree.tsx',
  'components/Settings/Registry/MeshModel.style.ts',
  'components/Settings/Registry/MeshModelComponent.tsx',
  'components/Settings/Registry/MeshModelDetails.tsx',
  'components/Settings/Registry/MesheryTreeView.tsx',
  'components/Settings/Registry/MesheryTreeViewModel.tsx',
  'components/Settings/Registry/MesheryTreeViewRegistrants.tsx',
  'components/Settings/Registry/RelationshipTree.tsx',
  'components/Settings/Registry/Stepper/CSVStepper.tsx',
  'components/Settings/Registry/Stepper/UrlStepper.tsx',
  'components/SpacesSwitcher/MainDesignsContent.tsx',
  'components/SpacesSwitcher/WorkspaceModal.tsx',
  'components/SpacesSwitcher/components.tsx',
  'components/TroubleshootingModalComponent.tsx',
  'components/UserPreferences/index.tsx',
  'components/ViewInfoModal.tsx',
  'components/ViewSwitch.tsx',
  'components/YamlDialog.tsx',
  'components/configuratorComponents/MeshModel/index.tsx',
  'components/configuratorComponents/NameToIcon.tsx',
  'components/connections/ConnectionChip.tsx',
  'components/connections/ConnectionTable.tsx',
  'components/icons/index.ts',
  'components/telemetry/grafana/GrafanaCustomChart.tsx',
  'components/telemetry/grafana/GrafanaDateRangePicker.tsx',
  'components/telemetry/prometheus/PrometheusSelectionComponent.tsx',
  'pages/_app.tsx',
];

const legacyLiteralColorOffenders = [
  'components/ConfirmationModal.tsx',
  'components/Dashboard/charts/ResourceUtilizationChart.tsx',
  'components/Dashboard/components.tsx',
  'components/Dashboard/images/info-icon.tsx',
  'components/Dashboard/images/meshery-icon.tsx',
  'components/Dashboard/style.ts',
  'components/DatabaseSummary.tsx',
  'components/DesignLifeCycle/DryRun.tsx',
  'components/DesignLifeCycle/ValidateDesign.tsx',
  'components/DesignLifeCycle/common.tsx',
  'components/ExportModal.tsx',
  'components/General/Modals/Information/InfoModal.tsx',
  'components/General/Modals/Modal.tsx',
  'components/General/TipsCarousel.tsx',
  'components/General/error-404/CurrentSession.tsx',
  'components/General/error-404/socials/styles.tsx',
  'components/General/error-404/styles.tsx',
  'components/Header.styles.tsx',
  'components/Header.tsx',
  'components/Lifecycle/Environments/environment-card.tsx',
  'components/Lifecycle/Environments/index.tsx',
  'components/Lifecycle/Environments/styles.tsx',
  'components/Lifecycle/General/empty-state/curvedArrowIcon.tsx',
  'components/Lifecycle/General/empty-state/index.tsx',
  'components/Lifecycle/General/flip-card/index.tsx',
  'components/Lifecycle/Workspaces/index.tsx',
  'components/LoadingComponents/Animations/AnimatedFilter.tsx',
  'components/LoadingComponents/Animations/AnimatedLightMeshery.tsx',
  'components/LoadingComponents/Animations/AnimatedMeshPattern.tsx',
  'components/LoadingComponents/Animations/AnimatedMeshery.tsx',
  'components/LoadingComponents/Animations/AnimatedMesheryCSS.tsx',
  'components/LoadingComponents/LoadingComponentServer.tsx',
  'components/MeshAdapterConfigComponent.tsx',
  'components/MesheryAdapterPlayComponent.tsx',
  'components/MesheryChart.tsx',
  'components/MesheryFilters/Filters.tsx',
  'components/MesheryFilters/FiltersGrid.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/ArrayFieldTemlate.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/CustomBaseInput.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/ObjectFieldTemplate.tsx',
  'components/MesheryMeshInterface/PatternServiceForm.tsx',
  'components/MesheryPatterns/MesheryPatternCard.tsx',
  'components/MesheryPatterns/MesheryPatternGridView.tsx',
  'components/MesheryPatterns/MesheryPatterns.tsx',
  'components/MesheryPatterns/style.tsx',
  'components/NavigatorExtension.tsx',
  'components/NotificationCenter/constants.tsx',
  'components/NotificationCenter/formatters/relationship_evaluation.tsx',
  'components/NotificationCenter/index.tsx',
  'components/NotificationCenter/notificationCenter.style.tsx',
  'components/Performance/PerformanceCard.tsx',
  'components/Performance/PerformanceResults.tsx',
  'components/Performance/assets/facebookIcon.tsx',
  'components/Performance/assets/linkedinIcon.tsx',
  'components/Performance/assets/twitterIcon.tsx',
  'components/Performance/index.tsx',
  'components/Performance/style.tsx',
  'components/Settings/MesherySettings.tsx',
  'components/Settings/Registry/MeshModelDetails.tsx',
  'components/Settings/Registry/helper.ts',
  'components/SpacesSwitcher/WorkspaceSwitcher.tsx',
  'components/SpacesSwitcher/styles.tsx',
  'components/StyledAccordion.tsx',
  'components/TroubleshootingModalComponent.tsx',
  'components/TypingFilter/style.tsx',
  'components/UserPreferences/index.tsx',
  'components/UserPreferences/style.tsx',
  'components/ViewInfoModal.tsx',
  'components/configuratorComponents/CustomBreadCrumb.tsx',
  'components/configuratorComponents/MeshModel/styledComponents/AppBar.tsx',
  'components/configuratorComponents/MeshModel/utils.tsx',
  'components/configuratorComponents/NameToIcon.tsx',
  'components/connections/meshSync/Stepper/Notification.tsx',
  'components/connections/meshSync/Stepper/StepperContent.tsx',
  'components/connections/meshSync/Stepper/StepperContentWrapper.tsx',
  'components/connections/meshSync/Stepper/index.tsx',
  'components/connections/styles.tsx',
  'components/load-test-timer-dialog.tsx',
  'components/multi-select-wrapper.tsx',
  'components/telemetry/grafana/GrafanaCustomGaugeChart.tsx',
  'components/telemetry/grafana/GrafanaDateRangePicker.tsx',
  'css/icons.styles.ts',
  'pages/extension/AccessMesheryModal.tsx',
  'utils/charts.ts',
  'utils/custom-search.tsx',
  'utils/lightenOrDarkenColor.ts',
];

const legacyMaxLineOffenders = [
  'components/Dashboard/resources/configuration/config.tsx',
  'components/Dashboard/resources/workloads/config.tsx',
  'components/MesheryAdapterPlayComponent.tsx',
  'components/MesheryFilters/Filters.tsx',
  'components/MesheryPatterns/MesheryPatterns.tsx',
  'components/Performance/index.tsx',
  'components/connections/ConnectionTable.tsx',
];

module.exports = [
  // Global ignores (replaces .eslintignore — not supported in flat config)
  {
    ignores: [
      'node_modules/**',
      'out/**',
      '.next/**',
      'static/**',
      'public/static/**',
      'lib/**',
      'tests/samples/**',
      '**/__generated__/**',
      'playwright-report/**',
      'playground/**',
      'test-results/**',
      // Non-JS/TS assets — ESLint 10 flat config processes every non-ignored file
      // in the directory tree when given `.` as the argument; these would be parsed
      // as JavaScript and cause hangs or parse errors.
      '**/*.svg',
      '**/*.png',
      '**/*.gif',
      '**/*.webp',
      '**/*.jpg',
      '**/*.jpeg',
      '**/*.wasm',
      '**/*.zip',
      '**/*.webm',
      '**/*.css',
      '**/*.html',
      '**/*.md',
      '**/*.json',
      '**/*.yml',
      '**/*.yaml',
      '**/*.txt',
      '**/*.csv',
      '**/*.otf',
      '**/*.woff',
      '**/*.woff2',
      '**/*.ttf',
    ],
  },

  // Globals via default parser (avoids babel parser / addGlobals incompatibility)
  {
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
        Atomics: 'readonly',
        SharedArrayBuffer: 'readonly',
        globalThis: 'readonly',
      },
      parserOptions: {
        ecmaFeatures: { jsx: true },
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
    },
  },

  // Next.js flat config (includes react, react-hooks, @next/next rules)
  ...patchedNextConfig,

  // Prettier integration (flat config format — disables conflicting style rules)
  prettierRecommended,

  // Custom overrides
  {
    plugins: {
      'unused-imports': unusedImports,
    },
    settings: {
      // eslint-plugin-react calls context.getFilename() during 'detect' (removed in ESLint 9+).
      // Provide an explicit version to skip detection entirely.
      react: { version: '19' },
    },
    rules: {
      '@next/next/no-img-element': 'off',
      'react-hooks/rules-of-hooks': 'warn',
      'react-hooks/exhaustive-deps': 'off',
      // Disabled: all React Compiler rules added by react-hooks v7 via eslint-config-next.
      // This project does not use the React Compiler; these rules are inapplicable and slow.
      'react-hooks/static-components': 'off',
      'react-hooks/use-memo': 'off',
      'react-hooks/set-state-in-effect': 'off',
      'react-hooks/component-hook-factories': 'off',
      'react-hooks/preserve-manual-memoization': 'off',
      'react-hooks/incompatible-library': 'off',
      'react-hooks/immutability': 'off',
      'react-hooks/globals': 'off',
      'react-hooks/refs': 'off',
      'react-hooks/error-boundaries': 'off',
      'react-hooks/purity': 'off',
      'react-hooks/set-state-in-render': 'off',
      'react-hooks/unsupported-syntax': 'off',
      'react-hooks/config': 'off',
      'react-hooks/gating': 'off',
      'jsx-a11y/alt-text': 'off',
      'valid-typeof': 'warn',
      'react/react-in-jsx-scope': 'off',
      'no-undef': 'error',
      'react/jsx-uses-vars': [2],
      'react/jsx-no-undef': 'error',
      'no-console': 0,
      'unused-imports/no-unused-imports': 'error',
      'react/jsx-key': 'warn',
      'no-dupe-keys': 'error',
      'react/prop-types': 'off',
      'prettier/prettier': ['error', { endOfLine: 'lf' }],

      // ---------------------------------------------------------------------
      // UI restructure guardrails (warn mode, phase 1).
      //
      // These rules encode the target architecture: one design system
      // (@sistent/sistent), one theme source (@/theme), and a size budget
      // for component files. They ship as warnings so CI stays green on
      // day one; a later phase will allowlist today's offenders and promote
      // the rules to errors.
      // ---------------------------------------------------------------------

      // Ban Material UI and legacy theme imports. @sistent/sistent is the
      // only UI kit; @/theme is the approved Phase 1 theme entry point.
      'no-restricted-imports': [
        'warn',
        {
          paths: [
            {
              name: '@/theme/index',
              message: 'Use @/theme; do not deep-import the local theme entry point.',
            },
            {
              name: '@mui/material',
              message: 'Use @sistent/sistent instead.',
            },
            {
              name: '@mui/icons-material',
              message: 'Use @sistent/sistent icons, or add an SVG component to ui/assets/icons.',
            },
            {
              name: '@mui/x-date-pickers',
              message:
                'Wrap @mui/x-date-pickers in a single shared primitive; do not import it directly.',
            },
            {
              name: '@mui/x-tree-view',
              message:
                'Wrap @mui/x-tree-view in a single shared primitive; do not import it directly.',
            },
            {
              name: '@rjsf/mui',
              message: 'Use the shared RJSF wrapper; do not import @rjsf/mui directly.',
            },
            {
              name: '@/themes',
              message: 'Use @/theme, the approved Phase 1 theme entry point.',
            },
            {
              name: '@/themes/app',
              message: 'Use @/theme and theme.palette.* instead of the legacy Colors object.',
            },
            {
              name: '@/themes/index',
              message: 'Use @/theme and theme.palette.* instead of NOTIFICATIONCOLORS.',
            },
            {
              name: '@/constants/colors',
              message: 'Use @/theme and theme.palette.* instead of legacy color constants.',
            },
          ],
          patterns: [
            {
              group: ['@mui/*'],
              message: 'Use @sistent/sistent instead.',
            },
            {
              group: ['@material-ui/*'],
              message: 'Material UI v4 is deprecated in this project — use @sistent/sistent.',
            },
          ],
        },
      ],

      // Size budget for component files. 1000 lines is the hard ceiling;
      // the plan is to drop this to 600 once the eight giant files are
      // broken up in phase 5.
      'max-lines': ['warn', { max: 1000, skipComments: true, skipBlankLines: true }],
    },
  },

  // ---------------------------------------------------------------------
  // Ban hex (#RRGGBB) and rgb()/rgba() literals in source files.
  //
  // Colors must come from theme.palette.* (or be composed with alpha() /
  // lighten() / darken() from @sistent/sistent). The only places allowed
  // to contain a literal color are:
  //
  //   - ui/theme/**       (the theme module itself)
  //   - ui/themes/**      (legacy theme module, scheduled for deletion)
  //   - ui/assets/**      (SVG icons encoded as React components)
  //   - ui/constants/**   (legacy color constants, scheduled for deletion)
  //   - ui/lib/**         (third-party integration helpers)
  //   - ui/public/**      (static assets)
  // ---------------------------------------------------------------------
  {
    files: ['**/*.{ts,tsx,js,jsx}'],
    ignores: [
      'theme/**',
      'themes/**',
      'assets/**',
      'constants/**',
      'lib/**',
      'public/**',
      'tests/**',
      'scripts/**',
      'eslint.config.js',
    ],
    rules: {
      'no-restricted-syntax': [
        'warn',
        {
          selector: 'Literal[value=/^#(?:[0-9a-fA-F]{3,4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$/]',
          message: 'Hex color literals are forbidden outside ui/theme/. Use theme.palette.*.',
        },
        {
          selector: 'Literal[value=/rgba?\\(/]',
          message:
            'rgb()/rgba() literals are forbidden outside ui/theme/. Use theme.palette.* (or alpha() from @sistent/sistent).',
        },
      ],
    },
  },

  // Current legacy violations that are being refactored incrementally.
  {
    files: legacyRestrictedImportOffenders,
    rules: {
      'no-restricted-imports': 'off',
    },
  },
  {
    files: legacyLiteralColorOffenders,
    rules: {
      'no-restricted-syntax': 'off',
    },
  },
  {
    files: legacyMaxLineOffenders,
    rules: {
      'max-lines': 'off',
    },
  },

  // no-unused-vars: JS/JSX only — TypeScript files should use @typescript-eslint/no-unused-vars
  {
    files: ['**/*.{js,jsx,mjs,cjs}'],
    rules: {
      'no-unused-vars': ['error', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
    },
  },
];
