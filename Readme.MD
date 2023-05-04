# Ticket Event Shop

<div align="center">
  <p>
    Ticket Event Shop
  </p>

  <p align="center">
    Group Project 3 Unit 2
    <br />
    <a href="https://app.swaggerhub.com/apis/ropel12/tes/1.0.0"><strong>| Open API Documentation |</strong></a>
    <br />
    <br />
  </p>
</div>

## 🧑‍💻 About the Project

<p align="justify">Event Ticket Shop is a web application-based online market place for selling event ticket. Users can also discuss with other users about the event to be held. Here you can sell paid or free tickets. We implement midtrans for payment gateway</p>


## 🛠 Technology Stack
<div align="center">

![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)
![Visual Studio Code](https://img.shields.io/badge/Visual%20Studio%20Code-0078d7.svg?style=for-the-badge&logo=visual-studio-code&logoColor=white)
![MySQL](https://img.shields.io/badge/mysql-%2300f.svg?style=for-the-badge&logo=mysql&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Google Cloud](https://img.shields.io/badge/GoogleCloud-%234285F4.svg?style=for-the-badge&logo=google-cloud&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)
![Cloudflare](https://img.shields.io/badge/Cloudflare-F38020?style=for-the-badge&logo=Cloudflare&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)
![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)
![Midtrans](https://img.shields.io/badge/Midtrans-FF6C37?style=for-the-badge&logo=midtrans&logoColor=white)

  </div>

# 🔗 ERD

<div align="center">
<img src="docs/erd.png" width="800" height="800">
  </div>

  # ⚡ Features

<details>
  <summary>🎫 Auth</summary>
  
| Method      | Endpoint            | Params      |q-Params            | JWT Token   | Function                                |
| ----------- | ------------------- | ----------- |--------------------| ----------- | --------------------------------------- |
| POST        | /register           | -           |-                   | NO         | Register a new User                |
| POST        | /login      | -           |-                   | NO         | Login to the system        |
  
</details>

<details>
  <summary>🙍‍♂️ Users</summary>
  
| Method      | Endpoint            | Params      |q-Params            | JWT Token   | Function                                |
| ----------- | ------------------- | ----------- |--------------------| ----------- | --------------------------------------- |
| GET        | /users           | -           |-                   | YES         | Show profile                |
| PUT        | /users      | -           |-                   | YES         | Update profile data        |
| DELETE        | /users      | -           |-                   | YES         | Delete user data        |
| GET        | /users/events      | -           |-                   | YES         | Show list user event's        |
| GET        | /users/history      | -           |-                   | YES         | Show history        |
| GET        | /users/transaction      | -           |-                   | YES         | show data of user ticket's        |
</details>

<details>
  <summary>🕺Events</summary>
  
| Method      | Endpoint            | Params      |q-Params            | JWT Token   | Function                                |
| ----------- | ------------------- | ----------- |--------------------| ----------- | --------------------------------------- |
| PUT        | /events      | -           |-                   | YES         | Update event data        |
| POST        | /events      | -           |-                   | YES         | Post new event        |
| GET        | /events           | -           |-                   | YES         | Show all events                |
| DELETE        | /events/{id_event}     | id_event           |-                   | YES         | Delete event by ID        |
| GET        | /events/{id_event}      | id_event           |-                   | YES         | Show detail event        |
| POST        | /events/participant      | -           |-                   | YES         | Join the event        |
</details>

<details>
  <summary>👋 Events Participants</summary>
  
| Method      | Endpoint            | Params      |q-Params            | JWT Token   | Function                                |
| ----------- | ------------------- | ----------- |--------------------| ----------- | --------------------------------------- |
| POST        | /comments           | -           |-                   | YES         | Create a new comment                |
  
</details>

<details>
  <summary>🛍️ Transactions</summary>
  
| Method      | Endpoint            | Params      |q-Params            | JWT Token   | Function                                |
| ----------- | ------------------- | ----------- |--------------------| ----------- | --------------------------------------- |
| POST        | /transaction/checkout           | -           |-                   | YES         | User checkout products                |
| POST        | /transactions/cart      | -           |-                   | YES         | User add products to cart        |
| GET        | /transactions/cart      | -           |-                   | YES         | User show their products in cart        |
| GET        | /transactions/{invoice}      | -           | invoice                   | YES         | Show User invoices        |
  
</details>