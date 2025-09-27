# Focus Timer

A simple timer with minimal distractions

The timer is designed for longer durations since loaded progress has 1 minute granularity. 
For example, starting a 1 minute timer and stopping it at 45s would log 75% progress.
On timer restart, a 1 minute timer at 75% progress completed is loaded. The remaining 25% 
of progress will take 1 minute. Therefore the total time was 1 minute 45 seconds rather than 
the desired 1 minute.

Pause functionality is not implemented. Completed time is logged in a 
file. If a break is needed, read the file, subtract completed minutes,
launch a new timer once break is complete.

## TODO

- Add optional terminal input to set time rather than force mouse wheel.