import { ImportModal } from './components/ImportModal';
import { DashboardPage } from '../DashboardPage';
import { DeployModal } from './components/DeployModal';

export class DesignPage {
  constructor(page) {
    this.page = page;
    this.DashboardPage = new DashboardPage(page);
    this.DeployModal = new DeployModal(page);
    this.ImportModal = new ImportModal(page);

    this.pageLoader = page.getByTestId('page-loader');
    this.createDesignBtn = page.getByTestId('meshery-patterns-create-design-btn');
    this.importDesignBtn = page.getByTestId('meshery-patterns-import-design-btn');
    this.searchBar = page.getByTestId('meshery-patterns-search-bar');
    this.searchBtn = page.getByTestId('meshery-patterns-search-btn');
    this.designGrid = page.getByTestId('meshery-patterns-grid');
    this.designCards = page.getByTestId('meshery-pattern-card-item');
    this.deleteDesignBtn = page
      .locator('[data-testid="pattern-btn-delete"], [data-testid="pattern-btn-action-delete"]')
      .first();

    this.deleteDesignModalHeader = page
      .locator(
        '[data-testid="modal-header"], [role="dialog"] h2, [role="dialog"] h6, [role="dialog"] [class*="title"]',
      )
      .first();
    this.deleteConfirmationBtn = page
      .locator('[data-testid="prompt-primary-button"], button:has-text("DELETE")')
      .first();

    this.universalFilter = page.getByTestId('meshery-patterns-universal-filter');
    this.visibilityFilter = page.getByTestId('meshery-patterns-universal-filter-select-visibility');
    this.filterApplyBtn = page.getByTestId('meshery-patterns-universal-filter-apply-btn');
  }

  async navigateTo() {
    await this.DashboardPage.navigateToDashboard();
    await this.DashboardPage.navigateToDesigns();
    await this.pageLoader.waitFor({ state: 'detached', timeout: 30000 }).catch(() => {});
  }

  async navigateToDesignConfigurator() {
    await this.applyVisibilityFilter('published');
    const firstCard = this.getFirstCardByVisibilityBadge('published');
    const cardElements = this.getCardElements(firstCard, 'published');
    await cardElements.actionElements.edit.click();
  }

  async clickImportDesignButton() {
    await this.importDesignBtn.click();
  }

  getBaseCardElements(card) {
    return {
      display: {
        card: card,
        name: card.getByTestId('pattern-card-name'),
        visibility: card.getByTestId('visibility-chip-menu'),
        modifiedDate: card.getByTestId('pattern-card-modified-on'),
        actionsBar: card.getByTestId('pattern-card-actions'),
      },
      actionToggleBtn: card.getByTestId('action-btn-toggle'),
      actionElements: {
        validate: this.page
          .locator(
            '[data-testid="pattern-btn-validate"], [data-testid="action-btn-option-Validate"]',
          )
          .first(),
        dryRun: this.page
          .locator('[data-testid="pattern-btn-dryrun"], [data-testid="action-btn-option-Dry Run"]')
          .first(),
        deploy: this.page
          .locator('[data-testid="pattern-btn-deploy"], [data-testid="action-btn-option-Deploy"]')
          .first(),
        unDeploy: this.page
          .locator(
            '[data-testid="pattern-btn-undeploy"], [data-testid="action-btn-option-Undeploy"]',
          )
          .first(),
        clone: card.getByTestId('pattern-btn-clone'),
        download: card.getByTestId('pattern-btn-download'),
        edit: card.getByTestId('pattern-btn-edit'),
        info: card.getByTestId('pattern-btn-info'),
      },
    };
  }

  async applyVisibilityFilter(visibility) {
    if (await this.universalFilter.isVisible()) {
      await this.universalFilter.click();
      if (await this.visibilityFilter.first().isVisible()) {
        await this.visibilityFilter.first().click();
        const option = this.page
          .getByTestId(`meshery-patterns-universal-filter-option-${visibility}`)
          .or(this.page.getByRole('option', { name: new RegExp(visibility, 'i') }));
        if (await option.isVisible()) {
          await option.click();
        }
      }
      if (await this.filterApplyBtn.isVisible()) {
        await this.filterApplyBtn.click();
      }
    }
  }

  getFirstCardByVisibilityBadge(visibility) {
    return this.designCards
      .filter({
        has: this.page.locator(
          `[data-testid="design-visibility-${visibility}"], [data-testid="visibility-chip-menu"]:has-text("${visibility}")`,
        ),
      })
      .or(this.designCards)
      .first();
  }

  getCardElements(card, visibility) {
    const baseElements = this.getBaseCardElements(card);

    const extraActions = {};
    const actionElements = this._getCardActionsMap()[visibility] || [];
    actionElements.forEach((element) => {
      extraActions[element] = card.getByTestId(`pattern-btn-${element}`);
    });

    return {
      ...baseElements,
      display: {
        ...baseElements.display,
      },
      actionElements: {
        ...baseElements.actionElements,
        ...extraActions,
      },
    };
  }

  _getCardActionsMap() {
    return {
      public: [],
      published: ['unpublish'],
    };
  }
}
