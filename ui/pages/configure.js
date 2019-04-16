import { NoSsr } from "@material-ui/core";
import MeshConfigComponent from "../components/MeshConfigComponent";
import MeshAdapterConfigComponent from "../components/MeshAdapterConfigComponent";

const Config = () => (
  <NoSsr>
    <MeshConfigComponent />
    <MeshAdapterConfigComponent />
  </NoSsr>
)

export default Config;