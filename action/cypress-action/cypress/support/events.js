Cypress.on("uncaught:exception", (err, runnable) => {
  console.error("Uncaught Application Error :", err, runnable);
  // returning false here prevents Cypress from
  // failing the test
  return false;
});
