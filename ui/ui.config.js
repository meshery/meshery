const {
  AnimatedMesheryDark,
} = require('./components/shared/LoadingState/Animations/AnimatedMesheryCSS');
const {
  default: AnimatedMeshery,
} = require('./components/shared/LoadingState/Animations/AnimatedMesheryCSS');

module.exports = {
  components: {
    navigator: true, // set false to disable the navigator component ( default: true )
  },
  AnimatedLogo: AnimatedMeshery,
  AnimatedLogoDark: AnimatedMesheryDark,
};
