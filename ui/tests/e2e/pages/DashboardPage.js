import { expect } from '@playwright/test';

export const NAVIGATION_TIMEOUT = 60000;
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

  async navigateToMenu(navItem) {
    const menuItem = this.page.getByTestId(navItem);
    await expect(menuItem).toBeVisible({ timeout: NAVIGATION_TIMEOUT });
    await menuItem.click();
  }

  async navigateToSubMenuItem(parentItem, childItem) {
    await this.navigateToMenu(parentItem);
    const submenuItem = this.page.getByTestId(childItem);
    await expect(submenuItem).toBeVisible({ timeout: NAVIGATION_TIMEOUT });
    await submenuItem.click();
  }

  async navigateToDashboard() {
    await this.page.goto(LEFT_NAV.DASHBOARD.path, { waitUntil: 'domcontentloaded' });
    await expect(this.navigationPanel).toBeVisible({ timeout: NAVIGATION_TIMEOUT });
    await expect(this.headerMenu).toBeVisible({ timeout: NAVIGATION_TIMEOUT });
  }

  async navigateToPerformance() {
    await this.navigateToMenu(LEFT_NAV.PERFORMANCE.name);
  }

  async navigateToExtensions() {
    await this.navigateToMenu(LEFT_NAV.EXTENSIONS.name);
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
    await expect(this.headerMenu).toBeVisible({ timeout: NAVIGATION_TIMEOUT });
    await this.headerMenu.click();
    const headerItem = this.page.getByTestId(navItem);
    await expect(headerItem).toBeVisible({ timeout: NAVIGATION_TIMEOUT });
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
