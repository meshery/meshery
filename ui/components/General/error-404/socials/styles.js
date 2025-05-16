import { styled } from '@layer5/sistent';
import Message from '@/assets/icons/general/message';
import Docker from '@/assets/icons/social/docker';
import Github from '@/assets/icons/social/github';
import Slack from '@/assets/icons/social/slack';
import Twitter from '@/assets/icons/social/twitter';
import Youtube from '@/assets/icons/social/youtube';

export const SocialMain = styled('div')(() => ({
  padding: '1.3rem 2rem',
  display: 'flex',
  justifyContent: 'center',
}));

export const SocialContainer = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
}));

export const MessageIcon = styled(Message)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: theme.palette.icon.default,
    color: theme.palette.icon.default,
  },
}));

export const SlackIcon = styled(Slack)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: '#4A154B',
    color: '#4A154B',
  },
}));

export const TwitterIcon = styled(Twitter)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: '#1da1f2',
    color: '#1da1f2',
  },
}));

export const TwitterHandleIcon = styled(Twitter)(({ theme }) => ({
  fill: theme.palette.keppelGreen,
  color: theme.palette.keppelGreen,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: '#1da1f2',
    color: '#1da1f2',
  },
}));

export const GithubHandleIcon = styled(Github)(({ theme }) => ({
  fill: theme.palette.keppelGreen,
  color: theme.palette.keppelGreen,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: '#24292e',
    color: '#24292e',
  },
}));

export const GithubIcon = styled(Github)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: '#24292e',
    color: '#24292e',
  },
}));

export const YoutubeIcon = styled(Youtube)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: '#ff0000',
    color: '#ff0000',
  },
}));

export const DockerIcon = styled(Docker)(({ theme }) => ({
  fill: theme.palette.icon.default,
  color: theme.palette.icon.default,
  cursor: 'pointer',
  transition: 'all .3s',
  '&:hover': {
    fill: '#2496ed',
    color: '#2496ed',
  },
}));
