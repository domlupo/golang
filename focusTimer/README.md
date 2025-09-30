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
- Fix percentage load bug. See below example.

  8.7% 2025-09-29T18:56:13-07:00  219 minutes remaining
    9% 2025-09-29T18:56:50-07:00  219 minutes remaining
  9.2% 2025-09-29T18:57:14-07:00  218 minutes remaining
  9.6% 2025-09-29T18:58:12-07:00  217 minutes remaining
Timer set for 217 minutes.
startTime: 2025-09-29T18:58:40-07:00
  0.5% 2025-09-29T18:59:40-07:00  216 minutes remaining
  0.9% 2025-09-29T19:00:40-07:00  215 minutes remaining
    9% 2025-09-29T19:00:50-07:00  215 minutes remaining
  1.4% 2025-09-29T19:01:40-07:00  214 minutes remaining
