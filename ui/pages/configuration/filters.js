import React, { useEffect } from 'react';
import { NoSsr } from '@mui/material';
import MesheryFilters from '../../components/Filters';
import { updatepagepath } from '../../lib/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { Box } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

function NewFilters(props) {
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>Filters | Meshery</title>
      </Head>
      <UsesSistent>
        <Box
          sx={{
            maxWidth: '90%',
            margin: 'auto',
            overflow: 'hidden',
          }}
        >
          <MesheryFilters />
        </Box>
      </UsesSistent>
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
});

export default connect(null, mapDispatchToProps)(NewFilters);
