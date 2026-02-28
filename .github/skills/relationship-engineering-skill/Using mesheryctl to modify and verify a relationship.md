### **1\. Export the Model**

Export the model as an OCI tar file:  
bash  
mesheryctl model export \[model-name\] \-o oci \-l \[output-directory\]

This creates a model-name.tar file containing the OCI artifact.

### **2\. Extract the OCI Tar File**

The OCI tar file contains a structured directory format. Extract it to access the model files:  
bash  
mkdir extracted-model  
tar \-xf model-name.tar \-C extracted-model/

### **3\. Navigate to the Model Directory**

The extracted structure follows this pattern:  
Code  
{model-name}/  
  v{schema-version}/  
    {model-version}/  
      model.json (or model.yaml)  
      components/  
        \*.json (or \*.yaml)  
      relationships/  
        \*.json (or \*.yaml)

For example: istio/v1.0.0/1.25.0/

### **4\. Modify Relationship Definitions**

Navigate to the relationships directory and edit the relationship definition files:  
bash  
cd extracted-model/{model-name}/v1.0.0/{version}/relationships/

Edit the relationship JSON/YAML files as needed. The files are written using WriteRelationshipDefinition() method (see server/handlers/component\_handler.go).

### **5\. Rebuild the OCI Image**

Use mesheryctl model build to repackage your modifications:  
bash  
\# Navigate back to the directory containing the model folder  
cd path/to/parent-directory

\# Build the OCI image  
mesheryctl model build {model-name}/{model-version} \--path .

This uses MeshKit's BuildImage() function to create a new OCI-compliant tar file.

### **6\. Import the Modified Model**

Import the rebuilt model back into Meshery:  
bash  
mesheryctl model import \-f {model-name}-{version}.tar

