import { NoSsr } from "@material-ui/core";
import Head from 'next/head';
import dynamic from 'next/dynamic';
import { getPath } from "../lib/path";
import { connect } from "react-redux";
import { updatepagepath } from "../lib/store";
import { bindActionCreators } from 'redux';

class Performance extends React.Component {
  constructor(props){
    super(props);
  }

  componentDidMount () {
    console.log(`path: ${getPath()}`);
    this.props.updatepagepath({ path: getPath() });
  }

  render () {
    const MesheryVisualizeComponentNoSsr = dynamic(
      () => import('../components/visualize/MesheryVisualizeComponent'),
      { ssr: false }
    )

    return (
      <NoSsr>
        <Head>
          <title>MeshMap | Meshery</title>
        </Head>
        <NoSsr>
          <MesheryVisualizeComponentNoSsr />
        </NoSsr>
      </NoSsr>
    );
  }
}

const mapDispatchToProps = dispatch => ({
  updatepagepath: bindActionCreators(updatepagepath, dispatch)
})

export default connect(
  null,
  mapDispatchToProps
)(Performance);
