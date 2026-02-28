# Focus Timer

A flow timer that emphasizes minimal distractions and an efficient feature set. Minimal distractions
because the timer only shows a progress bar. There is no visual distraction occurring each second
which is common since many timers display ticking seconds. Efficient feature set because there is 
only three commands.

A mouse click starts the timer. I found this better than directly starting the timer from any
command because I cleared timer progress by mistake multiple times. There is no pause command. If a
break is needed, close the timer. Progress will be loaded the next time the timer is launched via
`focus`.

## Usage 

`focus` load a timer from prior progress. If no prior progress is found a default time is loaded.

`focus 5` load a 5 minute timer without starting the timer.

`focus l`, `focus log` cat the timer progress.
