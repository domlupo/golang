# Focus Timer

A simple timer with minimal distractions

Pause functionality is not implemented. Completed time is logged in a 
file. If a break is needed, read the file, subtract completed minutes,
launch a new timer once break is complete.

## TODO

- Don't clear log file until timer is started. This prevents unneeded
  log file deletion. Most useful when a user forgets the exact time
  they want to continue at.
- Add ability for the timer to read the log file to continue the 
  session. e.g. I just stopped timer at 50 mins so next time I open
  the default time is set to 50 and I can opt into continue to log in
  the same file.
- Add optional terminal input to set time rather than force mouse wheel.
