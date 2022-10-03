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
  - [Execution](#execution)
  - [Directory Structure](#directory-structure)
    - [Tests Directory Structure](#tests-directory-structure)
    - [Goldens Directory Structure](#goldens-directory-structure)
  - [Test Specification Format](#test-specification-format)
    - [IncludeFiles](#includefiles)
    - [IgnoreFields](#ignorefields)

## Usage

There are two available commands within the tool:

- `equivalence-test update --goldens=examples/example_golden_files --tests=examples/example_test_cases --binary=terraform`
- `equivalence-test diff --goldens=examples/example_golden_files --tests=examples/example_test_cases --binary=terraform`

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

## Execution

Each test case executes the following Terraform commands in order:

1. `terraform init`
2. `terraform plan -out=equivalence_test_plan`
3. `terraform apply -json equivalence_test_plan`
4. `terraform show -json`
5. `terraform show -json equivalence_test_plan`

There is currently no way to execute a different set of commands.

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
    - `plan.json`
    - `state.json`
  - `test_case_two/`
      - `apply.json`
      - `plan.json`
      - `state.json`

Note, that if you are writing golden files out for the first time you do not 
need to set up the directory structure yourself. The tool will update and write 
out the directory structure from scratch.

## Test Specification Format

Currently, the test specification has two fields:

- `IncludeFiles`: This field specifies a set of files that should be included as 
                  golden files.
- `IgnoreFields`: This field specifies a map between output files and JSON 
                  fields that should be ignored when reading from or writing to 
                  the golden files.

### IncludeFiles

The `apply.json`, `state.json`, and `plan.json` golden files are included by all
tests automatically.

- The `apply.json` file contains the output of `terraform apply -json equivalence_test_plan`.
- The `state.json` file contains the output of `terraform show -json`.
- The `plan.json` file contains the output of `terraform show -json equivalence_test_plan`.

You can then use this field to specify any additional files that should also be 
considered golden files. **Any additional files must be JSON formatted.**

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