# RTUN - Reverse tunnel

RTUN is a tool for connecting back to your local computer when connected to a
remote server via SSH. It relies on a local daemon and a UNIX socket forwarded
to the remote server.

## What do I use RTUN for?

One problem that often happens when you work on remote servers is how to access data or files from the server from your local computer. Perhaps you're a data scientist and you generate a figure on the server. How do you see that figure locally? There are some tools for displaying images in a terminal, but they either have low resolution or are very difficult to use. Cloud storage like Dropbox can be used to sync data from servers, but it can be difficult to setup for individual users (it can require a developer account). And -- it can be slow as data is sent from the server to the cloud and then back to your computer.

RTUN aims to solve these problem.

RTUN has three major functions: 

1. Sending notifications from a server to your local computer
2. Sending files from a remote server to your local computer
3. Viewing files from a remote server on your local computer (which is basically saving the file as a temp file and then calling `open`).

So long as you have an SSH connection to the remote server (with the appropriate SSH reverse tunnel established), then you can send a notification or a file to your local computer.

## Supported operating systems

Currently, the only supported operating systems are macOS and Linux. BSD
systems should work, but are untested. 

I'm not sure this is a good mechanism to use with Windows clients or not, as it is untested.

## Security

RTUN relies on UNIX sockets as opposed to IP ports or user accounts. This is so that on the
remote server, we don't need to open a public port. All security to access the
reverse tunnel is managed through file permissions to the socket file. This
way, the only way to access the local computer is through a file owned by you.
If we opened an IP port, then anyone with access to the remote server would
also be able to connect to your local computer (to send files or
notifications). This way, without needing to setup yet another account or
login, you can connect back to your local computer securely. You already authenticated with SSH, so why add another login to the mix?

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

Connect to the remote server (note, you must provide absolute pathnames to the socket here, in the format: `-R remote_file:local_file`)

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

### So, how do I set this up?

Personally, I have a bash/zsh alias setup to connect to a remote server. This alias looks like this:

    alias remote='rtun server -d; ssh -R /home/username/.rtun/rtun.sock.$RANDOM:/Users/username/.rtun/rtun.sock remote-host`

*(Don't forget to create the `$HOME/.rtun` directory on both the local and remote systems!)*

Then I can run `remote`. On the first run, this starts the RTUN daemon process. Then it connects to the remote server with SSH, setting up a UNIX socket reverse tunnel, with the socket name on the server to a random file name (but easily discoverable). On the server, whenever `rtun` is run, then it figures out which connection is valid. When you sign off of the server (with `exit`), the socket file is not automatically removed.

The next time you run the `remote` command, it tries to start an RTUN daemon, but fails with a message saying that the socket is already in use. Then when the SSH connection is established, a new socket file is created. When you run `rtun`, the program looks for possible socket files. When the right one is found, it is used to connect back to the local RTUN daemon.

In theory, I think this could also be setup using only the `.ssh/config` file, but I have not tested this.