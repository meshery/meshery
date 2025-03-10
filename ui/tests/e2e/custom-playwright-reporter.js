import { readFileSync, writeFileSync } from 'fs';
import { template } from 'lodash';
import moment from 'moment';
import path from 'path';

class MyReporter {
  introMessage = '';
  totalTests = '';
  testData = [];
  passed = 0;
  failed = 0;
  skipped = 0;
  flaky = 0;
  countLog = 1;

  constructor() {
    this.countTestStatus = this.countTestStatus.bind(this);
  }

  onBegin(_config, suite) {
    this.introMessage = `- Testing started at: ${moment().format('MMMM Do YYYY, h:mm:ss a')}`;
    this.totalTests = `- Total test cases: ${suite.allTests().length}`;
  }

  onStdOut(chunk) {
    const text = chunk.toString('utf-8');
    process.stdout.write(text);
  }

  onStdErr(chunk) {
    const text = chunk.toString('utf-8');
    process.stderr.write(text);
  }

  onTestEnd(test, result) {
    const status = test.outcome();
    const project = test.parent?.project()?.name;
    const spec = test.parent?.title;

    this.displayLogs(project, test.title, test.tags, status, result);
    this.addTestData(project, spec, test.title, test.tags, status, result.retry, test.retries);
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
        ? ('\n' +
            `File Location: ${result.error.location?.file ?? 'Not Found'}` +
            '\n' +
            result.error?.snippet ?? 'No snippet' + '\n' + result.error?.message)
        : ''
    }\n`;

    process.stdout.write(logs);

    this.countLog++;
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

  addTestData(project, spec, title, tags, status, retry, retries) {
    this.countTestStatus(tags, status, retry, retries);

    if (status === 'expected') return;

    const lastRetriesRun = retry === retries;
    const isFail = status === 'unexpected';
    const isSkipped = status === 'skipped';

    if ((isFail || isSkipped) && !lastRetriesRun) {
      return;
    }

    const allTags = tags.map((item) => item.replace('@', '')).join(', ');

    this.testData.push({
      project,
      spec,
      title,
      tags: allTags,
      status: this.getStatusEmoji(tags, status),
    });
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
      testData: this.testData,
    });
  }
}

export default MyReporter;
