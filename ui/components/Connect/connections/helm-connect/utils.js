/**
 * Generates a selected Helm chart repository based on the user's choices.
 *
 * This function processes the selected Helm charts and their dependencies (sub charts), and
 * constructs an array of selected Helm charts with appropriate properties.
 *
 * @param {object} data - The user's data containing selected Helm charts.
 * @param {object} sharedData - All discovered helm Charts.
 * @returns {object} An object containing the updated user data and an array of selected Helm charts.
 */
export const generateSelectedHelmRepo = (data, sharedData) => {
  let selectedNameOfCharts = data.selectedHelmRepos.slice();
  if (data?.selectedHelmRepos.length > 0) {
    let arrayOfSelected = [];

    data.selectedHelmRepos.forEach((selectedChartName) => {
      if (selectedChartName.includes('/')) {
        // Handle repositories with dependencies
        const [parentName, dependencyName] = selectedChartName.split('/');

        // Find the parent chart by matching the name before "/"
        const parentChart = sharedData.helmRepoChartsData.find(
          (chartData) => chartData.Name === parentName,
        );
        if (parentChart) {
          // Find the parent chart in arrayOfSelected and add it if not already there
          let existingParent = arrayOfSelected.find((item) => item.name === parentChart.Name);
          if (!existingParent) {
            existingParent = {
              name: parentChart.Name,
              appVersion: parentChart.AppVersion,
              version: parentChart.Version,
              apiVersion: parentChart.APIVersion || '',
              dependencies: parentChart.Dependencies
                ? parentChart.Dependencies.map((dependency) => {
                    if (dependency.name === dependencyName) {
                      return {
                        ...dependency,
                        enabled: true,
                      };
                    } else {
                      return {
                        ...dependency,
                      };
                    }
                  })
                : [],
            };
            arrayOfSelected.push(existingParent);
            if (!selectedNameOfCharts.includes(existingParent.name)) {
              selectedNameOfCharts.push(existingParent.name);
            }
          } else {
            // If parentChart already exists, add enabled: true only to the selected dependency
            const existingDependency = existingParent.dependencies.find(
              (dependency) => dependency.name === dependencyName,
            );
            if (existingDependency) {
              existingDependency.enabled = true;
            }
          }
        }
      } else {
        // Handle repositories without dependencies
        const selectedChart = sharedData.helmRepoChartsData.find(
          (chartData) => chartData.Name === selectedChartName,
        );

        let existingParent = arrayOfSelected.find((item) => item.name === selectedChart.Name);

        if (selectedChart && !existingParent) {
          arrayOfSelected.push({
            name: selectedChart.Name,
            appVersion: selectedChart.AppVersion,
            version: selectedChart.Version,
            apiVersion: selectedChart.APIVersion || '',
            dependencies: selectedChart.Dependencies
              ? selectedChart.Dependencies.map((dependency) => {
                  // Remove the enabled property if it's false
                  const { enabled, ...rest } = dependency;
                  return enabled ? { ...rest } : { ...rest };
                })
              : [], // Initialize dependencies array
          });
          if (!selectedNameOfCharts.includes(selectedChart.Name)) {
            selectedNameOfCharts.push(selectedChart.Name);
          }
        }
      }
    });

    // Remove dependencies property if empty
    arrayOfSelected = arrayOfSelected.map((item) => {
      if (item.dependencies.length === 0) {
        // Remove the dependencies property if there are no dependencies
        // eslint-disable-next-line no-unused-vars
        const { dependencies, ...itemWithoutDependencies } = item;
        return itemWithoutDependencies;
      }
      return item;
    });

    const updatedData = { ...data, selectedHelmRepos: selectedNameOfCharts };
    return { updatedData, arrayOfSelected };
  }
  return { updatedData: data, arrayOfSelected: [] };
};

/**
 * Adds a 'status' field to each chart in the given array based on its digestion status.
 *
 * @param {Array} charts - An array of chart objects.
 * @param {Array} undigestedCharts - An array of undigested chart objects for comparison.
 * @returns {Array} - An array of chart objects with added 'status' field indicating digestion status.
 */
export const addStatusToCharts = (charts, undigestedCharts) => {
  const flattenedCharts = [];

  const flatten = (chart) => {
    const isUndigested = undigestedCharts.some((undigestedChart) => {
      return undigestedChart.name === chart.name;
    });

    // Create a copy of the chart without dependencies
    const chartWithoutDependencies = { ...chart };
    delete chartWithoutDependencies.dependencies;

    if (!isUndigested || chart.enabled) {
      flattenedCharts.push({ ...chartWithoutDependencies, status: !isUndigested });
    }

    if (chart.dependencies) {
      chart.dependencies.forEach((dependency) => {
        if (dependency.enabled === true) {
          flatten(dependency);
        }
      });
    }
  };

  charts.forEach((chart) => flatten(chart));

  return flattenedCharts;
};
