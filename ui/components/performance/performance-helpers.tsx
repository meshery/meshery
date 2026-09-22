import React from 'react';
import { Link } from '@sistent/sistent';
import { generateTestName } from './helper';

/**
 * generatePerformanceProfile takes in data and generate a performance
 * profile object from it
 * @param {*} data
 */
export function generatePerformanceProfile(data) {
  const {
    id,
    name,
    loadGenerator,
    additional_options,
    endpoint,
    serviceMesh,
    concurrentRequest,
    qps,
    duration,
    requestHeaders,
    requestCookies,
    requestBody,
    contentType,
    caCertificate,
  } = data;

  const performanceProfileName = generateTestName(name, serviceMesh);

  return {
    ...(id && { id }),
    name: performanceProfileName,
    loadGenerators: [loadGenerator],
    endpoints: [endpoint],
    serviceMesh,
    concurrentRequest,
    qps,
    duration,
    requestHeaders,
    requestBody,
    requestCookies,
    contentType,
    metadata: {
      additional_options: [additional_options],
      ca_certificate: {
        file: caCertificate.file,
        name: caCertificate.name,
      },
    },
  };
}

export const loadGenerators = ['fortio'];

export const infoFlags = <>Only .json files are supported.</>;

export const infoCRTCertificates = <>Only .crt files are supported.</>;

export const infoloadGenerators = (
  <>
    Which load generators does Meshery support?
    <ul>
      <li>
        fortio - Fortio load testing library, command line tool, advanced echo server and web UI in
        go (golang). Allows to specify a set query-per-second load and record latency histograms and
        other useful stats.{' '}
      </li>
    </ul>
    <Link
      style={{ textDecoration: 'underline' }}
      color="inherit"
      href="https://docs.meshery.io/functionality/performance-management"
    >
      {' '}
      Performance Management
    </Link>
  </>
);
