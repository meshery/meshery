import os
import json
import sys
import pathlib

#recursively merge two dictionaries
def merge(dict1, dict2):
    res = {**(dict1 or {}), **dict2}
    for key, value in res.items():
        if key in dict1 and key in dict2:
            if isinstance(value, dict):
                res[key] = merge(dict1[key], dict2[key])
    return res



shape_capabilities =  {
      "designer": {
        "edit": {
          "shape": {
            "convert-shape": True
          },
          "style": True,
          "config": False,
          "lock": True
        },
        "label": {
          "edit": True,
          "sync-with-config-property": "label",
          "show": True
        }
      }
}

DEFAULT_CAPABILITIES = {
    "SHAPE" : shape_capabilities
}


def add_capabilities(component_path:pathlib.Path, capabilities):
    print("Adding capabilities to " + component_path.name)

    data = component_path.read_text(encoding='utf-8')
    if data == '':
        print('File is empty, skipping')
        return
    data_json = json.loads(data)
    data_json['metadata']['capabilities'] = merge(data_json.get('capabilities') or {}, capabilities)
    component_path.write_text(json.dumps(data_json), encoding='utf-8')



def add_capabilities_all_files_in_dir(path , capabilities):
    # get all files in path
    files = pathlib.Path(path).glob('**/*')
    print("files",files)
    for file in files:
        if file.is_file() and file.suffix == '.json':
            add_capabilities(file, capabilities)


ACTION = {
  "add-to-model": "add capability to all components in model",
  "add-to-component": "add capability to a component in component folder",
}



def main(argv):
    if len(argv) < 3:
        print('Usage: python add-capabilities.py <action>  <path to component folder> <capabilities json string or key>')
        sys.exit(2)

    if argv[1] not in ACTION:
        print('Invalid action: ' + argv[1])
        print('Valid actions:')
        for key, value in ACTION.items():
            print(key + ': ' + value)
        sys.exit(2)

    capabilities = DEFAULT_CAPABILITIES.get(argv[3]) or json.loads(argv[3])

    action = argv[1]
    if action == 'add-to-model':
        path = pathlib.Path(argv[2])
        add_capabilities_all_files_in_dir(path, capabilities )
        return

    if action == 'add-to-component':
        path = pathlib.Path(argv[2])
        add_capabilities(path, capabilities)
        return




if __name__ == "__main__":
    main(sys.argv)