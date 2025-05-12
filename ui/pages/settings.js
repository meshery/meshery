import MesherySettings from '../components/MesherySettings';
import { NoSsr } from '@layer5/sistent';
import { updatepagepath } from '../lib/store';
import { connect } from 'react-redux';
import Head from 'next/head';
import { bindActionCreators } from 'redux';
import { getPath } from '../lib/path';
import React, { useEffect } from 'react';
import { useDispatchRtk } from '@/store/hooks';
import { updatePage, updateTitle } from '@/store/slices/mesheryUi';

function Settings(props) {
  const dispatch = useDispatchRtk();
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
    dispatch(updateTitle({ title: 'Settings' }));
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>Settings | Meshery</title>
      </Head>

      <MesherySettings />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default connect(null, mapDispatchToProps)(Settings);
