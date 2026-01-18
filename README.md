Introduction
============

This is a go project.  Claude Code helped to write the source code.

It controls 3 different types of light units, by sending various MQTT messages:

1. LED strip.  This is a length of RGB LEDs, where the colour is set by sending a JSON message to the topic "kevinoffice/ledstrip/sequence" with the following payload:
  {"sequence":"fill", "data":{"r":<int>,"g":<int>,"b":<int>}}
where "r", "g", and "b" contain integer values for red, green, and blue respectively.

2. LED bar.  This is a bar of RGBW LEDs, with some additional white-only LEDs.  There are 2 of these bars, and each one reads from a comma-separated list of values sent to the MQTT topic "kevinoffice/ledbar/0".  The values in the list are:
   * 6 sets of 4 values, for 6 RGBW LEDs
   * 13 sets of 1 value, for 13 white LEDs
   * 3 ignored values
   * 6 sets of 4 values, for 6 RGBW LEDs
   * 13 sets of 1 value, for 13 white LEDs

3. Video lights.  There are 2 of these, one at the topic "kevinoffice/videolight/1/command/light:0" and the other at the topic "kevinoffice/videolight/2/command/light:0".  They read a message in this format:
set,<on>,<brightness>
for example this message will turn the light on to half brightness:
set,true,50

Structure
=========

There is a folder "drivers", containing a folder for each of the types of light.  These drivers keep information about the current state of the relevant lights, and format the correct messages for publishing.  They get instantiated for each instance of that type of light.

There will be a user-interface in the future, containing buttons and dials to alter the state of any instantiated light, which will trigger sending the MQTT messages to change the lights.

The `main.go` file in the root contains the orchestration code.

Features
========

State storage
-------------

The state of all of the lights should be stored in a sqlite3 file in the current directory, called "lights.sqlite3".  The table/column structure is:

ledbars : id
ledbars_leds : id, ledbar_id, channel_num, value
ledstrips : id, red, green, blue
videolights : id, on, brightness

The state should be loaded on startup by querying the sqlite file, and saved back to the file every time a value changes and is published to MQTT.  Since there is only 1 LED bar and 1 LED strip, they are hard-coded as ID 0, and the 2 videolights are hard-coded as IDs 0 and 1.

User interfaces
===============

Multiple UIs are able to be run simultaneously.

TUI
---

One of the user-interfaces is a text user-interface.  The screen is split into 4 sections, one for each of the lights.  In each section there are controls for RGB, RGBW, W, or brightness as appropriate to that type of light.  The "TAB" key switches focus between the sections, while arrow keys move between the input controls.  Up and Down arrow keys change the values by small amounts, while holding shift with up and down changes the values in large amounts.

Web
---

One of the user-interfaces is a web interface, which runs in a spawned go func().  The web interfaces is composed of 2 separate parts:

1. An API at /api, which responds to GET requests by returning the complete status as a JSON structure, and responds to POST requests where the payload is a JSON doc that contains the complete status to change to.

2. An HTML page which makes an AJAX request to get the status from the API and renders some HTML UI components, and whenever the user changes something it sends the status back to the POST API endpoint.

**Run Web Interface:**
```bash
./office_lights web
# Access at http://localhost:8080
```

**Run both UIs simultaneously:**
```bash
./office_lights tui web
```

Streamdeck
----------

One of the user-interfaces is a Streamdeck+, which has:
* 8 buttons arranged into 2 rows of 4.  Each button has a configurable image behind it, of size 120x120 pixels.
* A long thin touchscreen, size 800x100, which can sense taps but nothing else.
* 4 physical rotary encoders, with detents so there are defined up or down actions, and they can also be pressed (or clicked).

The UI on the streamdeck should work like this:

* The first 3 buttons on the top row should act like a radio button set, selecting between:
  1. The LED strip
  2. The LED bar's RGBW lights
  3. The LED bar's plain white lights
  4. The video lights

* The touchscreen should be split into 4 areas, and show something different depending on which "radio button" from the top row of buttons is chosen:
  * For the LED strip, show 3 things: red, green, and blue
  * For the LED bar's RGBW lights, show 4 things: red, green, blue, and white
  * For the LED bar's plain white lights, show 2 things: the brightness of the plain white lights in the first section of 13 lights, and the brightness of the plain white lights in the second section of 13 lights.
  * For the video lights, show 2 things: the brightness of the first video light, and the brightness of the second video light.

* The 4 dials should allow the values shown on the touchscreen to be increased or decreased as the dial is turned, in increments of 5.  Clicking the dials will either toggle the value between "0" and the last-used value, or in the case of the video lights it will toggle the on/off state.  The other two dials will increase or decrease the 2 video lights in increments of 1.

**Run Stream Deck Interface:**
```bash
./office_lights streamdeck
```

**Run all UIs simultaneously:**
```bash
./office_lights tui web streamdeck
```
