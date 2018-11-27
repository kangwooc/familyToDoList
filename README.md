# FamilyToDo App - ZCO

## Members And Roles
<ul>
    <li>Kangwoo Choi - Mainly working on Microservice</li>
    <li>Juan Oh - Front-end / Microservice / Leaderboard </li>
    <li>Tina Zhuang - Authorization / Gateway / Microservice</li>
</ul>

## Description

Family to do list is a web application that allows parents to add tasks to other family members to do. The premise of this app is to organize the tasks that need to be finished in a family unit.

We want to build a handy app that organizes tasks in a family unit. The intention is to encourage family members(kids) to complete tasks with a sense of accomplishment through the leveling up system in the app. The app allows admin(parents) in a family to post tasks for family members to do. Members will level up as they complete more tasks. With the visualization, we are hoping to motivate family members to complete their tasks.

## Overview

![Alt text](/img/Overview.jpeg?raw=true "Overview of project")

## Priority for This Project

| Priority | User | Description | Technical Implementation Strategy |
| ------------- | ------------- | ------------- | ------------- |
| P0 | User | I want to sign up/sign in. | We will implement userhandler for authorization in gateway and save the user information in <strong>MySQL</strong>. |
| P1 | Admin | I want to create a family room. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and add the new room to FamilyRoom table in <strong>MySQL</strong>. |
| P2 | Admin | I want to add tasks in the to-do list. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and add the new task in <strong>MongoDB</strong>. |
| P3 | Member | I want to join a family room. | We will check for user authorization which is saved in user table in <strong>MySQL</strong> and implement requesthandler that allows user to send request. Finally, we will add the user as a member of the family in FamilyRoom table in <strong>MySQL</strong> if request is approved.|
| P4 | Admin | I want to receive notification when a user wants to join my family group. | We will implement requesthandler to allow room owner to receive and approve/disapprove request. |
| P5 | Member | I want to work on a certain task. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and change the task status saved in <strong>MongoDB</strong>. |
| P6 | Admin | I want to edit/delete tasks in the to-do list. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and update/delete the task status saved in <strong>MongoDB</strong>. |
| P7 | Member | I want to cancel a task that I am currently working on. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and update/delete the task status saved in <strong>MongoDB</strong>. |
| P8 | Admin | I want to delete a member from my family group. | We will check for user permission level which is saved in user table in <strong>MySQL</strong> and delete the member from the FamilyRoom table in <strong>MySQL</strong>. |
| P9 | User | I want to check top 5 users among all users. | We will check for user authorization which is saved in user table in <strong>MySQL</strong>. Also, we add the function handler for updating the ranking. |

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

| Request | type |
| ------------- | ------------- |
| req_id  | int |
| member_id | int |
| family_id | int |
| isPending | bool |


+ MongoDB for Task Microservice

| Task | type |
| ------------- | ------------- |
| task_id  | int |
| description | varchar |
| point | int |
| isProgress | bool |
| isDone | bool |

