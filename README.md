# FaizEye
[![License](https://img.shields.io/badge/license-MIT-_red.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/Ice3man543/hawkeye)](https://goreportcard.com/report/github.com/Ice3man543/hawkeye) 
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/Ice3man543/hawkeye/issues)

FaizEye is a fork of HawkEye,  but for Red Teaming. For those who don't know; HawkEye is a tool to crawl the filesystem looking for interesting stuff like SSH Keys, Log Files, Sqlite Database, password files, etc.

FaizEye is a modification that steals these files and uploads them to your C2. The only changes that are needed are to change the SSH credentials within main.go, FaizEye will detect the OS at runtime and upload interesting files found on the victims computer to your C2.
