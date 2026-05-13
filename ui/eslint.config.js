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
  'components/UserPreferences/index.tsx',
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
  'components/Dashboard/charts/ResourceUtilizationChart.tsx',
  'components/Dashboard/components.tsx',
  'components/Dashboard/images/info-icon.tsx',
  'components/Dashboard/images/meshery-icon.tsx',
  'components/Dashboard/style.ts',
  'components/DesignLifeCycle/DryRun.tsx',
  'components/DesignLifeCycle/ValidateDesign.tsx',
  'components/DesignLifeCycle/common.tsx',
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
  'components/TypingFilter/style.tsx',
  'components/UserPreferences/index.tsx',
  'components/UserPreferences/style.tsx',
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

// Files currently in the 600–1000 line "soft" range. They exceed the 600-line
// proactive warning threshold (§8.4) but stay under the hard 1000-line ceiling.
// Allowlisted so CI stays green; entries leave the list as files get split up.
const legacyMaxLineSoftOffenders = [
  'components/Dashboard/resources/network/config.tsx',
  'components/Lifecycle/Environments/index.tsx',
  'components/Navigator.tsx',
  'components/Performance/PerformanceResults.tsx',
  'components/Settings/Registry/Stepper/UrlStepper.tsx',
  'components/UserPreferences/index.tsx',
  'components/connections/meshSync/index.tsx',
  'components/telemetry/grafana/GrafanaCustomChart.tsx',
  'components/telemetry/grafana/GrafanaDateRangePicker.tsx',
  // This config file itself: the legacy allowlists above push it past 600
  // lines. Excluding it from the soft cap until the lists thin out naturally.
  'eslint.config.js',
  'pages/_app.tsx',
];

// Files that currently use inline `style={{ ... }}` props. The §8.3 guardrail
// nudges new code toward styled() from @sistent/sistent; existing offenders
// stay allowlisted until they are migrated.
const legacyInlineStyleOffenders = [
  'components/AppComponents.tsx',
  'components/BBChart.tsx',
  'components/ConfirmationModal.tsx',
  'components/Dashboard/UnsavedChangesModal.tsx',
  'components/Dashboard/charts/ConnectionCharts.tsx',
  'components/Dashboard/charts/DashboardMeshModelGraph.tsx',
  'components/Dashboard/charts/KubernetesConnectionChart.tsx',
  'components/Dashboard/charts/MesheryConfigurationCharts.tsx',
  'components/Dashboard/charts/WorkloadChart.tsx',
  'components/Dashboard/components.tsx',
  'components/Dashboard/debounceWidthProvider.tsx',
  'components/Dashboard/images/info-icon.tsx',
  'components/Dashboard/images/meshery-icon.tsx',
  'components/Dashboard/index.tsx',
  'components/Dashboard/overview.tsx',
  'components/Dashboard/resources/network/config.tsx',
  'components/Dashboard/resources/nodes/config.tsx',
  'components/Dashboard/resources/resources-table.tsx',
  'components/Dashboard/resources/security/config.tsx',
  'components/Dashboard/resources/sortable-table-cell.tsx',
  'components/Dashboard/tabpanel.tsx',
  'components/Dashboard/utils.tsx',
  'components/Dashboard/view-component.tsx',
  'components/Dashboard/view.tsx',
  'components/Dashboard/widgets/getting-started/data.tsx',
  'components/DataFormatter/index.tsx',
  'components/DatabaseSummary.tsx',
  'components/DesignLifeCycle/DeployStepper.tsx',
  'components/DesignLifeCycle/DeploymentSummary.tsx',
  'components/DesignLifeCycle/DryRun.tsx',
  'components/DesignLifeCycle/SelectDeploymentTarget.tsx',
  'components/DesignLifeCycle/ValidateDesign.tsx',
  'components/DesignLifeCycle/common.tsx',
  'components/DesignLifeCycle/finalizeDeployment.tsx',
  'components/DuplicatesDataTable.tsx',
  'components/ExportModal.tsx',
  'components/FlipCard.tsx',
  'components/General/ConnectClustersBtn.tsx',
  'components/General/CreateDesignBtn.tsx',
  'components/General/ErrorBoundary.tsx',
  'components/General/Modals/ConnectionModal.tsx',
  'components/General/Modals/Information/InfoModal.tsx',
  'components/General/Modals/Modal.tsx',
  'components/General/TipsCarousel.tsx',
  'components/General/error-404/index.tsx',
  'components/Header.tsx',
  'components/HeaderMenu.tsx',
  'components/Lifecycle/General/empty-state/index.tsx',
  'components/Lifecycle/Workspaces/WorkspaceActionList.tsx',
  'components/Lifecycle/Workspaces/WorkspaceDataTable.tsx',
  'components/Lifecycle/Workspaces/WorkspaceGridView.tsx',
  'components/Lifecycle/Workspaces/index.tsx',
  'components/LoadingComponents/Animations/AnimatedMeshSync.tsx',
  'components/LoadingComponents/LoadingComponent.tsx',
  'components/LoadingComponents/LoadingComponentServer.tsx',
  'components/MeshAdapterConfigComponent.tsx',
  'components/MesheryAdapterPlayComponent.tsx',
  'components/MesheryChart.tsx',
  'components/MesheryCredentialComponent.tsx',
  'components/MesheryFilters/CatalogFilter.tsx',
  'components/MesheryFilters/Filters.tsx',
  'components/MesheryFilters/FiltersCard.tsx',
  'components/MesheryFilters/FiltersGrid.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/Accordion.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/ArrayFieldTemlate.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/CustomBaseInput.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/CustomCheckboxWidget.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/CustomFileWidget.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/CustomSelectWidget.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/ObjectFieldTemplate.tsx',
  'components/MesheryMeshInterface/PatternService/RJSFCustomComponents/WrapIfAdditionalTemplate.tsx',
  'components/MesheryMeshInterface/PatternServiceForm.tsx',
  'components/MesheryPatterns/ActionButton.tsx',
  'components/MesheryPatterns/ActionPopover.tsx',
  'components/MesheryPatterns/CustomToolbarSelect.tsx',
  'components/MesheryPatterns/MesheryPatternCard.tsx',
  'components/MesheryPatterns/MesheryPatternGridView.tsx',
  'components/MesheryPatterns/MesheryPatterns.tsx',
  'components/MesheryPlayComponent.tsx',
  'components/MesheryProgressBar.tsx',
  'components/MesherySettingsEnvButtons.tsx',
  'components/Navigator.tsx',
  'components/NavigatorExtension.tsx',
  'components/NotificationCenter/formatters/common.tsx',
  'components/NotificationCenter/formatters/error.tsx',
  'components/NotificationCenter/formatters/meshsync_events.tsx',
  'components/NotificationCenter/formatters/model_registration.tsx',
  'components/NotificationCenter/formatters/relationship_evaluation.tsx',
  'components/NotificationCenter/index.tsx',
  'components/NotificationCenter/metadata.tsx',
  'components/NotificationCenter/notification.tsx',
  'components/Performance/Dashboard.tsx',
  'components/Performance/NodeDetails.tsx',
  'components/Performance/PerformanceCalendar.tsx',
  'components/Performance/PerformanceCard.tsx',
  'components/Performance/PerformanceProfileGrid.tsx',
  'components/Performance/PerformanceProfiles.tsx',
  'components/Performance/PerformanceResults.tsx',
  'components/Performance/assets/facebookIcon.tsx',
  'components/Performance/assets/linkedinIcon.tsx',
  'components/Performance/assets/twitterIcon.tsx',
  'components/Performance/index.tsx',
  'components/ReactSelectWrapper.tsx',
  'components/Registry/RegistryModal.tsx',
  'components/RelationshipBuilder/CreateRelationshipModal.tsx',
  'components/RelationshipBuilder/RelationshipFormStepper.tsx',
  'components/Settings/MesherySettings.tsx',
  'components/Settings/MesherySettingsPerformanceComponent.tsx',
  'components/Settings/Registry/ComponentTree.tsx',
  'components/Settings/Registry/CreateModelModal.tsx',
  'components/Settings/Registry/ImportModel.tsx',
  'components/Settings/Registry/ImportModelModal.tsx',
  'components/Settings/Registry/MeshModelComponent.tsx',
  'components/Settings/Registry/MeshModelDetails.tsx',
  'components/Settings/Registry/MesheryTreeView.tsx',
  'components/Settings/Registry/MesheryTreeViewItem.tsx',
  'components/Settings/Registry/MesheryTreeViewModel.tsx',
  'components/Settings/Registry/MesheryTreeViewRegistrants.tsx',
  'components/Settings/Registry/RelationshipTree.tsx',
  'components/Settings/Registry/Stepper/CSVStepper.tsx',
  'components/Settings/Registry/Stepper/UrlStepper.tsx',
  'components/Settings/Registry/StyledTreeItem.tsx',
  'components/SpacesSwitcher/DesignViewListItem.tsx',
  'components/SpacesSwitcher/MainDesignsContent.tsx',
  'components/SpacesSwitcher/MainViewsContent.tsx',
  'components/SpacesSwitcher/MenuComponent.tsx',
  'components/SpacesSwitcher/MobileViewSwitcher.tsx',
  'components/SpacesSwitcher/MyDesignsContent.tsx',
  'components/SpacesSwitcher/MyViewsContent.tsx',
  'components/SpacesSwitcher/RecentContent.tsx',
  'components/SpacesSwitcher/SharedContent.tsx',
  'components/SpacesSwitcher/SpaceSwitcher.tsx',
  'components/SpacesSwitcher/WorkspaceContent.tsx',
  'components/SpacesSwitcher/WorkspaceModal.tsx',
  'components/SpacesSwitcher/WorkspaceSwitcher.tsx',
  'components/SpacesSwitcher/components.tsx',
  'components/TroubleshootingComponent.tsx',
  'components/TypingFilter/index.tsx',
  'components/UserPreferences/index.tsx',
  'components/ViewInfoModal.tsx',
  'components/ViewSwitch.tsx',
  'components/configuratorComponents/MeshModel/LazyComponentForm.tsx',
  'components/configuratorComponents/MeshModel/index.tsx',
  'components/configuratorComponents/NameToIcon.tsx',
  'components/connections/ConnectionChip.tsx',
  'components/connections/ConnectionTable.tsx',
  'components/connections/common/index.tsx',
  'components/connections/index.tsx',
  'components/connections/meshSync/MeshSyncEmptyState.tsx',
  'components/connections/meshSync/RegisterConnectionModal.tsx',
  'components/connections/meshSync/Stepper/Notification.tsx',
  'components/connections/meshSync/Stepper/StepperContent.tsx',
  'components/connections/meshSync/Stepper/StepperContentWrapper.tsx',
  'components/connections/meshSync/index.tsx',
  'components/connections/metadata.tsx',
  'components/extensions/adapters/adapters.tsx',
  'components/layout.tsx',
  'components/multi-select-wrapper.tsx',
  'components/navigatorComponents.tsx',
  'components/shapes/Octagon.tsx',
  'components/telemetry/grafana/GrafanaComponent.tsx',
  'components/telemetry/prometheus/PrometheusSelectionComponent.tsx',
  'pages/_app.tsx',
  'pages/_document.tsx',
  'pages/extension/AccessMesheryModal.tsx',
  'pages/extensions.tsx',
  'utils/custom-search.tsx',
  'utils/utils.tsx',
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

      // Size budget for component files. 600 lines is the proactive warning
      // threshold; current files above 600 are allowlisted (legacyMaxLineSoft
      // Offenders for 600–1000, legacyMaxLineOffenders for >1000) so CI stays
      // green while the plan refactors them. The 1000-line hard ceiling from
      // the restructure plan is tracked separately by scripts/audit-size.js
      // (run via `npm run audit:size`) rather than by ESLint.
      'max-lines': ['warn', { max: 600, skipComments: true, skipBlankLines: true }],
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

  // ---------------------------------------------------------------------
  // Ban inline `style={{ ... }}` props in component code (§8.3).
  //
  // Styling belongs in styled() factories from @sistent/sistent. The inline
  // `style` prop is reserved for dynamic geometry (positions, sizes computed
  // at runtime) and should not be used for theme-derived values such as
  // colors, typography, or spacing tokens.
  //
  // Scoped to .tsx/.jsx component sources. The same dirs the hex-literal
  // guardrail ignores are ignored here (theme/themes/assets/lib/public are
  // not component code, and tests/scripts/the config itself are tooling).
  // ---------------------------------------------------------------------
  {
    files: ['**/*.{tsx,jsx}'],
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
      'react/forbid-dom-props': [
        'warn',
        {
          forbid: [
            {
              propName: 'style',
              message:
                'Use styled() from @sistent/sistent; inline style is reserved for dynamic geometry.',
            },
          ],
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
    files: legacyInlineStyleOffenders,
    rules: {
      'react/forbid-dom-props': 'off',
    },
  },
  // Hard 1000-line ceiling files: disable the 600-line warning entirely so
  // they are not double-reported. They are tracked separately in the giant-
  // files audit and will be refactored in phase 5.
  {
    files: legacyMaxLineOffenders,
    rules: {
      'max-lines': 'off',
    },
  },
  // 600–1000 line "soft" offenders: silence the 600-line warning for these
  // existing files only. New files crossing 600 lines will still warn.
  {
    files: legacyMaxLineSoftOffenders,
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
