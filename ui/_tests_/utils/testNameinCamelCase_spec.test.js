import { CamelCaseToSentenceCase } from "../../utils/camelCaseToSentanceCase";

describe('CamelCase String Utility', () => {
  test('make camelCase to sentansecase', () => {
    const val = CamelCaseToSentenceCase("mesheryComponents");
    expect(val).toBe("meshery  Components")
  })
})
