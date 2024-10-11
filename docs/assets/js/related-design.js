---
layout: blank
excluded_in_search: true
---

document.addEventListener("DOMContentLoaded", getRelatedDesigns);

async function getRelatedDesigns() {
    const relatedDesignsSection = document.querySelector(".related-designs-links");
    relatedDesignsSection.innerHTML = '';
    const docURL = relatedDesignsSection.getAttribute('data-doc-url');
    const docurl = docURL.split("/");
    const tech = docurl[docurl.length - 1];
    //console.log("techonolgy", `${tech}`);
    const url = `https://meshery.layer5.io/api/catalog/content/pattern?technology=` + `${tech}` + `&page=0&pagesize=5&search=&order=&metrics=true`; // to-do
  
    try {
      const res = await fetch(url);
      const response = await res.json(); 
      const patterns = response.patterns;

      patterns.map((pattern) => {
        const link = document.createElement("a");
        link.href = `https://meshery.layer5.io/catalog/content/catalog/${pattern.id}`;
        link.textContent = `Design ${pattern.name}`; 

        relatedDesignsSection.appendChild(link);
  
        relatedDesignsSection.appendChild(document.createElement("br"));
      });
    } catch (er) {
      console.error("Error fetching / processing patterns:", er);
    }
  }
  
  
  