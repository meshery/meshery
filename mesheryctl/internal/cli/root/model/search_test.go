func TestSearchModel(t *testing.T) {
    cmd := NewSearchCmd()
    cmd.SetArgs([]string{"sample-model"})
    err := cmd.Execute()

    assert.NoError(t, err, "Search command should not return an error")
}
