import { expect, test } from '@playwright/test';

// TEMPORARY - proof that the restored E2E gate actually fails the job.
// This file exists only on the throwaway branch fm/e2e-harness-gate-demo and is
// never merged. It asserts something false on purpose.
test.describe('E2E gate demonstration', () => {
  test('deliberately fails to prove the gate fires', async () => {
    expect(1, 'this failure is intentional; see PR #21140').toBe(2);
  });
});
