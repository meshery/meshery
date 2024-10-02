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
      const notification = await this.page.getByTestId('notification-success-bar');
      await expect(notification).toBeVisible();
      await expect(notification).toBeHidden({ timeout: 30000 });

      // Wait for the profile card to appear
      await this.page.waitForSelector(`text=${environmentName}`, {
        state: 'visible',
        timeout: 3 * 60 * 1000,
      });

      // Additional check: wait for the profile card to be clickable
      await this.page.waitForSelector(`text=${environmentName}`, {
        state: 'visible',
        timeout: 30000,
      });

      console.log(`Profile "${environmentName}" has been created and is visible.`);
    } else {
      console.log(`Profile "${environmentName}" already exists.`);
    }
  }

  async assignConnections(environmentName) {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/environments`);
    await this.page.waitForSelector(`text=${environmentName}`, { state: 'visible' });

    // Click the button using the newly added data-testid
    await this.page.getByTestId(`assigned-conn-${environmentName}`).click();

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
    await this.page.getByTestId('move-selected-right').click();

    // Verify if checkboxes are moved to the right list
    const rightList = await this.page.getByTestId('connections-rightList');

    for (const label of checkBoxAriaLabels) {
      const rightListItem = rightList.locator(`input[aria-labelledby="${label}"]`);
      await rightListItem.waitFor({ state: 'visible' });
      await expect(rightListItem).toBeVisible();
    }

    // Click save button
    const saveConnBtn = await this.page.getByTestId('assign-conn-save-button');
    await expect(saveConnBtn).toBeVisible();
    await saveConnBtn.click();

    // Verify if the assigned connection is visible
    await expect(this.page.getByTestId(`assigned-conn-${environmentName}`)).toBeVisible({
      state: 'visible',
      timeout: 3 * 60 * 1000,
    });
  }

  async moveAssignedConnectionsToAvailable(environmentName) {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/environments`);
    await this.page.waitForSelector(`text=${environmentName}`, { state: 'visible' });

    // Click the button using the newly added data-testid
    await this.page.getByTestId(`assigned-conn-${environmentName}`).click();

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

    await this.page.getByTestId('move-selected-left').click();

    const leftList = await this.page.getByTestId('connections-leftList');

    for (const label of AssignedConnections) {
      const rightListItem = leftList.locator(`input[aria-labelledby="${label}"]`);
      await rightListItem.waitFor({ state: 'visible' });
      await expect(rightListItem).toBeVisible();
    }

    const saveConnBtn = await this.page.getByTestId('assign-conn-save-button');
    await expect(saveConnBtn).toBeVisible();
    await saveConnBtn.click();

    // Verify if the assigned connection is visible
    await expect(this.page.getByTestId(`assigned-conn-${environmentName}`)).toBeVisible({
      state: 'visible',
      timeout: 3 * 60 * 1000,
    });
  }

  async addEnvironmentsToConnection(environmentName) {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/connections`);

    await this.page.waitForSelector('.MuiChip-label');
    await this.page.waitForSelector('tr[data-testid^="MUIDataTableBodyRow"]');

    const rows = await this.page.locator('tr[data-testid^="MUIDataTableBodyRow"]');
    const rowCount = await rows.count();
    let connectedRows = [];

    // Gather all connected rows
    for (let i = 0; i < rowCount; i++) {
      const connectionStatus = await rows
        .nth(i)
        .locator('td[data-colindex="10"] .MuiChip-label')
        .textContent();

      if (connectionStatus.trim().toLowerCase() === 'connected') {
        connectedRows.push(i);
      }

      // Break if we have found 5 connected rows
      if (connectedRows.length === 5) break;
    }

    // Iterate over the top 5 (or fewer) connected rows and assign the environment
    for (const index of connectedRows) {
      const dropdownIcon = rows
        .nth(index)
        .locator('td[data-colindex="4"] .css-1xc3v61-indicatorContainer')
        .first();

      // Set `data-testid` to the SVG icon in the current row
      await dropdownIcon.evaluate((node, index) => {
        node.setAttribute('data-testid', `dropdown-svg-${index}`);
      }, index);

      // Click on the dropdown icon
      await dropdownIcon.click();

      const checkboxSelector = `[data-testid="checkbox-${environmentName}"]`;

      // Check if the checkbox exists, then click it
      const checkbox = this.page.locator(checkboxSelector);
      await checkbox.waitFor({ state: 'visible' });

      // Check if the checkbox is already checked
      const isChecked = await checkbox.isChecked();

      if (!isChecked) {
        await checkbox.check();
      } else {
        // If already checked, close the dropdown without making changes
        await dropdownIcon.click();
      }
    }
  }

  async editEnvironmentCard(environmentName) {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/environments`);
    await this.page.click(
      `[data-testid="flip-card-${environmentName}"] >> :not(:has-text("Assigned Connections"))`,
    );

    await this.page.getByTestId(`edit-button-${environmentName}`).click();
    await this.page.waitForSelector('#Environment_name');
    await this.fillInput(this.page.locator('#Environment_name'), 'Sample-playwright-test-edited');
    await this.page.getByRole('button', { name: 'Update', exact: true }).click();
    const notification = await this.page.getByText(
      'Environment "Sample-playwright-test-edited" updated',
    );
    await expect(notification).toBeVisible();
  }

  // AFTER EVERY THING IS DONE TEST THIS FUNCTION :# NOT TESTED
  async deleteEnvironmentCard() {
    await expect(this.page).toHaveURL(`${ENV.MESHERY_SERVER_URL}/management/environments`);
    await this.page.click(
      `[data-testid="flip-card-Sample-playwright-test-edited"] >> :not(:has-text("Assigned Connections"))`,
    );

    await this.page.getByTestId('delete-button-Sample-playwright-test-edited').click();
    await this.page.getByRole('button', { name: 'DELETE', exact: true }).click();
    const notification = await this.page.getByText('Environment deleted');
    await expect(notification).toBeVisible();
  }
}
