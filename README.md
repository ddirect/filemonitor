# filemonitor
`filemonitor` monitors a tree of files for modifications and sends information about each file to immudb Vault.

Each file maps to an immudb document in this format:

```
{
  "Id": "<immudb_id>",
  "FileName": "<file_name>",
  "Hash": "<hash_of_file_contents>",
  "ModTime": "<last_modified_time>",
}
```

The hash is SHA256, base64 encoded. The timestamp is in nanoseconds since Unix Epoch, stored as string.

When a new file is found, a new immudb document is created; when a file is modified, the immudb document is updated with the new data.

It is possible to audit the immudb document using `viewer`.

## Building
Use the `./build_all.sh` script from the root of the repository. The binaries will be written to the `bin` directory.

## Synopsis
Run filemonitor on a directory:

`filemonitor -base <directory_path>`

Audit a file:

`viewer -audit -id <file_name>` or `viewer -audit -native -id <immudb_id>`

Additional options for `filemonitor`, `viewer` and `manage` (a tool to perform some operations on immudb) can be seen by running them with the `-help` argument.

## API key
The immudb Vault API key can be either set on the command line or with the `IMMUDB_API_KEY` environment variable. 

## Architecture
To parallelize the workload, `filemonitor` employs a pipelined architecture, using the following workers:

- **queuer**: buffers new file names which need to be processed by the pipeline and ensures a file is present a single time in the following stages
- **checker**: checks if modifications occurred on a file (calling `Lstat`); if the last modification time of the file changed:
	- for regular files it sends them to the hasher
	- for directories it reads the names of the files it contains and sends them to the queuer
- **hasher**: calculates the file hash and passes it to the sender
- **sender**: creates or stores the immudb document

## TODO
- Add tests for checker, hasher and sender
- Improve coverage of existing tests
- Tweak file checking intervals
- Tweak the number of workers
