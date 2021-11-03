# irods-compare
`iRODS-Compare` is a command-line tool that compares data between local disks and iRODS (CyVerse DataStore). It checks files out of sync using MD5 checksum and file size.

# Pre-built binaries
`iRODS-Compare` can run on any OSes. We provide pre-built binaries for some OSes. If you cannot find one for your OS, you can build it from source code by yourself.

Links for `iRODS-Compare` binaries `v0.1.1`:

For Linux:
- Intel 64Bit (amd64): [irods_compare_amd64_linux_v0.1.1.tar](https://github.com/cyverse/irods-compare/releases/download/v0.1.1/irods_compare_amd64_linux_v0.1.1.tar)
- Intel 32Bit (i386): [irods_compare_i386_linux_v0.1.1.tar](https://github.com/cyverse/irods-compare/releases/download/v0.1.1/irods_compare_i386_linux_v0.1.1.tar)
- ARM 64Bit: [irods_compare_arm64_linux_v0.1.1.tar](https://github.com/cyverse/irods-compare/releases/download/v0.1.1/irods_compare_arm64_linux_v0.1.1.tar)
- ARM 32Bit: [irods_compare_arm_linux_v0.1.1.tar](https://github.com/cyverse/irods-compare/releases/download/v0.1.1/irods_compare_arm_linux_v0.1.1.tar)

For Macos (OSX):
- Intel 64Bit (amd64): [irods_compare_amd64_darwin_v0.1.1.tar](https://github.com/cyverse/irods-compare/releases/download/v0.1.1/irods_compare_amd64_darwin_v0.1.1.tar)
- ARM 64Bit: [irods_compare_arm64_darwin_v0.1.1.tar](https://github.com/cyverse/irods-compare/releases/download/v0.1.1/irods_compare_arm64_darwin_v0.1.1.tar)

For Windows:
- Intel 64Bit (amd64): [irods_compare_amd64_windows_v0.1.1.tar](https://github.com/cyverse/irods-compare/releases/download/v0.1.1/irods_compare_amd64_windows_v0.1.1.tar)
- Intel 32Bit (i386): [irods_compare_i386_windows_v0.1.1.tar](https://github.com/cyverse/irods-compare/releases/download/v0.1.1/irods_compare_i386_windows_v0.1.1.tar)

# USAGE
Create config.yaml file to configure your iRODS account.
```yaml
host: data.cyverse.org
port: 1247
user: "your_iRODS_username"
zone: iplant
password: "your_password" or leave empty to type in later
colorize: true
```

You can put your iRODS password in the password field if you don't want to type the password every time you run the command. But this is insecure. iRODS-Compare will ask you to type in if it's empty.

```bash
./irods-compare -config config.yaml <local_file_or_dir_path> <irods_file_or_dir_path>
```

# Example
Compare two directories recursively.
```shell script
./irods-compare -config config.yaml . /iplant/home/iychoi
```

In this case, `iRODS-Compare` will first find files in the given local directory and compare them with iRODS ones.  

OUTPUT:
```shell script
~/irods-compare/bin$ ./irods-compare -config config.yaml . /iplant/home/iychoi
Password:
INFO[0004] Connecting to data.cyverse.org:1247           function=Connect package=connection struct=IRODSConnection
INFO[0004] Start up a connection without CS Negotiation  function=connectWithoutCSNegotiation package=connection struct=IRODSConnection
INFO[0005] Logging in using native authentication method  function=loginNative package=connection struct=IRODSConnection
INFO[0005] Checking local file .                         function=main package=main
ERRO[0005] failed to find irods file - /iplant/home/iychoi/irods-compare  function=main package=main
ERRO[0005] failed to find irods file - /iplant/home/iychoi/irods-compare.exe  function=main package=main
+-----+------------------------------------------------------------------+----------------------------------+----------------+-----------------------------------------+------------+
| #   | PATH                                                             | HASH                             | FILE SIZE      | MODIFIED TIME                           | CONSISTENT |
+-----+------------------------------------------------------------------+----------------------------------+----------------+-----------------------------------------+------------+
| 1   | /home/iychoi/Projects/irods-compare/bin/POV_L.Spr.O.10m_reads.fa | 8846629b71f41e8c97e81438ff6afd0d | 21047770       | 2021-10-23 15:51:26.193513367 -0700 MST | TRUE       |
| --> | /iplant/home/iychoi/POV_L.Spr.O.10m_reads.fa                     | 8846629b71f41e8c97e81438ff6afd0d | 21047770       | 2021-10-27 11:10:56 -0700 MST           | TRUE       |
+-----+------------------------------------------------------------------+----------------------------------+----------------+-----------------------------------------+------------+
| 2   | /home/iychoi/Projects/irods-compare/bin/config.yaml              | ad048c7fc69f47f1f3b5b8a901cc0f28 | 87             | 2021-11-02 16:18:45.297936406 -0700 MST | FALSE      |
| --> | /iplant/home/iychoi/config.yaml                                  | 27ea81374d137b016cf5b58380868b30 | 79             | 2021-10-23 19:36:37 -0700 MST           | FALSE      |
+-----+------------------------------------------------------------------+----------------------------------+----------------+-----------------------------------------+------------+
| 3   | /home/iychoi/Projects/irods-compare/bin/irods-compare            | SKIP                             | 6636104        | 2021-11-02 15:47:59.79808146 -0700 MST  | FALSE      |
| --> | /iplant/home/iychoi/irods-compare                                | FILE NOT FOUND                   | FILE NOT FOUND | FILE NOT FOUND                          | FALSE      |
+-----+------------------------------------------------------------------+----------------------------------+----------------+-----------------------------------------+------------+
INFO[0005] Disconnecting the connection                  function=Disconnect package=connection struct=IRODSConnection
```

# WIKI
For more information about downloading pre-built binaries, building from source code, and more examples, please refer WIKI.\
WIKI: [https://github.com/cyverse/irods-compare/wiki](https://github.com/cyverse/irods-compare/wiki)
https://github.com/cyverse/irods-compare/wiki
