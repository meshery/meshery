import { NoSsr } from '@material-ui/core';
import React from 'react';
import ProviderComponent from '../components/ProviderComponent';

class Index extends React.Component {
  componentDidMount() {
  }

  render() {
    return (
      <NoSsr>
        <ProviderComponent />
      </NoSsr>
    );
  }
}


export default Index;
