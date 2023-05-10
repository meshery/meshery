import { CamelCaseToSentanceCase } from "../../utils/camelCaseToSentanceCase";

describe('CamelCase String Utility', () => {
  test('make camelCase to sentansecase', () => {
    const val = CamelCaseToSentanceCase("mesheryComponents");
    expect(val).toBe("meshery  Components")
  })
})