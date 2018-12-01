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
var taskchannel;
// message queue struct
var buffer = {
    "type": "",
    "task": {},
    "tasks": []
};
var Task = require('./models/task');
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
// GET /tasks/:familyID
// If a user is authenticated(member/admin of this family),
// show the public to do list with all the in-progress tasks and undo tasks. (called to show the public task list)
app.get('/tasks/:id', (req, res, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        Task.find({"familyID": id}).lean().exec((err, tasks) => {
            if (err) {
                res.statusCode = 500;
                res.send("Error while finding tasks");
                return;
            }
            res.statusCode = 200;
            res.end(JSON.stringify(tasks));
            return;
        });
    } else {
        res.statusCode = 401;
        res.send("no X-User header in the request");
        return;
    }
});

// POST /tasks/:familyID
// If a user is authenticated(admin), post the new task in his/her private task list and the public task list. (called when an admin clicks create task in his/her private task page)
app.post("/tasks/:id", (res, req, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");
    // Check whether user is member or admin
    if (userJSON) {
        let user = JSON.parse(userJSON);
        var task;
        switch (user.role) {
            case "Admin":
                // If a user is authenticated(admin), post the new task in his/her private task list and the public task list.
                // (called when an admin clicks create task in his/her private task page)
                var task = new Task ({
                    description: req.Body.description,
                    point: req.Body.point,
                    familyID: id
                });
                // Create new task and push to task table
                task.save((err) => {
                    if (err) {
                        console.log(err);
                    }
                });
                // message queue
                buffer["type"] = "task-new";
                buffer["tasks"] = tasks;
                // Push to message queue
                taskchannel.sendToQueue(
                    "taskQueue",
                    Buffer.from(JSON.stringify(buffer)),
                    {persistent: true}
                );
            break;
            default:
                res.statusCode = 401;
                res.send("not proper roles in the request");  
                return;  
            break;
        }
        // Return 201 and application/json
        res.statusCode = 201;
        res.setHeader('Content-Type', 'application/json');
        res.end(JSON.stringify(task));
    } else {
        res.statusCode = 401;
        res.send("no X-User header in the request");
        return;
    }
});

// PATCH /tasks/:taskid
//   + If a user is authenticated(admin), update the task in his/her private task list and the public task list. (called when an admin clicks update in his/her private task page)
app.patch("/tasks/:id", (res, req, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");
    // Check whether user is member or admin
    if (userJSON) {
        let user = JSON.parse(userJSON);
        if (user.role != "Admin") {
            res.statusCode = 401;
            res.send("not proper role in the request");
            return;
        }
        // Update the task and return 200
        Task.findOne({"_id": id}).exec((err, task) => {
            if (err) {
                res.statusCode = 500;
                res.send("Error on execute finding family");
                return;
            }
            // Push to message queue
            task.description = req.Body.description;
            buffer["type"] = "task-edit";
            buffer["task"] = task;
            taskchannel.sendToQueue(
                "taskQueue",
                Buffer.from(JSON.stringify(buffer)),
                {persistent: true}
            );
        });
    } else {
        // If not return 401.
        res.statusCode = 401;
        res.send("no X-User header in the request");
        return;
    } 
});


// + DELETE /tasks/:taskId
//    + If a user is authenticated(admin), delete the task from his/her private task list and the public task list.
app.delete("/task/:id", (res, req, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");
    // Check whether user is member or admin
    if (userJSON) {
        let user = JSON.parse(userJSON);
        // If a user is authenticated(admin), delete the task from his/her private task list and the public task list.
        // If not return 401.
        if (user.role != "Admin") {
            res.statusCode = 401;
            res.send("not proper role in the request");
            return;
        }
        // Update the task and return 200
        Task.delete({"_id": id}).exec((err, task) => {
            if (err) {
                res.statusCode = 500;
                res.send("Error on execute finding family");
                return;
            }
            // Push to message queue
            buffer["type"] = "task-delete";
            buffer["memberID"]
            taskchannel.sendToQueue(
                "taskQueue",
                Buffer.from(JSON.stringify(buffer)),
                {persistent: true}
            );
            res.statusCode = 200;
            res.send("successfully delete!");
        });
    } else {
        // If not return 401.
        res.statusCode = 401;
        res.send("no X-User header in the request");
        return;
    }
    // Delete task
    // If it doesn’t have in task, return 500.
    // Push to message queue
})



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
        ch.assertQueue("taskQueue", {durable: true});
        taskchannel = ch;
        // start the server listening on host:port
        app.listen(portNum, host, () => {
            console.log(`server is listening at ${addr}`);
        });
    });
});