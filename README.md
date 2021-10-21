# FileCrypt
A command line program in Go that performs basic encryption and editing operations on files.
Can either encrypt or decrypt a file with a given key, or perform an all in one decrypt-append lines-encrypt with the key and some arguments.

## Setup
1. You'll need a copy of Go (ideally v1.17+) on your machine.
2. Clone this repo and navigate to its root folder
3. Run `go build filecrypt.go` to create the binary.
4. Use `./filecrypt <flags> <args>` to use the program.

e.g. `./filecrypt -f ./myfile.txt -p mypassword -r -a "line one" linetwo`
or `./filecrypt -d`

## Flags
### General
* `-f` Specify the filepath to the file you want to operate on - default is ./cryptfile.txt.
* `-k` Specify a pass phrase to encrypt or decrypt files with. **Don't lose this!** Default is 'secret'.
* `-r` Read mode. Will print the file's contents before and after job operations. Can be used alongside other modes.

### Job operations
Choose only one of these jobs per command.
* `-e` Encrypt mode. Encrypts the file given to `-f` with the key provided to `-k` or uses default values.
* `-d` Decrypt Mode. Attempts to decrypt the file given to `-f` with the key given to `-k` or uses default values.
* `-a` Add Mode. Currently takes encrypted files only. Decrypts the file, adds any string arguments given to the command each as a single new line then re-encrypts the file all with the key provided.

## Arguments / Add Mode
In Add Mode any string argument given to the command after flags and flag arguments will be appended each as a new line to the decrypted file. e.g.

myfile.txt (before)
```
...
old file contents
```
run FileCrypt...
`./filecrypt -f ./myfile.txt -a new lines "to add"`

myfile.txt (after)
```
...
old file contents
new
lines
to add
```

Args will be ignored in other modes.

## Warnings
- This program has no option to select what encryption algorithms or modes are used, so you would be extraordinarily lucky to decrypt a file with this tool that was not originally encrypted with it.
- You can keep encrypting a file over and over. To see plain text again you will need to successfully decrypt as many times as you encrypted with the correct key each time. You can use a different key for each encryption.
- You can of course always use the default pass key, but that wouldn't be very secure. However, there's no way to decrypt a file without the proper key so **Please remember it/ note it down** - or think carefully before encrypting unique and vital information.
- Recursively running this programme on a system's files could be very bad indeed.
- There are no checks on unsuccessful decryption attempts so a patient person could brute force your file. Usual password complexity advice applies.
- Attempting add mode on an unencrypted file gives the same error as providing an incorrect decryption key. If you're having cipher message authentication errors in add mode check your file manually or use read mode - it might be open to edit normally. If not check the password you're giving.

## Improvements to make
- better messaging (job progress, man page etc.).
- verbose/silent modes - To avoid printing sensitive data if necessary.
- Fix the above add mode error quirk.
- Packaging for cross-platform dl.
- Maybe even a custom complex lock/unlock script that can be stored separately.
