# DeliveryRL

The ancient dragon is hungry and you're short on time. You had better deliver that pizza quickly!

Instructions
------------

Welcome to DeliveryRL! You are a typical delivery person for a most unusual pizza shop. Sometimes your store gets orders from mythical creatures. Today an Ancient Dragon has ordered a pizza and expects it to be delivered promptly, and warm! Why does a dragon need delivery? It's not your job to ask those questions, but you drew the short straw this time.

Race to the depths of the caverns and deliver the pizza to the Dragon who is waiting. Survive by avoiding monsters, fighting for your life, and distracting them with some extra food you brought along. You also have a special trick up your sleeve, the ability to warp a short distance. It tires you out but don't let that slow you down; time is of the essence, after all! Maybe you can scrounge up something useful from past adventurers, but remember: even though you can teleport you're no warrior!

You will find the dragon on the 10th floor!


# Running It

Windows
-------

Running delivery-rl.exe in this directory should be sufficient to start.

Linux
-----

If you have SDL2, SDL2_image, and libpng installed you should be able to just run delivery-rl. If not, you can try running with the bundled shared objects by running delivery-rl.sh instead which will prepend them to your LD_LIBRARY_PATH and then launch delivery-rl. You can also install SDL2, SDL2_image, and libpng via your package manager and run delivery-rl as normal.

OSX
---

If you've downloaded the OSX package, you will need to install sdl2 and sdl2_image if you don't already have them installed.
```
brew install sdl2{,_image,} pkg-config
```
Then you can directly run `delivery-rl.sh`. If you can't get that to work, or you have downloaded the windows+linux archive you can run the windows version with wine.


Build from Source
-----------------

Source is distributed separately. Pull the latest version from https://github.com/thomas-holmes/delivery-rl

Development Instructions
=====

## Compilation

To compile/run you will need sdl2, sdl2_img, and gcc/mingw

The sdl2 bindings are vendored but you will still need to ensure that the shared libraries are installed and available on your system. Go to [veandco/go-sdl2](https://github.com/veandco/go-sdl2) and follow the SDL installation instructions for your platform.

After installing all the dependencies you can start the game by running `make run`. This will compile the game to `run_dir/delivery-rl` and start it.

This archive contains delivery-rl.exe (win64 build) and delivery-rl linux amd64 build.

The windows build should run in place, but if not follow instructions for installing SDL2 and SDL2_img on your system.

The buildwin makefile target is for cross compiling from linux->windows and should work if you have a full mingw-w64-gcc installation on your system.