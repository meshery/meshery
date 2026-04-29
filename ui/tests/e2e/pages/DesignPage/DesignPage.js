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
    this.deleteDesignBtn = page.getByTestId('pattern-btn-delete');

    this.deleteDesignModalHeader = page.getByTestId('modal-header');
    this.deleteConfirmationBtn = page.getByTestId('prompt-primary-button');

    this.universalFilter = page.getByTestId('meshery-patterns-universal-filter');
    this.visibilityFilter = page.getByTestId('meshery-patterns-universal-filter-select-visibility');
    this.filterApplyBtn = page.getByTestId('meshery-patterns-universal-filter-apply-btn');
  }

  async navigateTo() {
    await this.DashboardPage.navigateToDashboard();
    await this.DashboardPage.navigateToDesigns();
    await this.pageLoader.waitFor({ state: 'detached' });
  }

  async navigateToDesignConfigurator() {
    await this.applyVisibilityFilter('published');
    const firstCard = await this.getFirstCardByVisibilityBadge('published');
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
        validate: this.page.getByTestId('action-btn-option-Validate'),
        dryRun: this.page.getByTestId('action-btn-option-Dry Run'),
        deploy: this.page.getByTestId('action-btn-option-Deploy'),
        unDeploy: this.page.getByTestId('action-btn-option-Undeploy'),
        clone: card.getByTestId('pattern-btn-clone'),
        download: card.getByTestId('pattern-btn-download'),
        edit: card.getByTestId('pattern-btn-edit'),
        info: card.getByTestId('pattern-btn-info'),
      },
    };
  }

  async applyVisibilityFilter(visibility) {
    await this.universalFilter.click();
    await this.visibilityFilter.first().click();
    await this.page.getByTestId(`meshery-patterns-universal-filter-option-${visibility}`).click();

    await this.filterApplyBtn.click();
  }

  getFirstCardByVisibilityBadge(visibility) {
    return this.designCards
      .filter({
        has: this.page.getByTestId(`design-visibility-${visibility}`),
      })
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
