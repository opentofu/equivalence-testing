# terraform-equivalence-testing

The `terraform-equivalence-testing` repository provides a tool for comparing and
updating state files, plan files, and the JSON output of the `apply` command, produced by Terraform 
executions.

The framework uses a set of golden files to track outputs and verify changes
across different Terraform versions, provider versions, or even different 
Terraform configurations.

## Contents

- [terraform-equivalence-testing](#terraform-equivalence-testing)
  - [Usage](#usage)
    - [Optional Flags](#optional-flags)
  - [Execution](#execution)
  - [Directory Structure](#directory-structure)
    - [Tests Directory Structure](#tests-directory-structure)
    - [Goldens Directory Structure](#goldens-directory-structure)
  - [Test Specification Format](#test-specification-format)
    - [IncludeFiles](#includefiles)
    - [IgnoreFields](#ignorefields)
    - [Commands](#commands)

## Usage

There are two available commands within the tool:

- `./terraform-equivalence-testing update --goldens=examples/example_golden_files --tests=examples/example_test_cases`
- `./terraform-equivalence-testing diff --goldens=examples/example_golden_files --tests=examples/example_test_cases`

The first command will iterate through the test cases in 
`examples/example_test_cases`, run a set of Terraform commands while collecting
the Terraform output for these commands, and then write the outputs into a
directory within `examples/example_golden_files`. This command will overwrite 
any existing golden files that already exist.

The second command does the same as the first command, except instead of 
updating or overwriting the golden files it simply reports on any differences
found between the existing golden files and the outputs of the Terraform 
commands.

The above commands, when executed from the root of this repository, should be
successful using the examples provided in the `examples/` directory.

### Optional Flags

1. `--binary=terraform`
    - By default, the equivalence tests will look for the first binary named
      `terraform` within the path. 
    - This flag can be set to modify which Terraform binary is used to execute 
      these tests. 
2. `--filters=simple_resource,complex_resource`
    - By default, the equivalence tests will execute all the tests within the 
      specified `--tests` directory.
    - You can specify a subset of the tests to execute using this flag either by
      repeating the flag (eg. 
     `--filters=simple_resource --filters=complex_resource`), or with a comma
      separated list as in the original example.

## Execution

Each test case executes the following Terraform commands in order:

1. `terraform init`
2. `terraform plan -out=equivalence_test_plan`
3. `terraform apply -json equivalence_test_plan`
4. `terraform show`
5. `terraform show -json`
6. `terraform show -json equivalence_test_plan`

Consult the [Test Specification Format](#test-specification-format) section for
a run down on how to customise these commands using the `Commands` 
specification.

## Directory Structure

The tool reads in from and writes out to an expected directory structure. 

### Tests Directory Structure

The `--tests` flag specifies the input directory for the test cases.

Within the target directory there should be a set of subdirectories, with each 
subdirectory containing a single test case. Each test case is made up of a 
`spec.json` file, providing any customisations for the test, and then a set of
`.tf` Terraform files. The tool uses the name of each subdirectory to name the 
test case in any logs or output it produces.

Example input directory structure:

- `my_test_cases/`
  - `test_case_one/`
    - `spec.json`
    - `main.tf`
  - `test_case_two/`
    - `spec.json`
    - `main.tf`

### Goldens Directory Structure

The `--goldens` flag specifies the directory where the golden files should be
read from, when diffing, or written to, when updating.

The tool will write the golden files for a given test case into a subdirectory
using a name that matches the subdirectory in the input directory. You can use 
the subdirectory names to map between the input test cases and the output golden
files.

Example golden directory structure:

- `my_golden_files/`
  - `test_case_one/`
    - `apply.json`
    - `plan`
    - `plan.json`
    - `state.json`
  - `test_case_two/`
      - `apply.json`
      - `plan`
      - `plan.json`
      - `state.json`

Note, that if you are writing golden files out for the first time you do not 
need to set up the directory structure yourself. The tool will update and write 
out the directory structure from scratch.

## Test Specification Format

Currently, the test specification has three fields:

- `IncludeFiles`: This field specifies a set of files that should be included as 
                  golden files.
- `IgnoreFields`: This field specifies a map between output files and JSON 
                  fields that should be ignored when reading from or writing to 
                  the golden files.
- `Commands`: This field specifies a list of custom commands that should be 
              executed instead of the default set of commands.

### IncludeFiles

The `apply.json`, `state.json`, `plan.json`, and `plan`, golden files are 
included by all tests automatically.

- The `apply.json` file contains the output of 
  `terraform apply -json equivalence_test_plan`.
- The `state.json` file contains the output of `terraform show -json`.
- The `plan.json` file contains the output of
  `terraform show -json equivalence_test_plan`.
- The `plan` file contains the raw human-readable captured output of the 
  original `terraform plan` command.

You can then use this field to specify any additional files that should also be 
considered golden files.

### IgnoreFields

The following fields are ignored by default:

- In `apply.json`:
  - `0`: This is the first entry in the JSON list that comprises `apply.json`.
         It contains lots of execution specific information such as timing and 
         Terraform version which will change on every execution.
  - `*.@timestamp`: This removes the `@timestamp` field from every entry in the 
                    `apply.json` as the timestamp will change on every 
                    execution.
- In `state.json`:
  - `terraform_version`: The removes the Terraform version information from the 
                         state as it will create noise in our golden file diffs.
- In `plan.json`:
  - `terraform_version`: The removes the Terraform version information from the
                         plan as it will create noise in our golden file diffs.

If you need any other fields removed, either from the default golden files or
additional golden files, then you can specify them here as part of the test
specification.

Note, that you can only remove fields from JSON files. Other file types will not
be included when processing the `IgnoreFields` inputs.

### Commands

You can specify a custom list of terraform commands to execute instead of the 
default set specified in [Execution](#execution).

Each command has 5 required fields:
  
- `name`
- `arguments`
- `capture_output`
- `output_file_name`
- `has_json_output`
- `streams_json_output`

`name` (**required**) is a string only used for logging when reporting which 
commands might have failed, so you should make it unique and descriptive enough
that it can identify which part of the test failed when consulting the error 
log.

`arguments` (**required**) is a list of arguments that should be passed into the
Terraform binary for this command. For example, `[plan, -out=plan_output]` would
tell Terraform to perform a plan action and where to save the plan file.

`capture_output` (**optional**, defaults to `false`) is a boolean that tells the
equivalence tests to capture and save the output of this command as a golden 
file for diffing or updating.

`output_file_name` (**required** if `capture_output` is `true`) is a string 
that sets the filename that should be used for the output. If `capture_output` 
is `false`, this field is ignored.

`has_json_output` (**optional**, defaults to `false`) is a boolean that tells 
the equivalence tests that the output of this command will be in JSON format.
The framework will only use the `IgnoreFields` specification on JSON formatted
files so if you wish to remove any part of the output this must be true.

`streams_json_output` (**optional**, defaults to `false`) is a boolean
that tells the equivalence tests that the output is in the "structured JSON" 
format. Some Terraform commands, such as `terraform apply -json`, stream a list
of individual JSON objects to the output. This form of output is not a valid
JSON object when reading the output as a whole. When this value is true the 
framework will convert the output into a valid JSON object by replacing any `\n`
characters with `,` and putting the entire output in between `[` and `]`. If
`capture_output` or `has_json_output` is `false`, this field is ignored.

#### Examples

The following example demonstrates how to replicate the default commands using 
the custom `commands` entry in the test specification.

```json
{
  "commands": [
    {
      "name": "init",
      "arguments": ["init"],
      "capture_output": false
    },
    {
      "name": "plan",
      "arguments": ["plan", "-out=equivalence_test_plan", "-no-color"],
      "capture_output": true,
      "output_file_name": "plan",
      "has_json_output": false
    },
    {
      "name": "apply",
      "arguments": ["apply", "-json", "equivalence_test_plan"],
      "capture_output": true,
      "output_file_name": "apply.json",
      "has_json_output": true,
      "streams_json_output": true
    },
    {
      "name": "state",
      "arguments": ["show", "-no-color"],
      "capture_output": true,
      "output_file_name": "state",
      "has_json_output": false
    },
    {
      "name": "show_state",
      "arguments": ["show", "-json"],
      "capture_output": true,
      "output_file_name": "state.json",
      "has_json_output": true,
      "streams_json_output": false
    },
    {
      "name": "show_plan",
      "arguments": ["show", "-json", "equivalence_test_plan"],
      "capture_output": true,
      "output_file_name": "plan.json",
      "has_json_output": true,
      "streams_json_output": false
    }
  ]
}
```
