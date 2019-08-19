import React from 'react';
import { NoSsr } from '@material-ui/core';
import MesheryPlayComponent from '../components/MesheryPlayComponent';


const Manage = ({ query }) => {
    return(
    <NoSsr>
    <React.Fragment>
        <MesheryPlayComponent adapter={query.adapter} />
    </React.Fragment>
    </NoSsr>
    );
}

Manage.getInitialProps = ({query}) => {
    return {query}
  }

export default Manage;