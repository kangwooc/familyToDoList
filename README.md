# FamilyToDo App - ZCO

## Members And Roles
<ul>
    <li>Kangwoo Choi - Mainly working on Microservice</li>
    <li>Juan Oh - Front-end / Microservice</li>
    <li>Tina Zhuang - Authorization / Gateway / Microservice</li>
</ul>

## Description

Family to do list is a web application that allows parents to add tasks to other family members to do. The premise of this app is to organize the tasks that need to be finished in a family unit.

We want to build a handy app that organizes tasks in a family unit. The intention is to encourage family members(kids) to complete tasks with a sense of accomplishment through the leveling up system in the app. The app allows admin(parents) in a family to post tasks for family members to do. Members will level up as they complete more tasks. With the visualization, we are hoping to motivate family members to complete their tasks.

## Overview

![Alt text](/img/Overview.jpeg?raw=true "Overview of diagram")

## Priority for This Project

| Priority | User | Description |
| ------------- | ------------- | ------------- |
| P0 | User | I want to sign up/sign in |
| P1 | Admin | I want to create a family room |
| P2 | Admin | I want to add tasks in the to-do list |
| P3 | Member | I want to join a family room |
| P4 | Admin | I want to receive notification when a user wants to join my family group |
| P5 | Member | I want to work on a certain task |
| P6 | Admin | I want to edit/delete tasks in the to-do list |
| P7 | Member | I want to cancel a task that I am currently working on |
| P8 | Admin | I want to delete a member from my family group |

## Different Point (Methods for microservice)
+ GET /:id
    + If a user is authenticated(member/admin of this family), show the public to do list with all the in-progress tasks and undo tasks. (called to show the public task list)
+ POST /:id
  + If a user is authenticated(member), post the task that he/she started on his/her private task list. The public to do list should indicate that this member started on the chosen task. (called when a member clicks on a task in public task list)
   + If a user is authenticated(admin), post the new task in his/her private task list and the public task list. (called when an admin clicks create task in his/her private task page)
+ PATCH /:id
  + If a user is authenticated(admin), update the task in his/her private task list and the public task list. (called when an admin clicks update in his/her private task page)
+ DELETE /member/:id
   + If a user is authenticated(admin), delete the user from the family group. (called when an admin clicks on delete in his/her private task page)
+ DELETE /:taskId
   + If a user is authenticated(admin), delete the task from his/her private task list and the public task list.
+ POST /request
    + If a user is a family group admin, notify him/her about the request and let him/her approve/disapprove on the request.
## Appendix

+ MySQL for User Information
  
| FamilyRoom | type |
| ------------- | ------------- |
| family_id  | int |
| name | varchar |
| member_id  | int |

| Member | type |
| ------------- | ------------- |
| member_id  | int |
| firstName | varchar |
| lastName | varchar |
| userName | varchar |
| password | varchar |
| role_id | int |
| point | int |

| Role | type |
| ------------- | ------------- |
| role_id  | int |
| type | varchar |
| description | varchar |

+ MongoDB for Task Microservice

| Task | type |
| ------------- | ------------- |
| task_id  | int |
| description | varchar |
| point | int |
| isProgress | bool |
| isDone | bool |

