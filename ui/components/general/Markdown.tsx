import '@uiw/react-markdown-preview/markdown.css';
import '@uiw/react-md-editor/markdown-editor.css';

import dynamic from 'next/dynamic';
export const MDEditor = dynamic(() => import('@uiw/react-md-editor'), { ssr: false });
