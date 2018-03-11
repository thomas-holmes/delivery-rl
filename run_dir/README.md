# DeliveryRL

This archive contains delivery-rl.exe (win64 build) and delivery-rl linux amd64 build.

The windows build should run in place, but if not follow instructions for installing SDL2 and SDL2_img on your system.


Windows
-------

Running delivery-rl.exe in this directory should be sufficient to start.

Linux
-----

If you have SDL2, SDL2_image, and libpng installed you should be able to just run delivery-rl. If not, you can try running with the bundled shared objects by running delivery-rl.sh instead which will prepend them to your LD_LIBRARY_PATH and then launch delivery-rl. You can also install SDL2, SDL2_image, and libpng via your package manager and run delivery-rl as normal.

OSX
---

Run the windows version with wine

Build from Source
-----------------

Source is distributed separately. Pull the latest version from https://github.com/thomas-holmes/delivery-rl

