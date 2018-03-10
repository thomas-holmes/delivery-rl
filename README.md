# DeliveryRL

The ancient dragon is hungry and you're short on time. You had better deliver that pizza quickly!

## Compilation

To compile/run you will need sdl2, sdl2_img, and gcc/mingw

The sdl2 bindings are vendored but you will still need to ensure that the shared libraries are installed and available on your system. Go to [veandco/go-sdl2](https://github.com/veandco/go-sdl2) and follow the SDL installation instructions for your platform.

After installing all the dependencies you can start the game by running `make run`. This will compile the game to `run_dir/delivery-rl` and start it.