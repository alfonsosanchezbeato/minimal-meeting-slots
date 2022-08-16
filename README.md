# Scheduling meetings in the minimal number of time slots

This repository contains a small golang program that implements an
algorithm to find the minimal number of time slots that would be
required to schedule a set of meetings while avoiding
incompatibilities across meetings so all invitees can attend their
respective meetings.

More information on the problem and a description of the algorithm can
be found in this [blog
post](https://www.alfonsobeato.net/math/how-to-schedule-meetings-in-the-minimal-number-of-time-slots/).

## Build and test

Build with:

```
go build -o minimize-slots main.go
```

To run unit tests:

```
go test *.go
```

## Running

Run with

```
./minimize-slots <input_csv> <output_csv> [optional_dot_file]
```

Input and output files are [CSV
files](https://en.wikipedia.org/wiki/Comma-separated_values) that can
be easily exported and imported from spreadsheets. The third argument
is optional and if provided a description of the solution graph is
written in [dot language](https://graphviz.org/doc/info/lang.html).

The first column of the input file is a meeting title while
the second column is a list of assistants to that meeting, separated
by commas. For instance:

```
Discuss Q1 objectives,"Cide Hamete Benengeli, Teresa de Jesus"
Team building,"Piero della Francesca, Leopoldo Alas"
QA workshop,"Teresa de Jesus, Leopoldo Alas"
```

The output file's first column is the slot assigned to the meeting whose
title and assistants are in the second and third columns:

```
slot 0,Discuss Q1 objectives,"Cide Hamete Benengeli, Teresa de Jesus"
slot 0,Team building,"Piero della Francesca, Leopoldo Alas"
slot 1,QA workshop,"Teresa de Jesus, Leopoldo Alas"
```

To convert the dot file to a svg file it is recommended to use neato,
which is part of [graphviz](https://graphviz.org/):

```
neato -Tsvg <dot_file> > out.svg
```
