# RTUN - Reverse tunnel

RTUN is a tool for connecting back to your local computer when connected to a
remote server via SSH. It relies on a local daemon and a UNIX socket forwarded
to the remote server.

## Supported operating systems

Currently, the only supported operating systems are macOS and Linux. BSD
systems should work, but are untested. I'm not sure this is a good mechanism to
use with Windows clients or not, as it is untested.

## Security

RTUN relies on UNIX sockets as opposed to IP ports. This is so that on the
remote server, we don't need to open a public port. All security to access the
reverse tunnel is managed through file permissions to the socket file. This
way, the only way to access the local computer is through a file owned by you.
If we opened an IP port, then anyone with access to the remote server would
also be able to connect to your local computer (to send files or
notifications). This way, without needing to setup yet another account or
login, you can connect back to your local computer securely. 

*(This assumes you trust anyone that also has root access on the remote server,
but if you can't trust root on the remote server, then you have bigger issues
than protecting a file on the server.)*


## Usage

    rtun cmd


### Commands

`notify` - sends a desktop notification to the local computer

`send` - send a file (or directory) to the local computer

`view` - open a file on the local computer

`server` - start the server daemon

## Example workflow 

Start the local daemon

    [username@local ~] $ rtun server -d rtun.sock

Connect to the remote server

    [username@local ~]$ ssh -R /home/username/rtun.sock:$(realpath rtun.sock) username@remote-host

Send a message back to your local machine

    [username@remote-host ~]$ rtun notify -s ~/rtun.sock Hello from the server


