# FamilyToDo App - ZCO

## Members And Roles
<ul>
    <li>Kangwoo Choi - Mainly working on Microservice / Websocket</li>
    <li>Juan Oh - Front-end / Authorization / Websocket </li>
    <li>Tina Zhuang - Authorization / Gateway / Microservice</li>
</ul>

## Description

Family to do list is a web application that allows parents to add tasks for other family members to do. The premise of this app is to organize the tasks that need to be finished in a family unit.

We want to build a handy app that organizes tasks in a family unit. The intention is to encourage family members(kids) to complete tasks with a sense of accomplishment. The app allows admin(parents) in a family to post tasks for family members to do. We are hoping to motivate family members to complete their tasks.

## Overview

![Alt text](/img/final.jpeg?raw=true "Overview of project")

## Priority for This Project

| Priority | User | Description | Technical Implementation Strategy |
| ------------- | ------------- | ------------- | ------------- |
| P0 | User | I want to sign up/sign in. | We will implement userhandler for authorization in gateway and save the user information in <strong>MySQL</strong>. |
| P1 | Admin | I want to create a family room. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and add the new room to FamilyRoom table in <strong>MySQL</strong>. |
| P2 | Admin | I want to add tasks in the to-do list. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and add the new task in <strong>MongoDB</strong>. |
| P3 | Member | I want to join a family room. | We will check for user authorization which is saved in user table in <strong>MySQL</strong> and implement requesthandler that allows user to send request. Finally, we will add the user as a member of the family in FamilyRoom table in <strong>MySQL</strong> if request is approved.|
| P4 | Admin | I want to receive notification when a user wants to join my family group. | We will implement requesthandler to allow room owner to receive and approve/disapprove request. |
| P5 | Member | I want to work on a certain task. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and change the task status saved in <strong>MongoDB</strong>. |
| P6 | Member | I want to mark a task down (delete a task) when I am finished with it. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and delete the task status saved in <strong>MongoDB</strong>. |
| P7 | Admin | I want to delete a member from my family group. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and delete the member from the FamilyRoom table in <strong>MySQL</strong>. |

## Appendix

+ MySQL for User Information
   
| FamilyRoom | type |
| ------------- | ------------- |
| family_id  | int |
| name | varchar |


| User | type |
| ------------- | ------------- |
| member_id  | int |
| passhash  | binary |
| roomname | varchar |
| photourl | varchar |
| firstName | varchar |
| lastName | varchar |
| userName | varchar |
| personrole | varchar |
| score | int |


+ MongoDB for Task Microservice

| Task | type |
| ------------- | ------------- |
| task_id  | int |
| description | varchar |
| point | int |
| isProgress | bool |
| userID | int |
| familyRoomName | varchar |

+ Endpoints

<li>	GET	/tasks/:roomname</li>
<ul>
<li>    If a user is a member/admin of the family, shows the public to do list</li>
<li>        500, 200, 401</li>
</ul>
<li>	POST	/tasks/:roomname
<ul>
<li>	If a user is an admin of the family, post new task in task list</li>
<li>	400, 500, 201, 401</li>
</ul>
<li>	PATCH	/tasks/:id
<ul>
<li>	If a user is an admin of the family, update the task in task list</li>
<li>	400, 500, 200, 401</li>
</ul>
<li>	DELETE	/tasks/:id
<ul>
<li>	If a user is an admin of the family, delete the task in task list</li>
<li>	500, 200, 401</li>
</ul>
<li>	POST	/tasks/progress/:id
<ul>
<li>	If a user is a member of the family, add task to his/her private task list</li>
<li>	401, 500, 200</li>
</ul>
<li>	POST	/tasks/done/:id
<ul>
<li>	If a user is a member of the family, delete task from his/her private task list</li>
<li>	401, 500, 200</li>
