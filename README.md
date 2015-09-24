oku
=======

Command line tool which detects the encoding of a file and outputs to UTF-8.

### Usage examples ###

```bash
# detect file encoding only
oku -d my_file.txt

# detect current encoding and convert to utf-8
oku -o outfile.txt my_file.txt
# or
oku my_file.txt > outfile.txt

# read from stdin and convert to utf-8
... | oku | less
```

### Building ###

```bash
go install ./...
```

