import { readFileSync, writeFileSync } from 'fs';
import { template } from 'lodash';
import moment from 'moment';
import path from 'path';

class MyReporter {
  introMessage = '';
  totalTests = '';
  expectedTest = '';
  testTable = `| Test | Browser | Test Case | Tags | Result |
| :---: | :---: | :--- | :---: | :---: |`;
  passed = 0;
  failed = 0;
  skipped = 0;
  flaky = 0;
  countLog = 1;
  countTable = 1;

  onBegin(_config, suite) {
    this.introMessage = `- Testing started at: ${moment().format('MMMM Do YYYY, h:mm:ss a')}`;
    this.totalTests = `- Total tests cases: ${suite.allTests().length}`;
  }

  // eslint-disable-next-line no-unused-vars
  onStdOut(chunk, _test, _result) {
    const text = chunk.toString('utf-8');
    process.stdout.write(text);
  }

  // eslint-disable-next-line no-unused-vars
  onStdErr(chunk, _test, _result) {
    const text = chunk.toString('utf-8');
    process.stderr.write(text);
  }

  onTestEnd(test, result) {
    const status = test.outcome();
    const project = test.parent.project().name;

    this.displayLogs(project, test.title, test.tags, status, result);
    this.addTestTable(project, test.title, test.tags, status, result.retry, test.retries);
  }

  async onEnd(result) {
    const message = await this.buildMessage(result);

    try {
      writeFileSync('test-report.md', message);
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
        ? '\n' +
            `File Location: ${result.error.location?.file ?? 'Not Found'}` +
            '\n' +
            result.error?.snippet ?? 'No snippet' + '\n' + result.error?.message
        : ''
    }\n`;

    process.stdout.write(logs);

    this.countLog++;
  }

  addTestTable(project, title, tags, status, retry, retries) {
    this.countTestStatus(tags, status, retry, retries);

    if (status === 'expected') return;

    const lastRetriesRun = retry === retries;
    const isFail = status === 'unexpected';
    const isSkipped = status === 'skipped';

    if ((isFail || isSkipped) && !lastRetriesRun) {
      return;
    }

    const allTags = tags.map((item) => item.replace('@', '')).join(', ');

    const message = `| ${
      this.countTable
    } | ${project} | ${title} | ${allTags} | ${this.getStatusEmoji(tags, status)} |`;

    this.testTable += `\n${message}`;
    this.countTable++;
  }

  countTestStatus(tags, status, retry, retries) {
    const lastRetriesRun = retry === retries;
    const isFail = status === 'unexpected';
    const isSkipped = status === 'skipped';
    const isFlaky = status === 'flaky';

    const isUnstableTest = tags.includes('@unstable');

    if (status === 'expected') {
      this.passed++;
    }
    if (isFlaky || (isUnstableTest && isFail)) {
      this.flaky++;
    }
    if (isFail && lastRetriesRun && !isUnstableTest) {
      this.failed++;
    }
    if (isSkipped && lastRetriesRun && !isUnstableTest) {
      this.skipped++;
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
      testTable: this.testTable,
    });
  }
}

export default MyReporter;
