# Resume Wizard 🧙

CLI tool for converting YAML files to resume PDF's (Or html/markdown etc).

## CLI Workflows
- Use `build` to create the inital YAML file. (First time)
    - Also could be an `init` command.
- Use `wizard` (or `decorate`) decorate it for a job.
- Continue iterating with `build` until renders a good PDF.
- Use `save` to insert it into a db.
- Optional: `export` to view it as a csv

## V2 
- [ ] Support base resumes
    - [ ] Add
    - [ ] Delete
    - [ ] Edit
- [ ] Tailor workflow
    - [ ] Job info
    - [ ] Call to LLM
    - [ ] Edit/Review
    - [ ] Log event
- [ ] Store in a file store such as seaweedfs (S3 compatible)
        - [ ] For now can just dump them at a configurable directory
- [ ] Read only activity log (sqlite)


## TODO
- [X] Have selected tab be stored in local storage to make refresh easier
- [X] Allow editing of status via the UI
- [ ] Start working on metrics/analytics?

## High Level Design
- [ ] Create a CLI proof of concept
    - [X] Parse YAML/JSON into an IR
    - [X] Convert the IR to HTML
    - [ ] Wizard Command (Magic? Init?) 
        - [ ] What job is this for?
        - [ ] `$EDITOR copy_of_resume.yaml` (User deletes what they don't want
        - [ ] Which template (with default)
        - [ ] Done! New resume is saved. (Maybe update a spreadsheet with which resume used too?
        - [ ] Include a way to prune automatically?
- [X] Move to a web-based application if the complexity gets there

## Roadmap
- [ ] Support for JSON/TOML inputs and HTML/Markdown/PDF Outputs (QoL)
- [ ] Prune YAML to allow tailoring of resume's to job descriptions and website scrapping (Killer Feature)
- [ ] Create YAML by parsing existing resumes (Easy onboarding)
- [ ] Make the resume style customizable or add custom themes
