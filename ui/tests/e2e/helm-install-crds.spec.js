import { expect, test } from '@playwright/test';
import { ENV } from './env';
import { execSync } from 'child_process'; // Essential for running helm and kubectl commands
import os from 'os'; // Might be useful for kubeconfig paths, though not directly used in the Helm commands below

// Helper function to execute shell commands and log output
function execShellCommand(command, options = {}) {
    console.log(`Executing command: ${command}`);
    try {
        const output = execSync(command, {
            encoding: 'utf-8',
            stdio: 'inherit', // This makes the command output visible in the test console
            ...options
        });
        return output;
    } catch (error) {
        console.error(`Command failed: ${command}`);
        console.error(`Error: ${error.message}`);
        // If stdio was 'pipe', you could also log error.stderr here:
        // console.error(`Stderr: ${error.stderr}`);
        throw error; // Re-throw the error to ensure the test fails if the command fails
    }
}

test.describe('Meshery Helm CRD Installation Tests', () => {

    // Define the common CRDs that Meshery installs
    // IMPORTANT: You MUST verify these exact CRD names by checking the Meshery Helm chart's CRD folder (e.g., charts/meshery/crds/)
    // or by inspecting a live cluster with 'kubectl get crds' after a default Meshery install.
    const mesheryCRDs = [
        'brokers.meshery.io',
        'meshsyncs.meshery.io',
        // Add any other Meshery-specific CRD names here
    ];

    // Optional: beforeEach or beforeAll to ensure helm repo is added and updated
    // This depends on whether your CI environment already sets up the Helm repo.
    // If not, you might need something like this:
    /*
    test.beforeAll(() => {
        console.log('Ensuring Meshery Helm repo is added and updated...');
        try {
            execShellCommand('helm repo add meshery https://meshery.io/charts/ || true'); // `|| true` prevents error if already added
            execShellCommand('helm repo update');
        } catch (e) {
            console.error('Failed to setup Helm repository:', e.message);
            throw e; // Fail if repo setup fails
        }
    });
    */

    test('should not install Meshery CRDs when installCRDs=false', async () => {
        const releaseName = 'meshery-e2e-no-crds';
        const namespace = 'meshery-e2e-test-ns-no-crds'; // Use unique namespaces for each test if possible

        try {
            // 1. Install Meshery with installCRDs=false
            console.log(`--- Test: ${test.info().title} ---`);
            console.log(`Installing Meshery release '${releaseName}' in namespace '${namespace}' with installCRDs=false...`);
            execShellCommand(`helm install ${releaseName} meshery/meshery --namespace ${namespace} --create-namespace --set installCRDs=false`);
            console.log('Meshery installation command executed.');

            // 2. Verify CRDs are NOT present
            for (const crd of mesheryCRDs) {
                console.log(`Checking for absence of CRD: ${crd}`);
                try {
                    // Try to get the CRD. If it exists, execShellCommand will throw an error.
                    // We expect it to throw here because the CRD should NOT be present.
                    execShellCommand(`kubectl get crd ${crd}`, { stdio: 'pipe' }); // stdio: 'pipe' to capture output/error for specific checks
                    // If the above line executes without throwing, it means the CRD was found - this is a test failure.
                    expect(true, `CRD '${crd}' should NOT be present but was found.`).toBeFalsy();
                } catch (error) {
                    // This is the expected path: kubectl get crd should fail if the CRD is not there
                    expect(error.status).not.toBe(0); // Assert that the command failed (non-zero exit code)
                    console.log(`Confirmed CRD '${crd}' is NOT present (as expected).`);
                }
            }

        } finally {
            // 3. Cleanup: Uninstall Meshery and delete namespace regardless of test outcome
            console.log(`Cleaning up: Uninstalling ${releaseName} and deleting namespace ${namespace}...`);
            try {
                execShellCommand(`helm uninstall ${releaseName} --namespace ${namespace} || true`);
            } catch (e) {
                console.warn(`Warning: Uninstall failed for ${releaseName} (may not exist or other issue): ${e.message}`);
            }
            try {
                // Use --force and --grace-period=0 for immediate deletion in tests, and --timeout for robustness
                execShellCommand(`kubectl delete namespace ${namespace} --timeout=2m --force --grace-period=0 || true`);
            } catch (e) {
                console.warn(`Warning: Namespace deletion failed for ${namespace} (may not exist or other issue): ${e.message}`);
            }
        }
    });

    test('should install Meshery CRDs by default', async () => {
        const releaseName = 'meshery-e2e-with-crds';
        const namespace = 'meshery-e2e-test-ns-default';

        try {
            // 1. Install Meshery with default settings (CRDs installed)
            console.log(`--- Test: ${test.info().title} ---`);
            console.log(`Installing Meshery release '${releaseName}' in namespace '${namespace}' with default CRD installation...`);
            execShellCommand(`helm install ${releaseName} meshery/meshery --namespace ${namespace} --create-namespace`);
            console.log('Meshery default installation command executed.');

            // 2. Verify CRDs ARE present
            for (const crd of mesheryCRDs) {
                console.log(`Checking for presence of CRD: ${crd}`);
                // Implement a retry mechanism as CRDs might take a moment to become available
                let crdFound = false;
                const maxRetries = 15; // Increased retries for robustness
                const retryDelayMs = 2000; // 2 seconds delay

                for (let i = 0; i < maxRetries; i++) {
                    try {
                        execShellCommand(`kubectl get crd ${crd}`, { stdio: 'pipe' }); // Use pipe to avoid printing output on success
                        crdFound = true;
                        console.log(`Confirmed CRD '${crd}' is present.`);
                        break; // CRD found, exit retry loop
                    } catch (error) {
                        console.log(`CRD '${crd}' not yet present, retrying... (${i + 1}/${maxRetries})`);
                        await new Promise(resolve => setTimeout(resolve, retryDelayMs));
                    }
                }
                expect(crdFound, `CRD '${crd}' should be present but was not found after retries.`).toBeTruthy();
            }

        } finally {
            // 3. Cleanup: Uninstall Meshery and delete namespace
            console.log(`Cleaning up: Uninstalling ${releaseName} and deleting namespace ${namespace}...`);
            try {
                execShellCommand(`helm uninstall ${releaseName} --namespace ${namespace} || true`);
            } catch (e) {
                console.warn(`Warning: Uninstall failed for ${releaseName} (may not exist or other issue): ${e.message}`);
            }
            try {
                execShellCommand(`kubectl delete namespace ${namespace} --timeout=2m --force --grace-period=0 || true`);
            } catch (e) {
                console.warn(`Warning: Namespace deletion failed for ${namespace} (may not exist or other issue): ${e.message}`);
            }
        }
    });
});