import React, { useEffect } from 'react';
import MesherySMIResults from '../components/MesherySMIResults';
import { updatepagepath } from '../lib/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../lib/path';

/**
 * @deprecated This functionality has been deprecated and its child components can be left behind
 */
const SMIResults = (props) => {
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
  }, [props.updatepagepath]);

  return (
    <React.Fragment>
      <Head>
        <title>SMI Results | Meshery</title>
      </Head>
      <MesherySMIResults />
    </React.Fragment>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default connect(null, mapDispatchToProps)(SMIResults);
