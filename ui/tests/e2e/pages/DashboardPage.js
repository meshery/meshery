import { expect } from '@playwright/test';
const LEFT_NAV = {
  DASHBOARD: {
    name: 'Dashboard',
    path: '/',
    NAV_ITEMS: {},
  },
  LIFECYCLE: {
    name: 'lifecycle',
    NAV_ITEMS: {
      CONNECTIONS: 'connection',
      ENVIRONMENT: 'environment',
      WORKSPACE: 'Workspace',
      ADAPTERS: 'Adapters',
    },
  },
  CONFIGURATION: {
    name: 'configuration',
    NAV_ITEMS: {
      DESIGNS: 'design',
    },
  },
  TELEMETRY: {
    name: 'telemetry',
    NAV_ITEMS: {
      CHARTS: 'Grafana',
      METRICS: 'Prometheus',
    },
  },
  PERFORMANCE: {
    name: 'performance',
    NAV_ITEMS: {
      PROFILES: 'profiles',
    },
  },
  EXTENSIONS: {
    name: 'extensions',
    NAV_ITEMS: {},
  },
};

const HEADER_NAV = {
  name: 'header-nav',
  NAV_ITEMS: {
    SETTINGS: 'nav-item-settings',
    LOGOUT: 'nav-item-logout',
    PREFERENCES: 'nav-item-preferences',
    NOTIFICATIONS: 'nav-item-notifications',
  },
};

export class DashboardPage {
  constructor(page) {
    this.page = page;
    this.navigationPanel = this.page.getByTestId('navigation');
    this.notificationButton = this.page.getByTestId('notification-button');
    this.profileButton = this.page.getByTestId('profile-button');
    this.headerMenu = this.page.getByTestId('header-menu');
  }

  async navigateToMenu(navItem, options = {}) {
    const menuItem = this.page.getByTestId(navItem);
    await expect(menuItem).toBeVisible(options.timeout ? { timeout: options.timeout } : undefined);
    await menuItem.click();
  }

  // Expanding a left-nav section and clicking one of its children is not a
  // single settled interaction: the parent is a link that both navigates and
  // toggles the section, and a click that lands before the navigator has
  // hydrated is swallowed with no error. The submenu then never opens, so
  // waiting on the child is waiting on an element that will never appear -
  // which is why this timed out after resolving the child ~90 times as
  // "hidden" rather than failing fast (5 of the 20 CI runs to 2026-08-05).
  //
  // Retrying the parent click until the child is actually visible fixes the
  // race at its cause. A longer timeout cannot, because no amount of waiting
  // re-sends the swallowed click.
  async navigateToSubMenuItem(parentItem, childItem) {
    const submenuItem = this.page.getByTestId(childItem);
    // Bound BOTH waits inside the retry. With the parent check left on the
    // default 60s expect timeout, a single attempt could consume most of the
    // test budget and the case died on the 180s test timeout instead of its own
    // assertion - and a test-level timeout reports no locator, no snippet and no
    // file location, so the failure said nothing at all. Short per-attempt
    // timeouts make each retry cheap and let this fail legibly at ~60s.
    await expect(async () => {
      // Only click the parent when the submenu is not already open. The parent
      // both navigates and toggles, so an unconditional click on a retry can
      // CLOSE a submenu a previous attempt had just opened - the retry would
      // then fight itself, which is a plausible mechanism for this staying
      // flaky after the click was made retryable at all.
      if (!(await submenuItem.isVisible())) {
        await this.navigateToMenu(parentItem, { timeout: 5_000 });
      }
      await expect(submenuItem).toBeVisible({ timeout: 5_000 });
    }).toPass({ timeout: 60_000 });
    await submenuItem.click();
  }

  async navigateToDashboard() {
    await this.page.goto(LEFT_NAV.DASHBOARD.path, { waitUntil: 'domcontentloaded' });
    await expect(this.navigationPanel).toBeVisible();
    await expect(this.headerMenu).toBeVisible();
  }

  async navigateToPerformance() {
    await this.navigateToMenu(LEFT_NAV.PERFORMANCE.name);
  }

  async navigateToExtensions() {
    await this.navigateToMenu(LEFT_NAV.EXTENSIONS.name);
  }

  async navigateToTelemetry() {
    await this.navigateToMenu(LEFT_NAV.TELEMETRY.name);
  }

  async navigateToTelemetryCharts() {
    await this.navigateToSubMenuItem(LEFT_NAV.TELEMETRY.name, LEFT_NAV.TELEMETRY.NAV_ITEMS.CHARTS);
  }

  async navigateToTelemetryMetrics() {
    await this.navigateToSubMenuItem(LEFT_NAV.TELEMETRY.name, LEFT_NAV.TELEMETRY.NAV_ITEMS.METRICS);
  }

  async navigateToLifecycle() {
    await this.navigateToMenu(LEFT_NAV.LIFECYCLE.name);
  }

  async navigateToConfiguration() {
    await this.navigateToMenu(LEFT_NAV.CONFIGURATION.name);
  }

  async navigateToConnections() {
    await this.navigateToSubMenuItem(
      LEFT_NAV.LIFECYCLE.name,
      LEFT_NAV.LIFECYCLE.NAV_ITEMS.CONNECTIONS,
    );
  }

  async navigateToEnvironment() {
    await this.navigateToSubMenuItem(
      LEFT_NAV.LIFECYCLE.name,
      LEFT_NAV.LIFECYCLE.NAV_ITEMS.ENVIRONMENT,
    );
  }

  async navigateToWorkspace() {
    await this.navigateToSubMenuItem(
      LEFT_NAV.LIFECYCLE.name,
      LEFT_NAV.LIFECYCLE.NAV_ITEMS.WORKSPACE,
    );
  }

  async navigateToAdapters() {
    await this.navigateToSubMenuItem(
      LEFT_NAV.LIFECYCLE.name,
      LEFT_NAV.LIFECYCLE.NAV_ITEMS.ADAPTERS,
    );
  }

  async navigateToProfiles() {
    await this.navigateToSubMenuItem(
      LEFT_NAV.PERFORMANCE.name,
      LEFT_NAV.PERFORMANCE.NAV_ITEMS.PROFILES,
    );
  }

  async navigateToHeaderItem(navItem) {
    await expect(this.headerMenu).toBeVisible();
    await this.headerMenu.click();
    const headerItem = this.page.getByTestId(navItem);
    await expect(headerItem).toBeVisible();
    await headerItem.click();
  }

  async navigateToSettings() {
    await this.navigateToHeaderItem(HEADER_NAV.NAV_ITEMS.SETTINGS);
  }

  async navigateToPreferences() {
    await this.navigateToHeaderItem(HEADER_NAV.NAV_ITEMS.PREFERENCES);
  }

  async navigateToDesigns() {
    await this.navigateToSubMenuItem(
      LEFT_NAV.CONFIGURATION.name,
      LEFT_NAV.CONFIGURATION.NAV_ITEMS.DESIGNS,
    );
  }

  async navigateToLogout() {
    await this.navigateToHeaderItem(HEADER_NAV.NAV_ITEMS.LOGOUT);
  }
}
