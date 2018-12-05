'use strict';
//Import the mongoose module
var mongoose = require('mongoose');
var Schema = mongoose.Schema;

// Set the schema of task
// description: a short description of the task
// point: set the points for each task
// isProgress: check whether the task is progressing
// isDone: check whether the task finishes
// Reference: https://scotch.io/tutorials/using-mongoosejs-in-node-js-and-mongodb-applications
// https://stackoverflow.com/questions/17899750/how-can-i-generate-an-objectid-with-mongoose
// https://stackoverflow.com/questions/10006218/which-schematype-in-mongoose-is-best-for-timestamp
var taskSchema = new Schema({
    description: {type: String, unique: true},
    point: {type: Number, default: 5},
    isProgress: {type: Boolean, default: false},
    familyID: Number,
    familyRoomName: String,
    userID: Number
});
// function for update task
// https://stackoverflow.com/questions/16882938/how-to-check-if-that-data-already-exist-in-the-database-during-update-mongoose
taskSchema.statics.addTask = function(task, cb) {
    Task.find({description : task.description}).exec(function(err, docs) {
        if (docs.length) {
            cb('documents exists already', null);
        } else {
            task.save(function (err) {
                cb(err, docs);
            });
        }
    });
}

// create a model for our task
var Task = mongoose.model('Task', taskSchema);
// make this available to our users in our Node applications
module.exports = Task;