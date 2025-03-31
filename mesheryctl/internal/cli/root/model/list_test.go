func TestListModels(t *testing.T) {
    cmd := NewListCmd()
    err := cmd.Execute()

    assert.NoError(t, err, "List command should not return an error")
}
