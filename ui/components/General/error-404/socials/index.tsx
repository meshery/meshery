import { Tooltip, useTheme } from '@sistent/sistent';
import {
  DockerIcon,
  GithubIcon,
  MessageIcon,
  SlackIcon,
  SocialContainer,
  SocialMain,
  TwitterIcon,
  YoutubeIcon,
} from './styles';

export default function Socials() {
  const theme = useTheme();
  return (
    <SocialMain>
      <SocialContainer>
        <Tooltip title="Get connected with the Meshery community">
          <a href="mailto:maintainers@meshery.io">
            <MessageIcon height={45} width={45} className="" />
          </a>
        </Tooltip>

        <Tooltip title="Join the community Slack">
          <a href="https://slack.meshery.io">
            <SlackIcon height={45} width={45} className="" />
          </a>
        </Tooltip>

        <Tooltip title="Follow Meshery on X">
          <a href="https://x.com/mesheryio">
            <TwitterIcon height={40} width={40} fill={theme.palette.icon.default} className="" />
          </a>
        </Tooltip>

        <Tooltip title="Contribute to Meshery projects">
          <a href="https://github.com/meshery">
            <GithubIcon height={45} width={45} className="" />
          </a>
        </Tooltip>

        <Tooltip title="Watch community meeting recordings">
          <a href="https://www.youtube.com/playlist?list=PL3A-A6hPO2IMPPqVjuzgqNU5xwnFFn3n0">
            <YoutubeIcon height={45} width={45} className="" />
          </a>
        </Tooltip>

        <Tooltip title="Access Docker images">
          <a href="https://hub.docker.com/u/meshery/">
            <DockerIcon height={45} width={45} className="" />
          </a>
        </Tooltip>
      </SocialContainer>
    </SocialMain>
  );
}
