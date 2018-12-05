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
    "name": "",
    "task": {},
    "tasks": [],
    "point": 0,
    "User": {},
};
var Task = require('./models/task');
//add JSON request body parsing middleware
app.use(express.json());
//add the request logging middleware
app.use(morgan("dev"));

// GET /tasks/:familyRoom
// If a user is authenticated(member/admin of this family),
// show the public to do list with all the in-progress tasks and undo tasks. (called to show the public task list)
app.get('/tasks/:roomname', (req, res, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");
    if (userJSON) {
        var roomname = req.params.roomname;
        Task.find({"familyRoomName": roomname}).exec((err, tasks) => {
            if (err) {
                res.statusCode = 500;
                res.send("Error while finding tasks");
                return;
            }
            console.log(tasks)
            res.statusCode = 200;
            res.setHeader('Content-Type', 'application/json');
            res.end(JSON.stringify(tasks));
            return;
        });
    } else {
        res.statusCode = 401;
        res.send("no X-User header in the request");
        return;
    }
});


// POST /tasks/:roomname
// If a user is authenticated(admin), post the new task in his/her private task list and the public task list. (called when an admin clicks create task in his/her private task page)
app.post("/tasks/:roomname", (req, res, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");
    // Check whether user is member or admin
    if (userJSON) {
        let user = JSON.parse(userJSON);
        var roomname = req.params.roomname;
        var task;
        console.log(user);
        console.log(req.body)
        console.log("req.body.description: "+ req.body.description);
        console.log(roomname);
        // should return 400 if req.body.description is empty
        if (req.body.description === "" || req.body.description === undefined) {
            res.statusCode = 400;
            res.end("description is empty");
            return;
        }
        switch (user.personrole) {
            case "Admin":
                // If a user is authenticated(admin), post the new task in his/her private task list and the public task list.
                // (called when an admin clicks create task in his/her private task page)
                var task = new Task ({
                    description: req.body.description,
                    familyRoomName: user.roomname
                });
                // Create new task and push to task table
                Task.addTask(task, (err, task) => {
                    if (err) {
                        console.log("err: "+ err);
                        res.statusCode = 400;
                        res.end("duplicated document");
                        return;
                    }
                });
                Task.find({"familyRoomName": roomname}).exec((err, tasks) =>{
                    if (err) {
                        res.statusCode = 500;
                        res.send("Error while finding tasks");
                        return;
                    }
                    buffer["tasks"] = tasks;
                });
                // message queue
                buffer["name"] = "task-new";
                buffer["task"] = task;
                // Push to message queue
                taskchannel.sendToQueue(
                    "taskQueue",
                    Buffer.from(JSON.stringify(buffer)),
                    {persistent: true}
                );
                // Return 201 and application/json
                res.statusCode = 201;
                res.setHeader('Content-Type', 'application/json');
                res.end(JSON.stringify(task));
                return;
            break;
            default:
                res.statusCode = 401;
                res.send("not proper roles in the request");  
                return;  
            break;
        }
    } else {
        res.statusCode = 401;
        res.send("no X-User header in the request");
        return;
    }
});

// PATCH /tasks/:taskid
//  If a user is authenticated(admin), update the task in his/her private task list and the public task list. 
// (called when an admin clicks update in his/her private task page)
app.patch("/tasks/:id", (req, res, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");
    // Check whether user is member or admin
    if (userJSON) {
        let user = JSON.parse(userJSON);
        var id = req.params.id;
        console.log("Debug in patch //tasks/:id " + user);
        // if a user is not admin, the error should return 401
        if (user.personrole != "Admin") {
            res.statusCode = 401;
            res.send("not proper role in the request");
            return;
        }
        // should return 400 if req.body.description is empty
        if (req.body.description == "") {
            res.statusCode = 400;
            res.send("description is empty");
            return;
        }
        // Update the task and return 200
        Task.findOne({"_id": id}).exec((err, task) => {
            if (err) {
                res.statusCode = 500;
                res.send("Error on execute finding task");
                return;
            }
            // Push to message queue
            task.description = req.body.description;
            buffer["name"] = "task-edit";
            buffer["task"] = task;
            taskchannel.sendToQueue(
                "taskQueue",
                Buffer.from(JSON.stringify(buffer)),
                {persistent: true}
            );
            res.statusCode = 200;
            res.setHeader('Content-Type', 'application/json');
            res.end(JSON.stringify(task));
        });
    } else {
        // If not return 401.
        res.statusCode = 401;
        res.send("no X-User header in the request");
        return;
    }
});


// DELETE /tasks/:taskId
// If a user is authenticated(admin), delete the task from his/her private task list
// and the public task list.
app.delete("/tasks/:id", (req, res, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");
    // Check whether user is member or admin
    if (userJSON) {
        let user = JSON.parse(userJSON);
        var id = req.params.id;
        // If a user is authenticated(admin), 
        // delete the task from his/her private task list and the public task list.
        // If not return 401.
        if (user.personrole != "Admin") {
            res.statusCode = 401;
            res.send("not proper role in the request");
            return;
        }
        // Delete the task and return 200
        Task.deleteOne({"_id": id}).exec((err, task) => {
            if (err) {
                res.statusCode = 500;
                res.send("Error on execute finding family");
                return;
            }
            // Push to message queue
            buffer["name"] = "task-delete";
            buffer["task"] = task;
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
});

// POST /tasks/progress/:taskid
// If a user is authenticated(member), add task for him/her private task list
app.post('/tasks/progress/:id', (req, res, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");

    if (userJSON) {
        let user = JSON.parse(userJSON);
        var id = req.params.id;
        console.log("Debug: post /tasks/progress/:id " + user);
        console.log("Debug: post /tasks/progress/:id " + user.personrole);
        console.log("Debug: post /tasks/progress/:id " + id);
        if (user.personrole != "Member") {
            res.statusCode = 401;
            res.send("not proper role in the request");
            return;
        }
        Task.findByIdAndUpdate({_id: id}, {$set: {"userID": user.id, "isProgress": true}}, {$push: {"userID": user.id, "isProgress": true}}).exec((err, task) => {
            if (err) {
                res.statusCode = 500;
                res.send("Error on execute finding family");
                return;
            }
            console.log("Debug: post /tasks/progress/:id " + user);
            console.log("Debug: post /tasks/progress/:id task: " + task);
            // push to message queue
            buffer["name"] = "task-progress";
            buffer["user"] = user;
            buffer["task"] = task;
            console.log("Debug: post /tasks/progress/:id after task: " + task);
            taskchannel.sendToQueue(
                "taskQueue",
                Buffer.from(JSON.stringify(buffer)),
                {persistent: true}
            );
            res.statusCode = 200;
            res.setHeader('Content-Type', 'application/json');
            res.send(JSON.stringify(task));
        });
    } else {
        // If not return 401.
        res.statusCode = 401;
        res.send("no X-User header in the request");
        return;
    }
});

// POST /tasks/done/:taskid
// If a user is authenticated(member) and finished his/her task,
// delete task for him/her private task list and update user's point
app.post('/tasks/done/:id', (req, res, next) => {
    // Check whether user is authenticated using X-user header
    let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        var id = req.params.id;
        console.log("Debug: post /tasks/done/:id " + user);
        if (user.personrole != "Member") {
            res.statusCode = 401;
            res.send("not proper role in the request");
            return;
        }
        Task.findOneAndDelete({"_id": id}).exec((err, task) => {
            if (err) {
                res.statusCode = 500;
                res.send("Error on execute finding family");
                return;
            }
            // update the task
            task.userID = user.id;
            task.isProgress = true;
            // push to message queue
            buffer["name"] = "task-done";
            buffer["user"] = user;
            buffer["task"] = task;
            buffer["point"] = task.point; // increment point!
            taskchannel.sendToQueue(
                "taskQueue",
                Buffer.from(JSON.stringify(buffer)),
                {persistent: true}
            );
            res.statusCode = 200;
            res.send("Done!");
            return;
        });
    } else {
        // If not return 401.
        res.statusCode = 401;
        res.send("no X-User header in the request");
        return;
    }
});

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