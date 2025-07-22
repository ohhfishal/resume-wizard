#!/bin/bash

pandoc --standalone --css resume.css --output resume.html resume.md

# NOTE: I found it looks better if we just use the output

