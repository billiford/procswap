<img src="https://github.com/billiford/procswap/blob/media/procswap.png" width="125" align="left">

# procswap

Procswap is a simple application that allows your to prioritize processes on a Windows machine. This is great for mining when you want to mine all the time unless certain processes (looking at you Cyberpunk) are running.

It works by allowing the user to pass in their mining scripts and scanning a given directory for all `.exe` files and marking these as priorities. It polls the processes running in windows to check if one of these `.exe` priorities has started (e.g. `Hades.exe` has launched), and then stops all passed in scripts. Whenever there are no priority processes running, like when you get tired of playing Cyberpunk and close it, it launches all passed in scripts - these scripts are paths to your `.bat` [files](https://2miners.com/blog/phoenixminer-step-by-step-guide-for-beginners/#PhoenixMiner_Setup) that start your miners.

### Getting Started

- Download the latest version from the [releases page](https://github.com/billiford/procswap/releases) or build it from source if you don't like running random executables found on the internet.
- You'll have to [open a command prompt](https://www.howtogeek.com/235101/10-ways-to-open-the-command-prompt-in-windows-10/#:~:text=Open%20a%20Command%20Prompt%20in%20Admin%20Mode%20from,from%20the%20File%20Explorer%20Address%20Bar.%20More%20items) and you might want to run it as Administrator.
- Change directory to where you downloaded the `procswap.exe` (probably your Downloads folder).

### Examples

I like to have `procswap` scan my Steam games directory for `.exe` files and run two `.bat` files, one for [PheonixMiner](https://phoenixminer.org/) and another for [XMRig](https://xmrig.com/).

My command looks something like this:
```
procswap.exe --priority D:\Steam\steamapps\common --swap C:\Mining\PhoenixMiner\start_miner.bat --swap c:\Mining\xmrig\start.cmd
```
I store my steam games on drive `D:\` so I passed in the link to the common directory where all the games are held. I passed in two "swap" processes that will run whenever there's no games running. The two `.bat` scripts I passed in are for my PheonixMiner script and one for XMRig. You can pass in as many priority directories or swap processes that you want. Just keep in mind that **if any priority process starts, all swap processes are stopped**.

The output looks something like this:

<img src="https://github.com/billiford/procswap/blob/media/procswap-v0.1.0-output.png" align="left">

As you can see, the first thing that happened is it scanned my games folder, found all the executables and since none were running started up my `.bat` files to start mining! When I opened Hades, it logged that this priority process started and stopped all my scripts. When I quit playing Hades it started them up again. Efficient mining!

### Notes

This project was created casually over the course of a couple days. Any feedback or suggestions are welcome.
