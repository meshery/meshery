func TestImportModel(t *testing.T) {
    cmd := NewImportCmd()
    cmd.SetArgs([]string{"-f", "http://example.com/model"})
    err := cmd.Execute()

    assert.NoError(t, err, "Import command should not return an error")
}
