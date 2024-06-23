# streamline

#### Latest state
It's been copied from a different repository so it's probably failing to
compile.

Docker bit is under development so no idea what's there TBH.

#### Building
Use the following:

    cd cmd/ops-run
    go build .

#### Testing
Use the following:

    ./ops-run run-job -f ../../sample-files/job.yaml


#### Phase 0
Phase 0 is to get something running simple bash scripts, with some simple
dependencies between steps, some logic, piping standard output and error.
And all that should be possible to execute locally, in a specific
container or on a specified kubernetes cluster.

It should be possible to create job specification in YAML.

Once each phase is done, we shall ensure it's all covered with tests, even
if we have to write 50 or so of these guys.


## TODO

### Phase 0.1 - bash script blocks only
#### 0.1.4 pipe stdout and stderr properly
#### 0.1.5 add outputs
#### 0.1.6 add inputs

### Phase 0.2 - dependencies and handling logic
...
### Phase 0.3 - master <-?-> master (worker) <-> worker
...
### Phase 0.4 - logging
