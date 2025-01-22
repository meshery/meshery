export const getLegendTemplate = (title, color, data) => {
  const nodeCount = data.find((d) => d.id === title)?.value || 0;
  return `<div style='display:flex;flex-direction:column;gap:2px;align-items:center;'>
      <span style='color:${color};font-size:1.25rem;'>${nodeCount}</span>
      <span style='color:#fff;font-size:smaller;'>${title}</span>
    </div>`;
};
