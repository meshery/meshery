import PascalCaseToKebab from "../../utils/PascalCaseToKebab";

describe('PascalKebab String Utility', () => {
  test('make pascalcase to kebabcase', () => {
    const val = PascalCaseToKebab("TotalK8sConnections");
    expect(val).toBe("total-k8s-connections")
  })
})
