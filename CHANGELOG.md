<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [v0.5.0]

* [\#129](https://github.com/bianjieai/tibc-go/pull/129) Bump up irismod version to v1.8.0
* [\#128](https://github.com/bianjieai/tibc-go/pull/128) Bump up cosmos-sdk version to v0.47.4
* [\#127](https://github.com/bianjieai/tibc-go/pull/127) Support app wiring.

## [v0.4.3]

### Improvements

* [\#124](https://github.com/bianjieai/tibc-go/pull/124) Bump irismod version to v1.7.3
* [\#121](https://github.com/bianjieai/tibc-go/pull/121) Bump cosmos version to v0.46.9

### Bug Fixes

* [\#122](https://github.com/bianjieai/tibc-go/pull/122) Fix: solve proposal handler route conflict

## [v0.4.2] - 2022-11-28

* [\#119](https://github.com/bianjieai/tibc-go/pull/119) Bump up irismod to v1.7.2

## [v0.4.1] - 2022-11-18

## Improvements

* [\#118](https://github.com/bianjieai/tibc-go/pull/118) Bump up cosmos-sdk to v0.46.5
* [\#118](https://github.com/bianjieai/tibc-go/pull/118) Bump up irismod to v1.7.1

## [v0.4.0] - 2022-11-15

### Improvements

* [\#114](https://github.com/bianjieai/tibc-go/pull/114) Bump up cosmos-sdk to v0.46.4.
* [\#116](https://github.com/bianjieai/tibc-go/pull/116) Bump up irismod to v1.7.0
