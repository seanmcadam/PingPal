# PingPal
Curses based multi-ip monitor tool

Curses base - resizable text window
(See top, atop, btop, mtr for reference)

Takes command line arguments (TBD) and the last arguments are host names and IP addresses of the targets

pingpal -v 1.2.3.4 2.3.4.5 3.4.5.6 4.5.6.7

multiple screen display options
1)
Target | latency | ave latency
2) 
Target | visual line of latency times

For reference look at: mtr, but only presenting the target IP not all of the hops

