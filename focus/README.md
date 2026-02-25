# Focus Timer

A simple timer with minimal distractions

The timer is designed for longer durations since loaded progress has 1 minute granularity. For
example, starting a 1 minute timer and stopping it at 45s would log 75% progress. On timer restart,
a 1 minute timer that visually shows 75% progress completed is loaded. The remaining 25% of visual 
progress will take 1 minute to complete. Therefore the total time was 1 minute 45 seconds and the 
timer progress was visually inconsistent.

Pause functionality is not implemented. Completed time is logged in a file. If a break is needed,
close the timer. The next time the timer is launched via `focus`, progress is loaded from  the log
file.

## Usage 

`focus` load a timer with latest time from log file or a default time when no log file time found.

`focus 5` load a 5 minute timer without starting the timer.

`focus l`, `focus log` cat the timer progress from the log file.
