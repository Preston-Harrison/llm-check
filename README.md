# LLM Check

The purpose of this project is to allow programmers to make function documentation
enforcible. This is done by finding the call sites of functions with a specific tag,
and using a LLM to check whether the call site satisfies the requirements outlined
in the function documentation. 

Imagine the function `openFileForWriting` was written a while ago by
a developer who left the company. Later, another developer calls this function
without noticing that there is some implcit requirement that they must fulfill.

This project hopes to introduce a solution to this by running a check whenever
a call site of a documented function is edited. This check sends the function
body, call site, and requirement to an LLM and prompts it to check the requirement
is fulfilled.

For example:

```go
// sample/sample.go

// llm-check: make sure to close this file after use.
func openFileForWriting() (*os.File, error) {
    return os.OpenFile("/tmp/some_file.txt", os.O_CREATE, fs.ModeAppend)
}

// sample/other_file.go

func writeToFile(data []byte) error {
    file, err := openFileForWriting()
    if err != nil {
        return err
    }

    _, err = file.Write([]byte("hello world"))
    if err != nil {
        return err
    }

    return nil
}
```

Output:
```bash
$ go run ./cmd/cli --dir "." --key "<<open ai key>>"
```
```
Function openFileForWriting failed at writeToFile because:
The requirement states "make sure to close this file after use." The function 
`writeToFile` opens the file with the function `openFileForWriting`, writes 
data to it, but never explicitly closes the file. It should use the `file.Close()` 
method to close the file after its use.
```

This is a proof of concept for now. Ideally it would be run in CI, and only
re-check call sites that have been edited (using a git diff, or something similar).
