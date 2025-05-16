# Steps to build your own WASM binaries and use it in Meshery-UI or Meshery-Extensions

--- 

1. Write the rego file and its policy
2. Follow the hierarchy of defining policies, your polcies should be written in directory which should have same name as your policy file
3. install opa-cli tool
4. cd into the directory where you've written your policy
5. build-opa bundle with wasm target:
```
opa build -t wasm -e <package_name_of_rego_policy> <file_name_of_rego_policy>
```
6. There must be bundle created for you
7. cd into [temp-bundle-extract-dir](./policies/temp-bundle-extract-dir/)
8. extract the tar-file with
```
tar -xzf ../namespace_discovery_relationship_policy/bundle.tar.gz # example command
```
9. rename the wasm-binary
```
mv policy.wasm namespace_discovery_relationship_policy.wasm
```
10. move the new wasm binary to [ui/public/static/](../ui/public/static/)
11. reference the policy file from UI

## Please don't push the contents of  [temp-bundle-extract-dir](./policies/temp-bundle-extract-dir/)