import { NoSsr } from "@material-ui/core";
import Head from 'next/head';
import dynamic from 'next/dynamic';
import { getPath } from "../lib/path";

class Performance extends React.Component {
  componentDidMount () {
    console.log(`path: ${getPath()}`);
    // this.props.updatepagepath({ path: getPath() });
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

export default Performance;