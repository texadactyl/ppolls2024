Analysis of the 2024 USA Presidential Election State Polls
==========================================================

#### Overview

This project analyzes the 2024 USA Presididential Election polling data, available from https://electoral-vote.com/.  To date, there doesn't seem to be much in the way of polling data compared to other years.

The ```ppolls2024``` executable will create directories and files as needed on-the-fly. Make sure that the standard Go $HOME/go/bin directory is in the $PATH. Similar advice for Windows users (you would know better than me).

Free advice regarding poll data: Do not bother looking at **national** poll data. It would make more sense if the national popular vote determined the winner; see https://www.brennancenter.org/our-work/research-reports/national-popular-vote-explained. But, that is not the way Presidential elections currently work in the USA. I shall say no more on this subject!

#### Installation

```
git clone https://github.com/texadactyl/ppolls2024/
cd ppolls2024
go install ./...
```
The first time you run ```go install```, you will probably be warned to ```go get``` other modules. Do that. Then, continue.
```
go install ./...
```

#### Invocation from a terminal window command-line

```
cd ppolls2024
ppolls2024 -f # Download the latest poll data.
ppolls2024 -l # Load the database with the downloaded data.
ppolls2024 -r tx # Get detailed report for Texas. The string "TX" is also acceptable.
ppolls2024 -r ec # Get summary report for all states. The string "EC" is also acceptable.
```

#### Licensing

This is NOT commercial software; instead, usage is covered by the GNU General Public License version 3 (2007). In a nutshell, please feel free to use the project and share it as you will but please don't sell it. Thank you!

See the LICENSE file for the GNU licensing information.

Feel free to create an issue record for any questions, bugs, or enhancement requests. I'll respond as soon as I can.

Richard Elkins

Dallas, Texas, USA, 3rd Rock, Sol, ...
