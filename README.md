Analysis of the 2020 USA Presidential Election State Polls
==========================================================

#### Overview

This project analyzes the 2024 USA Presididential Election polling data, available from https://electoral-vote.com/.  To date, there doesn't seem to be much in the way of polling data compared to other years.  .

Free advice: Do not bother looking at national poll data.

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
ppolls2024 -r tx # Get detailed report for Texas.
ppolls2024 -r ec # Get summary report for all states.
```

#### Licensing

This is NOT commercial software; instead, usage is covered by the GNU General Public License version 3 (2007). In a nutshell, please feel free to use the project and share it as you will but please don't sell it. Thank you!

See the LICENSE file for the GNU licensing information.

Feel free to create an issue record for any questions, bugs, or enhancement requests. I'll respond as soon as I can.

Richard Elkins

Dallas, Texas, USA, 3rd Rock, Sol, ...
