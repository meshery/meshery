func TestGenerateModel(t *testing.T) {
    cmd := NewGenerateCmd()
    cmd.SetArgs([]string{"sample-model"})
    err := cmd.Execute()

    assert.NoError(t, err, "Generate command should not return an error")
}
