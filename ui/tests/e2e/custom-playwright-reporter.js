import { createHash } from 'crypto';
import { readFileSync, writeFileSync } from 'fs';
import { template } from 'lodash';
import moment from 'moment';
import path from 'path';

const ANSI_PATTERN =
  /[\u001B\u009B][[\]()#;?]*(?:(?:\d{1,4}(?:;\d{0,4})*)?[\dA-PR-TZcf-ntqry=><~]|(?:[^\u001B\u009B])*?[A-PR-TZcf-ntqry=><~])/g;
const MAX_ERROR_LINES = 4;
const MAX_ERROR_LENGTH = 320;
const MAX_SNIPPET_LINES = 8;
const MAX_SNIPPET_LENGTH = 600;

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
      test,
      result,
    );

    // Process relationship annotations if present
    const relationshipAnnotations =
      test.annotations?.filter((a) => a.type === 'relationship') || [];
    for (const annotation of relationshipAnnotations) {
      try {
        const data = JSON.parse(annotation.description);
        this.collectRelationshipData(data);
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

    let details = `| Test | Provider | Browser | File | Failure Summary |\n| :--- | :---: | :---: | :--- | :--- |\n`;
    for (const failedTest of this.failedTestDetails) {
      details += `\n| ${this.escapeTableValue(failedTest.title)} | ${this.escapeTableValue(
        failedTest.provider,
      )} | ${this.escapeTableValue(failedTest.project)} | \`${failedTest.fileLocation}\` | ${
        failedTest.errorSummary
      } |`;
    }
    return details;
  }

  async onEnd(result) {
    const message = await this.buildMessage(result);
    const jsonReport = this.buildStructuredReport(result);

    try {
      writeFileSync('test-report.md', message);
      writeFileSync('test-report.json', JSON.stringify(jsonReport, null, 2));
    } catch (e) {
      console.log('Cannot write test reporter ', e);
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

  addTestTable(project, provider, title, tags, status, retry, retries, test, result) {
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
      const fileLocation = this.getRelativePath(
        result.error.location?.file ?? test?.location?.file ?? 'Not Found',
      );
      this.failedTestDetails.push({
        id: this.buildTestId(project, provider, title, fileLocation),
        provider: provider,
        project: project,
        title: title,
        tags: allTags,
        fileLocation,
        errorSummary: this.summarizeError(result.error.message),
        errorMessage: this.cleanText(result.error.message ?? 'No error message available'),
        snippet: this.cleanSnippet(result.error.snippet),
        attachments: this.collectAttachments(result.attachments),
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

  buildStructuredReport(result) {
    const duration = moment.duration(result.duration, 'milliseconds');

    return {
      version: 1,
      generatedAt: new Date().toISOString(),
      summary: {
        startedAt: this.introMessage.replace('- Testing started at: ', ''),
        durationMs: result.duration,
        durationHuman: `${Math.floor(duration.asMinutes())} minutes and ${duration.seconds()} seconds`,
        total: this.passed + this.failed + this.flaky + this.skipped,
        passed: this.passed,
        failed: this.failed,
        flaky: this.flaky,
        skipped: this.skipped,
        overallStatus: this.failed > 0 ? 'failed' : 'passed',
      },
      tests: this.testTableData.map((test) => ({
        provider: test.provider,
        project: test.project,
        title: test.title,
        tags: test.tags,
        statusEmoji: test.statusEmoji,
      })),
      failedTests: this.failedTestDetails,
      relationshipTests: this.relationshipTestData,
    };
  }

  buildTestId(project, provider, title, fileLocation) {
    return createHash('sha1')
      .update([project, provider, title, fileLocation].join('::'))
      .digest('hex')
      .slice(0, 12);
  }

  collectAttachments(attachments = []) {
    return attachments
      .filter((attachment) => attachment.path)
      .map((attachment) => ({
        name: attachment.name,
        contentType: attachment.contentType ?? null,
        path: this.getRelativePath(attachment.path),
      }));
  }

  getRelativePath(filePath) {
    if (!filePath || filePath === 'Not Found') {
      return 'Not Found';
    }

    return path.relative(process.cwd(), filePath) || filePath;
  }

  cleanSnippet(snippet) {
    const cleanValue = this.cleanText(snippet);
    if (!cleanValue) {
      return null;
    }

    const lines = cleanValue.split('\n').slice(0, MAX_SNIPPET_LINES);
    const normalized = lines.join('\n');
    if (normalized.length <= MAX_SNIPPET_LENGTH) {
      return normalized;
    }

    return `${normalized.slice(0, MAX_SNIPPET_LENGTH - 1)}…`;
  }

  summarizeError(message) {
    const cleanValue = this.cleanText(message);
    if (!cleanValue) {
      return 'No failure summary available.';
    }

    const summary = cleanValue
      .split('\n')
      .map((line) => line.trim())
      .filter(Boolean)
      .slice(0, MAX_ERROR_LINES)
      .join(' ');

    const normalized = summary.replace(/\s+/g, ' ');
    if (normalized.length <= MAX_ERROR_LENGTH) {
      return this.escapeTableValue(normalized);
    }

    return `${this.escapeTableValue(normalized.slice(0, MAX_ERROR_LENGTH - 1))}…`;
  }

  cleanText(text) {
    if (!text) {
      return '';
    }

    return text
      .replace(ANSI_PATTERN, '')
      .replace(/\r\n/g, '\n')
      .replace(/\r/g, '\n')
      .split('\n')
      .map((line) => line.trimEnd())
      .join('\n')
      .trim();
  }

  escapeTableValue(value) {
    return String(value ?? '')
      .replace(/\|/g, '\\|')
      .replace(/\n/g, '<br />');
  }
}

export default MyReporter;
