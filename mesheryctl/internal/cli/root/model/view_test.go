func TestViewModel(t *testing.T) {
    cmd := NewViewCmd()
    cmd.SetArgs([]string{"sample-model"})
    err := cmd.Execute()

    assert.NoError(t, err, "View command should not return an error")
}
