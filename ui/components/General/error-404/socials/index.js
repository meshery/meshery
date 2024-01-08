import { Tooltip, Fade } from '@material-ui/core';
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
          title="Get connected with the Layer5 community"
        >
          <a href="mailto:community@layer5.io">
            <MessageIcon height={45} width={45} />
          </a>
        </Tooltip>

        <Tooltip
          TransitionComponent={Fade}
          TransitionProps={{ timeout: 600 }}
          title="Join the community Slack"
        >
          <a href="https://slack.layer5.io">
            <SlackIcon height={45} width={45} />
          </a>
        </Tooltip>

        <Tooltip
          TransitionComponent={Fade}
          TransitionProps={{ timeout: 600 }}
          title="Follow Layer5 on Twitter"
        >
          <a href="https://twitter.com/layer5">
            <TwitterIcon height={40} width={40} />
          </a>
        </Tooltip>

        <Tooltip
          TransitionComponent={Fade}
          TransitionProps={{ timeout: 600 }}
          title="Contribute to Layer5 projects"
        >
          <a href="https://github.com/layer5io">
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
          <a href="https://hub.docker.com/u/layer5/">
            <DockerIcon height={45} width={45} />
          </a>
        </Tooltip>
      </SocialContainer>
    </SocialMain>
  );
}
