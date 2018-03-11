package main

var topLeft = rune('┌')     // 0x250C
var horizontal = rune('─')  // 0x2500
var topRight = rune('┐')    // 0x2510
var vertical = rune('│')    // 0x2502
var bottomLeft = rune('└')  // 0x2514
var bottomRight = rune('┘') // 0x2518

// gterm font stuff is still kinda janky so I need to do this
var upArrow = rune(30)    // 0x25B2 ▲
var downArrow = rune(31)  // 0x25BC ▼
var rightArrow = rune(16) // 0x25BA ►

var partialBlockLeft = rune('▌')  // 0x258C
var partialBlockRight = rune('▐') // 0x2590
var fullBlock = rune('█')         // 0x2588

var grease = rune('≈')
