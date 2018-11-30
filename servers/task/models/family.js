'use strict';

var mongoose = require('mongoose');
var Schema = mongoose.Schema;

// Set the schema of task
// description: a short description of the task
// members: if the channel is private, an array containing the profiles of all channel members
// point: set the points for each task
// isProgress: check whether the task is progressing
// isDone: check whether the task finishes
// Reference: https://scotch.io/tutorials/using-mongoosejs-in-node-js-and-mongodb-applications
// https://stackoverflow.com/questions/17899750/how-can-i-generate-an-objectid-with-mongoose
// https://stackoverflow.com/questions/10006218/which-schematype-in-mongoose-is-best-for-timestamp
var familySchema = new Schema({
    familyName: String,
    familyMembers: [{type: Object}],
    creator: Object,
    createdAt: Date,
    taskList: [{type: Object}]
});
// create a model for our family
var Family = mongoose.model('Family', familySchema);
// make this available to our users in our Node applications
module.exports = Family;