import { createHash, randomUUID } from 'crypto';
import { mkdirSync, readFileSync, writeFileSync } from 'fs';
import { template } from 'lodash';
import moment from 'moment';
import os from 'os';
import path from 'path';

class MyReporter {
  introMessage = '';
  totalTests = '';
  expectedTest = '';
  testTableData = [];
  failedTestDetails = [];
  relationshipTestData = [];
  seenRelationships = new Set();
  passed = 0;
  failed = 0;
  skipped = 0;
  flaky = 0;
  countLog = 1;
  relationshipAllureEntries = new Map();

  onBegin(_config, suite) {
    this.introMessage = `- Testing started at: ${moment().format('MMMM Do YYYY, h:mm:ss a')}`;
    this.totalTests = `- Total test cases: ${suite.allTests().length}`;
  }

  onStdOut(chunk, _test, _result) {
    const text = chunk.toString('utf-8');
    process.stdout.write(text);
  }

  onStdErr(chunk, _test, _result) {
    const text = chunk.toString('utf-8');
    process.stderr.write(text);
  }

  onTestEnd(test, result) {
    const status = test.outcome();
    const project = test.parent?.project()?.name;
    const provider = test.parent?.project()?.use?.provider || 'Unknown';

    this.displayLogs(project, test.title, test.tags, status, result);
    this.addTestTable(
      project,
      provider,
      test.title,
      test.tags,
      status,
      result.retry,
      test.retries,
      result,
    );

    // Process relationship annotations if present
    const relationshipAnnotations =
      test.annotations?.filter((a) => a.type === 'relationship') || [];
    for (const annotation of relationshipAnnotations) {
      try {
        const data = JSON.parse(annotation.description);
        this.collectRelationshipData(data);
        this.collectRelationshipAllureEntry(test, result, project, provider, data);
      } catch {
        // Skip invalid annotations
      }
    }
  }

  collectRelationshipData(data) {
    // Create a unique key for deduplication
    const key = `${data.kind}|${data.type}|${data.subType}|${data.from}|${data.to}|${data.model}|${data.designName}`;

    // Skip if we've already seen this relationship
    if (this.seenRelationships.has(key)) {
      return;
    }
    this.seenRelationships.add(key);
    this.relationshipTestData.push(data);
  }

  buildRelationshipTestTable() {
    // Sort by kind, type, subType, model, from, to
    this.relationshipTestData.sort((a, b) => {
      return (
        a.kind.localeCompare(b.kind) ||
        a.type.localeCompare(b.type) ||
        a.subType.localeCompare(b.subType) ||
        a.model.localeCompare(b.model) ||
        a.from.localeCompare(b.from) ||
        a.to.localeCompare(b.to)
      );
    });

    let table = `| Kind | Type | SubType | From | To | Model | Design Name | Status |
| :---: | :---: | :---: | :---: | :---: | :---: | :--- | :---: |`;

    for (const data of this.relationshipTestData) {
      const statusEmoji = data.status === 'pass' ? '✅' : '❌';
      table += `\n| ${data.kind} | ${data.type} | ${data.subType} | ${data.from} | ${data.to} | ${data.model} | ${data.designName} | ${statusEmoji} |`;
    }

    return table;
  }

  buildTestTable() {
    // Sort by provider, then by project, then by title
    this.testTableData.sort((a, b) => {
      return (
        a.provider.localeCompare(b.provider) ||
        a.project.localeCompare(b.project) ||
        a.title.localeCompare(b.title)
      );
    });

    let table = `| Test | Provider | Browser | Test Case | Tags | Result |
| :---: | :---: | :---: | :--- | :---: | :---: |`;

    let count = 1;
    for (const row of this.testTableData) {
      table += `\n| ${count} | ${row.provider} | ${row.project} | ${row.title} | ${row.tags} | ${row.statusEmoji} |`;
      count++;
    }

    return table;
  }

  buildFailedTestDetails() {
    if (this.failedTestDetails.length === 0) {
      return '';
    }

    let details = '';
    for (const failedTest of this.failedTestDetails) {
      const escapedTitle = failedTest.title.replace(/[|]/g, '\\|');
      details += `
<details>
<summary>❌ ${escapedTitle} (${failedTest.provider} - ${failedTest.project})</summary>

**File Location:** \`${failedTest.fileLocation}\`

**Error Message:**
\`\`\`
${failedTest.errorMessage}
\`\`\`

${failedTest.snippet ? `**Code Snippet:**\n\`\`\`\n${failedTest.snippet}\n\`\`\`` : ''}

</details>
`;
    }
    return details;
  }

  async onEnd(result) {
    const message = await this.buildMessage(result);

    try {
      writeFileSync('test-report.md', message);
    } catch (e) {
      console.log('Cannot write test reporter ', e);
    }

    this.writeRelationshipAllureResults();
  }

  collectRelationshipAllureEntry(test, result, project, provider, data) {
    const fileName = path.basename(test.location.file);
    const suiteTitle = test.parent?.title || 'Relationship Evaluation';
    const key = this.createRelationshipAllureKey(project, provider, data);

    this.relationshipAllureEntries.set(key, {
      fileName,
      project,
      provider,
      suiteTitle,
      data,
      start: result.startTime?.getTime?.() ?? Date.now(),
      stop:
        (result.startTime?.getTime?.() ?? Date.now()) +
        (typeof result.duration === 'number' ? result.duration : 0),
      line: test.location.line,
      column: test.location.column,
    });
  }

  createRelationshipAllureKey(project, provider, data) {
    return [
      project,
      provider,
      data.kind,
      data.type,
      data.subType,
      data.from,
      data.to,
      data.model,
      data.designName,
    ].join('|');
  }

  writeRelationshipAllureResults() {
    if (this.relationshipAllureEntries.size === 0) {
      return;
    }

    const allureResultsDir = path.resolve(process.cwd(), 'allure-results');
    mkdirSync(allureResultsDir, { recursive: true });

    for (const entry of this.relationshipAllureEntries.values()) {
      const { fileName, project, provider, suiteTitle, data, start, stop, line, column } = entry;

      const historySource = JSON.stringify({
        fileName,
        project,
        provider,
        suiteTitle,
        kind: data.kind,
        type: data.type,
        subType: data.subType,
        from: data.from,
        to: data.to,
        model: data.model,
        designName: data.designName,
      });
      const testCaseId = createHash('md5').update(historySource).digest('hex');
      const attachmentSource = `${randomUUID()}-attachment.json`;
      const resultSource = `${randomUUID()}-result.json`;
      const relationshipName =
        `Relationship Result: ${data.kind}/${data.type}/${data.subType || '-'} ` +
        `${data.from} -> ${data.to} (${data.designName})`;

      writeFileSync(path.join(allureResultsDir, attachmentSource), JSON.stringify(data, null, 2));

      writeFileSync(
        path.join(allureResultsDir, resultSource),
        JSON.stringify(
          {
            uuid: randomUUID(),
            name: relationshipName,
            historyId: testCaseId,
            status: data.status === 'pass' ? 'passed' : 'failed',
            statusDetails: {},
            stage: 'finished',
            steps: [],
            attachments: [
              {
                name: 'relationship-data',
                source: attachmentSource,
                type: 'application/json',
              },
            ],
            parameters: [
              { name: 'Project', value: project },
              { name: 'Provider', value: provider },
              { name: 'Design', value: data.designName },
              { name: 'Model', value: data.model },
              { name: 'From', value: data.from },
              { name: 'To', value: data.to },
            ],
            labels: [
              { name: 'project', value: process.env.ALLURE_LABEL_project || 'Meshery' },
              { name: 'type', value: process.env.ALLURE_LABEL_type || 'e2e' },
              { name: 'language', value: 'javascript' },
              { name: 'framework', value: 'playwright' },
              { name: 'package', value: fileName },
              { name: 'tag', value: 'relationship' },
              { name: 'parentSuite', value: project },
              { name: 'suite', value: fileName },
              { name: 'subSuite', value: suiteTitle },
              { name: 'host', value: os.hostname() },
              { name: 'thread', value: `pid-${process.pid}` },
            ],
            links: [],
            start,
            stop,
            testCaseId,
            fullName: `${fileName}:${line}:${column}:${testCaseId}`,
            titlePath: [fileName, suiteTitle, relationshipName],
          },
          null,
          2,
        ),
      );
    }
  }

  displayLogs(project, title, tags, status, result) {
    const allTags = tags.map((item) => item.replace('@', '')).join(', ');

    const logs = `${
      this.countLog
    }. Project: ${project}, Test: ${title}, Status: ${this.getStatusEmoji(
      tags,
      status,
    )}, Tags: ${allTags} ${
      status === 'unexpected' && result.error !== undefined
        ? '\nFile Location: ' +
          (result.error.location?.file ?? 'Not Found') +
          '\n' +
          (result.error?.snippet ?? 'No snippet') +
          '\n' +
          (result.error?.message ?? 'No message')
        : ''
    }\n`;

    process.stdout.write(logs);

    this.countLog++;
  }

  addTestTable(project, provider, title, tags, status, retry, retries, result) {
    this.countTestStatus(tags, status, retry, retries);

    const lastRetriesRun = retry === retries;
    const isFail = status === 'unexpected';
    const isSkipped = status === 'skipped';
    const isUnstableTest = tags.includes('@unstable');

    // Skip intermediate retry attempts for failed tests
    if ((isFail || isSkipped) && !lastRetriesRun) {
      return;
    }

    // For flaky unstable tests that failed, skip intermediate retries
    if (isUnstableTest && isFail && !lastRetriesRun) {
      return;
    }

    const allTags = tags.map((item) => item.replace('@', '')).join(', ');
    const statusEmoji = this.getStatusEmoji(tags, status);

    // Add to test table data for sorting
    this.testTableData.push({
      provider: provider,
      project: project,
      title: title,
      tags: allTags,
      statusEmoji: statusEmoji,
    });

    // Collect failed test details for collapsible section
    if (isFail && lastRetriesRun && !isUnstableTest && result.error) {
      this.failedTestDetails.push({
        provider: provider,
        project: project,
        title: title,
        tags: allTags,
        fileLocation: result.error.location?.file ?? 'Not Found',
        errorMessage: result.error.message ?? 'No error message available',
        snippet: result.error.snippet ?? null,
      });
    }
  }

  countTestStatus(tags, status, retry, retries) {
    const lastRetriesRun = retry === retries;
    const isFail = status === 'unexpected';
    const isSkipped = status === 'skipped';
    const isFlaky = status === 'flaky';

    const isUnstableTest = tags.includes('@unstable');

    if (status === 'expected') {
      this.passed++;
      return;
    }
    if (isFlaky || (isUnstableTest && isFail)) {
      this.flaky++;
      return;
    }
    if (isFail && lastRetriesRun && !isUnstableTest) {
      this.failed++;
      return;
    }
    if (isSkipped && lastRetriesRun && !isUnstableTest) {
      this.skipped++;
      return;
    }
  }

  getStatusEmoji(tags, status) {
    const isFlakyTest = tags.includes('@unstable');

    if (status === 'expected') {
      return '✅';
    }
    if (status === 'unexpected' && !isFlakyTest) {
      return '❌';
    }
    if (status === 'flaky' || isFlakyTest) {
      return '⚠️';
    }
    if (status === 'skipped') {
      return '➖';
    }
    return '➖';
  }

  async buildMessage(result) {
    const duration = moment.duration(result.duration, 'milliseconds');
    const minutes = Math.floor(duration.asMinutes());
    const seconds = duration.seconds();
    const templateStr = readFileSync(path.join(__dirname, 'reporterSummary.md'), 'utf8');
    return template(templateStr)({
      introMessage: this.introMessage,
      minutes,
      seconds,
      passed: this.passed,
      failed: this.failed,
      flaky: this.flaky,
      skipped: this.skipped,
      totalTests: this.totalTests,
      testTable: this.buildTestTable(),
      failedTestDetails: this.buildFailedTestDetails(),
      failedTestCount: this.failedTestDetails.length,
      relationshipTestTable: this.buildRelationshipTestTable(),
      relationshipTestCount: this.relationshipTestData.length,
    });
  }
}

export default MyReporter;
