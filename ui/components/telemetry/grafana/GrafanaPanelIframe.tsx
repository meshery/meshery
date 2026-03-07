interface GrafanaPanelIframeProps {
  src: string;
  title: string;
}

const GrafanaPanelIframe = ({ src, title }: GrafanaPanelIframeProps) => (
  <iframe src={src} title={title} loading="lazy" />
);

export default GrafanaPanelIframe;
