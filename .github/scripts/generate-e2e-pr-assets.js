const fs = require("fs");
const path = require("path");

const QA_BASE_URL = "https://qa.meshery.io/meshery/prs";

function readJson(filePath) {
  return JSON.parse(fs.readFileSync(filePath, "utf8"));
}

function ensureDir(dirPath) {
  fs.mkdirSync(dirPath, { recursive: true });
}

function buildArtifactUrl(repository, runId, artifactId) {
  return `https://github.com/${repository}/actions/runs/${runId}/artifacts/${artifactId}`;
}

function escapeHtml(value) {
  return String(value ?? "")
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function formatLinks(test, artifactMap) {
  const links = [];
  if (test.qa?.shaUrl) {
    links.push(`<a href="${test.qa.shaUrl}">QA summary</a>`);
  }
  if (artifactMap["playwright-report"]) {
    links.push(
      `<a href="${artifactMap["playwright-report"].url}">Playwright report</a>`,
    );
  }
  if (artifactMap["playwright-traces"]) {
    links.push(
      `<a href="${artifactMap["playwright-traces"].url}">Trace bundle</a>`,
    );
  }
  if (artifactMap["meshery-server-log"]) {
    links.push(
      `<a href="${artifactMap["meshery-server-log"].url}">Server log</a>`,
    );
  }
  if (artifactMap["allure-results"]) {
    links.push(
      `<a href="${artifactMap["allure-results"].url}">Allure results</a>`,
    );
  }

  return links.join(" · ");
}

function classifyAttachments(attachments = []) {
  return attachments.reduce(
    (acc, attachment) => {
      const name =
        `${attachment.name ?? ""} ${attachment.path ?? ""}`.toLowerCase();
      if (name.includes("trace")) {
        acc.trace.push(attachment.path);
      } else if (name.includes("video")) {
        acc.video.push(attachment.path);
      } else if (name.includes("screenshot")) {
        acc.screenshot.push(attachment.path);
      } else {
        acc.other.push(attachment.path);
      }
      return acc;
    },
    { trace: [], video: [], screenshot: [], other: [] },
  );
}

function buildHtml(contextData) {
  const failedRows =
    contextData.failedTests.length === 0
      ? "<p>No failed tests were reported for this run.</p>"
      : `<table>
  <thead>
    <tr>
      <th>Failed test</th>
      <th>Provider</th>
      <th>Browser</th>
      <th>Failure summary</th>
      <th>Artifacts</th>
    </tr>
  </thead>
  <tbody>
    ${contextData.failedTests
      .map(
        (test) => `<tr id="${escapeHtml(test.id)}">
      <td><strong>${escapeHtml(test.title)}</strong><br /><code>${escapeHtml(test.fileLocation)}</code></td>
      <td>${escapeHtml(test.provider)}</td>
      <td>${escapeHtml(test.project)}</td>
      <td>${escapeHtml(test.errorSummary)}</td>
      <td>${formatLinks(test, contextData.artifacts)}</td>
    </tr>`,
      )
      .join("\n")}
  </tbody>
</table>`;

  return `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Meshery PR E2E Report - PR #${escapeHtml(contextData.pullRequest.number)}</title>
    <style>
      body { font-family: Arial, sans-serif; margin: 2rem auto; max-width: 1100px; line-height: 1.5; color: #111827; padding: 0 1rem; }
      code { background: #f3f4f6; padding: 0.1rem 0.3rem; border-radius: 4px; }
      .summary { display: flex; flex-wrap: wrap; gap: 0.75rem; margin: 1rem 0; }
      .pill { background: #f3f4f6; border-radius: 999px; padding: 0.4rem 0.75rem; }
      .links { margin: 1rem 0; }
      table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
      th, td { border: 1px solid #d1d5db; padding: 0.75rem; vertical-align: top; text-align: left; }
      th { background: #f9fafb; }
      .notice { background: #eff6ff; border-left: 4px solid #2563eb; padding: 1rem; margin: 1rem 0; }
    </style>
  </head>
  <body>
    <h1>Meshery PR E2E Report</h1>
    <p>Pull request <strong>#${escapeHtml(contextData.pullRequest.number)}</strong> · Commit <code>${escapeHtml(
      contextData.pullRequest.sha,
    )}</code></p>
    <div class="summary">
      <span class="pill">✅ Passed: ${escapeHtml(contextData.summary.passed)}</span>
      <span class="pill">❌ Failed: ${escapeHtml(contextData.summary.failed)}</span>
      <span class="pill">⚠️ Flaky: ${escapeHtml(contextData.summary.flaky)}</span>
      <span class="pill">⏩ Skipped: ${escapeHtml(contextData.summary.skipped)}</span>
      <span class="pill">⌛ Duration: ${escapeHtml(contextData.summary.durationHuman)}</span>
    </div>
    <div class="notice">
      This PR-specific summary is generated from the Meshery E2E workflow. If a QA page has not propagated yet,
      use the GitHub Actions artifact links below.
    </div>
    <div class="links">
      <a href="${escapeHtml(contextData.run.url)}">Workflow run</a>
      · <a href="${escapeHtml(contextData.qa.latestUrl)}">Canonical latest QA URL</a>
      · <a href="${escapeHtml(contextData.qa.shaUrl)}">Canonical SHA QA URL</a>
    </div>
    <h2>Failed tests</h2>
    ${failedRows}
  </body>
</html>`;
}

module.exports = async ({ github, context, core }) => {
  const workspace = process.env.GITHUB_WORKSPACE || process.cwd();
  const uiDir = path.join(workspace, "ui");
  const reportPath = path.join(uiDir, "test-report.json");

  if (!fs.existsSync(reportPath)) {
    core.setFailed(`Structured E2E report not found at ${reportPath}`);
    return;
  }

  const report = readJson(reportPath);
  const prNumber = process.env.PR_NUMBER;
  const prSha = process.env.PR_SHA || process.env.GITHUB_SHA;

  if (!prNumber) {
    core.setFailed("PR_NUMBER is required to generate PR E2E assets.");
    return;
  }

  const artifactResponse = await github.rest.actions.listWorkflowRunArtifacts({
    owner: context.repo.owner,
    repo: context.repo.repo,
    run_id: context.runId,
    per_page: 100,
  });

  const artifactMap = artifactResponse.data.artifacts.reduce(
    (acc, artifact) => {
      acc[artifact.name] = {
        id: artifact.id,
        name: artifact.name,
        url: buildArtifactUrl(
          process.env.GITHUB_REPOSITORY,
          context.runId,
          artifact.id,
        ),
        expired: artifact.expired,
      };
      return acc;
    },
    {},
  );

  const qaBasePath = `${QA_BASE_URL}/pr-${prNumber}`;
  const contextData = {
    version: 1,
    generatedAt: new Date().toISOString(),
    repository: process.env.GITHUB_REPOSITORY,
    pullRequest: {
      number: Number(prNumber),
      sha: prSha,
    },
    run: {
      id: context.runId,
      url: `${process.env.GITHUB_SERVER_URL}/${process.env.GITHUB_REPOSITORY}/actions/runs/${context.runId}`,
    },
    summary: report.summary,
    artifacts: artifactMap,
    qa: {
      latestUrl: `${qaBasePath}/latest/`,
      shaUrl: `${qaBasePath}/${prSha}/`,
    },
    failedTests: report.failedTests.map((test) => ({
      ...test,
      attachmentGroups: classifyAttachments(test.attachments),
      qa: {
        latestUrl: `${qaBasePath}/latest/#${test.id}`,
        shaUrl: `${qaBasePath}/${prSha}/#${test.id}`,
      },
    })),
    tests: report.tests,
    relationshipTests: report.relationshipTests,
  };

  fs.writeFileSync(
    path.join(uiDir, "e2e-pr-comment-context.json"),
    JSON.stringify(contextData, null, 2),
  );

  const qaReportDir = path.join(uiDir, "qa-pr-report");
  ensureDir(qaReportDir);
  fs.writeFileSync(
    path.join(qaReportDir, "manifest.json"),
    JSON.stringify(contextData, null, 2),
  );
  fs.writeFileSync(
    path.join(qaReportDir, "index.html"),
    buildHtml(contextData),
  );

  core.info(
    `Generated PR E2E assets for PR #${prNumber} with ${contextData.failedTests.length} failed test entries.`,
  );
};
