import React from 'react';
import CustomErrorMessage from '../components/general/ErrorPage';
import { MesheryPage } from '../components/general/MesheryPage';

const Error = () => (
  <MesheryPage title="Error" headTitle="404 - Page Not Found" noSuffix>
    <CustomErrorMessage />
  </MesheryPage>
);

export default Error;
