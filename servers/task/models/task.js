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
    user: Object
});
// create a model for our task
var Task = mongoose.model('Task', taskSchema);
// make this available to our users in our Node applications
module.exports = Task;