import { expect } from '@playwright/test';
import { ENV } from '../env';

export class EnvironmentPage {
  /**
   * @param {import('@playwright/test').Page} page
   */
  constructor(page) {
    this.page = page;
  }
  async navigate() {
    await this.page.goto(`${ENV.MESHERY_SERVER_URL}/management/environments`);
  }
  async fillInput(locator, value) {
    await locator.fill(value);
  }

  async goToProfiles() {
    await this.page.goto(`${ENV.MESHERY_SERVER_URL}/management/environments`);
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/environments`);
  }

  async goToConnections() {
    await this.page.goto(`${ENV.MESHERY_SERVER_URL}/management/connections`);
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/connections`);
  }

  async doesProfileExist(environmentName) {
    await this.goToProfiles();
    return this.page
      .locator(`text=${environmentName}`)
      .count()
      .then((count) => count === 1);
  }
  async createEnvironmentProfile(environmentName) {
    if (!(await this.doesProfileExist(environmentName))) {
      await this.page.getByRole('button', { name: 'Create' }).click();
      await this.fillInput(this.page.locator('#Environment_name'), `${environmentName}`);
      await this.fillInput(
        this.page.locator('#Environment_description'),
        'Created profile for playwright tests',
      );
      await this.page.getByRole('button', { name: 'Save', exact: true }).click();
      // Check for notification visibility
      const notification = await this.page.locator(`text=Environment "${environmentName}" created`);
      await expect(notification).toBeVisible();
      await this.page.waitForTimeout(8000);
    } else {
      console.log(`Profile "${environmentName}" already exists.`);
    }
  }

  async assignConnections(environmentName) {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/environments`);
    await this.page.waitForSelector(`text=${environmentName}`, { state: 'visible' });
    await this.page.evaluate(() => {
      const button = [...document.querySelectorAll('button')].find((b) =>
        b.textContent.includes('Assigned Connections'),
      );
      if (button) {
        button.setAttribute('data-testid', 'assigned-connections-button');
      }
    });

    // Click the button using the newly added data-testid
    await this.page.getByTestId('assigned-connections-button').click();

    // Wait for and locate the dialog
    const dialog = await this.page.locator('[aria-labelledby="alert-dialog-slide-title"]');
    await expect(dialog).toBeVisible();
    // Define the aria-labelledby values for checkboxes
    const checkBoxAriaLabels = [
      'transfer-list-item-kubernetes-mantra-9-label',
      'transfer-list-item-meshery-themis-7-label',
      'transfer-list-item-meshery-vader-9-label',
      'transfer-list-item-meshery-exogorth-0-label',
      'transfer-list-item-meshery-hydra-8-label',
    ];

    // Iterate over the aria-labels and check each checkbox
    for (const label of checkBoxAriaLabels) {
      const checkbox = dialog.locator(`input[aria-labelledby="${label}"]`);
      await checkbox.waitFor({ state: 'visible' });
      if (!(await checkbox.isChecked())) {
        await checkbox.check();
      }
    }

    await this.page.click('button[aria-label="move selected right"]');
    await this.page.getByRole('button', { name: 'Save', exact: true }).isEnabled();
    await this.page.getByRole('button', { name: 'Save', exact: true }).click();
  }

  async moveAssignedConnectionsToAvailable(environmentName) {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/environments`);
    await this.page.waitForSelector(`text=${environmentName}`, { state: 'visible' });
    await this.page.evaluate(() => {
      const button = [...document.querySelectorAll('button')].find((b) =>
        b.textContent.includes('Assigned Connections'),
      );
      if (button) {
        button.setAttribute('data-testid', 'assigned-connections-button');
      }
    });

    // Click the button using the newly added data-testid
    await this.page.getByTestId('assigned-connections-button').click();

    // Wait for and locate the dialog
    const dialog = await this.page.locator('[aria-labelledby="alert-dialog-slide-title"]');
    await expect(dialog).toBeVisible();

    // Define the aria-labelledby values for checkboxes
    const AssignedConnections = [
      'transfer-list-item-kubernetes-mantra-9-label',
      'transfer-list-item-meshery-themis-7-label',
      'transfer-list-item-meshery-vader-9-label',
      'transfer-list-item-meshery-exogorth-0-label',
      'transfer-list-item-meshery-hydra-8-label',
    ];

    // Iterate over the aria-labels and check each checkbox in #rightList
    for (const label of AssignedConnections) {
      const checkbox = dialog.locator(`input[aria-labelledby="${label}"]`);
      await checkbox.waitFor({ state: 'visible' });
      if (!(await checkbox.isChecked())) {
        await checkbox.check();
      }
    }

    // Click on the move button to transfer items from rightList to left
    await this.page.locator('button[aria-label="move selected left"]').click();

    // Wait for the Save button to be enabled and click it
    await this.page.getByRole('button', { name: 'Save', exact: true }).isEnabled();
    await this.page.getByRole('button', { name: 'Save', exact: true }).click();
  }

  async addEnvironmentsToConnection(environmentName) {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/connections`);

    await this.page.waitForSelector('.MuiChip-label');
    await this.page.waitForSelector('tr[data-testid^="MUIDataTableBodyRow"]');

    const rows = await this.page.locator('tr[data-testid^="MUIDataTableBodyRow"]');
    const rowCount = await rows.count();
    let connectedNames = [];

    // gather all connected rows' names
    for (let i = 0; i < rowCount; i++) {
      // Find the connection status cell in the current row
      const connectionStatus = await rows
        .nth(i)
        .locator('td[data-colindex="10"] .MuiChip-label')
        .textContent();

      // If the connection status is 'connected', extract the connection name
      if (connectionStatus.trim().toLowerCase() === 'connected') {
        const connectionName = await rows
          .nth(i)
          .locator('td[data-colindex="3"] .MuiChip-label')
          .textContent();
        connectedNames.push({ index: i, name: connectionName.trim() });
      }
    }

    // iterate over the connected rows and select the environment for each one
    for (const { index } of connectedNames) {
      const dropdownIcon = rows
        .nth(index)
        .locator('td[data-colindex="4"] .css-1xc3v61-indicatorContainer')
        .first();

      // `data-testid` to the SVG icon in the current row
      await dropdownIcon.evaluate((node, index) => {
        node.setAttribute('data-testid', `dropdown-svg-${index}`);
      }, index);

      // Click on the dropdown icon (now identified by `data-testid`)
      await dropdownIcon.click();

      // Type the environment name into the dropdown input to filter options
      await this.page.fill(`#react-select-${index}-input`, environmentName);

      // Hit the "Enter" key after selecting the environment option
      await this.page.press(`#react-select-${index}-input`, 'Enter');
    }
  }

  async editEnvironmentCard(environmentName) {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/environments`);
    await this.page.click(
      `.MuiCard-root:has-text("${environmentName}") >> :not(:has-text("Assigned Connections"))`,
    );

    await this.page.click(
      `.MuiCard-root:has-text("${environmentName}") >> :not(:has-text("Assigned Connections")) >> [data-testid="envCard-edit-button"]`,
    );
    await this.page.waitForSelector('#Environment_name');
    await this.fillInput(this.page.locator('#Environment_name'), 'Sample-playwright-test-edited');
    await this.page.getByRole('button', { name: 'Update', exact: true }).click();
    const notification = await this.page.getByText(
      'Environment "Sample-playwright-test-edited" updated',
    );
    await expect(notification).toBeVisible();
  }

  async deleteEnvironmentCard() {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/environments`);
    await this.page.click(
      '.MuiCard-root:has-text("Sample-playwright-test-edited") >> :not(:has-text("Assigned Connections"))',
    );

    await this.page.click(
      `.MuiCard-root:has-text("Sample-playwright-test-edited") >> :not(:has-text("Assigned Connections")) >> [data-testid="envCard-delete-button"]`,
    );
    await this.page.getByRole('button', { name: 'DELETE', exact: true }).click();
    const notification = await this.page.getByText('Environment deleted');
    await expect(notification).toBeVisible();
  }
}
