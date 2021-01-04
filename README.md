# RTUN - Reverse tunnel

RTUN is a tool for connecting back to your local computer when connected to a
remote server via SSH. It relies on a local daemon and a UNIX socket forwarded
to the remote server.

## Supported operating systems

Currently, the only supported operating systems are macOS and Linux. BSD
systems should work, but are untested. 

I'm not sure this is a good mechanism to use with Windows clients or not, as it is untested.

## Security

RTUN relies on UNIX sockets as opposed to IP ports. This is so that on the
remote server, we don't need to open a public port. All security to access the
reverse tunnel is managed through file permissions to the socket file. This
way, the only way to access the local computer is through a file owned by you.
If we opened an IP port, then anyone with access to the remote server would
also be able to connect to your local computer (to send files or
notifications). This way, without needing to setup yet another account or
login, you can connect back to your local computer securely. 

*Note: This all assumes you trust anyone that also has root access on the remote server, but if you 
can't trust root on the remote server, then you have bigger issues than protecting the socket file.*


## Usage

    rtun cmd


### Commands

`notify` - sends a desktop notification to the local computer

`send` - upload a file (or directory) to the local computer

`view` - open a file on the local computer

`server` - start the server daemon

## Example workflow 

Start the local daemon

    [username@local ~] $ rtun server -d -s rtun.sock

Connect to the remote server

    [username@local ~]$ ssh -R /home/username/rtun.sock:/Users/username/rtun.sock username@remote-host

Send a message back to your local machine

    [username@remote-host ~]$ rtun notify -s ~/rtun.sock Hello from the server


### Default socket names

By default, on the local machine, the server socket will be saved to `$HOME/.rtun/rtun.sock`. Similarly, if the server is started in background/daemon mode, then the output log will be written to `$HOME/.rtun/rtun.log` by default.

On the remote server, the socket file can be auto-discovered if it is named `$HOME/.rtun/rtun.sock.*`. The reason why we scan for multiple files is that when you disconnect from the remote server, any remote UNIX socket files will not be automatically removed. This can lead to stale file handles and will make it so that any subsequent connection attempts will fail to create the socket with an error.

Example: 

    [username@local ~]$ ssh -R /home/username/.rtun/rtun.sock:/Users/username/.rtun/rtun.sock remote-host

    [username@remote-host ~]$ rtun ping -s .rtun/rtun.sock
    OK PONG 

    [username@remote-host ~]$ exit 

    [username@local ~]$ ssh -R /home/username/.rtun/rtun.sock:/Users/username/.rtun/rtun.sock remote-host
    Warning: remote port forwarding failed for listen path /home/username/.rtun/rtun.sock

    [username@remote-host ~]$ rtun ping -s .rtun/rtun.sock
    panic: Unable to find a valid server socket (1): .rtun/rtun.sock, dial unix .rtun/rtun.sock: connect: connection refused

It is possible to setup SSH to automatically remove these stale files, but it requires support from the sshd configuration, which isn't always possible to change. So, in order to get around this, we will scan the `$HOME/.rtun/` directory for all files that start with `rtun.sock`. With this setup, you can now create multiple randomly named files, and RTUN will select the first valid connection it finds.

*Note: This does have the side effect that you cannot have multiple active sessions to the same remote host from multiple local clients. However, this should be a rare use-case. There is a mechanism in SSH to send the name of the socket using an ENV variable, but this too requires support on the sshd configuration.*

Here is one way to do this:

    [username@local ~]$ ssh -R /home/username/.rtun/rtun.sock.$RANDOM:/Users/username/.rtun/rtun.sock remote-host
