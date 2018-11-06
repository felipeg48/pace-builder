package resources

var DefaultConfig = `{
    "workshopSubject":"PACE",
    "workshopHomepage":"",
    "modules": [
    {
        "type": "concepts",
        "content": [
            {
            "name":"example-slide",
            "filename":"example/example-slide"
            }
        ]
    },
    {
        "type": "demos",
        "content": [
            {
            "name":"example-demo",
            "filename":"example/example-demo"
            }
        ]
    }
  ]
}`

var DefaultManifest = `---
applications:
- name: my-pace-workshop
  memory: 512M
  instances: 1
  buildpacks: 
  - staticfile_buildpack
  random-route: true
  path: public/`
