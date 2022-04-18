'use strict';
const express = require('express');
const uuid = require('uuid')
// Constants
const PORT = 8080;
const HOST = 'localhost';
// App
const app = express();
app.get('/v1/integrate', (req, res) => {
    console.log("v1 HI")
    res.json({ secret: uuid.v4()});
});
app.get('/integrate', (req, res) => {
    console.log("HI")
    res.json({ secret: uuid.v4()});
});
app.listen(PORT, function() {
    console.log(`Example app listening on port ${PORT}!`)
});