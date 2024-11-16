
# YouTube Video Downloader

This simple App can download any video or playlist on YouTube.
It downloads both the video with audio and audio file separetly.


The entire application is written in golang , however it uses ffmpeg.exe as an external dependency.

## Suggestion
It is advised to use cli tool for best performance as demonstrated

## Requirements
1)WINDOWS OS / LINUX based distribution

2)100MB free disk space

ps - Currently supports only windows os , other OS builds are in progress
## Development

To deploy this project run

1)Clone this repo into your local system

2)Extract the .zip file and get inside the repo

3)open ytdownloader.exe and paste video or playlist url and press enter.

4)Theres one more way to download the video using cmd or powershell .You can go thorugh the following example 

## Deployment
HELP:-
```bash
  ./ytdownloader.exe --help
```
example 1:-
```bash
  ./ytdownloader.exe --vidurl https://www.youtube.com/watch?v=Rtpu2cWz7W8
```
example 2:-
```bash
  ./ytdownloader.exe --urlplay https://www.youtube.com/playlist?list=PLMC9KNkIncKtPzgY-5rmhvj7fax8fdxoj
```


## Demo

Demo coming soon
## Appendix

Application is referenced from https://github.com/kkdai/youtube & https://www.ffmpeg.org




