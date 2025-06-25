import { Tooltip, Fade } from '@sistent/sistent';
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
  return (
    <SocialMain>
      <SocialContainer>
        <Tooltip
          TransitionComponent={Fade}
          TransitionProps={{ timeout: 600 }}
          title="Get connected with the Meshery community"
        >
          <a href="mailto:maintainers@meshery.io">
            <MessageIcon height={45} width={45} />
          </a>
        </Tooltip>

        <Tooltip
          TransitionComponent={Fade}
          TransitionProps={{ timeout: 600 }}
          title="Join the community Slack"
        >
          <a href="https://slack.meshery.io">
            <SlackIcon height={45} width={45} />
          </a>
        </Tooltip>

        <Tooltip
          TransitionComponent={Fade}
          TransitionProps={{ timeout: 600 }}
          title="Follow Meshery on Twitter"
        >
          <a href="https://twitter.com/mesheryio">
            <TwitterIcon height={40} width={40} />
          </a>
        </Tooltip>

        <Tooltip
          TransitionComponent={Fade}
          TransitionProps={{ timeout: 600 }}
          title="Contribute to Meshery projects"
        >
          <a href="https://github.com/meshery">
            <GithubIcon height={45} width={45} />
          </a>
        </Tooltip>

        <Tooltip
          TransitionComponent={Fade}
          TransitionProps={{ timeout: 600 }}
          title="Watch community meeting recordings"
        >
          <a href="https://www.youtube.com/playlist?list=PL3A-A6hPO2IMPPqVjuzgqNU5xwnFFn3n0">
            <YoutubeIcon height={45} width={45} />
          </a>
        </Tooltip>

        <Tooltip
          TransitionComponent={Fade}
          TransitionProps={{ timeout: 600 }}
          title="Access Docker images"
        >
          <a href="https://hub.docker.com/u/meshery/">
            <DockerIcon height={45} width={45} />
          </a>
        </Tooltip>
      </SocialContainer>
    </SocialMain>
  );
}
