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

const socialLinks = [
  {
    href: 'mailto:maintainers@meshery.io',
    title: 'Get connected with the Meshery community',
    icon: <MessageIcon height={45} width={45} />,
  },
  {
    href: 'https://slack.meshery.io',
    title: 'Join the community Slack',
    icon: <SlackIcon height={45} width={45} />,
  },
  {
    href: 'https://x.com/mesheryio',
    title: 'Follow Meshery on X',
    icon: <TwitterIcon height={40} width={40} />,
  },
  {
    href: 'https://github.com/meshery',
    title: 'Contribute to Meshery projects',
    icon: <GithubIcon height={45} width={45} />,
  },
  {
    href: 'https://www.youtube.com/playlist?list=PL3A-A6hPO2IMPPqVjuzgqNU5xwnFFn3n0',
    title: 'Watch community meeting recordings',
    icon: <YoutubeIcon height={45} width={45} />,
  },
  {
    href: 'https://hub.docker.com/u/meshery/',
    title: 'Access Docker images',
    icon: <DockerIcon height={45} width={45} />,
  },
] as const;

export default function Socials() {
  return (
    <SocialMain>
      <SocialContainer>
        {socialLinks.map(({ href, title, icon }) => (
          <Tooltip
            key={href}
            TransitionComponent={Fade}
            TransitionProps={{ timeout: 600 }}
            title={title}
          >
            <a
              href={href}
              rel="noreferrer"
              target={href.startsWith('mailto:') ? undefined : '_blank'}
            >
              {icon}
            </a>
          </Tooltip>
        ))}
      </SocialContainer>
    </SocialMain>
  );
}
