# Resume Wizard 🧙

CLI tool for converting YAML files to resume PDF's (Or html/markdown etc).

## Design
- [ ] Create a CLI proof of concept
    - [X] Parse YAML/JSON into an IR
    - [ ] Convert the IR to HTML
    - [ ] Wizard Command (Magic? Init?) 
        - [ ] What job is this for?
        - [ ] `$EDITOR copy_of_resume.yaml` (User deletes what they don't want
        - [ ] Which template (with default)
        - [ ] Done! New resume is saved. (Maybe update a spreadsheet with which resume used too?

- [ ] Move to a web-based application if the complexity gets there

## Roadmap
- [ ] Support for JSON/TOML inputs and HTML/Markdown/PDF Outputs (QoL)
- [ ] Way to use secrets. Maybe outline.yaml and secret.yaml? With secret last thing added
- [ ] Prune YAML to allow tailoring of resume's to job descriptions and website scrapping (Killer Feature)
- [ ] Create YAML by parsing existing resumes (Easy onboarding)
