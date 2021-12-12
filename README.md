# Auto-updates PaperMC projects

### Motivation

Paper has made it very difficult to auto update servers, because they discourage it. This means that you can't just request the latest version from their API.  
But then you see, yesterday a log4j vunl in minecraft allowed RCE.  
As long as you have a GOOD backup, IMO auto updates are fine.

### TODO
- [x] Working version
- [ ] Docker container
- [ ] Bash script for auto update

### Installation

Download the correct binary from Releases and install it to your $PATH.  

### Running

`./the-binary-file paper-project-name project-version`  

### Examples  

Download PaperMC 1.18.1  
`./paper-autoupdater.bin paper 1.18.1`  
Download Velocity 3.1.1  
`./paper-autoupdater.bin velocity 3.1.1`
Download & run PaperMC 1.18.1
`./paper-autoupdater.bin paper 1.18.1 && java -jar paper-1.18.1.jar`

### Auto update

Recommended way is autorestarting server and having auto updater in script when server is restarted
