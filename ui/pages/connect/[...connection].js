import { useEffect } from 'react';
import Head from 'next/head';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import ConnectionWizardNew from '../../components/Connect/index.js';
import { updatepagepath, updatepagetitle } from '../../lib/store';
import { getPath } from '../../lib/path';
import { NoSsr } from '@material-ui/core';

const ConnectionWizard = (props) => {
  useEffect(() => {
    props.updatepagepath({ path: getPath() });
    props.updatepagetitle({ title: 'Connection Wizard' });
  }, []);

  return (
    <NoSsr>
      <Head>
        {/**
         * CDN's for slick-carousel
         */}

        <link
          rel="stylesheet"
          type="text/css"
          // charSet="UTF-8"
          href="https://cdnjs.cloudflare.com/ajax/libs/slick-carousel/1.6.0/slick.min.css"
        />
        <link
          rel="stylesheet"
          type="text/css"
          href="https://cdnjs.cloudflare.com/ajax/libs/slick-carousel/1.6.0/slick-theme.min.css"
        />

        <title>Connection wizard | Meshery</title>
      </Head>
      <ConnectionWizardNew />
    </NoSsr>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch),
  updatepagetitle: bindActionCreators(updatepagetitle, dispatch),
});

export default connect(null, mapDispatchToProps)(ConnectionWizardNew);
