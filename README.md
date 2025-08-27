# ssi_race_condition_reproduction

This is a minimal setup for reproducing a bug in postgres' serialization isolation level
due to a race condition between `SetNewSxactGlobalXmin` and `PredicateLockPageSplit`

## Usage
* Apply `schema.sql` to create the table, index, and initial data.
* Adjust the connection string then run `go run cmd/ssi_stress_test/main.go` to run the test. 
  * A 5 minute run reproduced the bug several times on a build of postgres with a [sleep](https://github.com/joshcurtis/postgres/commit/36d4f447b8e12a3b3dba852ef85dbb8db542425d) added to `SetNewSxactGlobalXmin`.
  * I was able to inconsistently reproduce the bug on the latest build of postgres with a similar test setup, but there 
    were ~tens of hours between occurrences.

