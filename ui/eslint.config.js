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
  'components/general/ReactSelectWrapper.tsx',
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
  'components/SpacesSwitcher/components.tsx',
  'components/general/ViewSwitch.tsx',
  'components/general/YamlDialog.tsx',
  'components/configuratorComponents/MeshModel/index.tsx',
  'components/configuratorComponents/NameToIcon.tsx',
  'components/connections/ConnectionChip.tsx',
  'components/icons/index.ts',
  'components/telemetry/grafana/GrafanaCustomChart.tsx',
  'components/telemetry/grafana/GrafanaDateRangePicker.tsx',
  'components/telemetry/prometheus/PrometheusSelectionComponent.tsx',
  'pages/_app.tsx',
];

const legacyLiteralColorOffenders = [
  'components/dashboard/charts/ResourceUtilizationChart.tsx',
  'components/dashboard/components.tsx',
  'components/dashboard/images/info-icon.tsx',
  'components/dashboard/images/meshery-icon.tsx',
  'components/dashboard/style.ts',
  'components/designs/lifecycle/DryRun.tsx',
  'components/designs/lifecycle/ValidateDesign.tsx',
  'components/designs/lifecycle/common.tsx',
  'components/general/TipsCarousel.tsx',
  'components/general/error-404/CurrentSession.tsx',
  'components/general/error-404/socials/styles.tsx',
  'components/general/error-404/styles.tsx',
  'components/layout/Header/Header.styles.tsx',
  'components/layout/Header/Header.tsx',
  'components/environments/environment-card.tsx',
  'components/environments/index.tsx',
  'components/environments/styles.tsx',
  'components/lifecycle/general/empty-state/curvedArrowIcon.tsx',
  'components/lifecycle/general/empty-state/index.tsx',
  'components/lifecycle/general/flip-card/index.tsx',
  'components/workspaces/index.tsx',
  'components/shared/LoadingState/Animations/AnimatedFilter.tsx',
  'components/shared/LoadingState/Animations/AnimatedLightMeshery.tsx',
  'components/shared/LoadingState/Animations/AnimatedMeshPattern.tsx',
  'components/shared/LoadingState/Animations/AnimatedMeshery.tsx',
  'components/shared/LoadingState/Animations/AnimatedMesheryCSS.tsx',
  'components/shared/LoadingState/LoadingComponentServer.tsx',
  'components/MeshAdapterConfigComponent.tsx',
  'components/adapter-play-styled.tsx',
  'components/MesheryChart.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/ArrayFieldTemlate.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/CustomBaseInput.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/ObjectFieldTemplate.tsx',
  'components/meshery-mesh-interface/PatternServiceForm.tsx',
  'components/designs/patterns/MesheryPatternCard.tsx',
  'components/designs/patterns/MesheryPatternGridView.tsx',
  'components/designs/patterns/MesheryPatterns.columns.tsx',
  'components/designs/patterns/MesheryPatterns.tsx',
  'components/designs/patterns/design-lifecycle-handlers.tsx',
  'components/general/FlipCard.styles.tsx',
  'components/general/YamlDialog.styles.tsx',
  'components/layout/Navigator/NavigatorExtension.tsx',
  'components/layout/NotificationCenter/constants.tsx',
  'components/layout/NotificationCenter/formatters/relationship_evaluation.tsx',
  'components/layout/NotificationCenter/index.tsx',
  'components/layout/NotificationCenter/notificationCenter.style.tsx',
  'components/performance/PerformanceCard.tsx',
  'components/performance/PerformanceForm.tsx',
  'components/performance/PerformanceResults.tsx',
  'components/performance/assets/facebookIcon.tsx',
  'components/performance/assets/linkedinIcon.tsx',
  'components/performance/assets/twitterIcon.tsx',
  'components/performance/index.tsx',
  'components/performance/style.tsx',
  'components/settings/MesherySettings.tsx',
  'components/registry/MeshModelDetails.tsx',
  'components/workspaces/SpacesSwitcher/WorkspaceSwitcher.tsx',
  'components/workspaces/SpacesSwitcher/styles.tsx',
  'components/shared/FormFields/typing-filter/style.tsx',
  'components/user-preferences/index.tsx',
  'components/user-preferences/style.tsx',
  'components/designs/configurator/CustomBreadCrumb.tsx',
  'components/designs/configurator/MeshModel/styledComponents/AppBar.tsx',
  'components/designs/configurator/MeshModel/utils.tsx',
  'components/designs/configurator/NameToIcon.tsx',
  'components/connections/meshSync/Stepper/Notification.tsx',
  'components/connections/meshSync/Stepper/StepperContent.tsx',
  'components/connections/meshSync/Stepper/StepperContentWrapper.tsx',
  'components/connections/meshSync/Stepper/index.tsx',
  'components/filters/Filters.tsx',
  'components/filters/Filters.columns.tsx',
  'components/filters/FiltersGrid.tsx',
  'components/filters/ImportModal.tsx',
  'components/filters/PublishModal.tsx',
  'components/load-test-timer-dialog.tsx',
  'components/telemetry/grafana/GrafanaCustomGaugeChart.tsx',
  'components/telemetry/grafana/GrafanaDateRangePicker.tsx',
  'css/icons.styles.ts',
  'utils/custom-search.tsx',
];

const legacyMaxLineOffenders = [
  'components/designs/patterns/MesheryPatterns.tsx',
  'components/connections/ConnectionTable.tsx',
];

// Files currently in the 600–1000 line "soft" range. They exceed the 600-line
// proactive warning threshold (§8.4) but stay under the hard 1000-line ceiling.
// Allowlisted so CI stays green; entries leave the list as files get split up.
const legacyMaxLineSoftOffenders = [
  'components/environments/index.tsx',
  'components/layout/Navigator/Navigator.tsx',
  'components/performance/PerformanceResults.tsx',
  'components/registry/Stepper/UrlStepper.tsx',
  'components/user-preferences/index.tsx',
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
  'assets/icons/shapes/Octagon.tsx',
  'components/layout/AppShell/AppComponents.tsx',
  'components/general/BBChart.tsx',
  'components/connections/ConnectionTable.columns.tsx',
  'components/dashboard/charts/ConnectionCharts.tsx',
  'components/dashboard/charts/DashboardMeshModelGraph.tsx',
  'components/dashboard/charts/KubernetesConnectionChart.tsx',
  'components/dashboard/charts/MesheryConfigurationCharts.tsx',
  'components/dashboard/charts/WorkloadChart.tsx',
  'components/dashboard/components.tsx',
  'components/dashboard/debounceWidthProvider.tsx',
  'components/dashboard/images/info-icon.tsx',
  'components/dashboard/images/meshery-icon.tsx',
  'components/dashboard/index.tsx',
  'components/dashboard/overview.tsx',
  'components/dashboard/resources/network/service-columns.tsx',
  'components/dashboard/resources/nodes/config.tsx',
  'components/dashboard/resources/resources-table.tsx',
  'components/dashboard/resources/security/config.tsx',
  'components/dashboard/resources/sortable-table-cell.tsx',
  'components/dashboard/tabpanel.tsx',
  'components/dashboard/utils.tsx',
  'components/dashboard/view-component.tsx',
  'components/dashboard/view.tsx',
  'components/dashboard/widgets/getting-started/data.tsx',
  'components/data-formatter/index.tsx',
  'components/DatabaseSummary.tsx',
  'components/designs/lifecycle/DeployConfirmationModal.tsx',
  'components/designs/lifecycle/DeployStepper.tsx',
  'components/designs/lifecycle/DeploymentSummary.tsx',
  'components/designs/lifecycle/DryRun.tsx',
  'components/designs/lifecycle/SelectDeploymentTarget.tsx',
  'components/designs/lifecycle/ValidateDesign.tsx',
  'components/designs/lifecycle/common.tsx',
  'components/designs/lifecycle/finalizeDeployment.tsx',
  'components/general/FlipCard.tsx',
  'components/general/ConnectClustersBtn.tsx',
  'components/general/CreateDesignBtn.tsx',
  'components/shared/ErrorBoundary/ErrorBoundary.tsx',
  'components/shared/Modal/Information/LegacyInfoModal.tsx',
  'components/shared/Modal/LegacyRJSFModal.tsx',
  'components/general/TipsCarousel.tsx',
  'components/general/error-404/index.tsx',
  'components/layout/Header/Header.tsx',
  'components/layout/Header/HeaderMenu.tsx',
  'components/lifecycle/general/empty-state/index.tsx',
  'components/workspaces/WorkspaceActionList.tsx',
  'components/workspaces/WorkspaceDataTable.tsx',
  'components/workspaces/WorkspaceGridView.tsx',
  'components/workspaces/index.tsx',
  'components/shared/LoadingState/Animations/AnimatedMeshSync.tsx',
  'components/shared/LoadingState/LoadingComponent.tsx',
  'components/shared/LoadingState/LoadingComponentServer.tsx',
  'components/MeshAdapterConfigComponent.tsx',
  'components/MesheryAdapterPlayComponent.tsx',
  'components/adapter-play-addon-switches.tsx',
  'components/adapter-play-category-card.tsx',
  'components/MesheryChart.tsx',
  'components/MesheryCredentialComponent.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/Accordion.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/ArrayFieldTemlate.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/CustomBaseInput.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/CustomCheckboxWidget.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/CustomFileWidget.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/CustomSelectWidget.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/ObjectFieldTemplate.tsx',
  'components/meshery-mesh-interface/PatternService/RJSFCustomComponents/WrapIfAdditionalTemplate.tsx',
  'components/meshery-mesh-interface/PatternServiceForm.tsx',
  'components/designs/patterns/ActionButton.tsx',
  'components/designs/patterns/ActionPopover.tsx',
  'components/designs/patterns/CustomToolbarSelect.tsx',
  'components/designs/patterns/MesheryPatternCard.tsx',
  'components/designs/patterns/MesheryPatternGridView.tsx',
  'components/designs/patterns/MesheryPatterns.tsx',
  'components/designs/patterns/MesheryPatternsToolbar.tsx',
  'components/MesheryPlayComponent.tsx',
  'components/general/MesheryProgressBar.tsx',
  'components/MesherySettingsEnvButtons.tsx',
  'components/layout/Navigator/Navigator.tsx',
  'components/layout/Navigator/NavigatorExtension.tsx',
  'components/layout/NotificationCenter/formatters/common.tsx',
  'components/layout/NotificationCenter/formatters/error.tsx',
  'components/layout/NotificationCenter/formatters/meshsync_events.tsx',
  'components/layout/NotificationCenter/formatters/model_registration.tsx',
  'components/layout/NotificationCenter/formatters/relationship_evaluation.tsx',
  'components/layout/NotificationCenter/index.tsx',
  'components/layout/NotificationCenter/metadata.tsx',
  'components/layout/NotificationCenter/notification.tsx',
  'components/performance/Dashboard.tsx',
  'components/performance/NodeDetails.tsx',
  'components/performance/PerformanceCalendar.tsx',
  'components/performance/PerformanceCard.tsx',
  'components/performance/PerformanceForm.tsx',
  'components/performance/PerformanceFormActions.tsx',
  'components/performance/PerformanceProfileGrid.tsx',
  'components/performance/PerformanceProfiles.tsx',
  'components/performance/PerformanceResults.tsx',
  'components/performance/PerformanceTestResults.tsx',
  'components/performance/assets/facebookIcon.tsx',
  'components/performance/assets/linkedinIcon.tsx',
  'components/performance/assets/twitterIcon.tsx',
  'components/performance/index.tsx',
  'components/performance/performance-helpers.tsx',
  'components/general/ReactSelectWrapper.tsx',
  'components/relationship-builder/CreateRelationshipModal.tsx',
  'components/relationship-builder/RelationshipFormStepper.tsx',
  'components/settings/MesherySettings.tsx',
  'components/settings/MesherySettingsPerformanceComponent.tsx',
  'components/registry/ComponentTree.tsx',
  'components/registry/CreateModelModal.tsx',
  'components/registry/ImportModel.tsx',
  'components/registry/ImportModelModal.tsx',
  'components/registry/MeshModelComponent.tsx',
  'components/registry/MeshModelDetails.tsx',
  'components/registry/MesheryTreeView.tsx',
  'components/registry/MesheryTreeViewItem.tsx',
  'components/registry/MesheryTreeViewModel.tsx',
  'components/registry/MesheryTreeViewRegistrants.tsx',
  'components/registry/RegistryModal.tsx',
  'components/registry/RelationshipTree.tsx',
  'components/registry/Stepper/CSVStepper.tsx',
  'components/registry/Stepper/UrlStepper.tsx',
  'components/registry/StyledTreeItem.tsx',
  'components/workspaces/SpacesSwitcher/DesignViewListItem.tsx',
  'components/workspaces/SpacesSwitcher/MainDesignsContent.tsx',
  'components/workspaces/SpacesSwitcher/MainViewsContent.tsx',
  'components/workspaces/SpacesSwitcher/MenuComponent.tsx',
  'components/workspaces/SpacesSwitcher/MobileViewSwitcher.tsx',
  'components/workspaces/SpacesSwitcher/MyDesignsContent.tsx',
  'components/workspaces/SpacesSwitcher/MyViewsContent.tsx',
  'components/workspaces/SpacesSwitcher/RecentContent.tsx',
  'components/workspaces/SpacesSwitcher/SharedContent.tsx',
  'components/workspaces/SpacesSwitcher/SpaceSwitcher.tsx',
  'components/workspaces/SpacesSwitcher/WorkspaceContent.tsx',
  'components/workspaces/SpacesSwitcher/WorkspaceSwitcher.tsx',
  'components/workspaces/SpacesSwitcher/components.tsx',
  'components/TroubleshootingComponent.tsx',
  'components/shared/FormFields/typing-filter/index.tsx',
  'components/user-preferences/index.tsx',
  'components/general/ViewSwitch.tsx',
  'components/designs/configurator/MeshModel/LazyComponentForm.tsx',
  'components/designs/configurator/MeshModel/index.tsx',
  'components/designs/configurator/NameToIcon.tsx',
  'components/connections/ConnectionChip.tsx',
  'components/connections/ConnectionTable.tsx',
  'components/connections/common/index.tsx',
  'components/connections/index.tsx',
  'components/connections/meshSync/MeshSyncEmptyState.tsx',
  'components/connections/meshSync/Stepper/Notification.tsx',
  'components/connections/meshSync/Stepper/StepperContent.tsx',
  'components/connections/meshSync/Stepper/StepperContentWrapper.tsx',
  'components/connections/meshSync/index.tsx',
  'components/connections/metadata.tsx',
  'components/extensions/adapters/adapters.tsx',
  'components/filters/CatalogFilter.tsx',
  'components/filters/Filters.tsx',
  'components/filters/FiltersCard.tsx',
  'components/filters/FiltersGrid.tsx',
  'components/filters/ImportModal.tsx',
  'components/filters/PublishModal.tsx',
  'components/filters/YAMLEditor.tsx',
  'components/layout/AppShell/layout.tsx',
  'components/general/multi-select-wrapper.tsx',
  'components/layout/Navigator/navigatorComponents.tsx',
  'components/telemetry/grafana/GrafanaComponent.tsx',
  'components/telemetry/prometheus/PrometheusSelectionComponent.tsx',
  'pages/_app.tsx',
  'pages/_document.tsx',
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

      // Ban Material UI imports and deep imports of the local theme entry
      // point. @sistent/sistent is the only UI kit; @/theme is the approved
      // Phase 1 theme entry point.
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
  //   - ui/assets/**      (SVG icons encoded as React components)
  //   - ui/lib/**         (third-party integration helpers)
  //   - ui/public/**      (static assets)
  // ---------------------------------------------------------------------
  {
    files: ['**/*.{ts,tsx,js,jsx}'],
    ignores: [
      'theme/**',
      'assets/**',
      'lib/**',
      'public/**',
      'tests/**',
      'scripts/**',
      'eslint.config.js',
      '**/*.test.{ts,tsx,js,jsx}',
      '**/__tests__/**',
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
  // guardrail ignores are ignored here (theme/assets/lib/public are
  // not component code, and tests/scripts/the config itself are tooling).
  // ---------------------------------------------------------------------
  {
    files: ['**/*.{tsx,jsx}'],
    ignores: [
      'theme/**',
      'assets/**',
      'lib/**',
      'public/**',
      'tests/**',
      'scripts/**',
      'eslint.config.js',
      '**/*.test.{ts,tsx,js,jsx}',
      '**/__tests__/**',
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

  // Test files: relax rules that get in the way of mocks and ad-hoc test
  // doubles. Inline mock components rarely warrant a displayName, JSX in
  // .ts test files for hooks/RTK-Query providers, etc., are all expected
  // in a testing context.
  {
    files: ['**/*.test.{ts,tsx,js,jsx}', '**/__tests__/**/*.{ts,tsx,js,jsx}'],
    rules: {
      'react/display-name': 'off',
      'react/no-children-prop': 'off',
      'react/forbid-dom-props': 'off',
      'react-hooks/rules-of-hooks': 'off',
      'no-undef': 'off',
      'no-restricted-syntax': 'off',
    },
  },
];
