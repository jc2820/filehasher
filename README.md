# FileCrypt
A command line program in Go that performs basic encryption and editing operations on files.  
Can either encrypt or decrypt a file with a given key, or perform an all in one decrypt-append lines-encrypt with the key and some arguments.

## Installation (Go PATH)
1. You'll need a copy of Go (ideally v1.17+) on your machine.
2. run `go get github.com/jc2820/filecrypt` to build the latest binary in your Go bin.
3. Use `filecrypt <flags> <args>` to use the program.

e.g. `filecrypt -f ./path/to/myfile.txt -k mypasskey -r -a "line one" linetwo`
or `filecrypt -f myfile.txt -d`

## Build/edit and run
1. You need Go (ideally v1.17+).
2. Clone this repo and navigate to its root folder
3. Run `go build filecrypt.go` to build the binary for your OS.
4. Use `./filecrypt <flags> <args>` to use the program.
5. Edit and rebuild if you like.
6. You can move the binary file to somewhere in your PATH so you don't have to call it relatively.

## Flags
### General
* `-f` Specify the relative filepath to the file you want to operate on.
* `-k` Specify a pass phrase to encrypt or decrypt files with. **Don't lose this!** Default is 'secret'.
* `-r` Read mode. Will print the file's current contents. Can be used alongside other modes to print the end result.

### Job operations
Choose only one of these jobs per command.
* `-e` Encrypt mode. Encrypts the file given to `-f` with the key provided to `-k` or uses default key.
* `-d` Decrypt Mode. Attempts to decrypt the file given to `-f` with the key given to `-k` or uses default key.
* `-a` Add Mode. Decrypts the file, adds any string arguments given to the command each as a single new line then re-encrypts the file all with the key provided.

## Arguments / Add Mode
Works on files encrypted by this tool only, as it runs this tool's decrypt and encrypt modes as part of the job.
A quick one command method to add lines to a filecrypted file.
In Add Mode any string argument given to the command after flags and flag arguments will be appended each as a new line to the file. e.g.

myfile.txt (before)
```
...
old file contents
```
run FileCrypt...
`filecrypt -f myfile.txt -a new lines "to add"`

myfile.txt (after)
```
...
old file contents
new
lines
to add
```
The output will be encrypted.
Args will be ignored in other modes.

## Notes
- Paths are relative to the current directory. e.g. `./dir-inside-current-dir/filecrypt -f ../file-above-this-dir.txt -e` would work.
- This program has no option to select what encryption algorithms or modes are used, so you would be extraordinarily lucky to decrypt a file with this tool that was not originally encrypted with it.
- You can keep encrypting a file over and over. To see plain text again you will need to successfully decrypt as many times as you encrypted with the correct key each time. You can use a different key for each encryption.
- You can of course always use the default pass key if you don't bother to provide your own to commands, but that wouldn't be very secure. However, there's no way to decrypt a file without the proper key so **Please remember it/ note it down** - or think carefully before encrypting unique and vital information.
- Recursively running this programme on a system's files could be very bad indeed.
- There's no lockout after any number of unsuccessful decryption attempts so a patient person could brute force your file. Usual password complexity advice applies.
- Use Read Mode or look at your file manually if you just want to see its current content.

## Improvements to make
- better messaging (job progress, man page etc.).
- Modularise files/functions out of main.go
- verbose/silent modes - To avoid printing sensitive data if necessary.
- Packaging/release for cross-platform dl.
- Maybe even a custom complex lock/unlock script that can be stored separately.
