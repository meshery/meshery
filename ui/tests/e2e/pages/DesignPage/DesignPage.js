import { ImportModal } from './components/ImportModal';
import { DashboardPage } from '../DashboardPage';
import { DeployModal } from './components/DeployModal';

export class DesignPage {
  constructor(page) {
    this.page = page;
    this.DashboardPage = new DashboardPage(page);
    this.DeployModal = new DeployModal(page);
    this.ImportModal = new ImportModal(page);

    this.createDesignBtn = page.getByTestId('meshery-patterns-create-design-btn');
    this.importDesignBtn = this.page.getByTestId('meshery-patterns-import-design-btn');
    this.searchBar = this.page.getByTestId('meshery-patterns-search-bar');
    this.designGrid = this.page.getByTestId('meshery-patterns-grid');
    this.designCards = this.page.getByTestId('meshery-pattern-card-item');
    this.deleteDesignBtn = page.getByTestId('pattern-btn-delete');

    this.deleteDesignModalHeader = page.getByTestId('modal-header');
    this.deleteConfirmationBtn = page.getByTestId('prompt-primary-button');
  }

  async navigateTo() {
    await this.DashboardPage.navigateToDashboard();
    await this.DashboardPage.navigateToDesigns();
  }

  async navigateToDesignConfigurator() {
    const card = this.getPublishedDesign();
    await card.actionElements.edit.click();
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
        download: card.getByTestId('pattern-btn-download'),
        edit: card.getByTestId('pattern-btn-edit'),
        info: card.getByTestId('pattern-btn-info'),
      },
    };
  }

  getPublishedDesign() {
    const card = this.designCards
      .filter({
        has: this.page.getByTestId('design-visibility-published'),
      })
      .first();

    const base = this.getBaseCardElements(card);

    return {
      ...base,
      actionElements: {
        ...base.actionElements,
        clone: card.getByTestId('pattern-btn-clone'),
        unpublish: card.getByTestId('pattern-btn-unpublish'),
      },
    };
  }
}
