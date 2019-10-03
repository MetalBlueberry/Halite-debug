# Simple Canvas Server 

## How to use me

Start the server by running the executable 'Halite-debug' and it will start a server in the port 8888 of the local host. Perform a post with a list of Methods to the path "http://localhost:8888/img/{game}/{turn:[0-9]+}" where the turn is a string identifier of the game and turn is a positive integer representing the halite turn. This will automatically open the browser in the url http://localhost:8888/visor/{game} in the turn number 1. Now you can use the mouse to pan/zoom and the arrows left/right to change the current turn.

## Input format

The endpoint is expecting a json file representing a list of json object with a mandatory field called "Method" and the other fields will depend on it.  See tests/sample.json for reference.

### Circle

```json
{
    "Method":"Circle",
    "X":100.0,
    "Y":100.0,
    "R":10.0,
    "Class":["planet player1"]
}
```

- X Y and R are the position and radius. 
- Class will be used as a identifier for the final SVG image that can be referenced using CSS styles

### Line

```json
{
    "Method":"Line",
    "X1":100.0,
    "Y1":100.0,
    "X2":120.0,
    "Y2":100.0,
    "Class":["target"]
}
```

- X1 Y1 X2 Y2 are the two point that identify the line.
- Class will be used as a identifier for the final SVG image that can be referenced using CSS styles

## Examples

You can send the sample.json under the test folder to the server to see how it looks like.

```sh
curl --data @sample.json  -X POST http://localhost:8888/img/test/1
```

## Styling

Add a CSS file under static/styles.css relative to the executable. You can use the file from the repository as an initial point.

