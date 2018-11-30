'use strict';

const express = require("express");
const morgan = require("morgan");
// create a new express application
const app = express();
// get ADDR environment variable,
// defaulting to ":80"
const addr = process.env.ADDR || ":80";
//split host and port using destructuring
const [host, port] = addr.split(":");
const portNum = parseInt(port);

var amqp = require('amqplib/callback_api');

// Import the mongoose module
var mongoose = require('mongoose');

// Get environment variable
const mongoaddr = process.env.MONGOADDR || ":27017";
const rabbitaddr = process.env.RABBITADDR || ":5672";

// Set up default mongoose connection
var mongoDB = 'mongodb://'+mongoaddr+'/userdb';
mongoose.connect(mongoDB, { useCreateIndex: true, useNewUrlParser: true });
// Get Mongoose to use the global promise library
mongoose.Promise = global.Promise;

// Get the default connection
var db = mongoose.connection;
// Bind connection to error event (to get notification of connection errors)
db.on('error', console.error.bind(console, 'MongoDB connection error:'));
db.once('open', function callback () {
    console.log('Conntected To Mongo Database');
});
//add JSON request body parsing middleware
app.use(express.json());
//add the request logging middleware
app.use(morgan("dev"));

// compare the json object
// Reference:https://stackoverflow.com/questions/23826209/javascript-compare-the-structure-of-two-json-objects-while-ignoring-their-value
function compareObjects(obj1, obj2){
    var equal = true;
    for (var i in obj1) {
        if (!obj2.hasOwnProperty(i)) {
            equal = false;
            break;
        }
    }
    return equal;
}

// arrayRemove is the function which removes the value in the array
function arrayRemove(arr, value) {
    console.log(arr);
    return arr.filter(function (ele) {
        return ele !== value;
    });
}

app.use('tasks/:id')




var rabbiturl = 'amqp://' + rabbitaddr;
amqp.connect(rabbiturl, function (err, conn) {
    if (err) {
        console.log("Failed to connect to Rabbit Instance from API Server.");
        console.log(err);
        process.exit(1);
    }
    conn.createChannel((err, ch) => {
        if (err) {
            console.log("Failed to connect to create channel from API Server.");
            process.exit(1);
        }
        ch.assertQueue("msgQueue", {durable: true});
        megchannel = ch;
        // start the server listening on host:port
        app.listen(portNum, host, () => {
            console.log(`server is listening at ${addr}`);
        });
    });
});