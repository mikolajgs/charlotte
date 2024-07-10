# charlotte

### Job
Simple app that takes a YAML file that contains a set of steps which are bash scripts, wrapped into simple logic, inputs, outputs etc.
and executes it.

#### Running test suite

    make test

#### Building `job` binary

    cd cmd/job
    go build .

#### Running `job`

    # TODO: Gotta prepare the YAML file (see _test.go files for now)
    #cd cmd/job
    #./job run -f ../../sample-files/job.yaml

#### v0.1

- [x] pipe stdout and stderr to files
- [x] environment (global and in-step)
- [x] variables
- [x] job inputs
- [x] step outputs
- [x] `continue_on_error`
- [x] values using golang templates
- [x] `if` - conditional steps (value templated, must equal to string `'true'`)
- [ ] running step(s) on success
- [ ] running step(s) on failure
- [ ] running step(s) always
- [ ] job outputs (json or write to files in a specific directory)
- [ ] prepare sample yaml file
- [ ] add building docker image (ko?) and pushing to our registry

#### v0.2
- [ ] validation
- [ ] extract steps so that they can be included (include file with inputs) + proper validation for that
- [ ] ...

### Pipeline
Layer on top of Jobs.
