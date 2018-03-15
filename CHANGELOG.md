Changelog
---------
* 2.0.0 (Future *NOT 7DRL*)
	* UI
		* Fixed self branding.
		* Fixed incorrect config flag descriptions.
		* Added Start/Quit choices on main screen.
	* Engine
		* Redid input handling.
	* Gameplay
		* Can now restart from the title screen after losing or winning.
		* Victory dialog shows turns and heat.
* 1.0.5 (Latest 7DRL Version)
	* Update version of gterm rendering library, reducing CPU usage by approximately half.
* 1.0.4
	* Balance
		* Rearrange weapon progression slightly to be more sensible
		* Slightly reduced effectieness of armour
		* Make certain items only spawn deeper
	* UI
		* Examine cursor now works with numpad and arrows.
		* Examine interface shows if a monster is slowed or confused.
		* HUD shows if the player is slowed.
		* Throwing target selection inspect selection now more like the warp interface.
		* Stairs are now colored orange instead of white.
		* Inventory cursor can wrap around.
		* Fixed coloration of Inspect UI.
	* Gameplay
		* Grease slows monsters while in it instead of if they were in it last turn.
* 1.0.3
	* Update help screen with information about the dragon's location
* 1.0.2
	* Accidentally built 1.0.1 with an unsaved file missing a few minor balance tweaks. Hand warmers are now slightly less effective to offset the gain in power from adding their flat +10 bonus.
* 1.0.1
	* Added configuration flags for almost all of the balance tuning parameters (see -help)
		* Starting items
		* Item spawn rates
		* Monster spawn rates
		* Player regen rates
		* Heat degen rate
		* Fixed bug where bonus modifier on item rolls was not being included (Thanks tself55!)
	* Added equipment power to the HUD for items on the ground
	* Adjusts balance to make the game a little bit harder
		* Slightly reduce item density across the board
		* Slightly increase monster density
		* Reduce starting item counts from 5 to 3
* 1.0.0
	* Initial 7DRL submission release
